FROM golang:alpine AS base
WORKDIR /go/src/github.com/sul-dlss-labs/taco/
COPY . .
RUN apk add --no-cache --virtual .build-deps git && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    apk del .build-deps
RUN CGO_ENABLED=0 go build .

FROM scratch
COPY --from=base /go/src/github.com/sul-dlss-labs/taco/taco .
CMD ["/taco"]
