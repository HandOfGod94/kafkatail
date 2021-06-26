APP_NAME=kafkatail

all: clean fmt tidy build

clean:
	rm -rf out/

quality-check:
	staticcheck ./...
	gocyclo -over 15 .
	gocognit -over 15 .

.PHONY: integration-tests
integration-tests:
	go install
	go test -tags=integration ./integration-tests


test:
	go install
	go test -timeout 30s -cover -count=1 -coverprofile=coverage.out ./...
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

setup-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/vektra/mockery/v2@v2.5.1
	go install github.com/uudashr/gocognit/cmd/gocognit@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/mcubik/goverreport@latest
