# Make a release with
# make -j4 release

VERSION=$(shell git describe)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

.PHONY: build
build: bindata.go
	go build ${LDFLAGS}

STATICFILES := $(wildcard static/*)
TEMPLATES := $(wildcard templates/*)
bindata.go: $(STATICFILES) $(TEMPLATES)
	go-bindata -tags '!debug' static/... templates/...
	mv bindata.go src/
	sed -i 's/package main/package server/g' src/bindata.go