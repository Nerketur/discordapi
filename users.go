package discord

import (
	"fmt"
)

//type _cache map[string]interface{}

func (c Discord) PrivateChan(name string) (ret Channel, err error) {
	var chs []Channel
	chs, err = c.FindNamePrivChanCache(name)
	if err == nil {
		ret = chs[0]
	}
	return 
}
func (c Discord) PrivateChanFromID(ID string) (Channel, error) {
	return c.PrivChanCache(ID)
}

func (c Discord) PrivateChanID(name string) (string, error) {
	resp, err := c.PrivateChan(name)
	return resp.ID, err
}

/*func (c _privchan) Find(name string) (ele Channel, err error) {
	for _, ele = range c {
		if ele.Recipient.Username == name {
			return
		}
	}
	err = NameNotFoundError(name)
	return
}
func (c _privchan) FindID(ID string) (ele Channel, err error) {
	for _, ele = range c {
		if ele.Recipient.ID == ID {
			return
		}
	}
	err = IDNotFoundError(ID)
	return
}*/

func (c Discord) PrivateChansRest() (resp []Channel, err error) {
	resp = make([]Channel, 0)
	if err = c.Get(MyChansURL, &resp); err != nil {
		return
	}
	fmt.Println("got private channels successfully!")
	return
}

/*
 * Pehaps not 100% POST compliant, but creates channel if it doesn't already exist
 * on a different URL
 */
func (c *Discord) CreatePrivateChan(userID string) (resp Channel, err error) {
	req := struct{UserID string `json:"recipient_id"`}{
		UserID: userID,
	}
	//TODO: test other userIDs and see if it says forbidden
	err = c.Post(MyChansURL, req, &resp) // can also use our userID instead of @me
	if err != nil {
		return
	}
	
	//pcs := _chan(c.PrivChanCache)
	//pcs.AddChan(resp) //only needed if no websocket
	//fmt.Printf("%#v\n", pcs)
	fmt.Println("created (opened) private channel successfully!")
	return
}

/* func (c *Discord) RemPrivChanIdx(idx int) {
	if idx == 0 {
		c.MyChans = c.MyChans[1:]
	} else if idx == len(c.MyChans)-1 {
		c.MyChans = c.MyChans[:idx-1]
	} else {
		c.MyChans = append(c.MyChans[:idx-1], c.MyChans[idx+1:]...)
	}
} */

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
	if err = c.PrivateChanDelete(channel.ID); err != nil {
		return err
	}
	fmt.Println("removed private channel (from list) successfully!")
	return nil
}

func (c Discord) UserPres(guildID, userID string) (p WSPres, err error) {
/* 	gIdx, err := guilds(c.cache.Guilds).FindIdxID(guildID)
	if err != nil {
		return p, err
	} */
	g, err := c.GuildCache(guildID)
	if err != nil {
		return p, err
	}
	pIdx, err := _pres(g.Presences).FindIdx(userID)
	if err != nil {
		return p, err
	}
	fmt.Printf("presence: %#v\n", g.Presences[pIdx])
	return g.Presences[pIdx], err
}

func (p WSPres) Playing() *string {
	if p.Game == nil {
		return nil
	}
	tmp := *p.Game
	fmt.Println("game:", tmp)
	return &tmp.Name
}
