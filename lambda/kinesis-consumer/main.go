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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/fran11388/trendmicro-hw/model"
	sqsutils "github.com/fran11388/trendmicro-hw/utils/aws/sqs"
	"log"
)

type ClientEvent struct{
	ClientId string `json:"client_id"`
	Msg string `json:"msg"`
}

const SQS_Error_Queue_Name="consumer-error-queue"
var dynamodbClient *dynamodb.Client
var sqsClient *sqs.Client
var queueURL *string
func init(){
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
		return
	}

	// Using the Config value, create the DynamoDB client
	dynamodbClient = dynamodb.NewFromConfig(cfg)
	sqsClient=sqs.NewFromConfig(cfg)
	queue:=SQS_Error_Queue_Name
	// Get URL of queue
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queue,
	}
	result, err := sqsutils.GetQueueURL(context.TODO(), sqsClient, gQInput)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}
	queueURL = result.QueueUrl
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
			sendToSQS(clientevent,err.Error())
			continue
		}

		eventdoc:=model.NewEvent(clientevent.ClientId,clientevent.Msg)
		item, err := attributevalue.MarshalMap(eventdoc)
		if err != nil {
			sendToSQS(clientevent,err.Error())
			continue
		}
		_, err = dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(model.TABLE_NAME), Item: item,
		})
		if err != nil {
			errReason:=fmt.Sprintf("Couldn't add item to table. Here's why: %v",err)
			sendToSQS(clientevent,errReason)
			continue
		}

	}
	return nil
}

type SQSMsgBody struct{
	ClientEvent *ClientEvent `json:"client_event"`
	Error string `json:"error"`
}
func sendToSQS(ce *ClientEvent,errReason string){
	msgbody:=&SQSMsgBody{
		ClientEvent: ce,
		Error: errReason,
	}
	msgBytes, err := json.Marshal(msgbody)
	if err != nil {
		fmt.Println(err)
		return
	}

	sMInput := &sqs.SendMessageInput{
		MessageBody: aws.String(string(msgBytes)),
		QueueUrl:    queueURL,
	}

	_, err = sqsutils.SendMsg(context.TODO(), sqsClient, sMInput)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return
	}
}

func main() {
	lambda.Start(HandleLambdaEvent)
}