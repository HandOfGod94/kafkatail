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

clean-integration:
	docker compose stop
	docker compose rm -f -v
	docker volume rm -f kafkatail_kafka_data kafkatail_zookeeper_data
	docker compose up -d

.PHONY: integration-tests
integration-tests:
	go install
	go test -p 1 -count=1 -timeout 90s -tags=integration ./integration-tests

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

setup-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/vektra/mockery/v2@v2.5.1
	go install github.com/uudashr/gocognit/cmd/gocognit@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/mcubik/goverreport@latest

include scripts/cucumber.mk
