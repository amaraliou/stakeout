hello:
	@echo "Hello"

run:
	@go run main.go

test_handlers:
	@go test ./tests/handlerstest/... -v -coverpkg=./... -coverprofile=handlers.out

test_models:
	@go test ./tests/modelstest/... -v -coverpkg=./... -coverprofile=models.out

coverage: test_handlers test_models
	@gocovmerge ./handlers.out ./models.out > coverage.out
	@go tool cover -html=coverage.out