package streaming

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// Stream the interface for streaming pipeline
type Stream interface {
	SendMessage(string) error
}

type KinesisStream struct {
	StreamName string
	Connection *kinesis.Kinesis
}

func Connect(session *session.Session, kinesisEndpoint string) *kinesis.Kinesis {
	return kinesis.New(session, &aws.Config{Endpoint: aws.String(kinesisEndpoint)})
}
func (d KinesisStream) SendMessage(message string) error {
	streams, err := d.Connection.DescribeStream(&kinesis.DescribeStreamInput{StreamName: &d.StreamName})
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", streams)

	putOutput, err := d.Connection.PutRecord(&kinesis.PutRecordInput{
		Data:         []byte(message),
		StreamName:   &d.StreamName,
		PartitionKey: aws.String("key1"),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", putOutput)
	return nil
}
