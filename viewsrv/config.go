package main

import (
	"flag"
	"log"
)

type configData struct {
	port                  int
	apolloURL             string
	iiifURL               string
	wslsURL               string
	hostname              string
	rightsURL             string
	archivematicaS3Bucket string
}

// globals for the CFG
var config configData

func getConfiguration() {
	flag.IntVar(&config.port, "port", 8080, "Port to offer service on (default 8085)")
	flag.StringVar(&config.apolloURL, "apollo", "https://apollo.lib.virginia.edu/api", "Apollo URL")
	flag.StringVar(&config.iiifURL, "iiif", "https://iiifman.lib.virginia.edu", "IIIF manifest URL")
	flag.StringVar(&config.wslsURL, "fedora", "https://wsls.lib.virginia.edu", "WSLS Fedora URL")
	flag.StringVar(&config.rightsURL, "rights", "https://rights-wrapper.lib.virginia.edu/api/pid", "Rights wrapper URL")
	flag.StringVar(&config.archivematicaS3Bucket, "archivematicaBucket", "", "Archivematica S3 Bucket")
	flag.StringVar(&config.hostname, "host", "curio.lib.virginia.edu", "Curio hostname")
	flag.Parse()

	log.Printf("[CONFIG] port                  = [%d]", config.port)
	log.Printf("[CONFIG] apolloURL             = [%s]", config.apolloURL)
	log.Printf("[CONFIG] iiifURL               = [%s]", config.iiifURL)
	log.Printf("[CONFIG] wslsURL               = [%s]", config.wslsURL)
	log.Printf("[CONFIG] rightsURL             = [%s]", config.rightsURL)
	log.Printf("[CONFIG] archivematicaS3Bucket = [%s]", config.archivematicaS3Bucket)
	log.Printf("[CONFIG] hostname              = [%s]", config.hostname)
}

//
// end of file
//
