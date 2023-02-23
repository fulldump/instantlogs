

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

.PHONY: release
release: clean
	# daemon
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build $(FLAGS) -o bin/instantlogsd.linux.arm64 ./cmd/instantlogsd
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build $(FLAGS) -o bin/instantlogsd.linux.amd64 ./cmd/instantlogsd
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(FLAGS) -o bin/instantlogsd.win.arm64.exe ./cmd/instantlogsd
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(FLAGS) -o bin/instantlogsd.win.amd64.exe ./cmd/instantlogsd
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build $(FLAGS) -o bin/instantlogsd.mac.arm64 ./cmd/instantlogsd
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build $(FLAGS) -o bin/instantlogsd.mac.amd64 ./cmd/instantlogsd
	# client
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build $(FLAGS) -o bin/instantlogs.linux.arm64 ./cmd/instantlogs
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build $(FLAGS) -o bin/instantlogs.linux.amd64 ./cmd/instantlogs
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build $(FLAGS) -o bin/instantlogs.win.arm64.exe ./cmd/instantlogs
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(FLAGS) -o bin/instantlogs.win.amd64.exe ./cmd/instantlogs
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build $(FLAGS) -o bin/instantlogs.mac.arm64 ./cmd/instantlogs
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build $(FLAGS) -o bin/instantlogs.mac.amd64 ./cmd/instantlogs
	md5sum bin/instantlogs.* > bin/checksum-md5
	sha256sum bin/instantlogs.* > bin/checksum-sha256
