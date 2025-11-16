
dev:
	air

run-local:
	go run cmd/server/main.go

# Build commands
build:
	go build -o bin/server cmd/server/main.go

# Test commands
test:
	go test -v ./...

# Dependency management
deps:
	go mod download
	go mod tidy

# Clean
clean:
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

