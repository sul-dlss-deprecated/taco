package config

import (
	"log"
	"os"
	"strconv"
)

// Config is configuration for the TACO application
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
}

// NewConfig creates a new configuration with values from environment variables
// or defaults
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
	log.Printf("%s: Found setting [true].", envVar)
	return true
}

func getInteger(envVar string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(envVar))
	if err != nil || value == 0 {
		log.Printf("%s: Using default [%v].", envVar, defaultValue)
		return defaultValue
	}
	log.Printf("%s: Found setting [%v].", envVar, value)
	return value
}
