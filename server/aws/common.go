package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var sess *session.Session

func GetSession() *session.Session {
	if sess == nil {
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String("us-east-1"),
			},
			SharedConfigState: session.SharedConfigEnable,
		}))
	}
	return sess
}