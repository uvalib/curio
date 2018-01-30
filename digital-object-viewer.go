package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

const version = "0.0.1"

var db *sql.DB

type oEmbedData struct {
	Title  string
	HTML   string
	URL    string
	Width  int
	Height int
}

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

	// Init DB connection
	log.Printf("Init DB connection to %s...", viper.GetString("db_host"))
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", viper.GetString("db_user"), viper.GetString("db_pass"),
		viper.GetString("db_host"), viper.GetString("db_name"))
	db, err = sql.Open("mysql", connectStr)
	if err != nil {
		log.Fatal("Database connection failed: %s", err.Error())
		os.Exit(1)
	}
	_, err = db.Query("SELECT 1")
	if err != nil {
		log.Fatal("Database query failed: %s", err.Error())
		os.Exit(1)
	}
	defer db.Close()
	log.Printf("DB Connection established")

	// Set routes and start server
	mux := httprouter.New()
	mux.GET("/", loggingHandler(rootHandler))
	mux.GET("/images/:id", loggingHandler(imagesHandler))
	mux.GET("/oembed", loggingHandler(oEmbedHandler))
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
	urlStr := req.URL.Query().Get("url")
	if len(urlStr) == 0 {
		http.Error(rw, "URL is required!", http.StatusBadRequest)
		return
	}
	respFormat := req.URL.Query().Get("format")

	if len(respFormat) == 0 || strings.Compare(respFormat, "json") == 0 {
		log.Printf("JSON response requested")
		renderJSONResponse(urlStr, rw)
	} else if strings.Compare(respFormat, "xml") == 0 {
		log.Printf("XML response requested")
		renderXMLResponse(urlStr, rw)
	} else {
		// error: unsupported format
		err := fmt.Sprintf("Requested format '%s' is invalid.", respFormat)
		http.Error(rw, err, http.StatusBadRequest)
	}
}

func renderJSONResponse(rawURL string, rw http.ResponseWriter) {
	rw.Header().Set("content-type", "application/json; charset=utf-8")

	// The URL request must be of the expected format: http://[host]/images/[PID]
	// Extract the PID and generate the JSON data. First, make sure it is valid:
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		msg := fmt.Sprintf("Invalid URL: %s", err.Error())
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}

	// Now split out relatve path. This should be something like: /images/[PID]
	relPath := parsedURL.Path
	bits := strings.Split(relPath, "/")
	if len(bits) != 3 {
		msg := fmt.Sprintf("Invalid URL in request: %s", rawURL)
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}
	if strings.Compare(bits[1], "images") != 0 {
		msg := fmt.Sprintf("Invalid resource type in URL: %s", bits[1])
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "This would be JSON data")

	// tmpl := template.Must(template.ParseFiles("response.json"))
	// err := tmpl.ExecuteTemplate(rw, "iiif.json", data)
	// if err != nil {
	// 	logger.Printf("Unable to render IIIF metadata for %s: %s", data.MetadataPID, err.Error())
	// 	fmt.Fprintf(rw, "Unable to render IIIF metadata: %s", err.Error())
	// 	return
	// }
	// logger.Printf("IIIF Metadata generated for %s", data.MetadataPID)
}

func renderXMLResponse(url string, rw http.ResponseWriter) {
	rw.Header().Set("content-type", "text/xml; charset=utf-8")
}

/**
 * Handle a request for images from a specific ID (TrackSys PID)
 */
func imagesHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	url := fmt.Sprintf("%s/%s/manifest.json", viper.GetString("iiif_manifest_url"), params.ByName("id"))
	log.Printf("Target manifest: %s", url)
	template, err := template.ParseFiles("templates/view.html")
	if err != nil {
		msg := fmt.Sprintf("Unable to render viewer: %s", err.Error())
		http.Error(rw, msg, http.StatusInternalServerError)
	} else {
		template.Execute(rw, url)
	}
}
