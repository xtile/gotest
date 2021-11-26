#TODO: add env variables
DOCKER_LABEL=gotest:latest



build: 
  @echo "start building... "
  go build -0 bin/ .
  @echo "done! "
  
  
run: 
  @echo "running... "
  go run main.go
  
  
test: 
	@echo "start tests...  "
	go test  -v -race -timeout 10s ./...
	@echo "tests complete...  "
	
	
get: 
	@echo "not implemented"


makedir: 
	@echo "not implemented"


	
	
compile: 
	@echo "start building multiple architectures... "
	GOOS=darwin GOARCH=arm64 go build -o bin_m1/macos-arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin_x64/linux-amd64 .


clean: 
	@echo "not implemented"


all: makedir get build


docker build: 
	@echo "started docker build"
	docker build -t gotest:latest
	@echo "docker image complete"


docker run
	@echo "starting docker container"
	docker run -d --name gotest_running -e  . 
	@echo "complete!"


.DEFAULT_GOAL := build  
