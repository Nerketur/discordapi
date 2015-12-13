package discord

import (
	"fmt"
)

type _privchan []Channel

func (c Discord) PrivateChan(name string) Channel {
	resp, err := _privchan(c.MyChans).Find(name)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
func (c Discord) PrivateChanFromID(ID string) Channel {
	resp, err := _privchan(c.MyChans).FindID(ID)
	if err != nil {
		fmt.Println(err)
	}
	return resp
}
func (c Discord) PrivateChanID(name string) string {
	return c.PrivateChan(name).ID
}

func (c _privchan) Find(name string) (Channel, error) {
	for _, ele := range c {
		if ele.Recipient.Username == name {
			return ele, nil
		}
	}
	return Channel{}, NameNotFoundError(name)
}
func (c _privchan) FindID(ID string) (Channel, error) {
	for _, ele := range c {
		if ele.Recipient.ID == ID {
			return ele, nil
		}
	}
	return Channel{}, IDNotFoundError(ID)
}

func (c Discord) GetMyPrivateChans() ([]Channel, error) {
	
	resp := make([]Channel, 0)
	err := c.Get(MyChansURL, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("got private channels successfully!")
	return resp, nil
}

/*
 * Pehaps not 100% POST compliant, but creates channel if it doesn't already exist
 * on a different URL
 */
func (c Discord) CreatePrivateChan(userID string) (resp Channel, err error) {
	
	req := struct{UserID string `json:"recipient_id"`}{
		UserID: userID,
	}
	//TODO: test other userIDs and see if it says forbidden
	err = c.Post(MyChansURL, req, &resp) // can also use our userID instead of @me
	
	if err != nil {
		return
	}
	c.MyChans = append(c.MyChans, resp)
	fmt.Println("created (opened) private channel successfully!")
	return
}

/*
 * Note that this does not remove message history; only removes the channel from
 * the list of private chats you have open.
 * The current implementation is to find it in our cache, then delete it using DeleteChannel().
 * This could change in the future.
 */
func (c Discord) DeletePrivateChan(userID string) error { //API also returns deleted channel
	//use cache of private convos
	channel, err := c.PrivateChanFromID(userID)
	if err != nil {
		return err
	}
	
	if err = c.ChanDelete(channel.ID); err != nil { //..but we ignore it.
		return err
	}
	
	fmt.Println("removed private channel (from list) successfully!")
	return nil
}
