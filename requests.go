package discord

import (
	"fmt"
	"io"
	"bytes"
	"net/http"
	"encoding/json"
	"os"
)

const (
	LoginURL = "https://discordapp.com/api/auth/login"
	LogoutURL = "https://discordapp.com/api/auth/logout"
	ChanMsgURL = "https://discordapp.com/api/channels/%s/messages"
	MsgURL = "https://discordapp.com/api/channels/%s/messages/%s" //channel ID, message ID
)

func (c Discord) Send(method, url string, data, want interface{}) error {
	if c.Token == "" && !c.LoggingIn { // not logged in or logging in
		return TokenError("Not logged in!")
	}
	var req *http.Request
	var err error
	if data != nil {
		//fmt.Println("data not nil")

		b, err := json.Marshal(data)
		if err != nil {
			return EncodingError(fmt.Sprintf("to-json:%s", err))
		}
		//os.Stdout.Write(b)
		//os.Stdout.Write([]byte{10})
		req, err = http.NewRequest(method, url, bytes.NewBuffer(b))
		if err != nil {
			return UnknownError(fmt.Sprintf("%s", err))
		}
	} else {
		//fmt.Println("data nil")
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return UnknownError(fmt.Sprintf("%s", err))
		}
	}
	//return UnknownError("Testing!")
	if req == nil {
		return UnknownError("Why is req nil???")
	}
	if c.Token != "" {
		//fmt.Println("adding token")
		req.Header.Add("Authorization", c.Token)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	resp, err := c.Client.Do(req) // at this point, bytes buffer is closed if needed
	defer resp.Body.Close()
	if err != nil {
		return PostError(fmt.Sprintf("%s", err))
	}
	
	var raw []byte
	if resp.ContentLength > 0 {
		raw = make([]byte, resp.ContentLength)
	} else {
		raw = make([]byte, 1024)
	}
	num, err := resp.Body.Read(raw)
	if err == io.EOF {
		//reached the end, ignore.
	} else if err != nil {
		return ReadError(fmt.Sprintf("%s", err))
	}
	raw = raw[:num]
	os.Stdout.Write(raw)
	os.Stdout.Write([]byte{10})
	
	//loginResp := CredsResp{}
	//fmt.Println("before unmarshal")
	
	if want != nil {
		err = json.Unmarshal(raw, &want)
		//fmt.Println("after unmarshal")
		if err != nil {
			return EncodingError(fmt.Sprintf("from-json:%s", err))
		}
	}
	
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return HTTPError(*resp)
	}
	c.LoggingIn = false
	return nil
}
func (c Discord) Post(url string, req, resp interface{}) error {
	return c.Send("POST", url, req, resp)
}
func (c Discord) Get(url string, resp interface{}) error {
	return c.Send("GET", url, nil, resp)
}
