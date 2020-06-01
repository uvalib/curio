package main

// Check health of service
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func healthCheckHandler(c *gin.Context) {
	log.Printf("Checking Tracksys...")
	tsStatus := true
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	url := fmt.Sprintf("%s/pid/uva-lib:1157560/type", config.tracksysURL)
	log.Printf("Tracksys test URL: %s", url)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("ERROR: TrackSys service (%s)", err)
		tsStatus = false
	} else {
		b, errRead := ioutil.ReadAll(resp.Body)
		if errRead != nil {
			log.Printf("ERROR: TrackSys service (%s)", errRead)
			tsStatus = false
		} else {
			resp.Body.Close()
			if string(b) != "sirsi_metadata" {
				log.Printf("ERROR: TrackSys bad response (%s)", b)
				tsStatus = false
			}
		}
	}

	// make sure IIIF manifest service is alive
	log.Printf("Checking IIIF...")
	url = fmt.Sprintf("%s/version", config.iiifURL)
	log.Printf("IIIF test URL: %s", url)
	iiifStatus := true
	resp, err = client.Get(config.iiifURL)
	if err != nil {
		log.Printf("ERROR: IIIF service (%s)", err)
		iiifStatus = false
	} else {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ERROR: IIIF service (%s)", err)
			iiifStatus = false
		} else {
			resp.Body.Close()
			if strings.Contains(string(b), "IIIF metadata service") == false {
				log.Printf("ERROR: IIIF service reports unexpected version info (%s)", string(b))
				iiifStatus = false
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"alive": true, "iiif": iiifStatus, "tracksys": tsStatus})
}
