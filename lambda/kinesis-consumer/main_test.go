package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"testing"
)

func TestHandleLambdaEvent(t *testing.T) {
	type args struct {
		ctx          context.Context
		kinesisEvent events.KinesisEvent
	}

	kinesisEvent:=events.KinesisEvent{
		Records: []events.KinesisEventRecord{
			events.KinesisEventRecord{
				Kinesis:events.KinesisRecord{
					Data:[]byte("{\n    \"client_id\": \"client123\",\n    \"msg\": \"from golang ide some thing happen lol\"\n  }"),
				},
		}},
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:"case",
			args: args{
				ctx: context.TODO(),
				kinesisEvent:kinesisEvent,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HandleLambdaEvent(tt.args.ctx, tt.args.kinesisEvent); (err != nil) != tt.wantErr {
				t.Errorf("HandleLambdaEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
