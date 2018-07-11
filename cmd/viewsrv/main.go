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
const Version = "1.3.0"

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
	flag.StringVar(&config.fedoraURL, "fedora", os.Getenv("FEDORA"), "Fedora URL (required)")
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

// func renderOembedResponse(tgtURL *url.URL, format string, maxWidth int, maxHeight int, rw http.ResponseWriter, req *http.Request) {
// 	// init data used to render the oEmbed response
// 	var data oEmbedData
//
// 	// Get page param if any...
// 	qp, _ := url.ParseQuery(tgtURL.RawQuery)
// 	data.StartPage = 0
// 	if len(qp["page"]) > 0 {
// 		data.StartPage, _ = strconv.Atoi(qp["page"][0])
// 	}
//
// 	// accept 1 based page numbers from client, but use
// 	// 0-based canvas index in UV embed snippet
// 	if data.StartPage > 0 {
// 		data.StartPage--
// 		log.Printf("Requested starting page index %d", data.StartPage)
// 	}
//
// 	// URL for IIIF manifest
// 	data.PID = bits[2]
// 	data.SourceURI = fmt.Sprintf("%s/%s", config.iiifURL, data.PID)
//
// 	// Validate that the manifest has images
// 	if isManifestViewable(data.SourceURI) == false {
// 		log.Printf("Requested URL %s has no visible images", data.SourceURI)
// 		http.Error(rw, "Sorry, the requested resource is not available.", http.StatusNotFound)
// 		return
// 	}
//
// 	// scheme / host for UV javascript
// 	data.Scheme = "http"
// 	if strings.Contains(data.SourceURI, "https") {
// 		data.Scheme = "https"
// 	}
// 	data.EmbedHost = config.dovHost
// 	if len(data.EmbedHost) == 0 {
// 		data.EmbedHost = req.Host
// 	}
//
// 	// default embed size is 800x600. Params maxwidth and maxheight can override.
// 	data.Width = 800
// 	if maxWidth > 0 && maxWidth < data.Width {
// 		data.Width = maxWidth
// 	}
// 	data.Height = 600
// 	if maxHeight > 0 && maxHeight < data.Height {
// 		data.Height = maxHeight
// 	}
//
// 	// Hit Tracksys API to get brief metadata
// 	pidURL := fmt.Sprintf("%s/metadata/%s?type=brief", config.tracksysURL, data.PID)
// 	jsonResp, err := apisvc.GetAPIResponse(pidURL)
// 	if err != nil {
// 		rw.WriteHeader(http.StatusServiceUnavailable)
// 		fmt.Fprintf(rw, "Unable to connect with TrackSys to describe pid %s", data.PID)
// 		return
// 	}
//
// 	tsMetadata := apisvc.ParseTracksysResponse(jsonResp)
// 	data.Title = tsMetadata.Title
// 	data.Author = tsMetadata.Author
//
// 	// Render the <div> that will be included in the response, and used to embed the resource
// 	log.Printf("Rendering html snippet...")
// 	var renderedSnip bytes.Buffer
// 	snippet := htemplate.Must(htemplate.ParseFiles("templates/images/embed.html"))
// 	snipErr := snippet.Execute(&renderedSnip, data)
// 	if snipErr != nil {
// 		http.Error(rw, snipErr.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	rawHTML := strings.TrimSpace(renderedSnip.String())
//
// 	if strings.Compare(format, "json") == 0 {
// 		log.Printf("Rendering JSON output")
// 		data.HTML = strconv.Quote(rawHTML)
// 		rw.Header().Set("content-type", "application/json; charset=utf-8")
// 		jsonTemplate := template.Must(template.ParseFiles("templates/response.json"))
// 		jsonTemplate.Execute(rw, data)
// 	} else {
// 		rw.Header().Set("content-type", "text/xml; charset=utf-8")
// 		log.Printf("Rendering XML output")
// 		data.HTML = rawHTML
// 		var renderedSnip bytes.Buffer
// 		snippet := htemplate.Must(htemplate.ParseFiles("templates/response.xml"))
// 		snipErr := snippet.Execute(&renderedSnip, data)
// 		if snipErr != nil {
// 			log.Printf("Unable to render XML template: %s", snipErr.Error())
// 			http.Error(rw, snipErr.Error(), http.StatusInternalServerError)
// 		} else {
// 			fmt.Fprint(rw, renderedSnip.String())
// 		}
// 	}
// }
