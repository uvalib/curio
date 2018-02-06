package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//
// HealthCheck -- calls the service health check method
//
func HealthCheck(baseURI string) (int, string) {
	url := fmt.Sprintf("%s/healthcheck", baseURI)
	return getResponse(url)
}

//
// Oembed -- calls the oembed endpoint and returns results
//
func Oembed(baseURI string, iiifURL string, id string) (int, string) {
	url := fmt.Sprintf("%s/oembed?url=%s/%s&format=json", baseURI, iiifURL, id)
	return getResponse(url)
}

func getResponse(url string) (int, string) {
	resp, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}
