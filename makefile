.PHONY: all test integration-test ci-test clean build docker run pre-docker-build swagger sbom

BINARY_NAME=soarca
DIRECTORY = $(sort $(dir $(wildcard ./test/*/)))
VERSION = $(shell git describe --tags --dirty)
BUILDTIME := $(shell  date '+%Y-%m-%dT%T%z')

GOLDFLAGS += -X main.Version=$(VERSION)
GOLDFLAGS += -X main.Buildtime=$(BUILDTIME)
GOFLAGS = -ldflags "$(GOLDFLAGS)"

swagger:
	mkdir -p swaggerdocs
	swag init -o swaggerdocs

lint: swagger
	
	golangci-lint run  --timeout 5m -v

build: swagger
	CGO_ENABLED=0 go build -o ./build/soarca $(GOFLAGS) main.go

test:
	go test ./test/unittest/... -v

integration-test: swagger
	go test ./test/integration/... -v

ci-test: test integration-test

clean:
	rm -rf build/soarca* build/main
	rm -rf bin/*

compile: swagger
	echo "Compiling for every OS and Platform"
	
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-${VERSION}-linux-amd64 $(GOFLAGS) main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/${BINARY_NAME}-${VERSION}-darwin-arm64 $(GOFLAGS) main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/${BINARY_NAME}-${VERSION}-windows-amd64 $(GOFLAGS) main.go

sbom:
	echo "Generating SBOMs"
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 cyclonedx-gomod app -json -licenses -output bin/${BINARY_NAME}-${VERSION}-linux-amd64.bom.json
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 cyclonedx-gomod app -json -licenses -output bin/${BINARY_NAME}-${VERSION}-darwin-amd64.bom.json
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 cyclonedx-gomod app -json -licenses -output bin/${BINARY_NAME}-${VERSION}-windows-amd64.bom.json

pre-docker-build: swagger
	GOOS=linux GOARCH=amd64 go build -o bin/${BINARY_NAME}-${VERSION}-linux-amd64 $(GOFLAGS) main.go

docker: pre-docker-build
	docker build  -t soarca:${VERSION}  --build-arg="VERSION=${VERSION}" .

run: docker
	GIT_VERSION=${VERSION} docker compose up -d

all: build
