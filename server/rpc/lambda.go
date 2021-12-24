package rpc

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
)

var lambdaClient *lambda.Lambda

func GetPageSourceLambda(payload []byte) string {

	result, err := lambdaClient.Invoke(&lambda.InvokeInput{FunctionName: aws.String("lambda-prod-scrape-frontier"), Payload: payload})
	if err != nil {
		log.Fatalln(err)
	}

	return string(result.Payload)
}

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	lambdaClient = lambda.New(sess, &aws.Config{Region: aws.String("us-east-1")})
}
