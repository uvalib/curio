GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BASE_NAME=digital-object-viewer

build: darwin copy-web

all: darwin linux copy-web

linux-full: linux copy-web

darwin-full: darwin copy-web

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/$(BASE_NAME).darwin cmd/viewsrv/*.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/$(BASE_NAME).linux cmd/viewsrv/*.go

copy-web:
	mkdir -p bin/
	rm -rf bin/web
	rm -rf bin/templates
	cp -R web/ bin/web/
	cp -R templates bin/templates/

clean:
	rm -rf bin
