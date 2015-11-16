package discord

import (
	"fmt"
	"time"
	"net/url"
)

type _chan []Channel

func (c Discord) Chan(guild, name string) (Channel, error) {
	
	guild, err := c.Guild(guild)
	if err != nil {
		return Channel{}, err
	}
	chans, err := c.GuildChannels(guild)
	if err != nil {
		return Channel{}, err
	}
	return _chan(chans).Find(name, false)
}
func (c Discord) ChanID(guild, name string) (string, error) {
	channel, err := c.Chan(guild, name)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return channel.ID, nil
}
func (c Discord) PrivChan(name string) Channel {
	resp, err := _chan(c.MyChans).Find(name, true)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
func (c Discord) PrivChanID(name string) string {
	return c.PrivChan(name).ID
}
func (c _chan) Find(name string, private bool) (Channel, error) {
	for _, ele := range c {
		if private && (ele.Recipient.Username == name) || !private && (ele.Name == name) {
			return ele, nil
		}
	}
	return Channel{}, NotFoundError(name)
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
	err := c.Get(MyChansURL, &resp)
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

type userOrRole interface{
	GetID() string
	Type() string
}

func (x User) GetID() string {return x.ID}
func (x User) Type() string {return "member"}
func (x Role) GetID() string {return x.ID}
func (x Role) Type() string {return "role"}

func (c Discord) ChanReplacePerms(chanID string, ur userOrRole, allow, deny Perm) error {
	url := fmt.Sprintf(ChanPermIDURL, chanID, ur.GetID())
	fmt.Println(url)
	
	req := PermOver{
		Allow: allow,
		Deny: deny,
		ID: ur.GetID(),
		Type: ur.Type(),
	}
	err := c.Put(url, req)
	if err != nil {
		return err
	}
	
	fmt.Println("replaced chan perms!")
	return nil
}
func (c Discord) ChanDeletePerms(chanID string, ur userOrRole) error {
	url := fmt.Sprintf(ChanPermIDURL, chanID, ur.GetID())
	fmt.Println(url)
	
	err := c.Delete(url)
	if err != nil {
		return err
	}
	
	fmt.Println("deleted chan perms!")
	return nil
}