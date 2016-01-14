package discord

import (
	"fmt"
)

type _chan []Channel

func (c Discord) Chan(guildID, name string) (Channel, error) {
	chans, err := c.GuildChannels(guildID)
	if err != nil {
		return Channel{}, err
	}
	return _chan(chans).Find(name)
}
func (c Discord) ChanID(guildID, name string) (string, error) {
	channel, err := c.Chan(guildID, name)
	
	return channel.ID, err
}

func (c _chan) Find(name string) (ch Channel, err error) {
	idx, err := c.FindNameIdx(name)
	if err == nil {
		ch = c[idx]
	}
	return
}
func (c _chan) FindID(ID string) (ch Channel, err error) {
	idx, err := c.FindIdxID(ID)
	if err == nil {
		ch = c[idx]
	}
	return
}
func (c _chan) FindNameIdx(name string) (int, error) {
	for idx, ele := range c {
		if ele.Name == name {
			return idx, nil
		}
	}
	return -1, NameNotFoundError("name: " + name)
}
func (c _chan) FindIdxID(ID string) (int, error) {
	for idx, ele := range c {
		if ele.ID == ID {
			return idx, nil
		}
	}
	return -1, IDNotFoundError("id: " + ID)
}

func (c Discord) SendTyping(chanID string) (err error) {
	err = c.Post(fmt.Sprintf(ChanTypingURL, chanID), nil, nil)
	if err != nil {
		return
	}
	fmt.Println("typing sent!")
	return
}

type userOrRole interface{
	GetID() string
	Type() string
}

func (x User) GetID() string {return x.ID}
func (x User) Type() string {return "member"}
func (x Role) GetID() string {return x.ID}
func (x Role) Type() string {return "role"}

func (c Discord) ChanReplacePerms(chanID string, ur userOrRole, allow, deny Perm) (err error) {
	url := fmt.Sprintf(ChanPermIDURL, chanID, ur.GetID())
	fmt.Println(url)
	
	req := PermOver{
		Allow: allow,
		Deny: deny,
		ID: ur.GetID(),
		Type: ur.Type(),
	}
	
	if err = c.Put(url, req); err != nil {
		return
	}
	
	fmt.Println("replaced chan perms!")
	return
}
func (c Discord) ChanDeletePerms(chanID, permID string) (err error) {
	url := fmt.Sprintf(ChanPermIDURL, chanID, permID)

	if err = c.Delete(url); err != nil {
		return
	}
	
	fmt.Println("deleted chan perms!")
	return
}

func (c Discord) ChanEdit(chanID, name string, topic *string) (resp Channel, err error) {
	url := fmt.Sprintf(ChanIDURL, chanID)
	req := make(map[string]*string)
	if name != "" {req["name"] = &name} // name cannot be nil, empty ignore
	if *topic != "" {req["topic"] = topic} // empty ignore, nil remove
	if err = c.Patch(url, req, &resp); err != nil {
		return
	}
	
	fmt.Println("edited channel!")
	return
}
func (c Discord) ChanDelete(chanID string) (err error) {
	url := fmt.Sprintf(ChanIDURL, chanID)
	if err = c.Delete(url); err != nil {
		return
	}
	fmt.Println("deleted channel!")
	return
}
func (c *Discord) PrivateChanDelete(chanID string) (err error) {
	return c.ChanDelete(chanID)
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
		Uses *uint   `json:"max_uses,omitempty"`
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
func (c Discord) ChanInvites(chanID string) (resp []Invite, err error) {
	url := fmt.Sprintf(ChanInviteURL, chanID)
	resp = make([]Invite,0)
	if err = c.Get(url, &resp); err != nil {
		return
	}
	
	fmt.Println("got chan invites!")
	return
}