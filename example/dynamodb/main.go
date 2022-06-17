package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func main() {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	// Build the request with its input parameters
	resp, err := svc.ListTables(context.TODO(), &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	})
	if err != nil {
		log.Fatalf("failed to list tables, %v", err)
	}

	fmt.Println("Tables:")
	for _, tableName := range resp.TableNames {
		fmt.Println(tableName)
	}

	event := NewEvent("awd123213", "something happen lol")

	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		panic(err)
	}
	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}

}
