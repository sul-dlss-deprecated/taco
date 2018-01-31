package config

import (
	"log"
	"os"
)

type Config struct {
	AWS_Region  string
	Dynamo_Db   string
	Disable_SSL bool
	Table_Name  string
}

func NewConfig() *Config {
	return &Config{
		AWS_Region:  aws_region(),
		Dynamo_Db:   dynamo_db(),
		Disable_SSL: disable_ssl(),
		Table_Name:  table_name(),
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

func dynamo_db() string {
	var db string
	db = os.Getenv("DYNAMO_DB")
	if db == "" {
		db = "localhost:4569"
		log.Printf("DB_ENDPOINT: Using default [localhost:4569].")
	}
	log.Printf("DB_ENDPOINT: Found setting [%s]", db)
	return db
}

func disable_ssl() bool {
	var disablessl string
	disablessl = os.Getenv("DISABLE_SSL")
	if disablessl == "FALSE" || disablessl == "false" {
		log.Printf("DISABLE_SSL: Found setting [false].")
		return false
	} else {
		log.Printf("DISABLE_SSL: Using default [true].")
		return true
	}
}

func table_name() string {
	var tablename string
	tablename = os.Getenv("TABLE_NAME")
	if tablename == "" {
		tablename = "resources"
		log.Printf("TABLE_NAME: Using default [resources].")
	}
	return tablename
}
