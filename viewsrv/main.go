package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const Version = "3.1.1"

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
	router.StaticFile("/favicon.ico", "./web/favicon.ico")
	router.GET("/", versionHandler)
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
	log.Printf("INFO: start Curio on port %s with CORS support enabled", portStr)
	log.Fatal(router.Run(portStr))
}

// Handle a request for / and return version info
func versionHandler(c *gin.Context) {

	build := "unknown"

	// cos our CWD is the bin directory
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}

	vMap := make(map[string]string)
	vMap["version"] = Version
	vMap["build"] = build
	c.JSON(http.StatusOK, vMap)
}
