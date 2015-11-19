package discord

import (
	"fmt"
)

type _chan []Channel

func (c Discord) Chan(guild, name string) (Channel, error) {
	
	guild, err := c.GuildID(guild)
	if err != nil {
		return Channel{}, err
	}
	chans, err := c.GuildChannels(guild)
	if err != nil {
		return Channel{}, err
	}
	return _chan(chans).Find(name)
}
func (c Discord) ChanID(guild, name string) (string, error) {
	channel, err := c.Chan(guild, name)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return channel.ID, nil
}

func (c _chan) Find(name string) (Channel, error) {
	for _, ele := range c {
		if ele.Name == name {
			return ele, nil
		}
	}
	return Channel{}, NotFoundError(name)
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
	//fmt.Println(url)
	
	err := c.Delete(url)
	if err != nil {
		return err
	}
	
	fmt.Println("deleted chan perms!")
	return nil
}

type InvalidTypeError string

func (e InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type '%s'", string(e))
}

func (c Discord) ChanCreate(guildID, name, kind string) (Channel, error) {
	if kind != "text" && kind != "voice" {
		return Channel{}, InvalidTypeError(kind)
	}
	url := fmt.Sprintf(GuildChansURL, guildID)
	req := struct{
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		Name: name,
		Type: kind,
	}
	resp := Channel{}
	if err := c.Post(url, req, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("created channel!")
	return resp, nil
}

func (c Discord) ChanEdit(chanID, name string, topic *string) (Channel, error) {
	url := fmt.Sprintf(ChanIDURL, chanID)
	req := make(map[string]*string)
	if name != "" {req["name"] = &name} // name cannot be nil, empty ignore
	if *topic != "" {req["topic"] = topic} // empty ignore, nil remove
	resp := Channel{}
	if err := c.Patch(url, req, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("edited channel!")
	return resp, nil
}
func (c Discord) ChanDelete(chanID string) error {
	url := fmt.Sprintf(ChanIDURL, chanID)
	if err := c.Delete(url); err != nil {
		return err
	}
	
	fmt.Println("deleted channel!")
	return nil
}
//{"max_age":1800,"max_uses":0,"temporary":false,"xkcdpass":true}
//defaults: 1 day, 0, false, false
//TODO Disallow validate if anything else is there
func (c Discord) ChanInviteAllCreate(chanID string, age, uses *uint, temp, xkcd *bool, code *string) (Invite, error) {
	url := fmt.Sprintf(ChanInviteURL, chanID)
	req := struct{
		Age  *uint   `json:"max_age,omitempty"`
		Code *string `json:"validate,omitempty"`
		Temp *bool   `json:"temporary,omitempty"`
		Uses *uint    `json:"max_uses,omitempty"`
		XKCD *bool   `json:"xkcdpass,omitempty"`
	}{
		Age:  age,
		Code: code,
		Temp: temp,
		Uses: uses,
		XKCD: xkcd,
	}
	resp := Invite{}
	if err := c.Post(url, req, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("created (or validated) chan invite!")
	return resp, nil
}
func (c Discord) ChanInviteCreate(chanID string) (Invite, error) {
	return c.ChanInviteAllCreate(chanID, nil, nil, nil, nil, nil)
}
func (c Discord) ChanInviteValidate(chanID, code string) (Invite, error) {
	//validation, as of the time of this typing:
		//does not check the code
		//returns the first matching code (or the first code if none mtch)
		//guaranteed to return an active code (but can still be revoked)
		//code does not have to be created by us.
	return c.ChanInviteAllCreate(chanID, nil, nil, nil, nil, &code)
}
func (c Discord) ChanInvites(chanID string) ([]Invite, error) {
	url := fmt.Sprintf(ChanInviteURL, chanID)
	resp := make([]Invite,0)
	if err := c.Get(url, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("got chan invites!")
	return resp, nil
}
