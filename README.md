# Basecoat: Formula Tracking and Search Tool

Basecoat is a CRUD formula indexing tool meant to record formulas for certain colors and store them for future reference.

![Basecoat](https://i.imgur.com/ScgDBiZ.png) ![Basecoat Formula](https://i.imgur.com/nilixZL.png)


This project is slightly conveluted to get up and running
locally. I wish it was a simple go run but a few vue js
components, protobufs, grpc docs later and here we are
you'll need protoc: https://github.com/golang/protobuf
and go's plugin: go get github.com/golang/protobuf/protoc-gen-go
and grpc-web plugin: https://github.com/grpc/grpc-web/releases
get packr: go get -u github.com/gobuffalo/packr/packr
run go mod tidy
download npm and install webpack, webpack-cli
run npm install in the frontend folder
Generate localhost certs for TLS
go get -u github.com/FiloSottile/mkcert
$(go env GOBIN)/mkcert
run make run from root of project
You should be able to get to the frontend by localhost:8080
google datastore emulator(the current main database) can be downloaded and started by navagating to the docker folder and using docker-compose up

## Authors

- **Clint Edwards** - [Github](https://github.com/clintjedwards)
