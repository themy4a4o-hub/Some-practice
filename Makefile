lint:
	golangci-lint run
test: 
	go test -v ./...
check: lint test