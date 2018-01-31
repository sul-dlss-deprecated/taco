package config

import (
	"log"
	"os"
)

type Config struct {
	AWS_Region          string
	Dynamo_Db_Endpoint  string
	Dynamo_Disable_SSL  bool
	Resource_Table_Name string
}

func NewConfig() *Config {
	return &Config{
		AWS_Region:          aws_region(),
		Dynamo_Db_Endpoint:  dynamo_db_endpoint(),
		Dynamo_Disable_SSL:  dynamo_disable_ssl(),
		Resource_Table_Name: resource_table_name(),
	}
}

func aws_region() string {
	var region string
	region = os.Getenv("AWS_REGION")
	if region == "" {
		region = "localstack"
		log.Printf("AWS_REGION: Using default [localstack].")
	}
	log.Printf("AWS_REGION: Found setting [%s]", region)
	return region
}

func dynamo_db_endpoint() string {
	var db string
	db = os.Getenv("DYNAMO_DB_ENDPOINT")
	if db == "" {
		db = "localhost:4569"
		log.Printf("DYNAMO_DB_ENDPOINT: Using default [localhost:4569].")
	}
	log.Printf("DYNAMO_DB_ENDPOINT: Found setting [%s]", db)
	return db
}

func dynamo_disable_ssl() bool {
	var disablessl string
	disablessl = os.Getenv("DYNAMODB_DISABLE_SSL")
	if disablessl == "FALSE" || disablessl == "false" {
		log.Printf("DYNAMODB_DISABLE_SSL: Found setting [false].")
		return false
	} else {
		log.Printf("DYNAMODB_DISABLE_SSL: Using default [true].")
		return true
	}
}

func resource_table_name() string {
	var tablename string
	tablename = os.Getenv("RESOURCE_TABLE_NAME")
	if tablename == "" {
		tablename = "resources"
		log.Printf("RESOURCE_TABLE_NAME: Using default [resources].")
	}
	return tablename
}
