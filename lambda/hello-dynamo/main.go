package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"log"
	"time"
)

type MyTable struct {
	PK string `dynamodbav:"pk"`
	SK string `dynamodbav:"sk"`
}

type Event struct {
	MyTable
	Msg string `dynamodbav:"msg"`
}

func (e *Event) getPK(clientid string) string {
	return fmt.Sprintf("Event#CliendId#%s", clientid)
}

func (e *Event) getSK() string {
	u4 := uuid.New()
	uuidstr := u4.String()
	t := time.Now()
	timestamp := fmt.Sprintf("%d", t.Unix())
	return fmt.Sprintf("Timestamp#%s#UUID#%s", timestamp, uuidstr)
}

func NewEvent(clientid string, msg string) *Event {
	e := &Event{}
	e.PK = e.getPK(clientid)
	e.SK = e.getSK()
	e.Msg = msg
	return e
}

const TABLE_NAME = "mytable"

type ClientEvent struct{
	ClientId string`json:"client_id"`
	Msg string `json:"msg"`
}

type Response struct{
	Event  *Event `json:"event"`
}


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

func HandleLambdaEvent(ctx context.Context, clientEvent ClientEvent) (Response, error) {
	lc, _ := lambdacontext.FromContext(ctx)
	log.Print(lc)

	event := NewEvent(clientEvent.ClientId, clientEvent.Msg)

	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		panic(err)
	}
	_, err = dynamodbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}

	return Response{Event: event}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}