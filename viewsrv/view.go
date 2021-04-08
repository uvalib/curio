package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type viewerData struct {
	IIIFURI   string
	RightsURI string
	StartPage int
	PagePIDs  string
}

// viewHandler takes the initial viewer request and determines what type of resource it is and
// hands the rendering off to the appropriate handler - or returns 404
func viewHandler(c *gin.Context) {
	srcPID := c.Param("pid")
	if isIiifCandidate(srcPID) {
		unitID := c.Query("unit")
		iiifURL := fmt.Sprintf("%s/pid/%s", config.iiifURL, srcPID)
		if unitID != "" {
			iiifURL = fmt.Sprintf("%s?unit=%s", iiifURL, unitID)
		}

		log.Printf("INFO: render %s as image", srcPID)
		viewImage(c, iiifURL)
		return
	}

	// not an image; try Apollo for WSLS...
	log.Printf("INFO: %s is not image; check WSLS", srcPID)
	wslsData, err := getApolloWSLSMetadata(srcPID)
	if err == nil {
		log.Printf("INFO: render %s as WSLS", srcPID)
		viewWSLS(c, wslsData)
		return
	}

	// Nope; fail
	c.HTML(http.StatusOK, "not_available.html", nil)
}

// viewImage displays a series of images in the universalViewer
func viewImage(c *gin.Context, iiifURL string) {
	log.Printf("INFO: using iiif manifest %s", iiifURL)
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	manifestStr, err := getAPIResponse(iiifURL)

	var manifest struct {
		Sequences []struct {
			Canvases []struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"canvases"`
		} `json:"sequences"`
	}
	if jErr := json.Unmarshal([]byte(manifestStr), &manifest); jErr != nil {
		log.Printf("Unmarshal manifest failed: %s", jErr.Error())
		c.HTML(http.StatusOK, "not_available.html", nil)
		return
	}

	// https://iiif.lib.virginia.edu/iiif/tsm:2804870/full/!200,200/0/default.jpg
	re := regexp.MustCompile(`^.*iiif/|/full.*$`) // strip all but pid
	pids := make([]string, 0)
	for _, c := range manifest.Sequences[0].Canvases {
		pid := re.ReplaceAllString(c.Thumbnail, "")
		pids = append(pids, pid)
	}

	data := viewerData{RightsURI: config.rightsURL, IIIFURI: iiifURL, StartPage: page - 1, PagePIDs: strings.Join(pids, ",")}
	c.HTML(http.StatusOK, "image_view.html", data)
}

// viewWSLS renders a custom view of WSLS content that includes video clips, transcripts and a poster
func viewWSLS(c *gin.Context, wslsData *wslsMetadata) {
	if wslsData.HasVideo {
		// POSTER: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-poster.jpg
		// VIDEO (webm): http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.mp4
		wslsData.VideoURL = fmt.Sprintf("%s/%s/%s.mp4", config.wslsURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.PosterURL = fmt.Sprintf("%s/%s/%s-poster.jpg", config.wslsURL, wslsData.WSLSID, wslsData.WSLSID)
	}

	if wslsData.HasScript {
		// PDF: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.pdf
		// Thumb: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-script-thumbnail.jpg
		// Transcript: http://fedora01.lib.virginia.edu/wsls/0003_1/0003_1.txt
		wslsData.PDFURL = fmt.Sprintf("%s/%s/%s.pdf", config.wslsURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.PDFThumbURL = fmt.Sprintf("%s/%s/%s-script-thumbnail.jpg", config.wslsURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.TranscriptURL = fmt.Sprintf("%s/%s/%s.txt", config.wslsURL, wslsData.WSLSID, wslsData.WSLSID)
	}

	c.HTML(http.StatusOK, "wsls_view.html", wslsData)
}

// isIiifCandidate will call the IIIF manifest service exist endpoint to determine of the pid has IIIF data
func isIiifCandidate(pid string) bool {
	log.Printf("INFO: check if %s is a candidate for IIIF metadata...", pid)
	url := fmt.Sprintf("%s/pid/%s/exist", config.iiifURL, pid)
	resp, err := getAPIResponse(url)
	if err != nil {
		return false
	}
	var parsed struct {
		Exists bool   `json:"exists"`
		Cached bool   `json:"cached"`
		URL    string `json:"url"`
	}
	err = json.Unmarshal([]byte(resp), &parsed)
	if err != nil {
		log.Printf("ERROR: Unable to parse response from %s: %s", url, err.Error())
		return false
	}
	log.Printf("IIIF exist response %+v", parsed)
	return parsed.Exists
}
