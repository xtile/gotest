#TODO: add env variables
DOCKER_LABEL=gotest:latest

BUILDPATH = $(CURDIR)
GO=$(shell which go)
GOINSTALL=$(GO) install
GOCLEAN=$(GO) clean
GOGET=$(GO) get

#export GOPATH=$(CURDIR)


build: clean
	@echo "start building... "
	go build -o .bin/ ./cmd/gotest
	@echo "done! "
  
  
run: 
	@echo "running... "
	go run main.go
  
  
test: 
	@echo "start tests...  "
	go test  -v -race -timeout 10s ./...
	@echo "tests complete...  "
	


	
init: 
	@echo "init module"
	go mod init arbi
	go mod tidy


get: 
	@echo "Getting modules for the project..."
	@$(GOGET) "github.com/BurntSushi/toml"

	@$(GOGET) "github.com/xtile/gotest/internal/app/gotest"	
	@$(GOGET) "github.com/sacOO7/gowebsocket"

makedir: 
	@echo "creating dirs... "
	@if [ ! -d $(BUILDPATH)/bin ]; then mkdir -p $(BUILDPATH)/bin; fi
	@if [ ! -d $(BUILDPATH)/pkg ]; then mkdir -p $(BUILDPATH)/pkg; fi


	
	
compile: 
	@echo "start building multiple architectures... "
	@echo ".........................................................................................  "
	@echo "Start building for AMD64...  "
	GOOS=darwin GOARCH=arm64 go build -o ./bin/macos-arm64/ -v ./...
	@echo ".........................................................................................  "
	@echo "Start building for AMD64...  "
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux-amd64/ -v ./...
	@echo "All builds finished!  "


clean: 
	@echo "cleaning files... "
	@rm -rf $(BUILDPATH)/bin
	@rm -rf $(BUILDPATH)/pkg
	go clean -modcache


all: makedir get build


docker_build: 
	@echo "started docker build"
	@echo docker build -t gotest:latest
	docker build -t gotest:latest .
	@echo "docker image complete"


docker_run:
	@echo "starting docker container"
	docker run -d --name gotest_running -e  . 
	@echo "complete!"


.DEFAULT_GOAL := build  
