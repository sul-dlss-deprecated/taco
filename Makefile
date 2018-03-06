PROJECT                =taco
PROJECT_DIR            =$(shell pwd)
PATH_MAIN              =${PROJECT_DIR}/cmd/tacod/main.go
OS                     := $(shell go env GOOS)
ARCH                   := $(shell go env GOARCH)
BUILD_FLAGS            =-ldflags \"-s\" -a -installsuffix cgo
LOCALSTACK_SERVICES    =dynamodb,kinesis,s3
LOCAL_ENDPOINT         =--endpoint-url=http://localhost
DYNAMO_ENDPOINT        =${LOCAL_ENDPOINT}:4569
KINESIS_ENDPOINT       =${LOCAL_ENDPOINT}:4568
S3_ENDPOINT            =${LOCAL_ENDPOINT}:4572
PROJ_TABLE_NAME        =resources
PROJ_STREAM_NAME       =deposit
PROJ_BUCKET_NAME       =taco-deposited-files
PROJ_AWS_REGION        =us-west-2
PROJ_AWS_ACCESS_KEY_ID =999999
PROJ_AWS_SECRET_KEY    =999999
PROJ_ENV_VARS          =AWS_REGION=${PROJ_AWS_REGION} AWS_ACCESS_KEY_ID=${PROJ_AWS_ACCESS_KEY_ID} AWS_SECRET_KEY=${PROJ_AWS_SECRET_KEY}
PROJ_PORT              =8080
LOCAL_PORT             :=8080
LOCALSTACK_UI_PORT     =3000
LOCALSTACK_DATA_DIR    =/tmp/localstack/data
LOCALSTACK_DOCKER_PS   =$(shell docker ps -a | grep localstack/localstack | wc -l | xargs)
LOCALSTACK_DOCKER_Q    =$(shell docker ps -a | grep localstack/localstack | cut -d " " -f 1)

default: run

dependencies:
	go get github.com/golang/dep/cmd/dep
	dep ensure

create-resources: create-table create-stream create-bucket

create-table:
	$(eval TABLE=$(shell aws $(DYNAMO_ENDPOINT) dynamodb list-tables | jq '.TableNames[0] // ""'))
	@if [[ $(TABLE) != "" ]]; \
	  then echo "$(TABLE) table found"; \
		else aws $(DYNAMO_ENDPOINT) dynamodb create-table --table-name $(PROJ_TABLE_NAME) \
			--attribute-definitions "AttributeName=id,AttributeType=S" \
			--key-schema "AttributeName=id,KeyType=HASH" \
			--provisioned-throughput=ReadCapacityUnits=100,WriteCapacityUnits=100 ; \
	fi;

create-stream:
	$(eval STREAM=$(shell aws $(KINESIS_ENDPOINT) kinesis list-streams | jq '.StreamNames[0] // ""'))
	@if [[ $(STREAM) != "" ]]; \
    then echo "$(STREAM) stream found"; \
	  else aws $(KINESIS_ENDPOINT) kinesis create-stream --stream-name $(PROJ_STREAM_NAME) --shard-count 3 && \
		  echo "$(PROJ_STREAM_NAME) stream created"; \
	fi;

create-bucket:
	$(eval BUCKET=$(shell aws $(S3_ENDPOINT) s3api list-buckets | jq '.Buckets[0].Name // ""'))
	@if [[ $(BUCKET) != "" ]]; \
	  then echo "$(BUCKET) bucket found"; \
	  else aws $(S3_ENDPOINT) s3api create-bucket --bucket $(PROJ_BUCKET_NAME) && \
		  echo "$(PROJ_BUCKET_NAME) bucket created"; \
	fi;

docker-compose-up:
	docker-compose up -d
	sleep 5

docker-compose: docker-compose-up create-resources

localstack:
	PORT_WEB_UI=$(LOCALSTACK_UI_PORT) SERVICES=$(LOCALSTACK_SERVICES) DATA_DIR=$(LOCALSTACK_DATA_DIR) ENTRYPOINT=-d localstack --debug start --docker

run: localstack
	$(PROJ_ENV_VARS) TACO_PORT=$(PROJ_PORT) go run cmd/tacod/main.go

test: dependencies docker-compose
	$(PROJ_ENV_VARS) go test -v ./...
	docker-compose down

test-short: dependencies
	go test -v ./... -short

build-binary:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(PROJECT) $(BUILD_FLAGS) $(PATH_MAIN)

clean-localstack:
	@if [ $(LOCALSTACK_DOCKER_PS) -eq 1 ]; \
	  then docker stop "$(LOCALSTACK_DOCKER_Q)" && docker rm "$(LOCALSTACK_DOCKER_Q)"; \
		else echo "localstack docker is either not running or more than one is running"; \
	fi

clean-localstack-datadir:
	@if [ -d $(LOCALSTACK_DATA_DIR) ]; \
	  then echo "removing $(LOCALSTACK_DATA_DIR)" && rm -r $(LOCALSTACK_DATA_DIR); \
		else echo "localstack data dir does not exist"; \
	fi

# docker-compose down
# rm taco
# rm -rf /tmp/localstack/data
