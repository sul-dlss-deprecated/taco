# TACO Developer Guide

## Adding Environment Variables for Configuration

```shell
export YOUR_VAR='whatever'
```

### Add config and default value

1. Add to the config struct definition
```go
type Config struct {
	DynamodbEndpoint   string `envVariable:"DYNAMO_DB_ENDPOINT" defaultValue:"localhost:4569"`
	DynamodbDisableSSL bool   `envVariable:"DYNAMO_DISABLE_SSL" defaultValue:"true"`
	KinesisEndpoint    string `envVariable:"KINESIS_ENDPOINT" defaultValue:"localhost:4568"`
	KinesisDisableSSL  bool   `envVariable:"KINESIS_DISABLE_SSL" defaultValue:"true"`
	ResourceTableName  string `envVariable:"RESOURCE_TABLE_NAME" defaultValue:"resources"`
	DepositStreamName  string `envVariable:"DEPOSIT_STREAM_NAME" defaultValue:"deposit"`
  YourVar            string `envVariable:"YOUR_ENV_VAR" defaultValue:"yourDefaultValue"`
}
```
