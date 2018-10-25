package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type viewerData struct {
	URI       string
	StartPage int
}

// imagesHandler takes a request for images from a specific image PID and page offset (optional) and
// displays it in the universal viewer
func imagesHandler(c *gin.Context) {
	// pull page QP and use it for starting page. Any other params are ignored.
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	var data viewerData
	data.StartPage = page - 1
	data.URI = fmt.Sprintf("%s/%s", config.iiifURL, c.Param("pid"))

	// Make sure there are images visable for this PID.
	// Ahow an error page and bail if not
	if isManifestViewable(data.URI) == false {
		c.HTML(http.StatusNotFound, "not_available.html", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "image_view.html", data)
}

// Hit the target IIIF manifest URL and see if it contains any images
func isManifestViewable(url string) bool {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	log.Printf("Checking manifest URL %s for images", url)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("ERROR: IIIF URL: %s failed to return a response", url)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: IIIF URL: %s returned non-success status: %d", url, resp.StatusCode)
		return false
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	respStr := string(bytes)

	return strings.Contains(respStr, "dcTypes:Image")
}
