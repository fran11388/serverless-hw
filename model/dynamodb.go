package model

import (
	"fmt"
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

func (e *Event) getSK(timestamp time.Time,sn string) string {
	t := fmt.Sprintf("%d", timestamp.Unix())
	return fmt.Sprintf("Timestamp#%s#SN#%s", t, sn)
}

type NewEventInput struct{
	ClientId string
	Msg string
	Timestamp time.Time
	SN string
}

func NewEvent(input *NewEventInput) *Event {
	e := &Event{}
	e.PK = e.getPK(input.ClientId)
	e.SK = e.getSK(input.Timestamp,input.SN)
	e.Msg = input.Msg
	return e
}

type ErrorLog struct {
	MyTable
	Error string
	*ClientEvent
}

func (e *ErrorLog) getPK() string {
	return fmt.Sprintf("ErrorLog")
}

func (e *ErrorLog) getSK(timestamp string, SN string) string {
	return fmt.Sprintf("Timestamp#%s#SN#%s", timestamp, SN)
}

type NewErrorLogInput struct {
	Timestamp string
	SN      string
	Error     string
	ClientEvent *ClientEvent
}

func NewErrorLog(input *NewErrorLogInput) *ErrorLog {
	e := &ErrorLog{}
	e.PK = e.getPK()
	e.SK = e.getSK(input.Timestamp, input.SN)
	e.Error = input.Error
	e.ClientEvent=input.ClientEvent
	return e
}
