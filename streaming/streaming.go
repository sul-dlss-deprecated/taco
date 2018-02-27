package streaming

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// Stream the interface for streaming pipeline
type Stream interface {
	SendMessage(Message) error
}

type KinesisStream struct {
	StreamName *string
	Connection *kinesis.Kinesis
}

func Connect(session *session.Session, kinesisEndpoint string) *kinesis.Kinesis {
	return kinesis.New(session, &aws.Config{Endpoint: aws.String(kinesisEndpoint)})
}

// SendMessage publishes the given message to the kinesis stream
func (d *KinesisStream) SendMessage(message Message) error {
	strmsg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = d.Connection.PutRecord(&kinesis.PutRecordInput{
		Data:         strmsg,
		StreamName:   d.StreamName,
		PartitionKey: aws.String("key1"),
	})
	return err
}
