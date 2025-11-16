build:
	go build .

run:
	APP_ENV=dev go run main.go

lint:
	gofmt -w .
	goimports -w .
	golangci-lint run -c .golangci.yml

generate-mocks:
	mockgen -source=app/handler/handler.go -destination=app/mocks/mock_processor.go -package=mocks
	mockgen -source=app/processor/message_processor.go -destination=app/mocks/mock_service.go -package=mocks
	mockgen -source=app/service/message_service.go -destination=app/mocks/mock_repository.go -package=mocks

unit-test:
	go test -v ./app/handler/... ./app/processor/...  ./app/service/...  ./ -short

repository-test:
	go test -v ./app/repository -run TestRepository