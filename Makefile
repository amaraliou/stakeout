hello:
	@echo "Hello"

run:
	@go run main.go

test_handlers:
	@go test ./tests/handlerstest/... -v -coverpkg=./... -coverprofile=handlers.out

test_models:
	@go test ./tests/modelstest/... -v -coverpkg=./... -coverprofile=models.out

coverfile: test_handlers test_models
	@gocovmerge ./handlers.out ./models.out > coverage.out

test_junit:
	@rm -rf test-results
	@mkdir test-results
	@gotestsum --format standard-verbose --junitfile ./test-results/handlers-tests.xml ./tests/handlerstest/...
	@gotestsum --format standard-verbose --junitfile ./test-results/models-tests.xml ./tests/modelstest/...

coverage: coverfile
	@go tool cover -html=coverage.out