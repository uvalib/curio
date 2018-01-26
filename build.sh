#!/bin/bash
env GOOS=linux go build -o dist/digital-object-viewer.linux
cp config.yml.template dist/
