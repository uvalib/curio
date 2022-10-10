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

	// Check Archivematica
	log.Printf("INFO: %s is not WSLS; Checking Archivematica", srcPID)
	archivematicaData, err := getArchivematicaData(srcPID)
	if err == nil {
		log.Printf("INFO: render %s as Archivematica", srcPID)
		c.JSON(http.StatusOK, archivematicaData)
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

// ArchivematicaS3Node format for archivematica in s3
// Can be arbitrarily nested
type ArchivematicaS3Node struct {
	Name       string `json:"name,omitempty"`
	ID         string `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	Format     string `json:"format,omitempty"`
	SourceURL  string `json:"source_url,omitempty"`
	DisplayURL string `json:"display_url,omitempty"`
	MimeType   string `json:"mime_types,omitempty"`
	View       string `json:"view,omitempty"`

	Entries []ArchivematicaS3Node `json:"entries,omitempty"`
}

// TableNode format for PrimeVue TreeTable
// https://www.primefaces.org/primevue/treetable
type TableNode struct {
	Key        string      `json:"key"`
	Data       ColumnData  `json:"data"`
	Children   []TableNode `json:"children,omitempty"`
	StyleClass string      `json:"styleClass,omitempty"`
	//Leaf bool `json:"leaf"`
}

// ColumnData contains data to be displayed
type ColumnData struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Format string `json:"format"`
	Icon   string `json:"icon"`
	URL    string `json:"url"`
}

func getArchivematicaData(pid string) (viewResponse, error) {
	ArchivematicaResponse := viewResponse{Type: "archivematica"}

	// S3 retrieval
	fileName := fmt.Sprintf("%s.json", pid)

	resp, err := getS3Response(config.archivematicaBucket, fileName)
	if err != nil {
		return ArchivematicaResponse, err
	}

	var S3Format ArchivematicaS3Node
	err = json.Unmarshal(resp, &S3Format)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	// Convert to TreeNode

	ArchivematicaResponse.Data = transformNode(S3Format, 0)

	return ArchivematicaResponse, nil
}

// transformNode resursively converts the archivematica tree from S3 format to TreeNode
func transformNode(s3Node ArchivematicaS3Node, depth int) TableNode {
	var node TableNode
	var data ColumnData

	if len(s3Node.ID) > 0 {
		node.Key = fmt.Sprintf("%d-%s", depth, s3Node.ID)
	} else {
		node.Key = fmt.Sprintf("%d-%s", depth, s3Node.Name)
	}
	data.Name = s3Node.Name
	data.Type = s3Node.Type

	switch s3Node.Type {
	case "folder":
		data.Icon = "fa fa-folder"
		data.Format = "Folder"
	case "file":
		data.Icon = "fa fa-file"
		data.Format = s3Node.Format
		data.URL = s3Node.SourceURL
	}

	switch s3Node.View {
	case "image":
		data.Type = "image"
		data.Icon = "fa fa-file-image"
	}

	switch {
	case strings.Contains(s3Node.Format, "PDF"):
		data.Icon = "fa fa-file-pdf"
	case strings.Contains(s3Node.Format, "Word"):
		data.Icon = "fa fa-file-word"
	case strings.Contains(s3Node.Format, "Excel"):
		data.Icon = "fa fa-file-excel"
	}

	node.Data = data

	//recursively transform children
	for _, s3Child := range s3Node.Entries {
		node.Children = append(node.Children, transformNode(s3Child, depth+1))
	}

	return node
}
