package aws_session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func Connect(disableSSL bool) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		DisableSSL: aws.Bool(disableSSL),
	}))
}
