package main

// Check health of service
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func healthCheckHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")

	// TCheck TrackSys
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
				iiifStatus = false
			}
		}
	}
	out := fmt.Sprintf(`{"alive": true, "iiif": %t, "tracksys": %t}`, iiifStatus, tsStatus)
	if iiifStatus == false {
		http.Error(rw, out, http.StatusInternalServerError)
	} else {
		fmt.Fprintf(rw, out)
	}
}
