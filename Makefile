

run:
	go run ./cmd/instantlogsd

test:
	go test ./...

build:
	go build -o bin/ ./cmd/...

build-all:
	# https://golang.org/doc/install/source
	GOARCH=amd64 GOOS=linux   go build -o bin/instantlogs.linux64 ./cmd/instantlogs
	GOARCH=amd64 GOOS=darwin  go build -o bin/instantlogs.mac64 ./cmd/instantlogs
	GOARCH=amd64 GOOS=windows go build -o bin/instantlogs.win64.exe ./cmd/instantlogs

deps:
	go mod tidy
	go mod vendor

clean:
	rm -fr bin/*
