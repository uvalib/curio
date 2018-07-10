package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type viewerData struct {
	URI       string
	StartPage int
}

// Handle a request for images from a specific image PID and page offset (optional).
func imagesHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// pull page QP and use it for starting page. Any other params are ignored.
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	var data viewerData
	data.StartPage = page - 1
	data.URI = fmt.Sprintf("%s/%s", config.iiifURL, params.ByName("pid"))

	// Make sure there are images visable for this PID.
	// Ahow an error page and bail if not
	if isManifestViewable(data.URI) == false {
		rw.WriteHeader(http.StatusNotFound)
		bytes, _ := ioutil.ReadFile("web/not_available.html")
		fmt.Fprintf(rw, "%s", string(bytes))
		return
	}

	template, err := template.ParseFiles("templates/images/view.html")
	if err != nil {
		msg := fmt.Sprintf("Unable to render viewer: %s", err.Error())
		http.Error(rw, msg, http.StatusInternalServerError)
	} else {
		template.Execute(rw, data)
	}
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
