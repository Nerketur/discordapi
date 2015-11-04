package discord

import (
	"fmt"
	"io"
	"bytes"
	"net/http"
	"encoding/json"
	//"os"
)

const (
	APIURL = "https://discordapp.com/api"
	LoginURL = APIURL + "/auth/login"
	LogoutURL = APIURL + "/auth/logout"
	ChanMsgsURL = APIURL + "/channels/%s/messages" // chanID
	MsgURL = ChanMsgsURL + "/%s" //channel ID, message ID
	MsgAckURL = MsgURL + "/ack" //channel ID, message ID
)

func (c Discord) Send(method, url string, data, want interface{}) error {
	//note data and want are interfaces, so theoreticaly any object can be sent via JSON.
	//We don't have to worry about it, as an error will propogate if we send something unmarshalable
	if c.Token == "" && !c.LoggingIn { // not logged in or logging in
		return TokenError("Not logged in!")
	}
	var req *http.Request // important to define it before using
	var err error // define before use so we can use it later
	if data != nil {

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
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return UnknownError(fmt.Sprintf("%s", err))
		}
	}
	//err outside if would cause weird thing where err isn't populated
	if req == nil {
		return UnknownError("Why is req nil???") // should never happen
	}
	if c.Token != "" {
		req.Header.Add("Authorization", c.Token)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	resp, err := c.Client.Do(req) // at this point, bytes buffer is closed if needed
	if err != nil { //if theres an err, body couldbe nil
		return PostError(fmt.Sprintf("%s", err)) // body is nil here
	}
	defer resp.Body.Close() //after so panic wont occur for nil body
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return HTTPError(*resp)
	}
	
	var buff bytes.Buffer
	io.Copy(&buff, resp.Body)
	
	if want != nil {
		err = json.Unmarshal(buff.Bytes(), want)
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
