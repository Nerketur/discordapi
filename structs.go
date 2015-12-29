package discord

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"
)

type Creds struct{
	Email    []string `json:"email"`
	Pass     []string `json:"password"`
	Token    string `json:"token,omitempty"`
}

type Discord struct{
	Client    *http.Client
	Token     string
	LoggingIn bool
	MyGuilds  []Guild
	MyChans   []Channel
	sigStop   chan int
	sigSafe   chan int
	sigTime   chan int
}

/////////////////////////////

//Only reason not to remove this is its used evey time we create a message.
//Theoretically, we could use Message instead, and have the other fields of
//the struct set to omit from the json if empty, but that would just make
//everything annoying to read, so this stays seperate for now.
type MessageSend struct{ 
	Content  string   `json:"content"`
	Mentions []string `json:"mentions"`
	Nonce    int64    `json:"nonce,string"`
	Tts      bool     `json:"tts"`
}

type Member struct{
	GuildID string    `json:"guild_id,omitempty"`
	Joined  time.Time `json:"joined_at"`
	Deaf    bool      `json:"deaf"`
	User    User      `json:"user"`
	Roles   []string  `json:"roles"`
	Mute    bool      `json:"mute"`
}
	type User struct{
		Verified      bool    `json:"verified,omitempty"` //only for WS
		Username      string  `json:"username"`
		Email         string  `json:"email,omitempty"`    //only for WS
		Discriminator string  `json:"-"` //4 digits
		ID            string  `json:"id"`
		Avatar        *string `json:"avatar"` // hex string (can be null)
	}

func (u *User) UnmarshalJSON(raw []byte) (err error) {
	type user User
	u1, discInt, discStr := user{}, struct{D int `json:"discriminator"`}{}, struct{D string `json:"discriminator"`}{}
	err = json.Unmarshal(raw, &u1)
	if err != nil {
		return
	}
	err = json.Unmarshal(raw, &discInt)
	if err != nil {
		err = json.Unmarshal(raw, &discStr)
		if err != nil {
			return
		} else {
			u1.Discriminator = fmt.Sprintf("%v", discStr)
		}
	} else {
		u1.Discriminator = fmt.Sprintf("%v", discInt)
	}
	
	*u = User(u1)
	return
}

type Message struct{
	Nonce       int64        `json:"nonce,string,omitempty"` //only used when sending messages
	Attachments []Attachment `json:"attachments"`
	Tts         bool         `json:"tts"` 
	Embeds      []Embed      `json:"embeds"`
	Timestamp   time.Time    `json:"timestamp"`
	MentionAll  bool         `json:"mention_everyone"`
	ID          string       `json:"id"`
	EditedTime  *time.Time   `json:"edited_timestamp"` //can be null (not worth seperate struct member)
	Author      User         `json:"author"`
	Content     string       `json:"content"`
	ChanID      string       `json:"channel_id"`
	//including json tag magic to have it look for both, and ignore whichever one doesn't exist.
	MentionsNum []string     `json:"mentions,omitempty"` // Userids (usually only sent)
	MentionsUse []User       `json:"mentions,omitempty"` // Users (usually only receved)
}
	type Attachment struct{
		URL      string `json:"url"`      //URL of downloadable object
		ProxyURL string `json:"proxy_url"` //URL of ?
		Size     int    `json:"size"`     //size in bytes
		ID       string `json:"id"`       //id of attachment
		Filename string `json:"filename"` //filename
		Width    int    `json:"width,omitempty"`    //image width
		Height   int    `json:"height,omitempty"`   //image height
	}
	type Embed struct{
		Desc     *string   `json:"description"`
		Author   *Entity    `json:"author,omitempty"`
		URL      string    `json:"url"`
		Title    *string    `json:"title"`
		Provider *Entity   `json:"provider,omitempty"`
		Type     string    `json:"type"`
		Thumb    *Thumbnail `json:"thumbnail,omitempty"`
		Video    *Video     `json:"video,omitempty"`
	}
		type Entity struct{
			URL  *string `json:"url"` //can be null
			Name string  `json:"name"`
		}                
		type Thumbnail struct{
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Proxy  string `json:"proxy_url"`
			Height int    `json:"height"`
		}
		type Video struct{
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		}

