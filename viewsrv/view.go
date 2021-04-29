package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type viewResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type viewerData struct {
	IIIFURI   string `json:"iiif"`
	RightsURI string `json:"rights"`
	StartPage int    `json:"page"`
	PagePIDs  string `json:"page_pids"`
}

// viewHandler takes the initial viewer request and determines what type of resource it is and
// hands the rendering off to the appropriate handler - or returns 404
func viewHandler(c *gin.Context) {
	srcPID := c.Param("pid")
	unitID := c.Query("unit")
	log.Printf("INFO: Check if %s is an image...", srcPID)
	iiifManURL, iiifErr := getIIIFManifestURL(srcPID, unitID)
	if iiifErr == nil {
		log.Printf("INFO: render %s as image", srcPID)
		viewImage(c, iiifManURL)
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

	c.String(http.StatusNotFound, "not found")
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
		c.String(http.StatusNotFound, "not found")
		return
	}

	// https://iiif.lib.virginia.edu/iiif/tsm:2804870/full/!200,200/0/default.jpg
	re := regexp.MustCompile(`^.*iiif/|/full.*$`) // strip all but pid
	pids := make([]string, 0)
	for _, c := range manifest.Sequences[0].Canvases {
		pid := re.ReplaceAllString(c.Thumbnail, "")
		pids = append(pids, pid)
	}

	data := viewerData{RightsURI: config.rightsURL, IIIFURI: iiifURL, StartPage: page, PagePIDs: strings.Join(pids, ",")}
	out := viewResponse{Type: "iiif", Data: data}
	c.JSON(http.StatusOK, out)
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

	out := viewResponse{Type: "wsls", Data: wslsData}
	c.JSON(http.StatusOK, out)
}

// getIIIFManifestURL retrieves the cached IIIF manifest for an item. If a unit is specified,
// the manifest just needs to exist; cache does not matter as the manifest will be generated on the fly
func getIIIFManifestURL(pid string, unit string) (string, error) {
	log.Printf("INFO: check if %s, unitID [%s] is a candidate for IIIF metadata...", pid, unit)
	url := fmt.Sprintf("%s/pid/%s/exist", config.iiifURL, pid)
	resp, err := getAPIResponse(url)
	if err != nil {
		return "", err
	}
	var parsed struct {
		Exists bool   `json:"exists"`
		Cached bool   `json:"cached"`
		URL    string `json:"url"`
	}
	err = json.Unmarshal([]byte(resp), &parsed)
	if err != nil {
		return "", err
	}

	// when unit is present, dont care if it is cached or not, just care if the metadata exists
	if unit != "" {
		log.Printf("Unit %s present in request, not using IIIF cache", unit)
		if parsed.Exists {
			iiifURL := fmt.Sprintf("%s/pid/%s?unit=%s", config.iiifURL, pid, unit)
			log.Printf("INFO: IIIF manifest available at %s", iiifURL)
			return iiifURL, nil
		}
		return "", errors.New("manifest not found")
	}
	if !parsed.Exists || !parsed.Cached {
		return "", errors.New("manifest not found")
	}
	log.Printf("INFO: IIIF manifest cached at %s", parsed.URL)
	return parsed.URL, nil
}
