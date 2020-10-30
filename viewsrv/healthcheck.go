package main

// Check health of service
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func healthCheckHandler(c *gin.Context) {
	log.Printf("INFO: checking Tracksys...")
	url := fmt.Sprintf("%s/pid/uva-lib:1157560/type", config.tracksysURL)

	tsStatus := true
	_, err := getAPIResponse( url )
	if err != nil {
		tsStatus = false
	}

	// make sure IIIF manifest service is alive
	log.Printf("INFO: checking IIIF...")
	url = fmt.Sprintf( "%s/version", config.iiifURL)
	iiifStatus := true

	_, err = getAPIResponse( url )
	if err != nil {
		iiifStatus = false
	}

	c.JSON(http.StatusOK, gin.H{"alive": true, "iiif": iiifStatus, "tracksys": tsStatus})
}
