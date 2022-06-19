package model

type ClientEventReq struct {
	Data         ClientEvent `json:"Data"`
	PartitionKey string      `json:"PartitionKey"`
}

func NewClientEventReq(clientId string,msg string)*ClientEventReq{
	c:=&ClientEventReq{}
	c.Data.ClientId=&clientId
	c.Data.Msg=&msg
	c.PartitionKey=clientId
	return c
}
