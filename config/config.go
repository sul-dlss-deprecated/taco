package config

import (
	"log"
	"os"
	"reflect"
	"strconv"
)

// Config is configuration for the TACO application
type Config struct {
	DynamodbEndpoint  string `envVariable:"DYNAMO_DB_ENDPOINT" defaultValue:"localhost:4569"`
	AwsDisableSSL     bool   `envVariable:"AWS_DISABLE_SSL" defaultValue:"true"`
	KinesisEndpoint   string `envVariable:"KINESIS_ENDPOINT" defaultValue:"localhost:4568"`
	ResourceTableName string `envVariable:"RESOURCE_TABLE_NAME" defaultValue:"resources"`
	DepositStreamName string `envVariable:"DEPOSIT_STREAM_NAME" defaultValue:"deposit"`
	S3Endpoint        string `envVariable:"S3_ENDPOINT" defaultValue:"localhost:4572"`
	S3BucketName      string `envVariable:"S3_BUCKET_NAME" defaultValue:"taco-deposited-files"`
	Port              int    `envVariable:"TACO_PORT" defaultValue:"8080"`
}

// NewConfig creates a new configuration with values from environment variables
// or defaults
func NewConfig() *Config {
	configuration := Config{}
	typ := reflect.TypeOf(configuration)
	fieldsList := reflect.ValueOf(&configuration)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldSetter := fieldsList.Elem().FieldByName(field.Name)
		envVariable := field.Tag.Get("envVariable")
		defaultValue := field.Tag.Get("defaultValue")

		switch fieldSetter.Kind() {
		case reflect.Int:
			setInteger(fieldSetter, envVariable, defaultValue)
		case reflect.Bool:
			setBool(fieldSetter, envVariable, defaultValue)
		default:
			setString(fieldSetter, envVariable, defaultValue)
		}
	}
	return &configuration
}

func setString(field reflect.Value, envVariable string, defaultValue string) {
	value := os.Getenv(envVariable)
	if value == "" {
		value = defaultValue
		log.Printf("%s: Using default [%s].", envVariable, defaultValue)
	} else {
		log.Printf("%s: Found setting [%s].", envVariable, value)
	}
	field.SetString(value)
}

func setBool(field reflect.Value, envVariable string, defaultValue string) {
	value := os.Getenv(envVariable)
	if value == "" {
		value = defaultValue
		log.Printf("%s: Using default [%s].", envVariable, value)
	} else {
		log.Printf("%s: Found setting [%s].", envVariable, value)
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("%s: Error setting [%s] as a boolean, using [%s] instead", envVariable, value, defaultValue)
		boolVal, _ = strconv.ParseBool(defaultValue)
	}
	field.SetBool(boolVal)
}

func setInteger(field reflect.Value, envVariable string, defaultValue string) {
	value, err := strconv.ParseInt(os.Getenv(envVariable), 10, 64)
	if err != nil || value == 0 {
		defaultInteger, _ := strconv.ParseInt(defaultValue, 10, 64)
		log.Printf("%s: Using default [%v].", envVariable, defaultInteger)
		value = defaultInteger
	} else {
		log.Printf("%s: Found setting [%v].", envVariable, value)
	}
	field.SetInt(value)
}
