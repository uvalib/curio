package main

import (
	"flag"
	"log"
)

type configData struct {
	port          int
	apolloURL     string
	iiifURL       string
	iiifCacheURL  string
	cacheDisabled bool
	wslsURL       string
	hostname      string
	rightsURL     string
}

// globals for the CFG
var config configData

func getConfiguration() {
	flag.IntVar(&config.port, "port", 8085, "Port to offer service on (default 8085)")
	flag.StringVar(&config.apolloURL, "apollo", "https://apollo.lib.virginia.edu/api", "Apollo URL")
	flag.StringVar(&config.iiifURL, "iiif", "https://iiifman.lib.virginia.edu", "IIIF manifest URL")
	flag.StringVar(&config.iiifCacheURL, "iiifcache", "https://s3.us-east-1.amazonaws.com/iiif-manifest-cache-staging", "IIIF manifest cache URL")
	flag.StringVar(&config.wslsURL, "fedora", "https://wsls.lib.virginia.edu", "WSLS Fedora URL")
	flag.StringVar(&config.rightsURL, "rights", "https://rights-wrapper.lib.virginia.edu/api/pid", "Rights wrapper URL")
	flag.StringVar(&config.hostname, "host", "curio.lib.virginia.edu", "Curio hostname")
	flag.BoolVar(&config.cacheDisabled, "nocache", false, "Local dev mode flag to disable IIIF cache")
	flag.Parse()

	log.Printf("[CONFIG] port          = [%d]", config.port)
	log.Printf("[CONFIG] apolloURL     = [%s]", config.apolloURL)
	log.Printf("[CONFIG] iiifURL       = [%s]", config.iiifURL)
	log.Printf("[CONFIG] iiifCacheURL   = [%s]", config.iiifCacheURL)
	log.Printf("[CONFIG] cacheDisabled = [%t]", config.cacheDisabled)
	log.Printf("[CONFIG] wslsURL       = [%s]", config.wslsURL)
	log.Printf("[CONFIG] rightsURL     = [%s]", config.rightsURL)
	log.Printf("[CONFIG] hostname      = [%s]", config.hostname)
}

//
// end of file
//
