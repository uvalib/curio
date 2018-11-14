package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const Version = "2.0.0"

type configData struct {
	port        int
	tracksysURL string
	apolloURL   string
	iiifURL     string
	fedoraURL   string
	hostname    string
}

// golbals for DB and CFG
var config configData

func main() {
	// Load cfg
	log.Printf("===> Curio is staring up <===")
	getConfiguration()

	// Set routes and start server
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()

	// load all of the templates
	router.LoadHTMLFiles("./templates/image_view.html", "./templates/wsls_view.html", "./templates/not_available.html")

	// Set routes and start server
	router.Use(cors.Default())
	router.GET("/", versionHandler)
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/version", versionHandler)
	router.GET("/healthcheck", healthCheckHandler)
	router.Use(static.Serve("/public", static.LocalFile("./web", true)))
	router.Use(static.Serve("/view/uv", static.LocalFile("./web/uv", true)))
	router.GET("/view/:pid", viewHandler)
	router.GET("/oembed", oEmbedHandler)
	api := router.Group("/api")
	{
		api.GET("/aries", ariesPing)
		api.GET("/aries/:id", ariesLookup)
	}

	portStr := fmt.Sprintf(":%d", config.port)
	log.Printf("Start HTTP service on port %s with CORS support enabled", portStr)
	log.Fatal(router.Run(portStr))
}

func getConfiguration() {
	defTracksysURL := os.Getenv("TRACKSYS_URL")
	if defTracksysURL == "" {
		defTracksysURL = "https://tracksys.lib.virginia.edu/api"
	}

	defApolloURL := os.Getenv("APOLLO_URL")
	if defApolloURL == "" {
		defApolloURL = "https://apollo.lib.virginia.edu/api"
	}

	defIIIFURL := os.Getenv("CURIO_IIIF_MAN_URL")
	if defIIIFURL == "" {
		defIIIFURL = "https://iiifman.lib.virginia.edu/pid"
	}

	defFedoraURL := os.Getenv("WSLS_FEDORA_URL")
	if defFedoraURL == "" {
		defFedoraURL = "http://fedora01.lib.virginia.edu/wsls"
	}

	defHost := os.Getenv("CURIO_HOST")
	if defHost == "" {
		defHost = "curio.lib.virginia.edu"
	}

	// FIRST, try command line flags. Fallback is ENV variables
	flag.IntVar(&config.port, "port", 8085, "Port to offer service on (default 8085)")
	flag.StringVar(&config.tracksysURL, "tracksys", defTracksysURL, "TrackSys URL")
	flag.StringVar(&config.apolloURL, "apollo", defApolloURL, "Apollo URL")
	flag.StringVar(&config.iiifURL, "iiif", defIIIFURL, "IIIF Manifest URL")
	flag.StringVar(&config.fedoraURL, "fedora", defFedoraURL, "WSLS Fedora URL")
	flag.StringVar(&config.hostname, "host", defHost, "Curio Hostname")
	flag.Parse()
}

// Handle a request for / and return version info
func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Curio version %s", Version)
}

// faviconHandler is a dummy handler to silence browser API requests that look for /favicon.ico
func faviconHandler(c *gin.Context) {
}
