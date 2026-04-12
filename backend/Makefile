APP_NAME=seneprojects
CMD_PATH=./cmd/api

# Run app
run:
	go run $(CMD_PATH)/main.go

# Build binary
build:
	go build -o $(APP_NAME) $(CMD_PATH)/main.go

# Run tests
test:
	go test ./internal/...

# Verbose tests
test-verbose:
	go test -v ./internal/...

# Coverage
test-coverage:
	go test -coverpkg=./internal/auth/...,./internal/users/... -coverprofile=coverage.out ./internal/auth/... ./internal/auth/test/... ./internal/users/... ./internal/users/test/...
	go tool cover -func=coverage.out
