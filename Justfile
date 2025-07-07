run:
  go run cmd/commitmsg/main.go

build:
  go build -o ./bin/commitmsg cmd/commitmsg/main.go

format:
  go fmt ./...

lint:
  golangci-lint run ./...
