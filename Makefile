SHELL = /bin/bash
VERSION=$(shell date +%s)

GO_LDFLAGS := '-X "github.com/clintjedwards/basecoat/cmd.appVersion=$(VERSION)"'

build-dev:
	npx webpack --config="./frontend/webpack.config.js" --mode="development"
	packr build -ldflags $(GO_LDFLAGS) -o $(path)

build-backend: check-path-included
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	go test ./utils
	go build -ldflags $(GO_LDFLAGS) -o $(path)

build: check-path-included
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	go test ./utils
	npm run --prefix ./frontend build:production
	packr build -ldflags $(GO_LDFLAGS) -o $(path)

run:
	protoc --go_out=plugins=grpc:. api/*.proto
	protoc --js_out=import_style=commonjs,binary:./frontend/src/ --grpc-web_out=import_style=typescript,mode=grpcwebtext:./frontend/src/ -I ./api/ api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:development
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/basecoat && /tmp/basecoat server

install:
	protoc --go_out=plugins=grpc:. api/*.proto
	go mod tidy
	npm run --prefix ./frontend build:production
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/basecoat
	sudo mv /tmp/basecoat /usr/local/bin/
	chmod +x /usr/local/bin/basecoat

check-path-included:
ifndef path
	$(error path is undefined; ex. path=/tmp/basecoat)
endif
