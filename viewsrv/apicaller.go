package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

// getAPIResponse calls a JSON endpoint and returns the resoponce
func getAPIResponse(url string) (string, error) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	respString := string(bodyBytes)
	if resp.StatusCode != 200 {
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

// getTracksysMetadata pulls title and author data from a tracksys API call
func getTracksysMetadata(pid string) (*tracksysMetadata, error) {
	// Hit Tracksys API to get brief metadata
	metadataURL := fmt.Sprintf("%s/metadata/%s?type=brief", config.tracksysURL, pid)
	jsonResp, err := getAPIResponse(metadataURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect with TrackSys to describe pid %s", pid)
	}

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonResp), &jsonMap)
	var data tracksysMetadata
	data.Title = jsonMap["title"].(string)
	data.Author, _ = jsonMap["creator"].(string)
	return &data, nil
}
