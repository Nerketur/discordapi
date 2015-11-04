package discord

import (
	"fmt"
	"time"
)

func (c Discord) SendMsg(message, chanID string) (Message, error) {
	
	req := MessageSend{
		Content: message,
		Mentions: make([]string, 0),
		Nonce: fmt.Sprintf("%v", time.Now().Unix()), //almost always different.
		Tts: false,
	}
	resp := Message{}
	err := c.Post(fmt.Sprintf(ChanMsgsURL, chanID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Sent message successfully!")
	return resp, nil
}
func (c Discord) GetMsgs(chanID, before string, limit int) ([]Message, error) {
	
	resp := make([]Message, limit)
	url := fmt.Sprintf(ChanMsgsURL+"?", chanID)
	if before != "" {
		url += fmt.Sprintf("before=%v", before)
	}
	if limit > 1 {
		if before != "" { url += "&" }
		url += fmt.Sprintf("limit=%v", limit)
	}
	
	fmt.Println(url)
	//return resp.Msgs, UnknownError("Testing!")
	err := c.Get(url, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Got messages successfully!")
	return resp, nil
}
func (c Discord) EditMsg(msg Message, newMsg string) (Message, error) {
	//need messageID and channelID
	req := MessageSend{
		Content: newMsg,
		Mentions: make([]string, 0),
	}
	resp := Message{}
	err := c.Send("PATCH", fmt.Sprintf(MsgURL, msg.ChanID, msg.ID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("edited message successfully!")
	return resp, nil
}
func (c Discord) AckMsg(msg Message) error {
	//need messageID and channelID
	err := c.Post(fmt.Sprintf(MsgAckURL, msg.ChanID, msg.ID), nil, nil)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	
	fmt.Println("message acknowledged successfully!")
	return nil
}
func (c Discord) DelMsg(msg Message, newMsg string) (Message, error) {
	//need messageID and channelID
	req := MessageSend{
		Content: newMsg,
		Mentions: make([]string, 0),
		//Nonce: fmt.Sprintf("%v", time.Now().Unix()),
		//Tts: false,
	}
	resp := Message{}
	err := c.Send("DELETE", fmt.Sprintf(MsgURL, msg.ChanID, msg.ID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("deleted message successfully!")
	return resp, nil
}
