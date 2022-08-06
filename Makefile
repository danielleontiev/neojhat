all: format test build

format:
	gofmt -s -w .

test:
	go test -count=1 ./...

build:
	go build .

cover:
	go test -coverprofile coverage.out ./...
	go tool cover -html=coverage.out

artifacts: darwin linux

darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./out/darwin/amd64/neojhat .
	GOOS=darwin GOARCH=arm64 go build -o ./out/darwin/arm64/neojhat .

linux:
	GOOS=linux GOARCH=amd64 go build -o ./out/linux/amd64/neojhat .
	GOOS=linux GOARCH=arm64 go build -o ./out/linux/arm64/neojhat .

release: artifacts
	tar -cavf ./out/neojhat.darwin.amd64.tar.gz -C ./out/darwin/amd64 neojhat
	tar -cavf ./out/neojhat.darwin.arm64.tar.gz -C ./out/darwin/arm64 neojhat
	tar -cavf ./out/neojhat.linux.amd64.tar.gz -C ./out/linux/amd64 neojhat
	tar -cavf ./out/neojhat.linux.arm64.tar.gz -C ./out/linux/arm64 neojhat

