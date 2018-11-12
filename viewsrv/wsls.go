package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// wslsHandler accepts a request for a WSLS item and renders in in a custom viewer
func wslsHandler(c *gin.Context) {
	srcPID := c.Param("pid")
	wslsData, parseErr := getApolloWSLSMetadata(srcPID)
	if parseErr != nil {
		log.Printf("Unable to get Apollo response for %s: %s", srcPID, parseErr.Error())
		c.String(http.StatusNotFound, "%s not found", srcPID)
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
