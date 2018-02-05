# UVA Digital Object Viewer

This is a oEmbed-enabled, standalone service to render views for digitized objects.
It supports the following endpoints:

* / : returns version information
* /healthcheck : returns a JSON object with details about the health of the service
* /images/[identifier] : display a digital object. Identifier is currently a TrackSys PID.
* /oembed : implementation of the oEmbed spec described here: https://oembed.com/
