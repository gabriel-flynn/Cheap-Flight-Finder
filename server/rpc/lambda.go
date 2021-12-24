package rpc

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	aws2 "github.com/gabriel-flynn/Cheap-Flight-Finder/server/aws"
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

	lambdaClient = lambda.New(aws2.GetSession(), &aws.Config{Region: aws.String("us-east-1")})
}
