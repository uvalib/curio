package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handle a request for a WSLS item
func wslsHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprintf(rw, "WSLS support is under construction")
}
