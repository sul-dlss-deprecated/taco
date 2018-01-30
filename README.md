# taco 🌮🌮🌮 [![CircleCI](https://circleci.com/gh/sul-dlss-labs/taco.svg?style=svg)](https://circleci.com/gh/sul-dlss-labs/taco)
The next generation repository system for DLSS
![taco](https://user-images.githubusercontent.com/92044/34897877-016a4e36-f7b6-11e7-80e3-4edecfb2f89d.gif)

## Swagger API

This configuration is for AWS API Gateway.  It was retrieved by going to the API, selecting the "prod" under "Stages" and then doing "Export" and selecting "Export as Swagger + API Gateway Extensions"

## Go Local Development Setup

1. Install go (grab binary from here or use `brew install go` on Mac OSX).
2. Setup your Go workspace (where your Go code, binaries, etc. are kept together. TO DO: add info on Go workspace FYI.):
      ```bash
      $ mkdir -p ~/go
      $ export GOPATH=~/go
      $ export PATH=~/go/bin:$PATH
      $ cd ~/go
      ```
      Your Go code repositories will reside within `~/go/src/...` in the `$GOPATH`. Name these paths to avoid library clash, for example TACO Go code could be in `~/go/src/github.com/sul-dlss-labs/taco`. This should be where your Github repository resides too.
3. In order to download the project code to `~/go/src/github.com/sul-dlss-labs/taco`, from any directory, run:
```bash
$ go get github.com/sul-dlss-labs/taco
```
4. Handle dependencies with the Go Dep package:
    * Install Go Dep via `brew install dep` then `brew upgrade dep`.
    * If your project's `Gopkg.toml` has not yet been populated (i.e. there should be libraries not commented out), you need to add an inferred list of your dependencies by running `dep init`.
    * If your project has that, make sure your dependencies are synced via running `dep ensure`.
    * If you need to add a new dependency, run `dep ensure -add github.com/pkg/errors`. This should add the dependency and put the new dependency in your `Gopkg.*` files.

## Running the Go Code locally without a build


```shell
$ go run main.go
```

## Building to TACO Binary

### Building for Docker
```shell
$ docker build -t taco  .
$ docker run -p 8080:8080 taco
```

### Build for the local OS
```shell
% cd cmd/tacod
% dep ensure -vendor-only
% go build -o tacod main.go
```

## Testing

```shell
$ go test -v ./...
```

## Running the TACO Binary

First start up DynamoDB:
```shell
$ SERVICES=dynamodb localstack start
```

Then create the table:
```shell
$ awslocal dynamodb create-table --table-name resources \
  --attribute-definitions "AttributeName=id,AttributeType=S" \
  --key-schema "AttributeName=id,KeyType=HASH" \
  --provisioned-throughput=ReadCapacityUnits=100,WriteCapacityUnits=100
```

And add a stub record:
```
$ awslocal dynamodb put-item --table-name resources --item '{"id": {"S":"99"}, "title":{"S":"Ta-da!"}}'
```

```shell
% TACO_ENV=production AWS_ACCESS_KEY_ID=999999 AWS_SECRET_KEY=1231 ./tacod
```

Now visit: http://localhost:8080/v1/resource/99

## API Code Structure

We use `go-swagger` to generate the API code within `generated/`, and we connect that API code to our own handlers defined with `handlers/`. The handlers are where we add our own logic for processing requests. Our handlers and the generated API code is connected within `main.go`, which is the file to start the API.

### To generate the API code

The API code is generated from `swagger.yml` using `go-swagger` library. TBD: best way to handle regeneration (i.e. currently you're recommended to delete the generated code before re-running):

```shell
$ swagger generate server -t generated --exclude-main
```

#### ACHTUNG!
We've seen problems with the swagger at the HEAD of github. So we've used this workaround:

```shell
cd $GOPATH/src/github.com/go-swagger/go-swagger/cmd/swagger
git checkout 5ade92aa47f4b45e197e97b05f36e761fab9bbf0
go install
```

Do this prior to generating code.

### To run the API code

```shell
$ go run main.go
```

### Non-generated code

The API code generation does **not** touch the following, which we are writing locally:
- `main.go`
- `handlers/`

(basically anything outside of `generated`).


## SWAGGER Generated Documentation

To see the SWAGGER generated documentation, run the following:

```shell
$ swagger serve swagger.yml
```

This should prompt you to your web browser for the HTML generated docs. TBD: how we can have this consistently running on our servers de facto at a URL for the documentation.
