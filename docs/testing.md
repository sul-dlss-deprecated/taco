# Testing Guidelines & Framework

## Testing Framework

[Go](https://golang.org) has extensive support for automated [testing](https://golang.org/pkg/testing/) built into the language.


Additionally, we utilize the [Gofight](https://github.com/appleboy/gofight) framework to test API handler mock responses.

## Go Unit Tests
The unit tests have no external dependencies and can be run like so:
```shell
$ go test -v ./... -short
```

## Go Package Tests
The package tests are dependent on [Localstack](docs/localstack.md) running as they are testing individual connectivity to
external services.

```shell
AWS_REGION=localstack go test -v ./[PACKAGE NAME]/
```

## Go Integration Tests

The integration test depends on the taco binary and [Localstack](docs/localstack.md) running.  Once these conditions are met you can run the integration tests using:

```shell
$ go test test/integration_test.go
```

## Troubleshooting Common Errors

### Missing Region

```
&awserr.baseError{code:"MissingRegion", message:"could not find region configuration", errs:[]error(nil)}
```

This error indicates that the command line aws region key (`AWS_REGION`) are missing.

### Missing AWS Resource

```
&awserr.requestError{awsError:(*awserr.baseError)(0xc420060a00), statusCode:400, requestID:"8e442552-4370-46d1-8226-73a0db496dcd"}
```

With a corrisponding
```
WARNING:localstack.services.dynamodb.dynamodb_listener: Unknown table: resources not found in {}
```

Is an indication that the particular [Localstack](docs/localstack.md) resource is missing. Follow the instructions in the document
to ensure the resource is available.
