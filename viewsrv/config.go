package main

import (
	"flag"
	"log"
	"os"
)

type configData struct {
	port          int
	apolloURL     string
	iiifURL       string
	iiifRootURL   string
	cacheDisabled bool
	wslsURL       string
	hostname      string
}

// globals for the CFG
var config configData

func getConfiguration() {

	defApolloURL := os.Getenv("APOLLO_URL")
	if defApolloURL == "" {
		defApolloURL = "https://apollo.lib.virginia.edu/api"
	}

	defIIIFURL := os.Getenv("CURIO_IIIF_MAN_URL")
	if defIIIFURL == "" {
		defIIIFURL = "https://iiifman.lib.virginia.edu"
	}

	defIIIFRootURL := os.Getenv("CURIO_IIIF_MAN_ROOT_URL")
	if defIIIFRootURL == "" {
		defIIIFRootURL = "https://s3.us-east-1.amazonaws.com/iiif-manifest-cache-staging"
	}

	defFedoraURL := os.Getenv("WSLS_FEDORA_URL")
	if defFedoraURL == "" {
		defFedoraURL = "https://wsls.lib.virginia.edu"
	}

	defHost := os.Getenv("CURIO_HOST")
	if defHost == "" {
		defHost = "curio.lib.virginia.edu"
	}

	// FIRST, try command line flags. Fallback is ENV variables
	flag.IntVar(&config.port, "port", 8085, "Port to offer service on (default 8085)")
	flag.StringVar(&config.apolloURL, "apollo", defApolloURL, "Apollo URL")
	flag.StringVar(&config.iiifURL, "iiif", defIIIFURL, "IIIF manifest URL")
	flag.StringVar(&config.iiifRootURL, "iiifroot", defIIIFRootURL, "IIIF manifest root URL")
	flag.StringVar(&config.wslsURL, "fedora", defFedoraURL, "WSLS Fedora URL")
	flag.StringVar(&config.hostname, "host", defHost, "Curio hostname")
	flag.BoolVar(&config.cacheDisabled, "nocache", false, "Local dev mode flag to disable IIIF cache")
	flag.Parse()

	log.Printf("[CONFIG] port          = [%d]", config.port)
	log.Printf("[CONFIG] apolloURL     = [%s]", config.apolloURL)
	log.Printf("[CONFIG] iiifURL       = [%s]", config.iiifURL)
	log.Printf("[CONFIG] iiifRootURL   = [%s]", config.iiifRootURL)
	log.Printf("[CONFIG] cacheDisabled = [%t]", config.cacheDisabled)
	log.Printf("[CONFIG] wslsURL     = [%s]", config.wslsURL)
	log.Printf("[CONFIG] hostname      = [%s]", config.hostname)
}

//
// end of file
//
