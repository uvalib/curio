GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT = $(GOCMD) fmt
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOVET = $(GOCMD) vet

BASE_NAME=curio

build: darwin web

all: darwin linux web

linux-full: linux web

darwin-full: darwin web

web:
	mkdir -p bin/
	cd frontend && yarn install && yarn build
	rm -rf bin/public
	mv frontend/dist bin/public

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/$(BASE_NAME).darwin viewsrv/*.go
	cp -r templates/ bin/templates

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/$(BASE_NAME).linux viewsrv/*.go
	cp -r templates/ bin/templates

clean:
	rm -rf bin

fmt:
	cd viewsrv; $(GOFMT)

vet:
	cd viewsrv; $(GOVET)

dep:
	cd frontend && yarn upgrade
	$(GOGET) -u ./viewsrv/...
	$(GOMOD) tidy
	$(GOMOD) verify
