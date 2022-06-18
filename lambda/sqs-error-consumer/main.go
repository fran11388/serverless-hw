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

var dynamodbClient *dynamodb.Client
func init(){
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	dynamodbClient = dynamodb.NewFromConfig(cfg)
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		sqsErrMsg:=&model.SQSErrorMsgBody{}
		err := json.Unmarshal([]byte(message.Body), sqsErrMsg)
		if err != nil {
			fmt.Println("Unmarshal error:")
			fmt.Println(err)
			continue
		}

		errordoc:=model.NewErrorLog(&model.NewErrorLogInput{
			Timestamp :message.Attributes["ApproximateFirstReceiveTimestamp"],
			SN      :message.MessageId,
			Error     :sqsErrMsg.Error,
			ClientEvent :sqsErrMsg.ClientEvent,
		})


		item, err := attributevalue.MarshalMap(errordoc)
		if err != nil {
			fmt.Println("MarshalMap error:")
			fmt.Println(err)
			continue
		}
		_, err = dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(model.TABLE_NAME), Item: item,
		})
		if err != nil {
			log.Printf("Couldn't add item to table. Here's why: %v\n", err)
			continue
		}

	}

	return nil
}

func main() {
	lambda.Start(handler)
}