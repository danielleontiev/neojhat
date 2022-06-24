all: format test build

format:
	gofmt -s -w .

test:
	go test -count=1 ./...

build:
	go build -o ./neojhat ./cmd

cover:
	go test -coverprofile coverage.out ./...
	go tool cover -html=coverage.out

artifacts: darwin linux windows

darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./out/darwin/neojhat ./cmd

linux:
	GOOS=linux GOARCH=amd64 go build -o ./out/linux/neojhat ./cmd

windows:
	GOOS=windows GOARCH=amd64 go build -o ./out/windows/neojhat.exe ./cmd
