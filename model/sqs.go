package model


type SQSErrorMsgBody struct{
	ClientEvent *ClientEvent `json:"client_event"`
	Error string `json:"error"`
}
