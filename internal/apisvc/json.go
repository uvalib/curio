package apisvc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

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

// ParseTracksysResponse pulls title and author data from a tracksys API call
func ParseTracksysResponse(jsonStr string) (title string, author string) {
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &jsonMap)
	title = jsonMap["title"].(string)
	author, _ = jsonMap["creator"].(string)
	return title, author
}
