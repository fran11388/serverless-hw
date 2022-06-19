package main

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestHandleLambdaEvent(t *testing.T) {
	type args struct {
		ctx   context.Context
		event LambdaEvent
	}
	err := os.Setenv("SERVICE_END_POINT", "https://px3bcgybfh.execute-api.ap-northeast-1.amazonaws.com/test/streams/my-event-stream/record")
	if err != nil {
		fmt.Println(err)
		return
	}
	clientid:="local_test_client1"
	eventCount:=5
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:"case1",
			args: args{
				ctx: context.TODO(),
				event: LambdaEvent{
					ClientId: &clientid,
					EventCount: &eventCount,
				},
			},

		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HandleLambdaEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("HandleLambdaEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
