.PHONY: swagger run build install-swag dev clean

# Cài đặt swag CLI
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger documentation
swagger:
	$(shell go env GOPATH)/bin/swag init

# Run application
run:
	go run main.go

# Build application
build:
	go build -o phikhanh main.go

# Generate swagger và run (development)
dev: swagger run

# Clean build files
clean:
	rm -rf phikhanh docs/

# Install dependencies
deps:
	go mod download
