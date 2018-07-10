package apisvc

// TrackSysMetadata contains the basic metadata returned from the Tracksys API
type TrackSysMetadata struct {
	Title  string
	Author string
}

// WSLSMetadata contains the Apollo metadata supporting WSLS
type WSLSMetadata struct {
	HasVideo      bool
	HasScript     bool
	WSLSID        string
	Title         string
	Description   string
	VideoURL      string
	PosterURL     string
	PDFURL        string
	PDFThumbURL   string
	TranscriptURL string
}
