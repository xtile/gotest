build: 
  go build -0 bin/ .
  
  
run: 
  go run main.go
  
  
test: 
	go test  -v -race -timeout 10s ./...

clean: 
	@echo "not implemented"

.DEFAULT_GOAL := build  
