package main

// Check health of service
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func healthCheckHandler(c *gin.Context) {

	type healthcheck struct {
		Healthy bool   `json:"healthy"`
		Message string `json:"message"`
	}

	log.Printf("INFO: checking Tracksys...")
	url := fmt.Sprintf("%s/pid/uva-lib:1157560/type", config.tracksysURL)

	tsStatus := healthcheck{true, ""}
	_, err := getAPIResponse(url)
	if err != nil {
		tsStatus.Healthy = false
		tsStatus.Message = err.Error()
	}

	// make sure IIIF manifest service is alive
	log.Printf("INFO: checking IIIF...")
	url = fmt.Sprintf("%s/version", config.iiifURL)
	iiifStatus := healthcheck{true, ""}

	_, err = getAPIResponse(url)
	if err != nil {
		iiifStatus.Healthy = false
		iiifStatus.Message = err.Error()
	}

	httpStatus := http.StatusOK
	if tsStatus.Healthy == false || iiifStatus.Healthy == false {
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, gin.H{"tracksys": tsStatus, "iiifmanifest": iiifStatus})
}
