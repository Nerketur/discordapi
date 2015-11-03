package discord

import (
	"fmt"
	"time"
)

func (c Discord) SendMsg(message, chanID string) (MessageResp, error) {
	
	req := Message{
		Content: message,
		Mentions: make([]string, 0),
		Nonce: fmt.Sprintf("%v", time.Now().Unix()),
		Tts: false,
	}
	resp := MessageResp{}
	err := c.Post(fmt.Sprintf(ChanMsgURL, chanID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Sent message successfully!")
	return resp, nil
}
func (c Discord) EditMsg(msg MessageResp, newMsg string) (MessageResp, error) {
	//need messageID and channelID
	req := Message{
		Content: newMsg,
		Mentions: make([]string, 0),
		//Nonce: fmt.Sprintf("%v", time.Now().Unix()),
		//Tts: false,
	}
	resp := MessageResp{}
	err := c.Send("PATCH", fmt.Sprintf(MsgURL, msg.ChanID, msg.ID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("edited message successfully!")
	return resp, nil
}
