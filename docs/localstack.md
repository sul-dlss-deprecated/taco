# Localstack

We use [Localstack](https://github.com/localstack/localstack) as a stub implementation of many AWS services that we can run locally for development and testing of Taco.

## Installing
`localstack`'s dependencies are found here:
https://github.com/localstack/localstack#requirements

Follow the directions on installing `localstack` here: https://github.com/localstack/localstack#installing.


## Running

Start localstack by running the command:
```shell
$ SERVICES=dynamodb,s3 localstack start
```

Next we have to set up the resources:

```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 make resources
```
