
.PHONY: build build-alpine clean test help default

BIN_NAME=ip2location-pfsense
VERSION :=$(shell grep "const Version" backend/version/version.go | sed -E 's/.*"([^"]+)".*/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME :="proget.terra-net.io:443/images/ip2location-pfsense"

default: test

help:
	@echo 'Management commands for IP2Location-pfSense:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make get-deps        runs dep ensure, mostly used for ci.'
	@echo '    make build-alpine    Compile optimized for alpine linux.'
	@echo '    make package         Build final docker image with just the go binary inside'
	@echo '    make tag             Tag image created by package with latest, git commit and version'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make push            Push tagged images to registry'
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -C ./backend -ldflags "-X github.com/jpmchia/ip2location-pfsense/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/jpmchia/ip2location-pfsense/version.BuildDate=${BUILD_DATE}" -o bin/${BIN_NAME}

get-deps:
	@echo "getting dependencies"
	@echo "GOPATH=${GOPATH}"
	cd backend && go get . && go mod tidy && go mod vendor


build-alpine:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -C ./backend -ldflags '-w -linkmode external -extldflags "-static" -X github.com/jpmchia/ip2location-pfsense/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/jpmchia/ip2location-pfsense/version.BuildDate=${BUILD_DATE}' -o bin/${BIN_NAME}

package:
	@echo "building image ${BIN_NAME} ${VERSION} $(GIT_COMMIT)"
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(IMAGE_NAME):local .

tag: 
	@echo "Tagging: latest ${VERSION} $(GIT_COMMIT)"
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):$(GIT_COMMIT)
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):${VERSION}
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):latest

push: tag
	@echo "Pushing docker image to registry: latest ${VERSION} $(GIT_COMMIT)"
	docker push $(IMAGE_NAME):$(GIT_COMMIT)
	docker push $(IMAGE_NAME):${VERSION}
	docker push $(IMAGE_NAME):latest

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test -C ./backend ./...

