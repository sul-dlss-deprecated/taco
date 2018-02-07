# Benchmarking

## Install Apache Bench

OsX has Apache Bench installed by default at `/usr/sbin/ab`

## Start the Server

```
$ cd cmd/tacod
$ AWS_ACCESS_KEY_ID=999999 AWS_SECRET_KEY=1231 go run main.go
```

## Benchmark 1000 create resource requests:
```
ab -T application/json -p examples/request.json -n 1000 http://localhost:8080/v1/resource
```

## Benchmark 1000 GET requests:

```
ab -n 1000 http://localhost:8080/v1/resource/44af74ec-d8e2-4a7c-aa5c-874b86da350d
```


## Comparing to other things:

### Fedora 4

Note: this goes faster than TACO as it does no shape validation and writes to external services.  It would require several extra requests to accommodate object permissions.

```
gem install fcrepo_wrapper
fcrepo_wrapper
```

Then start the benchmark in a new shell
```
ab -p examples/request.json -n 1000 http://localhost:8984/rest
```

### DOR Services / Fedora 3
Note: this goes much slower than taco.  It writes to several external services.

Use the repository at https://github.com/sul-dlss-labs/dor-services-benchmark
