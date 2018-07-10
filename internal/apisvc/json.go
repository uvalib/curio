package apisvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// This is a minimal mapping to get the name/value of an apollo node
// Without the leading caps and json mapping, the unmarshall doesn't work
type apolloItem struct {
	ItemType struct {
		Name string
	} `json:"type"`
	Value string
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
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &jsonMap)
	itemData, ok := jsonMap["item"].(map[string]interface{})
	if !ok {
		return data, errors.New("Missing item data in response")
	}
	children, ok := itemData["children"].([]interface{})
	if !ok {
		return data, errors.New("Missing item children in response")
	}
	for _, c := range children {
		var item apolloItem
		itemBytes, err := json.Marshal(c)
		if err != nil {
			return data, fmt.Errorf("Unable to extract child item data: %s", err.Error())
		}
		err = json.Unmarshal(itemBytes, &item)
		if err != nil {
			return data, fmt.Errorf("Unable to extract child item data: %s", err.Error())
		}
		log.Printf("===== Item STRUCT: %+v", item)
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
