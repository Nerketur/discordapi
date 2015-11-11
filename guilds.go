package discord

import (
	"fmt"
	"encoding/json"
)

type NotFoundError string
type guild []Guild

func (c NotFoundError) Error() string {
	return fmt.Sprintf("name not found: %s.  returning \"\"", string(c))
}

func (c Discord) Guild(name string) (string, error) {
	resp, err := guild(c.MyGuilds).Find(name)
	if err != nil {
		return "", err
	}
	return resp, nil
}
func (c guild) Find(name string) (string, error) {
	for _, ele := range c {
		if ele.Name == name {
			return ele.ID, nil
		}
	}
	return "", NotFoundError(name)
}

func (c Discord) GuildMembers(guildID string) ([]Member, error) {
	
	resp := make([]Member, 0)
	err := c.Get(fmt.Sprintf(GuildMembersURL, guildID), &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Got memberss successfully!")
	return resp, nil
}
func (c Discord) GuildChannels(guildID string) ([]Channel, error) {
	
	resp := make([]Channel, 0)
	err := c.Get(fmt.Sprintf(GuildChansURL, guildID), &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Got channels successfully!")
	return resp, nil
}
func (c Discord) Guilds() ([]Guild, error) {
	
	resp := make([]Guild, 0)
	err := c.Get(fmt.Sprintf(UserIDURL, "@me")+"/guilds", &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Got guilds successfully!")
	return resp, nil
}

type RolesMsg struct{
	roles   []Role
	msg     struct{
		Message *string `json:"message"`
	}
}
type BansMsg struct{
	bans    []User
	msg     struct{
		Message *string `json:"message"`
	}
}

func (s *RolesMsg) UnmarshalJSON(raw []byte) error {
	//first we try unmarshaling into roles
	err1 := json.Unmarshal(raw, &s.roles)
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
func (s *BansMsg) UnmarshalJSON(raw []byte) error {
	//first we try unmarshaling into roles
	err1 := json.Unmarshal(raw, &s.bans)
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

type MissingPermissionError string

func (e MissingPermissionError) Error() string {
	return fmt.Sprintln("permission error:", string(e))
}

func (c Discord) GuildRoles(guildID string) ([]Role, error) {
	
	resp := RolesMsg{}
	err := c.Get(fmt.Sprintf(GuildRolesURL, guildID), &resp)
	if err != nil {
		return resp.roles, err
	}
	if resp.msg.Message != nil {
		return resp.roles, MissingPermissionError(*resp.msg.Message)
	}
	
	fmt.Println("Got roles successfully!")
	return resp.roles, nil
}
func (c Discord) GuildBans(guildID string) ([]User, error) {
	
	resp := BansMsg{}
	err := c.Get(fmt.Sprintf(GuildBansURL, guildID), &resp)
	if err != nil {
		return resp.bans, err
	}
	if resp.msg.Message != nil {
		return resp.bans, MissingPermissionError(*resp.msg.Message)
	}
	
	fmt.Println("Got bans successfully!")
	return resp.bans, nil
}


/*
{
    "afk_timeout": 300,
    "joined_at": "2012-12-21T12:34:56.789012+00:00",
    "afk_channel_id": null,
    "id": "111222333444555666",
    "icon": null,
    "name": "Name",
    "roles": [
        {
            "managed": false,
            "name": "@everyone",
            "color": 0,
            "hoist": false,
            "position": -1,
            "id": "111222333444555666",
            "permissions": 12345678
        }
    ],
    "region": "us-west",
    "embed_channel_id": null,
    "embed_enabled": false,
    "owner_id": "111222333444555666"
}
*/