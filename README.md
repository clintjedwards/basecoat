# Basecoat: Formula Tracking and Search Tool

Basecoat is a CRUD formula indexing tool meant to record formulas for certain colors and store them for future reference.

![Basecoat Logo](https://i.imgur.com/ScgDBiZ.png)
![Basecoat](https://i.imgur.com/PUleJk9.png)
![Basecoat Formula](https://i.imgur.com/dB3MJUW.png)

## Application Structure

- /api: GRPC protobuf definition and resulting auto-generated protobuf code
- /app: Connector package for rest API code and frontend code. Used to initialize the server as a whole.
- /cmd: Command line tooling. Essentially the root of the application. Allows starting of server and also can be used as a client.
- /config: Parses necessary configuration variables from the environment
- /docker: Docker assets used to run test database
- /frontend: Frontend vue js/grpc-web assets and code.
- /scripts: Scripts not compiled into the main binary, used for development
- /service: Logic for backend GRPC api service
- /storage: Storage interface for backend database functions
- /tests: Testing code which runs against a live local server to verify api functionality is correct

## How to run locally

This takes a few steps

### Download golang packages

`go mod tidy` or `go mod download`

### Generating go and javascript code using protobuf definition

This project uses grpc and grpc-web which use generated functions from a protobuf definition

You'll need:

- protoc: `https://github.com/golang/protobuf`
- go's proto plugin: `go get github.com/golang/protobuf/protoc-gen-go`
- grpc-web proto plugin: `https://github.com/grpc/grpc-web/releases`

### Building the frontend

This requires installing vue js and a bunch of other packages that the frontend uses. Basecoat uses go generate and vfsgen to bake frontend files into the binary.

- install npm packages: `cd frontend; npm install`

### Generate certificates

This tool will autogenerate certs for your and add them as trusted

`go get -u github.com/FiloSottile/mkcert`
`$(go env GOBIN)/mkcert`

### Start application

`make run`
`curl https://localhost:8080`

## How to run tests

The only tests are integration tests that must be run against a live server. Before running the next command you must bring up a local version of the application

`go test tests/integration_test.go -v -count 1`

## Authors

- **Clint Edwards** - [Github](https://github.com/clintjedwards)
