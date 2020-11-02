GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT = $(GOCMD) fmt
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOVET = $(GOCMD) vet

BASE_NAME=curio

build: darwin copy-web

all: darwin linux copy-web

linux-full: linux copy-web

darwin-full: darwin copy-web

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/$(BASE_NAME).darwin viewsrv/*.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/$(BASE_NAME).linux viewsrv/*.go

copy-web:
	mkdir -p bin/
	rm -rf bin/web
	rm -rf bin/templates
	cp -R web/ bin/web/
	cp -R templates bin/templates/

clean:
	rm -rf bin

fmt:
	cd viewsrv; $(GOFMT)

vet:
	cd viewsrv; $(GOVET)

dep:
	$(GOGET) -u ./viewsrv/...
	$(GOMOD) tidy
	$(GOMOD) verify
