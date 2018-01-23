# taco 🌮🌮🌮
The next generation repository system for DLSS
![taco](https://user-images.githubusercontent.com/92044/34897877-016a4e36-f7b6-11e7-80e3-4edecfb2f89d.gif)

## Swagger API

This configuration is for AWS API Gateway.  It was retrieved by going to the API, selecting the "prod" under "Stages" and then doing "Export" and selecting "Export as Swagger + API Gateway Extensions"

## Building

### Build for the local OS

```shell
% go get -t
% go build
```

### Building for Docker
```shell
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
docker build -t taco
docker run -p 8080:8080 taco
```

## Running

```shell
% taco
```

Now visit: http://localhost:8080/v1/resource/99

## Testing

```shell
% go test -v ./...
```
