package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handle a request for a WSLS item
func wslsHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// First, find the ApolloID for this PID...
	srcPID := params.ByName("pid")
	log.Printf("Get Apollo PID for %s", srcPID)
	pidURL := fmt.Sprintf("%s/external/%s", config.apolloURL, srcPID)
	apolloPID, err := GetAPIResponse(pidURL)
	if err != nil {
		log.Printf("ERROR: unable to get apollo pid for %s: %s", srcPID, err.Error())
		rw.WriteHeader(http.StatusNotFound)
		bytes, _ := ioutil.ReadFile("web/not_available.html")
		fmt.Fprintf(rw, "%s", string(bytes))
		return
	}

	// Use the ApolloPID to get metadata describing the items...
	metadataURL := fmt.Sprintf("%s/items/%s", config.apolloURL, apolloPID)
	metadataJSON, err := GetAPIResponse(metadataURL)
	if err != nil {
		log.Printf("ERROR: unable to parse apollo response for %s: %s", apolloPID, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Unable to connect with Apollo get metadata for Apollo PID %s", apolloPID)
		return
	}

	// ... and parse it into the necessary data for the viewer
	wslsData, parseErr := ParseApolloWSLSResponse(metadataJSON)
	if parseErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Unable to connect parse Apollo response for PID %s: %s", apolloPID, parseErr.Error())
		return
	}

	if wslsData.HasVideo {
		// POSTER: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-poster.jpg
		// VIDEO (webm): http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.webm
		wslsData.VideoURL = fmt.Sprintf("%s/%s/%s.webm", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.PosterURL = fmt.Sprintf("%s/%s/%s-poster.jpg", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
	}

	if wslsData.HasScript {
		// PDF: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.pdf
		// Thumb: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-script-thumbnail.jpg
		// Transcript: http://fedora01.lib.virginia.edu/wsls/0003_1/0003_1.txt
		wslsData.PDFURL = fmt.Sprintf("%s/%s/%s.pdf", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.PDFThumbURL = fmt.Sprintf("%s/%s/%s-script-thumbnail.jpg", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
		wslsData.TranscriptURL = fmt.Sprintf("%s/%s/%s.txt", config.fedoraURL, wslsData.WSLSID, wslsData.WSLSID)
	}

	template, err := template.ParseFiles("templates/wsls/view.html")
	if err != nil {
		msg := fmt.Sprintf("Unable to render viewer: %s", err.Error())
		http.Error(rw, msg, http.StatusInternalServerError)
	} else {
		template.Execute(rw, wslsData)
	}
}
