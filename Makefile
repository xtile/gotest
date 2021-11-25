build: 
  go build -0 bin/ .
  
  
run: 
  go run main.go
  
  
test: 
	go test  -v -race -timeout 10s ./...
	
	
compile: 
	GOOS=darwin GOARCH=arm64 go build -o bin_m1/macos-arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin_x64/linux-amd64 .


clean: 
	@echo "not implemented"

.DEFAULT_GOAL := build  
