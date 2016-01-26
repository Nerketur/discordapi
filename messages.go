package discord

import (
	"fmt"
	"time"
	"net/url"
	"regexp"
)

func (c Discord) SendRawMsg(message, chanID string, tts bool) (resp Message, err error) {
	if len(message) > 2000 {
		//TODO: change into error returned
		message = fmt.Sprintf("Content too long! (by %d chars)", len(message)-2000)
	}
	req := MessageSend{
		Content: message,
		Nonce: time.Now().Unix(), //almost always different.
		Tts: tts,
	}
	err = c.Post(fmt.Sprintf(ChanMsgsURL, chanID), req, &resp)
	if err != nil {
		//fmt.Println(err)
		return
	}
	fmt.Println("sent message successfully!")
	return
}

func (c Discord) replName2ID(match string) string {
	kind := string(match[:1])
	name := string(match[1:len(match)])
	switch kind {
	case "@":
		//does the user exist?
		if us, err := c.FindNameUserCache(name); err == nil {
			return fmt.Sprintf("<@%s>", us[0].ID)
		}
	case "#": //TODO: fix channels, requires guildID
		//does the channel exist?
		if chs, err := c.FindNameChanCache(name); err == nil {
			return fmt.Sprintf("<#%s>", chs[0].ID)
		}
	}
	return match
}
func (c Discord) replID2Name(match string) (ret string) {
	kind := string(match[1:2])
	ID := string(match[2:len(match)-1])
	switch kind {
	case "@":
		//does the user exist?
		if u, err := c.UserCache(ID); err != nil {
			return "@" + u.Username
		}
	case "#":
		//does the channel exist?
		if ch, err := c.ChanCache(ID); err != nil {
			return "#" + ch.Name
		}
	}
	return match
}

/*
func main() {
*/
func (c Discord) SendMsg(message, chanID string, tts bool) (Message, error) {
	//look for @name and see if any users match the name
	//look for #name and see if any channels match the name
	
	re := regexp.MustCompile(`[@#](\w*)`)
	message = re.ReplaceAllStringFunc(message, c.replName2ID)
	
	return c.SendRawMsg(message, chanID, tts)
}
func (c Discord) SendTextMsg(message, chanID string) (Message, error) {
	return c.SendMsg(message, chanID, false)
}
func (c Discord) SendSpeechMsg(message, chanID string) (Message, error) {
	return c.SendMsg(message, chanID, true)
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
func (c Discord) AckMsg(msg Message) (err error) {
	//need messageID and channelID
	err = c.Post(fmt.Sprintf(MsgAckURL, msg.ChanID, msg.ID), nil, nil)
	if err != nil {
		//fmt.Println(err)
		return
	}
	
	fmt.Println("message acknowledged successfully!")
	return
}
func (c Discord) DelMsg(msg Message) (err error) {
	//need messageID and channelID
	err = c.Delete(fmt.Sprintf(MsgURL, msg.ChanID, msg.ID))
	if err != nil {
		//fmt.Println(err)
		return
	}
	
	fmt.Println("deleted message successfully!")
	return
}