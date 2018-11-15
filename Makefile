SHELL = /bin/bash
GIT_COMMIT := $(shell git rev-parse --short HEAD)
VERSION_FILE=VERSION.md
VERSION=$(shell cat ${VERSION_FILE})

GO_LDFLAGS := '-X "github.com/clintjedwards/basecoat/cmd.appVersion=$(VERSION) $(GIT_COMMIT)"'

build-dev:
	npx webpack --config="./frontend/webpack.config.js" --mode="development"
	packr build -ldflags $(GO_LDFLAGS) -o $(path)

build: check-path-included
	npx webpack --config="./frontend/webpack.config.js" --mode="production"
	packr build -ldflags $(GO_LDFLAGS) -o $(path)

run:
	npx webpack --config="./frontend/webpack.config.js" --mode="development"
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/basecoat && /tmp/basecoat server

install:
	npx webpack --config="./frontend/webpack.config.js" --mode="production"
	packr build -ldflags $(GO_LDFLAGS) -o /tmp/basecoat
	sudo mv /tmp/basecoat /usr/local/bin/
	chmod +x /usr/local/bin/basecoat

check-path-included:
ifndef path
	$(error path is undefined; ex. path=/tmp/basecoat)
endif
