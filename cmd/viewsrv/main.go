package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	htemplate "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

// Version of the service
const Version = "1.0.0"

type oEmbedData struct {
	PID       string
	Title     string
	Author    sql.NullString
	HTML      string
	URL       string
	Width     int
	Height    int
	SourceURI string
	Scheme    string
	EmbedHost string
	StartPage int
}

type configData struct {
	port    int
	dbHost  string
	dbName  string
	dbUser  string
	dbPass  string
	iiifURL string
	dovHost string
}

type viewerData struct {
	URI       string
	StartPage int
}

// golbals for DB and CFG
var mysqlDB *sql.DB
var config configData

func main() {
	// Load cfg
	log.Printf("===> viewer staring up <===")
	getConfiguration()

	// Init DB connection
	var err error
	log.Printf("Init DB connection to %s...", config.dbHost)
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.dbUser, config.dbPass, config.dbHost, config.dbName)
	mysqlDB, err = sql.Open("mysql", connectStr)
	if err != nil {
		log.Printf("FATAL: Database connection failed: %s", err.Error())
		os.Exit(1)
	}
	_, err = mysqlDB.Query("SELECT 1")
	if err != nil {
		log.Printf("FATAL: Database query failed: %s", err.Error())
		os.Exit(1)
	}
	defer mysqlDB.Close()
	log.Printf("DB Connection established")

	// Set routes and start server
	mux := httprouter.New()
	mux.GET("/", loggingHandler(rootHandler))
	mux.GET("/images/:id", loggingHandler(imagesHandler))
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
	flag.StringVar(&config.dbHost, "dbhost", os.Getenv("DB_HOST"), "DB Host (required)")
	flag.StringVar(&config.dbName, "dbname", os.Getenv("DB_NAME"), "DB Name (required)")
	flag.StringVar(&config.dbUser, "dbuser", os.Getenv("DB_USER"), "DB User (required)")
	flag.StringVar(&config.dbPass, "dbpass", os.Getenv("DB_PASS"), "DB Password (required)")
	flag.StringVar(&config.iiifURL, "iiif", os.Getenv("IIIF"), "IIIF URL (required)")
	flag.StringVar(&config.dovHost, "dovhost", os.Getenv("DOV_HOST"), "DoViewer Hostname (optional)")
	flag.Parse()

	// if anything is still not set, die
	if len(config.dbHost) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if len(config.dbName) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if len(config.dbUser) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if len(config.dbPass) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if len(config.iiifURL) == 0 {
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

/**
 * Handle a request for / and return version info
 */
func rootHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprintf(rw, "UVA Viewer version %s", Version)
}

/**
 * Handle a request for oembed data
 */
func oEmbedHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Get some optional params; format, maxWidth and maxHeight
	respFormat := req.URL.Query().Get("format")
	maxWidth, err := strconv.Atoi(req.URL.Query().Get("maxwidth"))
	if err != nil {
		maxWidth = 0
	}
	maxHeight, err := strconv.Atoi(req.URL.Query().Get("maxheight"))
	if err != nil {
		maxHeight = 0
	}

	// Next, get the required URL and see if a page is requested
	urlStr, _ := url.QueryUnescape(req.URL.Query().Get("url"))
	if len(urlStr) == 0 {
		http.Error(rw, "URL is required!", http.StatusBadRequest)
		return
	}
	log.Printf("Target resource URL: %s", urlStr)

	if len(respFormat) == 0 || strings.Compare(respFormat, "json") == 0 {
		log.Printf("JSON response requested")
		renderOembedResponse(urlStr, "json", maxWidth, maxHeight, rw, req)
	} else if strings.Compare(respFormat, "xml") == 0 {
		log.Printf("XML response requested")
		renderOembedResponse(urlStr, "xml", maxWidth, maxHeight, rw, req)
	} else {
		// error: unsupported format
		err := fmt.Sprintf("Requested format '%s' is invalid.", respFormat)
		http.Error(rw, err, http.StatusBadRequest)
	}
}

