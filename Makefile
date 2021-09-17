
GIT_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
BUILD_DATE := $(shell date "+%Y-%m-%d:%H:%M")
GIT_TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
ifeq ($(GIT_TAG),)
VERSION := ${GIT_COMMIT}-${BUILD_DATE}
else
VERSION := $(GIT_TAG:v%=%)-${GIT_COMMIT}-${BUILD_DATE}
endif

build:
	go mod tidy
	go build -ldflags "-X main.Version=${VERSION}"

install: build
	go install -ldflags "-X main.Version=${VERSION}" .

test:
	go test -v ./...

create-schemas: build
	./magda-cli schema create --id cse-order --name cse-order --schema-file example/schema/order.json
	./magda-cli schema create --id cse-service --name cse-service --schemaFile example/schema/service.json

load-services:
	./magda-cli record update --id ffdi --name ffdi -a cse-service -f example/record/ffdi_service.json