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
- /utils: Utility folder which encompasses functions that really don't belong anywhere else

## Interacting with the API

Basecoat uses gRPC and serves requests through a gRPC-proxy. This means that it can receive http _and_ grpc requests. You can find the API endpoints in the [proto files](./api). The admin routes require a pre-negotiated password.

You can send requests using a utility like grpcurl:

`grpcurl -H "Authorization: Bearer lolwut" -d {} basecoat.clintjedwards.com:443 api.Basecoat/ListAccounts`

## How to run locally

This takes a few steps

### Download golang packages

`go mod tidy` or `go mod download`

### Generating go and javascript code using protobuf definition

This project uses grpc and grpc-web which use generated functions from a protobuf definition

You'll need:

- protoc: `https://github.com/golang/protobuf`
- go's proto plugin: `go install github.com/golang/protobuf/protoc-gen-go@latest`
- grpc-web proto plugin: `https://github.com/grpc/grpc-web/releases`

### Building the frontend

This requires installing vue js and a bunch of other packages that the frontend uses. Basecoat uses go generate and vfsgen to bake frontend files into the binary.

- install npm packages: `cd frontend; npm install`

### Start application

`make run`
`curl https://localhost:8080`

## How to run tests

The only tests are integration tests that must be run against a live server. The testing command will bring up the testing server for you.

`go test -v ./tests`

## Authors

- **Clint Edwards** - [Github](https://github.com/clintjedwards)

This software is provided as-is. It's a hobby project, done in my free time, and I don't get paid for doing it.
