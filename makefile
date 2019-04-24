	# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=jotun
BINARY_UNIX=$(BINARY_NAME)_unix
OS=$(shell uname -s)
VERSION=1.0.0-alpha


all: runtest build run
.PHONY: build
build:
	$(info Building for: $(OS))
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/jotun/... 
runtest: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f ./bin/$(BINARY_NAME)
	rm -f ./bin/$(BINARY_UNIX)
run:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/jotun/...
	./bin/$(BINARY_NAME)
runprd:
	$(info Building for: $(OS))
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/jotun/...
	tar cvzf "./bin/jotun-$(VERSION).tar.gz" ./LICENSE ./jotun.1 ./bin/jotun ./installer.sh
makedeb:
	$(info Creating DEBIAN package)
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/jotun/...
	if [ ! -d "./deb-release/jotun-${VERSION}" ]; then echo Creating deb dir..; mkdir ./deb-release; mkdir ./deb-release/jotun-${VERSION}; mkdir ./deb-release/jotun-${VERSION}/usr/; mkdir ./deb-release/jotun-${VERSION}/usr/local; mkdir ./deb-release/jotun-${VERSION}/usr/share; mkdir ./deb-release/jotun-${VERSION}/usr/share/man; mkdir ./deb-release/jotun-${VERSION}/usr/share/man/man1; mkdir ./deb-release/jotun-${VERSION}/usr/local/bin; mkdir ./deb-release/jotun-${VERSION}/DEBIAN; fi
	cp control ./deb-release/jotun-${VERSION}/DEBIAN/
	sed -i -e 's/versionhere/${VERSION}/g' ./deb-release/jotun-${VERSION}/DEBIAN/control
	cp ./bin/jotun ./deb-release/jotun-$(VERSION)/usr/local/bin/
	gzip -c jotun.1 > jotun.1.gz
	mv jotun.1.gz ./deb-release/jotun-$(VERSION)/usr/share/man/man1/
	dpkg-deb --build ./deb-release/jotun-${VERSION}/
	cp ./deb-release/*.deb ./bin/
	echo Deb package created!!