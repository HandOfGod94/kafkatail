APP_NAME=kafkatail

all: clean fmt tidy build install

clean:
	rm -rf out/

install:
	go install

quality-check:
	staticcheck ./...
	gocyclo -over 15 -ignore "testdata/" .
	gocognit -over 15 .

test:
	go test -timeout 60s -cover -count=1 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

tidy:
	go mod tidy -v

run: dev
	./out/$(APP_NAME)

dev:
	go build -v -o out/$(APP_NAME) main.go

build:
	go build -v -ldflags "-w" -o out/$(APP_NAME) main.go

fmt:
	go fmt ./...
	gofmt -s -w .
