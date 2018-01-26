package streaming

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/spf13/viper"
)

// Stream the interface for streaming pipeline
type Stream interface {
	SendMessage(string)
	GetIterator(shardID *string) *string
	GetRecords(iterator *string) ([]string, *string)
	// Do we want to expose these methods?
	// createStream()
	// waitUntilStreamExists()
}

type kinesisStream struct {
	streamName *string
	connection *kinesis.Kinesis
}

func (d kinesisStream) GetRecords(iterator *string) ([]string, *string) {

	// get records use shard iterator for making request
	messages, err := d.connection.GetRecords(&kinesis.GetRecordsInput{
		ShardIterator: iterator,
	})
	if err != nil {
		panic(err)
	}

	records := make([]string, len(messages.Records))
	for i, message := range messages.Records {
		records[i] = string(message.Data)
	}
	return records, messages.NextShardIterator
}

func (d kinesisStream) GetIterator(shardID *string) *string {
	// retrieve iterator
	iteratorOutput, err := d.connection.GetShardIterator(&kinesis.GetShardIteratorInput{
		// Shard Id is provided when making put record(s) request.
		ShardId:           shardID,
		ShardIteratorType: aws.String("TRIM_HORIZON"),
		// ShardIteratorType: aws.String("AT_SEQUENCE_NUMBER"),
		// ShardIteratorType: aws.String("LATEST"),
		StreamName: d.streamName,
	})
	if err != nil {
		panic(err)
	}
	return iteratorOutput.ShardIterator
}

// NewKinesisStream create a new kinesis stream
func NewKinesisStream(config viper.Viper) (Stream, error) {
	s := session.New(&aws.Config{Region: aws.String(config.GetString("kinesis.region")),
		Endpoint:   aws.String(config.GetString("kinesis.endpoint")),
		DisableSSL: aws.Bool(config.GetBool("kinesis.disable_ssl"))})

	kc := kinesis.New(s)

	streamName := aws.String(config.GetString("kinesis.stream"))

	return &kinesisStream{streamName: streamName, connection: kc}, nil
}

func (d kinesisStream) SendMessage(message string) {
	// d.createStream()
	// d.waitUntilStreamExists()

	streams, err := d.connection.DescribeStream(&kinesis.DescribeStreamInput{StreamName: d.streamName})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", streams)

	putOutput, err := d.connection.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte(message),
		StreamName:   d.streamName,
		PartitionKey: aws.String("key1"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", putOutput)
}

func (d kinesisStream) createStream() {
	out, err := d.connection.CreateStream(&kinesis.CreateStreamInput{
		ShardCount: aws.Int64(1),
		StreamName: aws.String(*d.streamName),
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", out)
}

func (d kinesisStream) waitUntilStreamExists() {
	if err := d.connection.WaitUntilStreamExists(&kinesis.DescribeStreamInput{StreamName: d.streamName}); err != nil {
		panic(err)
	}
}
