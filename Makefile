GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get

# project specific definitions
BASE_NAME=digital-object-viewer
SRC_TREE=cmd/viewsrv

build: build-darwin build-linux copy-web

build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/$(BASE_NAME).darwin $(SRC_TREE)/*

copy-web:
	cp -R web/ bin/web/

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/$(BASE_NAME).linux $(SRC_TREE)/*

fmt:
	$(GOFMT) $(SRC_TREE)/*

vet:
	$(GOVET) $(SRC_TREE)/*

clean:
	$(GOCLEAN)
	rm -rf bin/

deps:
	dep ensure
	dep status
