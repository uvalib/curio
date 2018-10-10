package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// Version of the service
const Version = "1.4.0"

type configData struct {
	port        int
	tracksysURL string
	apolloURL   string
	iiifURL     string
	fedoraURL   string
	dovHost     string
}

// golbals for DB and CFG
var config configData

func main() {
	// Load cfg
	log.Printf("===> viewer staring up <===")
	getConfiguration()

	// Set routes and start server
	mux := httprouter.New()
	mux.GET("/", loggingHandler(rootHandler))
	mux.GET("/images/:pid", loggingHandler(imagesHandler))
	mux.GET("/wsls/:pid", loggingHandler(wslsHandler))
	mux.GET("/oembed", loggingHandler(oEmbedHandler))
	mux.GET("/healthcheck", loggingHandler(healthCheckHandler))
	mux.ServeFiles("/web/*filepath", http.Dir("web/"))
	log.Printf("Start service on port %d with CORS support enabled", config.port)
	port := fmt.Sprintf(":%d", config.port)
	http.ListenAndServe(port, cors.Default().Handler(mux))
}

func getConfiguration() {
	// FIRST, try command line flags. Fallback is ENV variables
	flag.IntVar(&config.port, "port", 8085, "Port to offer service on (default 8085)")
	flag.StringVar(&config.tracksysURL, "tracksys", os.Getenv("TRACKSYS_URL"), "TrackSys URL (required)")
	flag.StringVar(&config.apolloURL, "apollo", os.Getenv("APOLLO_URL"), "Apollo URL (required)")
	flag.StringVar(&config.iiifURL, "iiif", os.Getenv("IIIF"), "IIIF URL (required)")
	flag.StringVar(&config.fedoraURL, "fedora", os.Getenv("WSLS_FEDORA_URL"), "WSLS Fedora URL (required)")
	flag.StringVar(&config.dovHost, "dovhost", os.Getenv("DOV_HOST"), "DoViewer Hostname (optional)")
	flag.Parse()

	// if anything is still not set, die
	if config.tracksysURL == "" || config.iiifURL == "" || config.apolloURL == "" || config.fedoraURL == "" {
		flag.Usage()
		os.Exit(1)
	}
	if len(config.dovHost) == 0 {
		log.Printf("DOV host not set; this info will be extracted from request headers")
	} else {
		log.Printf("DOV host set to: %s", config.dovHost)
	}
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

// Handle a request for / and return version info
func rootHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprintf(rw, "UVA Viewer version %s", Version)
}
