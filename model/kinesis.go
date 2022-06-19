package model

import "errors"

type ClientEvent struct{
	ClientId *string `json:"client_id"`
	Msg *string `json:"msg"`
	IssueError bool `json:"issue_error"`
}

func (c *ClientEvent)Validate()error{
	if c.ClientId==nil{
		return errors.New("client_id missing")
	}
	if c.Msg==nil{
		return errors.New("msg missing")
	}

	return nil
}

