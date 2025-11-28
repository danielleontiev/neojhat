.PHONY: check format vet tidy test build cover

check: format vet tidy test build

format:
	go fmt ./...
	git diff --exit-code

vet:
	go vet ./...

tidy:
	go mod tidy -diff

test:
	go test -count=1 ./...

build:
	go build .

cover:
	go test -coverprofile coverage.out ./...
	go tool cover -html=coverage.out
