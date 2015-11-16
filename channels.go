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
	fmt.Println(url)
	
	err := c.Delete(url)
	if err != nil {
		return err
	}
	
	fmt.Println("deleted chan perms!")
	return nil
}