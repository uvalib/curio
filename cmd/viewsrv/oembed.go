package main

import (
	"bytes"
	"errors"
	"fmt"
	htemplate "html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"github.com/uvalib/digital-object-viewer/internal/apisvc"
)

type oEmbedData struct {
	Title  string
	Author string
	HTML   string
	Width  int
	Height int
}

type embedImageData struct {
	PID       string
	Width     int
	Height    int
	SourceURI string
	Scheme    string
	EmbedHost string
	StartPage int
}

// Handle a request for oembed data
func oEmbedHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get some optional params; format, maxWidth and maxHeight
	respFormat := req.URL.Query().Get("format")
	if respFormat == "" {
		respFormat = "json"
	}

	maxWidth, err := strconv.Atoi(req.URL.Query().Get("maxwidth"))
	if err != nil {
		maxWidth = 0
	}

	maxHeight, err := strconv.Atoi(req.URL.Query().Get("maxheight"))
	if err != nil {
		maxHeight = 0
	}

	// Next, get the required URL and see if a page is requested
	urlStr, _ := url.QueryUnescape(req.URL.Query().Get("url"))
	if len(urlStr) == 0 {
		http.Error(rw, "URL is required!", http.StatusBadRequest)
		return
	}

	// The raw URL requested must be of the expected format: [http|https]://[host]/[images|wsls]/[PID][?page=n]
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		msg := fmt.Sprintf("Invalid URL: %s", err.Error())
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}

	// Now split out relatve path to find PID. This should be something like: /[images|wsls]/[PID]
	// NOTE: that this wil strip out all query params
	relPath := parsedURL.Path
	bits := strings.Split(relPath, "/")
	if len(bits) != 3 {
		msg := fmt.Sprintf("Invalid URL in request: %s", urlStr)
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}

	pid := bits[2]
	resourceType := bits[1]

	// See what type of resource is being requested
	var respData oEmbedData
	if resourceType == "images" {
		respData, err = getImageData(parsedURL, pid, maxWidth, maxHeight, req.Host)
	} else if resourceType == "wsls" {
		respData, err = getWSLSData(parsedURL, pid, respFormat, maxWidth, maxHeight)
	} else {
		err = fmt.Errorf("invalid resource type: %s", bits[1])
	}

	if err != nil {
		log.Printf("ERROR: Unable to render oEmbed response: %s", err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if respFormat == "json" {
		log.Printf("Rendering JSON output")
		rw.Header().Set("content-type", "application/json; charset=utf-8")
		jsonTemplate := template.Must(template.ParseFiles("templates/response.json"))
		jsonTemplate.Execute(rw, respData)
	} else {
		rw.Header().Set("content-type", "text/xml; charset=utf-8")
		log.Printf("Rendering XML output")
		var renderedSnip bytes.Buffer
		snippet := htemplate.Must(htemplate.ParseFiles("templates/response.xml"))
		snipErr := snippet.Execute(&renderedSnip, respData)
		if snipErr != nil {
			log.Printf("Unable to render XML template: %s", snipErr.Error())
			http.Error(rw, snipErr.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprint(rw, renderedSnip.String())
		}
	}
}

func getImageData(tgtURL *url.URL, pid string, maxWidth int, maxHeight int, reqHost string) (oEmbedData, error) {
	var respData oEmbedData
	var data embedImageData
	data.PID = pid
	data.SourceURI = fmt.Sprintf("%s/%s", config.iiifURL, data.PID)

	// Get page param if any...
	qp, _ := url.ParseQuery(tgtURL.RawQuery)
	data.StartPage = 0
	if len(qp["page"]) > 0 {
		data.StartPage, _ = strconv.Atoi(qp["page"][0])
	}

	// accept 1 based page numbers from client, but use
	// 0-based canvas index in UV embed snippet
	if data.StartPage > 0 {
		data.StartPage--
		log.Printf("Requested starting page index %d", data.StartPage)
	}

	// Validate that the manifest has images
	if isManifestViewable(data.SourceURI) == false {
		log.Printf("Requested URL %s has no visible images", data.SourceURI)
		return respData, errors.New("requested resource is not available")
	}

	// scheme / host for UV javascript
	data.Scheme = "http"
	if strings.Contains(data.SourceURI, "https") {
		data.Scheme = "https"
	}
	data.EmbedHost = config.dovHost
	if len(data.EmbedHost) == 0 {
		data.EmbedHost = reqHost
	}

	// default embed size is 800x600. Params maxwidth and maxheight can override.
	data.Width = 800
	if maxWidth > 0 && maxWidth < data.Width {
		data.Width = maxWidth
	}
	data.Height = 600
	if maxHeight > 0 && maxHeight < data.Height {
		data.Height = maxHeight
	}

	// Render the <div> that will be included in the response, and used to embed the resource
	log.Printf("Rendering html snippet...")
	var renderedSnip bytes.Buffer
	snippet := htemplate.Must(htemplate.ParseFiles("templates/images/embed.html"))
	snipErr := snippet.Execute(&renderedSnip, data)
	if snipErr != nil {
		return respData, snipErr
	}
	rawHTML := strings.TrimSpace(renderedSnip.String())

	// Hit Tracksys API to get brief metadata
	metadataURL := fmt.Sprintf("%s/metadata/%s?type=brief", config.tracksysURL, data.PID)
	jsonResp, err := apisvc.GetAPIResponse(metadataURL)
	if err != nil {
		return respData, fmt.Errorf("Unable to connect with TrackSys to describe pid %s", data.PID)
	}
	tsMetadata := apisvc.ParseTracksysResponse(jsonResp)

	respData.Title = tsMetadata.Title
	respData.Author = tsMetadata.Author
	respData.HTML = strconv.Quote(rawHTML)
	respData.Width = data.Width
	respData.Height = data.Height
	return respData, nil
}

func getWSLSData(tgtURL *url.URL, pid string, respFormat string, maxWidth int, maxHeight int) (oEmbedData, error) {
	var respData oEmbedData
	return respData, nil
}
