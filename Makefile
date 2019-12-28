APP_NAME = basecoat
BUILD_PATH = /tmp/${APP_NAME}
EPOCH_TIME = $(shell date +%s)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GO_LDFLAGS = '-X "github.com/clintjedwards/${APP_NAME}/cmd.appVersion=$(VERSION)" \
			   -X "github.com/clintjedwards/${APP_NAME}/service.appVersion=$(VERSION)"'
SEMVER = 0.0.1
SHELL = /bin/bash
VERSION = ${SEMVER}_${EPOCH_TIME}_${GIT_COMMIT}


## backup: backup production database using gcp
backup:
	gcloud datastore export gs://clintjedwardsbackups/basecoat-${EPOCH_TIME}

## build-prod: run tests and compile full app in production mode
build-prod:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:production
	go generate
	go build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

## build-backend: build backend without frontend assets
build-backend:
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	go generate
	go build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

## build-dev: build application in dev mode
build-dev:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:development
	go generate
	go build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

## build-protos: build required protobuf files
build-protos:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto

## help: prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## install: build application and install on system
install: build-prod
	sudo mv BUILD_PATH /usr/local/bin/
	chmod +x /usr/local/bin/${APP_NAME}

## run: build application and run server; useful for dev
run: export DEBUG=true
run:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:development
	go generate
	go build -ldflags $(GO_LDFLAGS) -o /tmp/${APP_NAME} && /tmp/${APP_NAME} server

## run-backend: build backend only and run server; useful for dev
run-backend:
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	go build -ldflags $(GO_LDFLAGS) -o /tmp/${APP_NAME} && /tmp/${APP_NAME} server
