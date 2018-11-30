GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_DIR=build
BINARY_PUBLISH=publish
BINARY_PROCESSENGINE=processengine
BINARY_ACTIONHANDLER=actionhandler
    
all: build
build: publish processengine actionhandler

publish:
		$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_PUBLISH) -v cmd/publish/publish.go 
processengine:
		$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_PROCESSENGINE) -v cmd/processengine/processengine.go
actionhandler:
		$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_ACTIONHANDLER) -v cmd/actionhandler/actionhandler.go

clean: 
		$(GOCLEAN)
		rm -rf $(BUILD_DIR)

# Cross compilation
build-linux: export CGO_ENABLED=0
build-linux: export GOOS=linux
build-linux: export GOARCH=amd64
build-linux: export BUILD_DIR=build/linux
build-linux: build