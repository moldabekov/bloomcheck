# MAKEFILE
#
# Author: Margulan Moldabekov <moldabekov [ a t ] gmail.com>
# Home: https://github.com/moldabekov/bloomcheck
#

# Program Name
PROGRAM = bloomcheck

# Go files wildcard
SOURCE = *.go

# Go defs
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Go options
CGO=CGO_ENABLED=0
LDFLAGS=-ldflags='-extldflags "-static -s -w"'
OPTIONS=-a -installsuffix cgo

BINARY_NAME=$(PROGRAM)
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WINDOWS=$(BINARY_NAME).exe
DOCKER_NAME=mldbk/$(PROGRAM)

# GO lang path
ifneq ($(GOPATH),)
	ifeq ($(findstring $(GOPATH),$(CURRENTDIR)),)
		# the defined GOPATH is not valid
		GOPATH=
	endif
endif
ifeq ($(GOPATH),)
	# extract the GOPATH
	GOPATH=$(firstword $(subst /src/, ,$(CURRENTDIR)))
endif


# Targets
all: build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
	strip $(PROGRAM)

fmt:
	gofmt -w $(SOURCE)

vet:
	go vet $(SOURCE)

run:
	go run $(SOURCE)

tag:
	docker tag mldbk/$(PROGRAM):latest

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_WINDOWS)
	rm -rf *.pem
	docker system prune --all

run:
	$(GOBUILD) -o $(BINARY_NAME) -v .
	./$(BINARY_NAME)

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(OPTIONS) -o $(BINARY_NAME) -v .

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) $(OPTIONS) -o $(BINARY_WINDOWS) -v .

docker: linux
	upx -9 $(BINARY_NAME)
	docker build -t $(DOCKER_NAME) .
	#docker run --rm -p 9876:9876 $(DOCKER_NAME)

deps:
	GOPATH=$(GOPATH) $(GOGET) -u github.com/gorilla/mux
	GOPATH=$(GOPATH) $(GOGET) -u github.com/willf/bloom
	GOPATH=$(GOPATH) $(GOGET) -u github.com/gorilla/handlers
