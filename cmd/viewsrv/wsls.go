package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/uvalib/digital-object-viewer/internal/apisvc"
)

// Handle a request for a WSLS item
func wslsHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	srcPID := params.ByName("pid")
	log.Printf("Get Apollo PID for %s", srcPID)
	pidURL := fmt.Sprintf("%s/external/%s", config.apolloURL, srcPID)
	apolloPID, err := apisvc.GetAPIResponse(pidURL)
	if err != nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(rw, "Unable to connect with Apollo get info for PID %s", srcPID)
		return
	}

	metadataURL := fmt.Sprintf("%s/items/%s", config.apolloURL, apolloPID)
	metadataJSON, err := apisvc.GetAPIResponse(metadataURL)
	if err != nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(rw, "Unable to connect with Apollo get metadata for Apollo PID %s", apolloPID)
		return
	}

	wslsData, parseErr := apisvc.ParseApolloWSLSResponse(metadataJSON)
	if parseErr != nil {
		log.Printf("FAILED PARSE: %s", parseErr.Error())
	}
	log.Printf("GOT %#v", wslsData)
	/* NOTES:
			   HIT apollo item APIs to get JSON for the item;
			      /api/external/uva-lib:2220355  : get Apollo PID
			      /api/items/uva-an109886 : get JSON for collection / item

		      Check hasVideo and hasScript properties to determine what to show
		      VIDEO:
		         POSTER: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-poster.jpg
		         VIDEO (webm): http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.webm
		      SCRIPT:
		         PDF: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}.pdf
		         Thumb: http://fedora01.lib.virginia.edu/wsls/{wslsID}/{wslsID}-script-thumbnail.jpg
	            Transcript: http://fedora01.lib.virginia.edu/wsls/0003_1/0003_1.txt
	*/
	fmt.Fprintf(rw, "WSLS support is under construction. Apollo PID: %s", apolloPID)
}
