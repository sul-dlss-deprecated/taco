package streaming

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/sul-dlss-labs/taco/config"
)

// Stream the interface for streaming pipeline
type Stream interface {
	SendMessage(string) error
}

type kinesisStream struct {
	streamName *string
	connection *kinesis.Kinesis
}

// NewKinesisStream create a new kinesis stream
func NewKinesisStream(config *config.Config) Stream {
	s, err := session.NewSession(&aws.Config{
		Endpoint:   aws.String(config.KinesisEndpoint),
		DisableSSL: aws.Bool(config.KinesisDisableSSL),
	})
	if err != nil {
		panic(err)
	}
	kc := kinesis.New(s)

	streamName := aws.String(config.DepositStreamName)

	return &kinesisStream{streamName: streamName, connection: kc}
}

func (d kinesisStream) SendMessage(message string) error {
	streams, err := d.connection.DescribeStream(&kinesis.DescribeStreamInput{StreamName: d.streamName})
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", streams)

	putOutput, err := d.connection.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte(message),
		StreamName:   d.streamName,
		PartitionKey: aws.String("key1"),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", putOutput)
	return nil
}
