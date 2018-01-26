package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

const version = "0.0.1"

func main() {
	// Load cfg
	log.Printf("===> viewer staring up <===")
	log.Printf("Load configuration...")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Unable to read config: %s", err.Error())
		os.Exit(1)
	}

	// Set routes and start server
	mux := httprouter.New()
	mux.GET("/", loggingHandler(rootHandler))
	mux.GET("/images/:id", imagesHandler)
	mux.GET("/oembed/:id", oEmbedHandler)
	mux.ServeFiles("/static/*filepath", http.Dir("static/"))
	log.Printf("Start service on port %s", viper.GetString("port"))
	http.ListenAndServe(":"+viper.GetString("port"), mux)
}

/**
 * Function Adapter for httprouter handlers that will log start and complete info.
 * A bit different than usual looking adapter because of the httprouter library used. Foun
 * this code here:
 *   https://stackoverflow.com/questions/43964461/how-to-use-middlewares-when-using-julienschmidt-httprouter-in-golang
 */
func loggingHandler(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		start := time.Now()
		log.Printf("Started %s %s", req.Method, req.RequestURI)
		next(w, req, ps)
		log.Printf("Completed %s %s in %s", req.Method, req.RequestURI, time.Since(start))
	}
}

/**
 * Handle a request for / and return version info
 */
func rootHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprintf(rw, "UVA Viewer version %s", version)
}

/**
 * Handle a request for oembed data
 */
func oEmbedHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprintf(rw, "FAKE OEMBED DATA")
}

/**
 * Handle a request for images from a specific ID
 */
func imagesHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// url := fmt.Sprintf("/static/test.html?id=%s", params.ByName("id"))
	url := "/static/viewer/app.html?manifestUri=http://search.lib.virginia.edu/catalog/tsb:18652/iiif/manifest.json"
	log.Printf("Redirecting to: %s", url)
	http.Redirect(rw, req, url, 301)
}
