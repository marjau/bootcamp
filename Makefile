APP_BIN = capstone

# Test controller layer
test-controller:
	go test -count=1 ./internal/controller

# Test service layer
test-service:
	go test -count=1 ./internal/service

# Test repository layer
test-repository:
	go test --race -count=1 ./internal/repository

# Test the clean architecture's layers: controller, service and repository.
test-clean-architecture:test-controller test-service test-repository

# Run all tests
test-all:
	go test --race -count=1 ./...

# Build the application
build:
	go mod download
	go mod verify
	go mod tidy -v
	go build -o $(APP_BIN) cmd/main.go

run: build
	./$(APP_BIN)

clean:
	@[ -f '$(APP_BIN)' ] && rm -v $(APP_BIN) || true

