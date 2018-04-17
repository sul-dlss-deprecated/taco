PROJECT                =taco
PROJECT_DIR            =$(shell pwd)
PATH_MAIN              =${PROJECT_DIR}/main.go
OS                     := $(shell go env GOOS)
ARCH                   := $(shell go env GOARCH)
PROJECT_UNIT_TESTS     =$(shell go list ./... | grep -v test | grep -v db)
BUILD_FLAGS            =-ldflags \"-s\" -a -installsuffix cgo
LOCALSTACK_SERVICES    =dynamodb,kinesis,s3
LOCAL_ENDPOINT_HOST    :=${LOCAL_ENDPOINT_HOST}
LOCAL_ENDPOINT         =--endpoint-url=http://${LOCAL_ENDPOINT_HOST}
DYNAMO_ENDPOINT        =${LOCAL_ENDPOINT}:4569
KINESIS_ENDPOINT       =${LOCAL_ENDPOINT}:4568
S3_ENDPOINT            =${LOCAL_ENDPOINT}:4572
PROJ_TABLE_NAME        =resources
PROJ_STREAM_NAME       =deposit
PROJ_BUCKET_NAME       =taco-deposited-files
PROJ_AWS_REGION        =us-east-1
PROJ_AWS_ACCESS_KEY_ID =999999
PROJ_AWS_SECRET_KEY    =999999
PROJ_ENV_VARS          =AWS_REGION=${PROJ_AWS_REGION} AWS_ACCESS_KEY_ID=${PROJ_AWS_ACCESS_KEY_ID} AWS_SECRET_KEY=${PROJ_AWS_SECRET_KEY}
PROJ_PORT              =8080
LOCAL_PORT             :=8080
LOCALSTACK_UI_PORT     =3000
LOCALSTACK_DATA_DIR    =/tmp/localstack/data
LOCALSTACK_DOCKER_PS   =$(shell docker ps -a | grep localstack/localstack | wc -l | xargs)
LOCALSTACK_DOCKER_Q    =$(shell docker ps -a | grep localstack/localstack | cut -d " " -f 1)

default: test

dependencies:
	go get github.com/golang/dep/cmd/dep
	dep ensure

resources: table stream bucket

table:
	$(eval TABLE=$(shell aws $(DYNAMO_ENDPOINT) dynamodb list-tables | jq '.TableNames[0] // ""'))
	@if [[ $(TABLE) != "" ]]; \
	  then echo "$(TABLE) table found"; \
		else aws $(DYNAMO_ENDPOINT) dynamodb create-table --table-name $(PROJ_TABLE_NAME) \
				--attribute-definitions AttributeName=tacoIdentifier,AttributeType=S \
				AttributeName=externalIdentifier,AttributeType=S \
				AttributeName=version,AttributeType=N \
				--key-schema "AttributeName=tacoIdentifier,KeyType=HASH" \
				--provisioned-throughput=ReadCapacityUnits=10,WriteCapacityUnits=10 \
				--global-secondary-indexes "IndexName=ResourceByExternalIDAndVersion, \
					KeySchema=[{AttributeName=externalIdentifier,KeyType=HASH}, \
					           {AttributeName=version,KeyType=RANGE}], \
					Projection={ProjectionType=ALL}, \
					ProvisionedThroughput={ReadCapacityUnits=10,WriteCapacityUnits=10}" ; \
	fi;

stream:
	$(eval STREAM=$(shell aws $(KINESIS_ENDPOINT) kinesis list-streams | jq '.StreamNames[0] // ""'))
	@if [[ $(STREAM) != "" ]]; \
    then echo "$(STREAM) stream found"; \
	  else aws $(KINESIS_ENDPOINT) kinesis create-stream --stream-name $(PROJ_STREAM_NAME) --shard-count 3 && \
		  echo "$(PROJ_STREAM_NAME) stream created"; \
	fi;

bucket:
	$(eval BUCKET=$(shell aws $(S3_ENDPOINT) s3api list-buckets | jq '.Buckets[0].Name // ""'))
	@if [[ $(BUCKET) != "" ]]; \
	  then echo "$(BUCKET) bucket found"; \
	  else aws $(S3_ENDPOINT) s3api create-bucket --bucket $(PROJ_BUCKET_NAME) && \
		  echo "$(PROJ_BUCKET_NAME) bucket created"; \
	fi;

localstack:
	PORT_WEB_UI=$(LOCALSTACK_UI_PORT) SERVICES=$(LOCALSTACK_SERVICES) DATA_DIR=$(LOCALSTACK_DATA_DIR) ENTRYPOINT=-d localstack --debug start --docker

dev: localstack resources
	$(PROJ_ENV_VARS) go run main.go

test: dependencies
	$(PROJ_ENV_VARS) go test -v $(PROJECT_UNIT_TESTS)

e2e_test:
	$(PROJ_ENV_VARS) go test test/integration_test.go
	$(PROJ_ENV_VARS) go test db/database_test.go
