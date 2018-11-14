# Curio

Curio is a oEmbed-enabled, standalone service to render views for digitized objects.
It supports the following endpoints:

* /healthcheck : returns a JSON object with details about the health of the service
* /version : returns the version of the service
* /view/[identifier] : display a digital object. Identifier is currently a TrackSys PID.
* /oembed : implementation of the oEmbed spec described here: https://oembed.com/
* /api/aries/:ID : implementation of the Aries API. Returns information about the ID if known

### System Requirements
* GO version 1.11.0 or greater
