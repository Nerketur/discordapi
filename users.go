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
	return Channel{}, NotFoundError(name)
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
