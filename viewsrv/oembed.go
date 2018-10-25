package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type oembed struct {
	Version     string `json:"version,omitempty" xml:"version,omitempty"`
	Type        string `json:"type,omitempty" xml:"type,omitempty"`
	Title       string `json:"title,omitempty" xml:"title,omitempty"`
	Author      string `json:"author,omitempty" xml:"author,omitempty"`
	HTML        string `json:"html,omitempty" xml:"html,omitempty"`
	Width       int    `json:"width,omitempty" xml:"width,omitempty"`
	Height      int    `json:"height,omitempty" xml:"height,omitempty"`
	Provider    string `json:"provider,omitempty" xml:"provider,omitempty"`
	ProviderURL string `json:"provider_url,omitempty" xml:"provider_url,omitempty"`
}

// custom marshal that doesn't do the weird escaling of < >
func (o *oembed) marshalJSON() string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "   ")
	encoder.Encode(o)
	return buffer.String()
}

// embedImageData is the data needed to render the HTML snippet fot embedded images
type embedImageData struct {
	Width     int
	Height    int
	SourceURI string
	Scheme    string
	EmbedHost string
	StartPage int
}

type embedWSLSData struct {
	Width     int
	Height    int
	SourceURI string
}

// oEmbedHandler returns the oEmbed data for a view
func oEmbedHandler(c *gin.Context) {
	// Get some optional params; format, maxWidth and maxHeight
	respFormat := c.Query("format")
	if respFormat == "" {
		respFormat = "json"
	}

	maxWidth, err := strconv.Atoi(c.Query("maxwidth"))
	if err != nil {
		maxWidth = 0
	}

	maxHeight, err := strconv.Atoi(c.Query("maxheight"))
	if err != nil {
		maxHeight = 0
	}

	// Next, get the required URL and see if a page is requested
	urlStr, _ := url.QueryUnescape(c.Query("url"))
	if len(urlStr) == 0 {
		c.String(http.StatusBadRequest, "A URL param is required!")
		return
	}

	// The raw URL requested must be of the expected format: [http|https]://[host]/[images|wsls]/[PID][?page=n]
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid URL: %s", err.Error())
		return
	}

	// Now split out relatve path to find PID. This should be something like: /[images|wsls]/[PID]
	// NOTE: that this wil strip out all query params
	relPath := parsedURL.Path
	bits := strings.Split(relPath, "/")
	if len(bits) != 3 {
		c.String(http.StatusBadRequest, "Invalid URL in request: %s", urlStr)
		return
	}

	pid := bits[2]
	resourceType := bits[1]

	// See what type of resource is being requested
	var respData oembed
	if resourceType == "images" {
		respData, err = getImageData(parsedURL, pid, maxWidth, maxHeight)
	} else if resourceType == "wsls" {
		respData, err = getWSLSData(parsedURL, pid, maxWidth, maxHeight)
	} else {
		err = fmt.Errorf("invalid resource type: %s", bits[1])
	}

	if err != nil {
		log.Printf("ERROR: Unable to render oEmbed response: %s", err.Error())
		c.String(http.StatusInternalServerError, "Unable to render oEmbed response: %s", err.Error())
		return
	}

	if respFormat == "json" {
		log.Printf("Rendering JSON output")
		c.Header("content-type", "application/json; charset=utf-8")
		c.String(http.StatusOK, respData.marshalJSON())
	} else {
		c.XML(http.StatusOK, respData)
	}
}

func getImageData(tgtURL *url.URL, pid string, maxWidth int, maxHeight int) (oembed, error) {
	respData := oembed{Version: "1.0", Type: "rich", Provider: "UVA Library", ProviderURL: "http://www.library.virginia.edu/"}
	var imgData embedImageData
	imgData.EmbedHost = config.dovHost
	imgData.SourceURI = fmt.Sprintf("%s/%s", config.iiifURL, pid)

	// Get page param if any...
	qp, _ := url.ParseQuery(tgtURL.RawQuery)
	imgData.StartPage = 0
	if len(qp["page"]) > 0 {
		imgData.StartPage, _ = strconv.Atoi(qp["page"][0])
	}

	// accept 1 based page numbers from client, but use
	// 0-based canvas index in UV embed snippet
	if imgData.StartPage > 0 {
		imgData.StartPage--
		log.Printf("Requested starting page index %d", imgData.StartPage)
	}

	// Validate that the manifest has images
	if isManifestViewable(imgData.SourceURI) == false {
		log.Printf("Requested URL %s has no visible images", imgData.SourceURI)
		return respData, errors.New("requested resource is not available")
	}

	// scheme / host for UV javascript
	imgData.Scheme = "http"
	if strings.Contains(imgData.SourceURI, "https") {
		imgData.Scheme = "https"
	}

	// default embed size is 800x600. Params maxwidth and maxheight can override.
	imgData.Width = 800
	if maxWidth > 0 && maxWidth < imgData.Width {
		imgData.Width = maxWidth
	}
	imgData.Height = 600
	if maxHeight > 0 && maxHeight < imgData.Height {
		imgData.Height = maxHeight
	}

	// Render the <div> that will be included in the response, and used to embed the resource
	log.Printf("Rendering html snippet...")
	var renderedSnip bytes.Buffer
	snippet := template.Must(template.ParseFiles("templates/image_embed.html"))
	snipErr := snippet.Execute(&renderedSnip, imgData)
	if snipErr != nil {
		return respData, snipErr
	}
	rawHTML := strings.TrimSpace(renderedSnip.String())

	tsMetadata, err := getTracksysMetadata(pid)
	if err != nil {
		return respData, err
	}

	respData.Title = tsMetadata.Title
	respData.Author = tsMetadata.Author
	respData.HTML = rawHTML
	respData.Width = imgData.Width
	respData.Height = imgData.Height
	return respData, nil
}

func getWSLSData(tgtURL *url.URL, pid string, maxWidth int, maxHeight int) (oembed, error) {
	respData := oembed{Version: "1.0", Type: "rich", Provider: "UVA Library", ProviderURL: "http://www.library.virginia.edu/"}
	var snipData embedWSLSData

	snipData.SourceURI = tgtURL.String()
	snipData.Width = 670
	if maxWidth > 0 && maxWidth < snipData.Width {
		snipData.Width = maxWidth
	}
	snipData.Height = 800
	if maxHeight > 0 && maxHeight < snipData.Height {
		snipData.Height = maxHeight
	}

	log.Printf("Rendering html snippet...")
	var renderedSnip bytes.Buffer
	snippet := template.Must(template.ParseFiles("templates/wsls_embed.html"))
	snipErr := snippet.Execute(&renderedSnip, snipData)
	if snipErr != nil {
		return respData, snipErr
	}
	rawHTML := strings.TrimSpace(renderedSnip.String())

	wslsData, parseErr := getApolloWSLSMetadata(pid)
	if parseErr != nil {
		return respData, parseErr
	}
	respData.Title = wslsData.Title
	respData.HTML = rawHTML
	respData.Width = snipData.Width
	respData.Height = snipData.Height
	return respData, nil
}
