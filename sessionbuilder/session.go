package sessionbuilder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sul-dlss-labs/taco/config"
)

// NewAwsSession creates a new session given the passed in config values
func NewAwsSession(config *config.Config) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		DisableSSL: aws.Bool(config.AwsDisableSSL),
	}))
}
