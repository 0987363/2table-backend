PACKAGE = .

#export GO15VENDOREXPERIMENT=1
#export GO111MODULE=on

BUILD_VERSION=$(shell git describe --tags --abbrev=0)
BUILD_NUMBER=$(strip $(if $(TRAVIS_BUILD_NUMBER), $(TRAVIS_BUILD_NUMBER), 0))
BUILD_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
BUILD_OWNER=$(shell whoami).$(shell uname -n)
ifeq ($(shell whoami ), "runner")
else
	BUILD_OWNER=$(shell git log -1 --pretty=format:'%an')
endif

IMAGE_DATE=$(shell date -u +%Y-%m-%d.%H%M%S)
IMAGE_SERVER_NAME=$(shell basename $(PWD))
IMAGE_TAG := $(IMAGE_DATE)-$(BUILD_COMMIT)-$(BUILD_VERSION)

image-test: doc release
	docker build --platform linux/arm64  -t registry.druid.company/ecotopia/service/$(IMAGE_SERVER_NAME):$(IMAGE_TAG) .
	docker push registry.druid.company/ecotopia/service/$(IMAGE_SERVER_NAME):$(IMAGE_TAG)
	docker rmi registry.druid.company/ecotopia/service/$(IMAGE_SERVER_NAME):$(IMAGE_TAG)

all: build

clean:
	rm -f backend

release:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v -ldflags "-X main.BuildVersion=$(BUILD_VERSION) -X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildDate=$(BUILD_DATE) -X main.BuildOwner=$(BUILD_OWNER)"

build:
	go build -v -ldflags "-X main.BuildVersion=$(BUILD_VERSION) -X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildDate=$(BUILD_DATE) -X main.BuildOwner=$(BUILD_OWNER)"

run: build
	./backend serve

