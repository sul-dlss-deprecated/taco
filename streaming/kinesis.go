package streaming

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/sul-dlss-labs/taco/datautils"
)

// KinesisStream represents the kinesis implementation of the Stream interface
type KinesisStream struct {
	StreamName *string
	Connection *kinesis.Kinesis
}

// Connect initialize the kinesis struct
func Connect(session *session.Session, kinesisEndpoint string) *kinesis.Kinesis {
	return kinesis.New(session, &aws.Config{Endpoint: aws.String(kinesisEndpoint)})
}

// Send creates a message and publishes it to the stream
func (d *KinesisStream) Send(action string, resource *datautils.Resource) error {
	return d.SendMessage(NewMessage(action, resource))
}

// SendMessage publishes the given message to the kinesis stream
func (d *KinesisStream) SendMessage(message *Message) error {
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
