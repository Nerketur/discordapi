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
type TokenStr struct{
	Token    string `json:"token"`
}

type Discord struct{
	Client    *http.Client
	Token     string
	LoggingIn bool
}

/////////////////////////////

type Message struct{
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
type Embed struct{}
type Attachment struct{
	URL      string `json:"url"`      //URL of downloadable object
	ProxyURL string `json:"poxy_url"` //URL of ?
	Size     int    `json:"size"`     //size in bytes
	ID       string `json:"id"`       //id of attachment
	Filename string `json:"filename"` //filename
}

type MessageResp struct{
	Nonce       string         `json:"nonce"`
	Attachments []Attachment   `json:"attachments"`
	Tts         bool           `json:"tts"` 
	Embeds      []Embed        `json:"embeds"`
	Timestamp   time.Time      `json:"timestamp"`
	MentionAll  bool           `json:"mention_everyone"`
	ID          string         `json:"id"`
	EditedTime  *time.Time     `json:"edited_timestamp"` //can be null
	Author      User           `json:"author"`
	Content     string         `json:"content"`
	ChanID      string         `json:"channel_id"`
	Mentions    []string       `json:"mentions"` // Userids
}

/*
{
	"attachments": [
		{
			"url": <url>,
			"proxy_url": <url>,
			"size": 82,
			"id": <id>,
			"filename": "hello.go"
		}
	],
	"tts": false,
	"embeds": [],
	"timestamp": "2015-11-02T22:44:36.580000+00:00",
	"mention_everyone": false,
	"id": <id>,
	"edited_timestamp": null,
	"author": {"username": "Nerketur", "discriminator": "1468", "id": "94473980307570688", "avatar": "fdfc0122012c75b3a164433856050b18"},
	"content": "",
	"channel_id": <id>,
	"mentions": []
}
*/
