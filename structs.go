package discord

import (
	"net/http"
	"time"
)

type Creds struct{
	Email    string `json:"email"`
	Pass     string `json:"password"`
}
type CredsResp struct{
	Email    []string `json:"email"`
	Pass     []string `json:"password"`
	Token    string `json:"token,omitempty"`
}

type Discord struct{
	Client    *http.Client
	Token     string
	LoggingIn bool
}

/////////////////////////////

type MessageSend struct{
	Content  string   `json:"content"`
	Mentions []string `json:"mentions"`
	Nonce    string   `json:"nonce"`
	Tts      bool     `json:"tts"`
}
type User struct{
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"` //4 digits
		ID            string `json:"id"`
		Avatar        string `json:"avatar"` // hex string (can be null?)
}
/*
{
	"username": "",
	"discriminator": "0000",
	"id": <authid>,
	"avatar": <hex> // 32 hex digits, 128-bit
},
*/
type Embed struct{
	Desc     string    `json:"description"`
	Author   Entity    `json:"author"`
	URL      string    `json:"url"`
	Title    string    `json:"title"`
	Provider Entity    `json:"provider"`
	Type     string    `json:"type"`
	Thumb    Thumbnail `json:"thumbnail"`
}
type Entity struct{
	URL  *string `json:"url"` //can be null
	Name string `json:"name"`
}                
type Thumbnail struct{
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Proxy  string `json:"proxy_url"`
	Height int    `json:"height"`
}
/*               
"embeds": [
	{
		"description": "",
		"author": {"url": "", "name": ""},
		"url": "",
		"title": "",
		"provider": {"url": null, "name": "GitHub"},
		"type": "link",
		"thumbnail": {
			"url": <url>,
			"width": 96,
			"proxy_url": <url>,
			"height": 96
		}
	}
],
*/
type Attachment struct{
	URL      string `json:"url"`      //URL of downloadable object
	ProxyURL string `json:"poxy_url"` //URL of ?
	Size     int    `json:"size"`     //size in bytes
	ID       string `json:"id"`       //id of attachment
	Filename string `json:"filename"` //filename
}

/*
"attachments": [
	{
		"url": <url>,
		"proxy_url": <url>,
		"size": 82, // bytes
		"id": <id>,
		"filename": "hello.go"
	}
],
*/

type Message struct{
	Nonce       string         `json:"nonce,omitempty"`
	Attachments []Attachment   `json:"attachments"`
	Tts         bool           `json:"tts"` 
	Embeds      []Embed        `json:"embeds"`
	Timestamp   time.Time      `json:"timestamp"`
	MentionAll  bool           `json:"mention_everyone"`
	ID          string         `json:"id"`
	EditedTime  *time.Time     `json:"edited_timestamp"` //can be null (not worth seperate struct member)
	Author      User           `json:"author"`
	Content     string         `json:"content"`
	ChanID      string         `json:"channel_id"`
	//including json tag magic to have it look for both, and ignore whichever one doesn't exist.
	MentionsNum []string       `json:"mentions,omitempty"` // Userids (usually only sent)
	MentionsUse []User         `json:"mentions,omitempty"` // Userids (usually only receved)
}

/*
[
	{
		"attachments": [],
		"tts": false,
		"embeds": [],
		"timestamp": "2015-11-03T20:07:16.292000+00:00",
		"mention_everyone": false,
		"id": <msgid>,
		"edited_timestamp": null,
		"author": {
			"username": "",
			"discriminator": "0000",
			"id": <authid>,
			"avatar": <hex> // 32 hex digits, 128-bit
		},
		"content": <url>,
		"channel_id": <chanid>,
		"mentions": []
	},
*/
