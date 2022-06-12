package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"reflect"
	"testing"
)

func TestHandleLambdaEvent(t *testing.T) {
	type args struct {
		event MyEvent
	}
	tests := []struct {
		name    string
		args    args
		want    MyResponse
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				event: MyEvent{
					Name: "frank",
					Age: 20,
				},
			},
			want: MyResponse{
				Message:"frank is 20 years old!",
			},
			wantErr: false,
		},
	}
	lc:=&lambdacontext.LambdaContext{
		AwsRequestID:"AwsRequestID123",
		InvokedFunctionArn:"InvokedFunctionArn123",
	}

	ctx:=context.Background()
	ctx=lambdacontext.NewContext(ctx,lc)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleLambdaEvent(ctx,tt.args.event)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleLambdaEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleLambdaEvent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