type PrivateChannel struct{ // not curently used
	LastMsgID   *string    `json:"last_message_id"` // can be null
	Recipient   *User      `json:"recipient,omitempty"` //only exists if private
	ID          string     `json:"id"`
	Private     bool       `json:"is_private"`
}
type Channel struct{
	GuildID     string     `json:"guild_id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Permissions []PermOver `json:"permission_overwrites,omitempty"`
	Topic       *string    `json:"topic,omitempty"` // can be null
	Position    int        `json:"position,omitempty"`
	LastMsgID   *string    `json:"last_message_id"` // can be null
	Recipient   *User      `json:"recipient,omitempty"` //only exists if private
	Type        string     `json:"type,omitempty"` //only exists if not private (text|voice)
	ID          string     `json:"id"`
	Private     bool       `json:"is_private"`
}
	type PermOver struct{
		Deny  Perm   `json:"deny"`
		Type  string `json:"type"`
		ID    string `json:"id"`
		Allow Perm   `json:"allow"`
	}

type Guild struct{
	VoiceStates  []WSVoiceState `json:"voice_states"` //only READY
	Roles        []Role         `json:"roles"`
	Region       string         `json:"region"`
	Presences    []WSPres       `json:"presences"` // only READY
	OwnerID      string         `json:"owner_id"`
	Name         string         `json:"name"`
	//Large        bool           `json:"large"` //only READY
	Members      []Member       `json:"members"` //only READY
	JoinedAt     time.Time      `json:"joined_at"`
	ID           string         `json:"id"`
	Icon         *string        `json:"icon"`
	Channels     []Channel      `json:"channels"` // only READY
	AfkTimeout   uint64         `json:"afk_timeout"`
	AfkChannelID *string        `json:"afk_channel_id"`
	Unavailable  *bool          `json:"unavailable,omitempty"` // if present, state changed
}
	type Role struct{
		Managed     bool   `json:"managed,omitempty"`
		Name        string `json:"name"`
		Color       int    `json:"color"`
		Hoist       bool   `json:"hoist"`
		Position    int    `json:"position,omitempty"`
		ID          string `json:"id,omitempty"`
		Permissions Perm   `json:"permissions"`
	}
type Invite struct{
    MaxAge    int       `json:"max_age"`
    Code      string    `json:"code"`
    Guild     struct{
        ID        string    `json:"id"`
        Name      string    `json:"name"`
    }                   `json:"guild"`
    Revoked   bool      `json:"revoked"`
    CreatedAt time.Time `json:"created_at"`
    Temporary bool      `json:"temporary"`
    Uses      int       `json:"uses"`
    MaxUses   int       `json:"max_uses"`
    Inviter   User      `json:"inviter"`
    XKCDpass  *string   `json:"xkcdpass"`
    Chan      Channel   `json:"channel"`
}
type Region struct{
	Hostname string `json:"sample_hostname"`
	Port     uint   `json:"sample_port"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}
type Connection struct{
	Integrations []Integration `json:"integrations"`
	Revoked      bool          `json:"revoked"`
	Type         string        `json:"type"`
	ID           string        `json:"id"`
	Name         string        `json:"name"`
}
	type Integration string
type Ice struct {
	TTL     string   `json:"ttl"`
	Servers []Server `json:"servers"`
}
	type Server struct {
		URL      string `json:"url"`
		Username string `json:"username,omitempty"`
		Cred     string `json:"credential,omitempty"`
	}
type Setting struct{
	RenderEmbeds      *bool    `json:"render_embeds,omitempty"`
	InlineEmbedMedia  *bool    `json:"inline_embed_media,omitempty"`
	EnableTTSCmd      *bool    `json:"enable_tts_command,omitempty"`
	MsgDispCompact    *bool    `json:"message_display_compact,omitempty"`
	Locale            string   `json:"locale,omitempty"`
	ShowCurrentGame   *bool    `json:"show_current_game,omitempty"`
	Theme             string   `json:"theme,omitempty"` // emptystring should not be sent
	MutedChanIDs      []string `json:"muted_channels,omitempty"`
	InlineAttachMedia *bool    `json:"inline_attachment_media,omitempty"`
	
}
type Tutorial struct{
	Confirmed  []string `json:"indicators_confirmed"`
	Suppressed bool     `json:"indicators_suppressed"`
}
type EmbedInfo struct{
	ChanID  *string `json:"channel_id"`
	Enabled bool    `json:"enabled"`
}
