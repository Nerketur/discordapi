package discord

import (
	"fmt"
	"time"
	"net/url"
	"encoding/json"
)

type _chan []Channel

func (c Discord) Channel(guild, name string) (string, error) {
	
	guild, err := c.Guild(guild)
	if err != nil {
		return "", err
	}
	chans, err := c.GuildChannels(guild)
	if err != nil {
		return "", err
	}
	resp, err := _chan(chans).Find(name, false)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return resp, nil
}
func (c Discord) PivChannel(name string) string {
	resp, err := _chan(c.MyChans).Find(name, true)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
func (c _chan) Find(name string, private bool) (string, error) {
	for _, ele := range c {
		if private && (ele.Recipient.Username == name) || !private && (ele.Name == name) {
			return ele.ID, nil
		}
	}
	return "", NotFoundError(name)
}

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
	
	fmt.Println("sent message successfully!")
	return resp, nil
}
func (c Discord) GetMsgs(chanID, before, after string, limit int) ([]Message, error) {
	
	resp := make([]Message, limit)
	baseURL := fmt.Sprintf(ChanMsgsURL+"?", chanID)
	params := url.Values{}
	if before != "" {
		params.Add("before", before)
	}
	if after != "" {
		params.Add("after", after)
	}
	if limit > 1 {
		params.Add("limit", string(limit))
	}
	
	fullURL := baseURL + params.Encode()
	fmt.Println(fullURL)
	err := c.Get(fullURL, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("got messages successfully!")
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
func (c Discord) DelMsg(msg Message) error {
	//need messageID and channelID
	err := c.Send("DELETE", fmt.Sprintf(MsgURL, msg.ChanID, msg.ID), nil, nil)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	
	fmt.Println("deleted message successfully!")
	return nil
}
func (c Discord) PrivateChannels() ([]Channel, error) {
	
	resp := make([]Channel, 0)
	err := c.Get(fmt.Sprintf(UserChansURL, "@me"), &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("got private channels successfully!")
	return resp, nil
}
func (c Discord) SendTyping(chanID string) error {
	
	err := c.Post(fmt.Sprintf(ChanTypingURL, chanID), nil, nil)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	
	fmt.Println("typing sent!")
	return nil
}

// interface PermsNeeded {
//	
// }

type ChanPermsMsg struct{
	perms []string
	msg   struct{
		Message *string
	}
}

func (s *ChanPermsMsg) UnmarshalJSON(raw []byte) error {
	//first we try unmarshaling into roles
	err1 := json.Unmarshal(raw, &s.perms)
	if err1 != nil {
		//retry with message
		err2 := json.Unmarshal(raw, &s.msg)
		if err2 != nil {
			//error, invalid
			return err2
		}
	}
	return nil
}

func (c Discord) ChanPerms(chanID string) ([]string, error) {
	resp := ChanPermsMsg{}
	err := c.Get(fmt.Sprintf(ChanPermsURL, chanID), &resp)
	if err != nil {
		return resp.perms, err
	}
	if resp.msg.Message != nil {
		return resp.perms, MissingPermissionError(*resp.msg.Message)
	}
	
	fmt.Println("got chan perms!")
	return resp.perms, nil
}