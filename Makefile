.PHONY: test

PROJECT_UNIT_TESTS     =$(shell go list ./... | grep -v test | grep -v db)
PROJECT_INT_TESTS      =$(shell go list ./db)
PROJECT_E2E_TESTS      =$(shell go list ./test)
LOCALSTACK_SERVICES    =dynamodb,s3
LOCAL_ENDPOINT_HOST    ?=localhost
LOCAL_ENDPOINT         =--endpoint-url=http://${LOCAL_ENDPOINT_HOST}
DYNAMO_ENDPOINT        =${LOCAL_ENDPOINT}:4569
S3_ENDPOINT            =${LOCAL_ENDPOINT}:4572
PROJ_TABLE_NAME        =resources
PROJ_BUCKET_NAME       =taco-deposited-files
PROJ_AWS_REGION        =us-east-1
PROJ_AWS_ACCESS_KEY_ID =999999
PROJ_AWS_SECRET_KEY    =999999
PROJ_ENV_VARS          =AWS_REGION=${PROJ_AWS_REGION} AWS_ACCESS_KEY_ID=${PROJ_AWS_ACCESS_KEY_ID} AWS_SECRET_ACCESS_KEY=${PROJ_AWS_SECRET_KEY}

default: dependencies test

test:
	$(PROJ_ENV_VARS) go test -v $(PROJECT_UNIT_TESTS)

compose_tests:
	$(PROJ_ENV_VARS) go test -v $(PROJECT_INT_TESTS)
	$(PROJ_ENV_VARS) go test -v $(PROJECT_E2E_TESTS)

dependencies:
	go get github.com/golang/dep/cmd/dep
	dep ensure

resources: table bucket

table:
	$(eval TABLE=$(shell aws $(DYNAMO_ENDPOINT) dynamodb list-tables | jq '.TableNames[0] // ""'))
	@if [[ $(TABLE) != "" ]]; \
	  then echo "$(TABLE) table found"; \
		else aws $(DYNAMO_ENDPOINT) dynamodb create-table --table-name $(PROJ_TABLE_NAME) \
				--attribute-definitions AttributeName=tacoIdentifier,AttributeType=S \
				AttributeName=externalIdentifier,AttributeType=S \
				AttributeName=version,AttributeType=N \
				AttributeName=dedupeIdentifier,AttributeType=S \
				--key-schema "AttributeName=tacoIdentifier,KeyType=HASH" \
				--provisioned-throughput=ReadCapacityUnits=10,WriteCapacityUnits=10 \
				--global-secondary-indexes "[{\"IndexName\":\"ResourceByExternalIDAndVersion\", \
				          \"KeySchema\":[{\"AttributeName\":\"externalIdentifier\",\"KeyType\":\"HASH\"}, \
				                     {\"AttributeName\":\"version\",\"KeyType\":\"RANGE\"}], \
				          \"Projection\":{\"ProjectionType\":\"ALL\"}, \
				          \"ProvisionedThroughput\":{\"ReadCapacityUnits\":10,\"WriteCapacityUnits\":10}}, \
				          {\"IndexName\":\"ResourceByDedupeIdentifier\", \
				          \"KeySchema\":[{\"AttributeName\":\"dedupeIdentifier\",\"KeyType\":\"HASH\"}], \
				          \"Projection\":{\"ProjectionType\":\"KEYS_ONLY\"}, \
				          \"ProvisionedThroughput\":{\"ReadCapacityUnits\":10,\"WriteCapacityUnits\":10}}]" ; \
	fi;

bucket:
	$(eval BUCKET=$(shell aws $(S3_ENDPOINT) s3api list-buckets | jq '.Buckets[0].Name // ""'))
	@if [[ $(BUCKET) != "" ]]; \
	  then echo "$(BUCKET) bucket found"; \
	  else aws $(S3_ENDPOINT) s3api create-bucket --bucket $(PROJ_BUCKET_NAME) && \
		  echo "$(PROJ_BUCKET_NAME) bucket created"; \
	fi;
