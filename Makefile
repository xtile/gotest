#TODO: add env variables
DOCKER_LABEL=gotest:latest

BUILDPATH = $(CURDIR)
GO=$(shell which go)
GOINSTALL=$(GO) install
GOCLEAN=$(GO) clean
GOGET=$(GO) get

export GOPATH=$(CURDIR)


build: 
	@echo "start building... "
	git config --global url."https://xtile:89188e6ef3a334cc8d29bc857e6bf48a90dee192@github.com".insteadOf "https://github.com"	
	go build -o bin/ ./cmd/arbilogger
	@echo "done! "
  
  
run: 
	@echo "running... "
	go run main.go
  
  
test: 
	@echo "start tests...  "
	go test  -v -race -timeout 10s ./...
	@echo "tests complete...  "
	

#  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	
#  89188e6ef3a334cc8d29bc857e6bf48a90dee192
	
get: 
	@echo "Getting modules for the project..."
	echo git config --global url."https://xtile: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@github.com".insteadOf "https://github.com"	
	@$(GOGET) "github.com/BurntSushi/toml"
	echo @$(GOGET) "github.com/xtile/gotest/internal/app/arbilogger"	
	GOPRIVATE=github.com/xtile go get -u github.com/xtile/gotest/internal/app/arbilogger
	@$(GOGET) "github.com/sacOO7/gowebsocket"

makedir: 
	@echo "not implemented"


	
	
compile: 
	@echo "start building multiple architectures... "
	GOOS=darwin GOARCH=arm64 go build -o bin_m1/macos-arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin_x64/linux-amd64 .


clean: 
	@echo "cleaning files... "
	@rm -rf $(BUILDPATH)/bin
	@rm -rf $(BUILDPATH)/pkg


all: makedir get build


docker_build: 
	@echo "started docker build"
	docker build -t gotest:latest
	@echo "docker image complete"


docker_run:
	@echo "starting docker container"
	docker run -d --name gotest_running -e  . 
	@echo "complete!"


.DEFAULT_GOAL := build  
