package discord

import (
	"fmt"
	"time"
	"net/url"
)

func (c Discord) SendMsg(message, chanID string, usrs []User) (Message, error) {
	//way 1.) look for @name and see if any users match the name
	//way 2.) use a passed in []User to fill mentions array
	//for now useway 2
	size := 0
	
	if usrs != nil {
		size = len(usrs)
	}
	ment := make([]string, size)
	if usrs != nil {
		for i, u := range usrs {
			ment[i] = u.ID
		}
	}
	req := MessageSend{
		Content: message,
		Mentions: ment,
		Nonce: fmt.Sprintf("%v", time.Now().Unix()), //almost always different.
		Tts: false,
	}
	resp := Message{}
	err := c.Post(fmt.Sprintf(ChanMsgsURL, chanID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("sent message successfully!")
	return resp, nil
}
func (c Discord) SendTextMsg(message, chanID string) (Message, error) {
	return c.SendMsg(message, chanID, nil)
}
func (c Discord) GetMsgs(chanID, before, after string, limit int) ([]Message, error) {
	resp := make([]Message, 0)
	baseURL := fmt.Sprintf(ChanMsgsURL+"?", chanID)
	params := url.Values{}
	if before != "" {
		params.Add("before", before)
	}
	if after != "" {
		params.Add("after", after)
	}
	if limit >= 0 {
		params.Add("limit", fmt.Sprintf("%v", limit))
	}
	
	fullURL := baseURL + params.Encode()
	//fmt.Println(fullURL)
	err := c.Get(fullURL, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("got messages successfully!")
	return resp, nil
}
func (c Discord) EditMsg(msg Message, newMsg string, usrs []User) (Message, error) {
	size := 0
	
	if usrs != nil {
		size = len(usrs)
	}
	ment := make([]string, size)
	if usrs != nil {
		for i, u := range usrs {
			ment[i] = u.ID
		}
	}
	//need messageID and channelID
	req := MessageSend{
		Content: newMsg,
		Mentions: ment,
	}
	resp := Message{}
	err := c.Patch(fmt.Sprintf(MsgURL, msg.ChanID, msg.ID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("edited message successfully!")
	return resp, nil
}
func (c Discord) EditTextMsg(msg Message, newMsg string) (Message, error) {
	return c.EditMsg(msg, newMsg, nil)
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
func (c Discord) DelMsg(msg Message) error {
	//need messageID and channelID
	err := c.Delete(fmt.Sprintf(MsgURL, msg.ChanID, msg.ID))
	if err != nil {
		//fmt.Println(err)
		return err
	}
	
	fmt.Println("deleted message successfully!")
	return nil
}