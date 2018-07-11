package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TrackSysMetadata contains the basic metadata returned from the Tracksys API
type TrackSysMetadata struct {
	Title  string
	Author string
}

// WSLSMetadata contains the Apollo metadata supporting WSLS
type WSLSMetadata struct {
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

// GetAPIResponse calls a JSON endpoint and returns the resoponce
func GetAPIResponse(url string) (string, error) {
	resp, err := http.Get(url)
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

// ParseApolloWSLSResponse pulls title and author data from a tracksys API call
func ParseApolloWSLSResponse(jsonStr string) (WSLSMetadata, error) {
	var data WSLSMetadata
	var respStruct apolloResp
	err := json.Unmarshal([]byte(jsonStr), &respStruct)
	if err != nil {
		return data, fmt.Errorf("Unable parse response: %s", err.Error())
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
	return data, nil
}

// ParseTracksysResponse pulls title and author data from a tracksys API call
func ParseTracksysResponse(jsonStr string) TrackSysMetadata {
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &jsonMap)
	var data TrackSysMetadata
	data.Title = jsonMap["title"].(string)
	data.Author, _ = jsonMap["creator"].(string)
	return data
}
