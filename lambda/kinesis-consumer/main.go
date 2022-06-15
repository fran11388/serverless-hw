package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fran11388/trendmicro-hw/model"
	"log"
)

type ClientEvent struct{
	ClientId string`json:"client_id"`
	Msg string `json:"msg"`
}


var dynamodbClient *dynamodb.Client
func init(){
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	dynamodbClient = dynamodb.NewFromConfig(cfg)
}



func HandleLambdaEvent(ctx context.Context, kinesisEvent events.KinesisEvent) error {
	for _, record := range kinesisEvent.Records {
		kinesisRecord := record.Kinesis
		dataBytes := kinesisRecord.Data
		dataText := string(dataBytes)
		fmt.Printf("%s Data = %s \n", record.EventName, dataText)

		clientevent:=&ClientEvent{}
		err := json.Unmarshal(dataBytes, clientevent)
		if err != nil {
			fmt.Println("some err happen", err)
			continue
		}

		eventdoc:=model.NewEvent(clientevent.ClientId,clientevent.Msg)
		item, err := attributevalue.MarshalMap(eventdoc)
		if err != nil {
			fmt.Println("some err happen", err)
			continue
		}
		_, err = dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(model.TABLE_NAME), Item: item,
		})
		if err != nil {
			log.Printf("Couldn't add item to table. Here's why: %v\n", err)
			continue
		}

		//todo error send to SQS

	}
	return nil

}

func main() {
	lambda.Start(HandleLambdaEvent)
}