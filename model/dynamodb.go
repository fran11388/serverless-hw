package model

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const TABLE_NAME = "mytable"

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

type ErrorLog struct {
	MyTable
	Error string
}

func (e *ErrorLog) getPK() string {
	return fmt.Sprintf("ErrorLog")
}

func (e *ErrorLog) getSK(timestamp int, uuidstr string) string {
	return fmt.Sprintf("Timestamp#%dUUID#%s", timestamp, uuidstr)
}

type NewErrorLogInput struct {
	Timestamp int
	UUID      string
	Error     string
}

func NewErrorLog(input *NewErrorLogInput) *ErrorLog {
	e := &ErrorLog{}
	e.PK = e.getPK()
	e.SK = e.getSK(input.Timestamp, input.UUID)
	e.Error = input.Error
	return e
}
