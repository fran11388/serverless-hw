package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fran11388/trendmicro-hw/model"
	"net/http"
	"os"
	"runtime"
)

type LambdaEvent struct {
	ClientId   *string `json:"client_id"`
	EventCount *int    `json:"event_count"`
}

func (l *LambdaEvent) Validate() error {
	errstr := "%s missing"
	if l.ClientId == nil {
		return errors.New(fmt.Sprintf(errstr, "client_id"))
	}
	if l.EventCount == nil {
		return errors.New(fmt.Sprintf(errstr, "event_count"))
	}
	return nil
}

var serviceEndPoint string

func init() {
	serviceEndPoint = os.Getenv("SERVICE_END_POINT")
}

func HandleLambdaEvent(ctx context.Context, event LambdaEvent) error {
	if serviceEndPoint == "" {
		err := "env SERVICE_END_POINT not been setting"
		fmt.Println(err)
		return errors.New(err)
	}

	err := event.Validate()
	if err != nil {
		return err
	}

	err = concurrentSendEventToService(event)
	if err != nil {
		return err
	}

	return nil
}

func concurrentSendEventToService(event LambdaEvent) error {
	workerSize := runtime.NumCPU()
	clientId := *event.ClientId
	eventCount := *event.EventCount
	finishChan := make(chan bool, eventCount)
	errChan := make(chan error, workerSize)

	//init job
	jobs := make(chan *model.ClientEventReq, eventCount)
	for i := 0; i < eventCount; i++ {
		msg := fmt.Sprintf("This event is generated by lambda client-event-producer, event: %d", i)
		e := model.NewClientEventReq(clientId, msg)
		jobs <- e
	}
	close(jobs)

	//init worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < workerSize; i++ {
		go func(ctx context.Context, jobs chan *model.ClientEventReq) {
			defer func() {
				if err := recover(); err != nil {
					s := err.(string)
					errChan <- errors.New(s)
				}
			}()

			for job:=range jobs{
				select {
				case <-ctx.Done():
					return
				default:
					err := sendEventToService(job)
					if err != nil {
						errChan <- err
						return
					} else {
						finishChan <- true
					}
				}
			}
		}(ctx,jobs)
	}

	for i := 0; i < eventCount; i++ {
		select {
		case <-finishChan:
		case err := <-errChan:
			cancel()
			return err
		}
	}
	return nil
}

func sendEventToService(e *model.ClientEventReq) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(data)
	req, err := http.NewRequest("PUT", serviceEndPoint, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
