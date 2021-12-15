package database

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gabriel-flynn/Cheap-Flight-Finder/server/models"
	"log"
	"sync"
)

var DDBClient *dynamodb.DynamoDB

func SaveOneWayFlight(flight *models.OneWayFlight) {
	av, err := dynamodbattribute.MarshalMap(flight)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}
	_, err = DDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("Flights"),
		Item:      av,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to put Record to DynamoDB, %v", err))
	}
}

func SaveOneWayFlightBatch(flights []*models.OneWayFlight) {
	var wg sync.WaitGroup
	for i := 0; i < len(flights); i+=25 {
		wg.Add(1)
		rightBound := i+25
		if rightBound > len(flights) {
			rightBound = len(flights)
		}

		go func(i int) {
			saveOneWayFlightBatch(flights[i:rightBound])
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func saveOneWayFlightBatch(flights []*models.OneWayFlight) {
	batchWriteItemInput := make(map[string][]*dynamodb.WriteRequest)
	for _, flight := range flights {
		if flight == nil {
			continue
		}
		av, err := dynamodbattribute.MarshalMap(flight)
		if err != nil {
			log.Fatal("failed to DynamoDB marshal Records, ", err)
		}
		writeRequest := &dynamodb.WriteRequest{
			DeleteRequest: nil,
			PutRequest:    &dynamodb.PutRequest{Item: av},
		}
		batchWriteItemInput["Flights"] = append(batchWriteItemInput["Flights"], writeRequest)
	}

	_, err := DDBClient.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems:                batchWriteItemInput,
		ReturnConsumedCapacity:      nil,
		ReturnItemCollectionMetrics: nil,
	})
	if err != nil {
		fmt.Println(batchWriteItemInput)
		log.Fatal("Failed to batch write records: ", err)
	}
}

func init() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	DDBClient = dynamodb.New(sess)
}
