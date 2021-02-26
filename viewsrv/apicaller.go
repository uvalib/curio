package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// This is a minimal mapping of the apollo items API request to the
// data needed by the digital object viewer. Note that alll property names
// must be leading caps and match the json repspone field (case insensitive)
// or be mapped with a json attribute
type apolloResp struct {
	Item struct {
		Children []struct {
			ItemType struct {
				Name string
			} `json:"type"`
			Value string
		}
	}
}

// tracksysMetadata contains the basic metadata returned from the Tracksys API
type tracksysMetadata struct {
	Title  string
	Author string
}

// wslsMetadata contains the Apollo metadata supporting WSLS
type wslsMetadata struct {
	HasVideo      bool
	HasScript     bool
	WSLSID        string
	Title         string
	Description   string
	VideoURL      string
	PosterURL     string
	PDFURL        string
	PDFThumbURL   string
	TranscriptURL string
	Duration      string
}

// use a shared client, 5 second connect, 15 second read timeout
var httpClient = httpClientWithTimeouts(5, 15)

// getAPIResponse calls a JSON endpoint and returns the response
func getAPIResponse(url string) (string, error) {

	log.Printf("INFO: GET %s", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("ERROR: %s returns %s", url, err.Error())
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	respString := string(bodyBytes)
	if resp.StatusCode != http.StatusOK {
		logLevel := "ERROR"
		// some errors are expected
		if resp.StatusCode == http.StatusNotFound {
			logLevel = "INFO"
		}
		log.Printf("%s: %s returns %d (%s)", logLevel, url, resp.StatusCode, respString)
		return "", errors.New(respString)
	}
	return respString, nil
}

func getApolloWSLSMetadata(pid string) (*wslsMetadata, error) {
	metadataURL := fmt.Sprintf("%s/items/%s", config.apolloURL, pid)
	metadataJSON, err := getAPIResponse(metadataURL)
	if err != nil {
		return nil, err
	}

	// ... and parse it into the necessary data for the viewer
	var data wslsMetadata
	var respStruct apolloResp
	err = json.Unmarshal([]byte(metadataJSON), &respStruct)
	if err != nil {
		return nil, fmt.Errorf("Unable parse response: %s", err.Error())
	}
	for _, item := range respStruct.Item.Children {
		switch val := item.ItemType.Name; val {
		case "wslsID":
			data.WSLSID = item.Value
		case "title":
			data.Title = item.Value
		case "hasVideo":
			data.HasVideo = (item.Value == "true")
		case "hasScript":
			data.HasScript = (item.Value == "true")
		case "abstract":
			data.Description = item.Value
		case "duration":
			data.Duration = item.Value
		}
	}
	return &data, nil
}

func httpClientWithTimeouts(connTimeout int, readTimeout int) *http.Client {

	client := &http.Client{
		Timeout: time.Duration(readTimeout) * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(connTimeout) * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return client
}

//
// end of file
//
