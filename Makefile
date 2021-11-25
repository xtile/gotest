build: 
  go build -0 bin/ .
  
  
run: 
  go run main.go
  
  
test: 
	go test  -v -race -timeout 10s ./...



.DEFAULT_GOAL := build  
