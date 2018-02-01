package config

import (
	"log"
	"os"
)

// Config is configuration for the TACO application
type Config struct {
	AWSRegion          string
	DynamodbEndpoint   string
	DynamodbDisableSSL bool
	ResourceTableName  string
	DepositStreamName  string
}

// NewConfig creates a new configuration with values from environment variables
// or defaults
func NewConfig() *Config {
	return &Config{
		AWSRegion:          getString("AWS_REGION", "localstack"),
		DynamodbEndpoint:   getString("DYNAMO_DB_ENDPOINT", "localhost:4569"),
		DynamodbDisableSSL: getBool("DYNAMODB_DISABLE_SSL", true),
		ResourceTableName:  getString("RESOURCE_TABLE_NAME", "resources"),
		DepositStreamName:  getString("DEPOSIT_STREAM_NAME", "deposit"),
	}
}

func getString(envVar string, defaultValue string) string {
	var value string
	value = os.Getenv(envVar)
	if value == "" {
		value = defaultValue
		log.Printf("%s: Using default [%s].", envVar, defaultValue)
		return defaultValue
	}
	log.Printf("%s: Found setting [%s].", envVar, value)
	return value
}

func getBool(envVar string, defaultValue bool) bool {
	var value string
	value = os.Getenv(envVar)
	if value == "FALSE" || value == "false" {
		log.Printf("%s: Using default [%s].", envVar, value)
		return false
	}
	log.Printf("%s: defaulting to true.", envVar)
	return true
}
