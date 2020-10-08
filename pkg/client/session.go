package client

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// GetAwsSession creates new session for AWS
func GetAwsSession() *session.Session {
	mySession := session.Must(session.NewSession())
	return mySession
}