func renderOembedResponse(rawURL string, format string, maxWidth int, maxHeight int, rw http.ResponseWriter, req *http.Request) {
	// init data used to render the oEmbed response
	var data oEmbedData

	// The raw URL requested must be of the expected format: [http|https]://[host]/images/[PID][?page=n]
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		msg := fmt.Sprintf("Invalid URL: %s", err.Error())
		http.Error(rw, msg, http.StatusInternalServerError)
		return
	}

	// Get page param if any...
	qp, _ := url.ParseQuery(parsedURL.RawQuery)
	data.StartPage = 0
	if len(qp["page"]) > 0 {
		data.StartPage, _ = strconv.Atoi(qp["page"][0])
	}

	// accept 1 based page numbers from client, but use
	// 0-based canvas index in UV embed snippet
	if data.StartPage > 0 {
		data.StartPage--
		log.Printf("Requested starting page index %d", data.StartPage)
	}

	// Now split out relatve path to find PID. This should be something like: /images/[PID]
	// NOTE: that this wil strip out all query params
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

	// URL for IIIF manifest
	data.PID = bits[2]
	data.SourceURI = fmt.Sprintf("%s/%s", config.iiifURL, data.PID)

	// Validate that the manifest has images
	if isManifestViewable(data.SourceURI) == false {
		log.Printf("Requested URL %s has no visible images", data.SourceURI)
		http.Error(rw, "Sorry, the requested resource is not available.", http.StatusNotFound)
		return
	}

	// scheme / host for UV javascript
	data.Scheme = "http"
	if strings.Contains(data.SourceURI, "https") {
		data.Scheme = "https"
	}
	data.EmbedHost = config.dovHost
	if len(data.EmbedHost) == 0 {
		data.EmbedHost = req.Host
	}

	// default embed size is 800x600. Params maxwidth and maxheight can override.
	data.Width = 800
	if maxWidth > 0 && maxWidth < data.Width {
		data.Width = maxWidth
	}
	data.Height = 600
	if maxHeight > 0 && maxHeight < data.Height {
		data.Height = maxHeight
	}

	log.Printf("Retrieving metadata for %s...", data.PID)
	qs := `select m.title, m.creator_name from metadata m where m.pid = ? group by m.id`
	queryErr := mysqlDB.QueryRow(qs, data.PID).Scan(&data.Title, &data.Author)
	if queryErr != nil {
		log.Printf("Request failed: %s", queryErr.Error())
		if strings.Contains(queryErr.Error(), "no rows") {
			msg := fmt.Sprintf("Invalid ID %s", data.PID)
			http.Error(rw, msg, http.StatusBadRequest)
		} else {
			msg := fmt.Sprintf("Unable to retreive oEmbed response: %s", queryErr.Error())
			http.Error(rw, msg, http.StatusInternalServerError)
		}
		return
	}

	// Render the <div> that will be included in the response, and used to embed the resource
	log.Printf("Rendering html snippet...")
	var renderedSnip bytes.Buffer
	snippet := htemplate.Must(htemplate.ParseFiles("templates/embed.html"))
	snipErr := snippet.Execute(&renderedSnip, data)
	if snipErr != nil {
		http.Error(rw, snipErr.Error(), http.StatusInternalServerError)
		return
	}
	rawHTML := strings.TrimSpace(renderedSnip.String())

	if strings.Compare(format, "json") == 0 {
		log.Printf("Rendering JSON output")
		data.HTML = strconv.Quote(rawHTML)
		rw.Header().Set("content-type", "application/json; charset=utf-8")
		jsonTemplate := template.Must(template.ParseFiles("templates/response.json"))
		jsonTemplate.Execute(rw, data)
	} else {
		rw.Header().Set("content-type", "text/xml; charset=utf-8")
		log.Printf("Rendering XML output")
		data.HTML = rawHTML
		var renderedSnip bytes.Buffer
		snippet := htemplate.Must(htemplate.ParseFiles("templates/response.xml"))
		snipErr := snippet.Execute(&renderedSnip, data)
		if snipErr != nil {
			log.Printf("Unable to render XML template: %s", snipErr.Error())
			http.Error(rw, snipErr.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprint(rw, renderedSnip.String())
		}
	}
}

/**
 * Handle a request for images from a specific ID (TrackSys PID) and page offset (optional)
 */
func imagesHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// pull page QP and use it for starting page. Any other params are ignored.
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	var data viewerData
	data.StartPage = page - 1
	data.URI = fmt.Sprintf("%s/%s", config.iiifURL, params.ByName("id"))

	// Make sure there are images visable for this PID.
	// Ahow an error page and bail if not
	if isManifestViewable(data.URI) == false {
		rw.WriteHeader(http.StatusNotFound)
		bytes, _ := ioutil.ReadFile("web/not_available.html")
		fmt.Fprintf(rw, "%s", string(bytes))
		return
	}

	template, err := template.ParseFiles("templates/view.html")
	if err != nil {
		msg := fmt.Sprintf("Unable to render viewer: %s", err.Error())
		http.Error(rw, msg, http.StatusInternalServerError)
	} else {
		template.Execute(rw, data)
	}
}

// Hit the target IIIF manifest URL and see if it contains any images
//
func isManifestViewable(url string) bool {
	timeout := time.Duration(15 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	log.Printf("Checking manifest URL %s for images", url)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("ERROR: IIIF URL: %s failed to return a response", url)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("ERROR: IIIF URL: %s returned non-success status: %d", url, resp.StatusCode)
		return false
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	respStr := string(bytes)

	return strings.Contains(respStr, "dcTypes:Image")
}

/**
 * Check health of service
 */
func healthCheckHandler(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")

	// make sure DB is connected
	log.Printf("Checking DB...")
	dbStatus := true
	_, err := mysqlDB.Query("SELECT 1")
	if err != nil {
		log.Printf("ERROR: DB access (%s)", err)
		dbStatus = false
	}

	// make sure IIIF manifest service is alive
	log.Printf("Checking IIIF...")
	iiifStatus := true
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(config.iiifURL)
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
	out := fmt.Sprintf(`{"alive": true, "mysql": %t, "iiif": %t}`, dbStatus, iiifStatus)
	if dbStatus == false || iiifStatus == false {
		http.Error(rw, out, http.StatusInternalServerError)
	} else {
		fmt.Fprintf(rw, out)
	}
}
