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

## build: run tests and compile full app in production mode
build:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:production
	packr build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

## build-backend: build backend without frontend assets
build-backend:
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	go build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

## run: build application and run server; useful for dev
build-dev:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:development
	packr build -race -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)

## build-protos: build required protobuf files
build-protos:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto

## deploy: deploy the application to production
deploy: backup
	wget -O /tmp/${APP_NAME} https://github.com/clintjedwards/${APP_NAME}/releases/download/v${SEMVER}/${APP_NAME}
	scp /tmp/${APP_NAME} ${SERVER_USERNAME}@${APP_NAME}.clintjedwards.com:/tmp/${APP_NAME}
	ssh -t ${SERVER_USERNAME}@${APP_NAME}.clintjedwards.com ' \
	sudo mv /tmp/${APP_NAME} /usr/local/bin/; \
	sudo chmod +x /usr/local/bin/${APP_NAME}; \
	sudo chown ${SERVER_USERNAME}:${SERVER_USERNAME} /usr/local/bin/${APP_NAME}; \
	sudo service ${APP_NAME} stop && sudo setcap CAP_NET_BIND_SERVICE=+eip /usr/local/bin/${APP_NAME} && sudo service ${APP_NAME} start; \
	'

## help: prints this help message
help:
	@echo "Usage: "
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## install: build application and install on system
install: build
	sudo mv BUILD_PATH /usr/local/bin/
	chmod +x /usr/local/bin/${APP_NAME}

## run: build application and run server; useful for dev
run:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:development
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/${APP_NAME} && /tmp/${APP_NAME} server

## run-backend: build backend only and run server; useful for dev
run-backend:
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/${APP_NAME} && /tmp/${APP_NAME} server
