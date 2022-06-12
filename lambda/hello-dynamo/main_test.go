package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"testing"
)

func TestHandleLambdaEvent(t *testing.T) {
	type args struct {
		ctx         context.Context
		clientEvent ClientEvent
	}

	lc:=&lambdacontext.LambdaContext{
		AwsRequestID:"AwsRequestID123",
		InvokedFunctionArn:"InvokedFunctionArn123",
	}

	ctx:=context.Background()
	ctx=lambdacontext.NewContext(ctx,lc)
	tests := []struct {
		name    string
		args    args
		want    Response
		wantErr bool
	}{
		{
			name:"case 1",
			args: args{
				ctx: ctx,
				clientEvent: ClientEvent{
					ClientId: "clientid123",
					Msg: "some thing happen lol",
				},
			},
			want: Response{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleLambdaEvent(tt.args.ctx, tt.args.clientEvent)
			fmt.Println(got,err)
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("HandleLambdaEvent() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("HandleLambdaEvent() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
