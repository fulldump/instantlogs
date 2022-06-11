

run:
	go run ./cmd/instantlogsd

test:
	go test ./...

build:
	go build -o bin/ ./cmd/...

deps:
	go mod tidy
	go mod vendor
