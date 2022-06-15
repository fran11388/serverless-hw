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

type Response struct{
	Event  *Event `json:"event"`
}

func NewEvent(clientid string, msg string) *Event {
	e := &Event{}
	e.PK = e.getPK(clientid)
	e.SK = e.getSK()
	e.Msg = msg
	return e
}
