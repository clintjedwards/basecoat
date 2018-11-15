# Basecoat: Formula Tracking and Search Tool

Basecoat is a CRUD formula indexing tool meant to record formulas for certain colors and store them for future reference.

![Basecoat](https://i.imgur.com/ScgDBiZ.png) ![Basecoat Formula](https://i.imgur.com/nilixZL.png)

[API Documentation](https://basecoat.clintjedwards.com)

## Installation

### Build from source

```bash
git clone git@github.com:clintjedwards/basecoat.git
make build path=somepathhere/basecoat
```

### Build from go

This method does **not** bake the frontend files into the binary

```bash
go get github.com/clintjedwards/basecoat
```

## Getting Started

### Docker Test Environment

For development all you need is a postgres server to test against.

The current docker compose file will attempt to launch a postgres database for you located at localhost:5432

```bash
cd docker
cp env.sample .env
docker-compose up
```

#### Testing using the included Frontend

You can turn the frontend on by the environment variable `frontend=true`.

### Application Configuration

You can set configuration options via environment variables. An up to date list of the environment variables can be found in the config folder. It is usually not necessary to change these for development as sane configurations based on the docker environment are already set.

### Application Structure

- /api: Rest API code.
- /app: Connector package for rest API code and frontend code. Used to initialize the server as a whole.
- /cmd: Command line tooling. Essentially the root of the application. Allows starting of server and also can be used as a client.
- /config: Parses necessary configuration variables from the environment
- /docker: Docker assets
- /frontend: Frontend assets and code.
- /tests: Testing code which runs against a live local server to verify api functionality is correct

### How to run

Basecoat users [packr](https://github.com/gobuffalo/packr) to build static assets into the Golang binary.
You can use the `packr` command in lieu of the standard go command. So to build a quick version you can use `packr build -o /tmp/basecoat && /tmp/basecoat server` to quickly build and run the binary.

From the project's root directory

### How to generate API Documentation

The API documentation is generated via the API Blueprint standard. The file is simple markdown file that is
formatted in a way that other tools can understand. To generate the html file for in app documentation
you'll need to do the following:

You'll need to download the go program [snowboard]("https://github.com/bukalapak/snowboard"). A quick and painless way to do it is through the docker container:

```
docker pull quay.io/bukalapak/snowboard
#change to directory of API.apib file
docker run -it --rm -v $PWD:/doc:z quay.io/bukalapak/snowboard html -o api.html API.apib
```

This will generate the html file from the API.apib markdown file.

### Running tests

1.  Spin up the docker compose environment: `docker-compose up`
2.  Run tests through go's built in test command: `go test -v ./tests`

The docker environment must be purged after each test run as it creates assets that aren't removed(This is due to database safety reasons)

## Authors

- **Clint Edwards** - [Github](https://github.com/clintjedwards)
