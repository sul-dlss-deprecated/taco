package config

import (
	"os"
	"log"
)

type Config struct {
    AWS_Region	string
		DB_Endpoint string
		Disable_SSL bool
}

func NewConfig() *Config {
	return &Config{
		AWS_Region: aws_region(),
		DB_Endpoint: db_endpoint(),
		Disable_SSL: disable_ssl(), // os.Getenv("DISABLE_SSL"),
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

func db_endpoint() string {
	var endpoint string
	endpoint = os.Getenv("DB_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4569"
		log.Printf("DB_ENDPOINT: Using default [localhost:4569].")
	}
	log.Printf("DB_ENDPOINT: Found setting [%s]", endpoint)
	return endpoint
}

func disable_ssl() bool {
	var disablessl string
  disablessl = os.Getenv("DISABLE_SSL")
	if (disablessl == "FALSE" || disablessl == "false") {
		log.Printf("DISABLE_SSL: Found setting [false].")
		return false
	} else {
		log.Printf("DISABLE_SSL: Using default [true].")
		return true
	}
}
