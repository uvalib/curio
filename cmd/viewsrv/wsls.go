package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uvalib/digital-object-viewer/pkg/apisvc"
)

// wslsHandler accepts a request for a WSLS item and renders in in a custom viewer
func wslsHandler(c *gin.Context) {
	srcPID := c.Param("pid")
	metadataURL := fmt.Sprintf("%s/items/%s", config.apolloURL, srcPID)
	metadataJSON, err := apisvc.GetAPIResponse(metadataURL)
	if err != nil {
		log.Printf("ERROR: Unable to connect with Apollo get metadata for Apollo PID %s: %s", srcPID, err.Error())
		c.String(http.StatusServiceUnavailable, "Unable to retrieve metadata for %s", srcPID)
		return
	}

	// ... and parse it into the necessary data for the viewer
	wslsData, parseErr := apisvc.ParseApolloWSLSResponse(metadataJSON)
	if parseErr != nil {
		log.Printf("ERROR: Unable to parse Apollo response for %s: %s", srcPID, parseErr.Error())
		c.String(http.StatusInternalServerError, "Unable to retrieve metadata for %s", srcPID)
		return
	}

	if wslsData.HasVideo {
		// POSTER: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-poster.jpg
		// VIDEO (webm): http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.webm
		wslsData.VideoURL = fmt.Sprintf("%s/%s/%s.webm", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.PosterURL = fmt.Sprintf("%s/%s/%s-poster.jpg", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
	}

	if wslsData.HasScript {
		// PDF: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.pdf
		// Thumb: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-script-thumbnail.jpg
		// Transcript: http://fedora01.lib.virginia.edu/wsls/0003_1/0003_1.txt
		wslsData.PDFURL = fmt.Sprintf("%s/%s/%s.pdf", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.PDFThumbURL = fmt.Sprintf("%s/%s/%s-script-thumbnail.jpg", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.TranscriptURL = fmt.Sprintf("%s/%s/%s.txt", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
	}

	c.HTML(http.StatusOK, "wsls_view.html", wslsData)
}
