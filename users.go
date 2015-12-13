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
func (c Discord) CreatePrivateChan(userID string) (Channel, error) {
	
	resp := Channel{}
	req := struct{UserID string `json:"recipient_id"`}{
		UserID: userID,
	}
	//TODO: test other userIDs and see if it says forbidden
	err := c.Post(MyChansURL, req, &resp) // can also use our userID instead of @me
	
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	c.MyChans = append(c.MyChans, resp)
	fmt.Println("created (opened) private channel successfully!")
	return resp, nil
}

/*
 * Note that this does not remove message history; only removes the channel from
 * the list of private chats you have open.
 * The current implementation is to create it then delete it using DeleteChannel().
 * This will change in the future.  Do NOT rely on this behavior.
 */
func (c Discord) DeletePrivateChan(userID string) error { //API also returns deleted channel
	
	//TODO: when better caching is implememnted, change this to use it.
	channel, err := c.CreatePrivateChan(userID)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	
	
	if err = c.ChanDelete(channel.ID); err != nil { //..but we ignore it.
		//fmt.Println(err)
		return err
	}
	
	fmt.Println("removed private channel (from list) successfully!")
	return nil
}
