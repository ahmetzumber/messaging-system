build:
	go build .

run:
	go run .

lint:
	gofmt -w .
	goimports -w .
	golangci-lint run -c .golangci.yml