# TACO Developer Guide

## Adding Environment Variables for Configuration

```shell
export YOUR_VAR='whatever'
```

### Add config and default value

1. Add to the config struct definition
```go
type Config struct {
  DynamodbEndpoint   string
	DynamodbDisableSSL bool
	KinesisEndpoint    string
	KinesisDisableSSL  bool
	ResourceTableName  string
	DepositStreamName  string
	S3Endpoint         string
	S3BucketName       string
	S3DisableSSL       bool
	Port               int
  YourVar            string
}
```

2. Load var into config struct
```go
func NewConfig() *Config {
  return &Config{
    DynamodbEndpoint:   getString("DYNAMO_DB_ENDPOINT", "localhost:4569"),
		DynamodbDisableSSL: getBool("DYNAMODB_DISABLE_SSL", true),
		ResourceTableName:  getString("RESOURCE_TABLE_NAME", "resources"),
		KinesisEndpoint:    getString("KINESIS_ENDPOINT", "localhost:4568"),
		KinesisDisableSSL:  getBool("KINESIS_DISABLE_SSL", true),
		DepositStreamName:  getString("DEPOSIT_STREAM_NAME", "deposit"),
		S3Endpoint:         getString("S3_ENDPOINT", "localhost:4572"),
		S3BucketName:       getString("S3_BUCKET_NAME", "taco-deposited-files"),
		S3DisableSSL:       getBool("S3_DISABLE_SSL", true),
		Port:               getInteger("TACO_PORT", 8080),
    YourVar:            getString("YOUR_VAR", "Some Default Value"),
	}
}
```
