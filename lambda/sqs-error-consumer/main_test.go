package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func Test_handler(t *testing.T) {
	type args struct {
		ctx      context.Context
		sqsEvent events.SQSEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:"case1",
			args: args{
				ctx:context.TODO(),
				sqsEvent: events.SQSEvent{
					Records: []events.SQSMessage{
						events.SQSMessage{
							MessageId:"MessageId123",
							//ReceiptHandle
							Body:"{\"client_event\":{\"client_id\":\"client831209\",\"msg\":\"send from postman 8\",\"issue_error\":true},\"error\":\"client issue error\"}",
							//Md5OfBody
							//Md5OfMessageAttributes
							Attributes:map[string]string{"ApproximateFirstReceiveTimestamp":"1234567891"},
							//MessageAttributes
							//EventSourceARN
							//EventSource
							//AWSRegion
						},
					}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handler(tt.args.ctx, tt.args.sqsEvent); (err != nil) != tt.wantErr {
				t.Errorf("handler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
