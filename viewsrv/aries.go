package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// AriesService contains details for a service reference
type AriesService struct {
	URL      string `json:"url"`
	Protocol string `json:"protocol"`
}

// Aries is a structure containing the data to be returned from an aries query
type Aries struct {
	Identifiers []string       `json:"identifier,omitempty"`
	AccessURL   []string       `json:"access_url,omitempty"`
	Services    []AriesService `json:"service_url,omitempty"`
}

// AriesPing handles requests to the aries endpoint with no params.
// Just returns and alive message
func ariesPing(c *gin.Context) {
	c.String(http.StatusOK, "Curio Aries API")
}

// AriesLookup will query apollo for information on the supplied identifer
func ariesLookup(c *gin.Context) {
	passedPID := c.Param("id")
	out := Aries{}
	out.Identifiers = append(out.Identifiers, passedPID)

	// easy check; see if there is an IIIF manifest with this PID visible
	_, err := getIIIFManifestURL(passedPID, "")
	if err == nil {
		// yes; this is an image asset. Return the oEbmed and viewer URLs
		publicURL := fmt.Sprintf("https://%s/view/%s", config.hostname, passedPID)
		oEmbedURL := fmt.Sprintf("https://%s/oembed?url=%s", config.hostname, url.QueryEscape(publicURL))
		svc := AriesService{URL: oEmbedURL, Protocol: "oembed"}
		out.AccessURL = append(out.AccessURL, publicURL)
		out.Services = append(out.Services, svc)

		c.JSON(http.StatusOK, out)
		return
	}

	// Not an image... see if it is WSLS
	_, err = getApolloWSLSMetadata(passedPID)
	if err == nil {
		publicURL := fmt.Sprintf("https://%s/view/%s", config.hostname, passedPID)
		out.AccessURL = append(out.AccessURL, publicURL)
		oEmbedURL := fmt.Sprintf("https://%s/oembed?url=%s", config.hostname, url.QueryEscape(publicURL))
		svc := AriesService{URL: oEmbedURL, Protocol: "oembed"}
		out.Services = append(out.Services, svc)
		c.JSON(http.StatusOK, out)
		return
	}

	c.String(http.StatusNotFound, "%s not found", passedPID)
}
