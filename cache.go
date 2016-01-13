package discord

import "fmt"

/* type DiscordCache interface{
	Ele(Discord, string, interface{}) error
	SetEle(Discord, interface{}) error
	FindKey(Discord, string) (string, error)
}

type UserCache Discord

func (c UserCache) Ele(userID string, val interface{}) (err error) {
	if user, ok := val.(*User); !ok {
		return NotPointerErr("element must be type *User")
	}
	*user = c.usrCache[userID]
	return
}
func (u UserCache) SetEle(val interface{}) (err error) {
	if id, ok := val.(User); !ok {
		return PointerErr("element must be type User")
	}
	c.usrCache[user.ID] = user
	return
}
func (u UserCache) FindKey(name string) (ret string, err error) {
	err = NameNotFoundError(name)
	for key, val := range c.usrCache {
		if val.Name == name {
			ret, err = key, nil
		}
	}
	return
} */
func (c Discord) GuildCacheGuilds() (ret []Guild) {
	for _, val := range c.gldCache {
		ret = append(ret, val)
	}
	return
}

func (c Discord) UserCache(userID string) (ret User, err error) {
	var ok bool
	if ret, ok = c.usrCache[userID]; !ok {
		err = IDNotFoundError(userID)
	}
	return
}
func (c Discord) MessageCache(msgID string) (ret Message, err error) {
	var ok bool
	if ret, ok = c.msgCache[msgID]; !ok {
		err = IDNotFoundError(msgID)
	}
	return
}
func (c Discord) GuildCache(guildID string) (ret Guild, err error) {
	var ok bool
	if ret, ok = c.gldCache[guildID]; !ok {
		err = IDNotFoundError(guildID)
	}
	return
}
func (c Discord) ChanCache(chanID string) (ret Channel, err error) {
	var ok bool
	if ret, ok = c.chnCache[chanID]; !ok {
		err = IDNotFoundError(chanID)
	}
	return
}
func (c Discord) PrivChanCache(userID string) (ret Channel, err error) {
	var ok bool
	if ret, ok = c.priCache[userID]; !ok {
		err = IDNotFoundError(userID)
	}
	return
}

func (c Discord) SetUserCache(u User) {
	c.usrCache[u.ID] = u
}
func (c Discord) SetMessageCache(m Message) {
	c.msgCache[m.ID] = m
}
func (c Discord) SetGuildCache(g Guild) {
	c.gldCache[g.ID] = g
}
func (c Discord) SetChanCache(ch Channel) {
	c.chnCache[ch.ID] = ch
}
func (c Discord) SetPrivChanCache(ch Channel) {
	c.priCache[ch.Recipient.ID] = ch
}

func (c Discord) DelUserCache(u User) {
	delete(c.usrCache, u.ID)
}
func (c Discord) DelMessageCache(m Message) {
	delete(c.msgCache, m.ID)
}
func (c Discord) DelGuildCache(g Guild) {
	delete(c.gldCache, g.ID)
}
func (c Discord) DelChanCache(ch Channel) {
	delete(c.chnCache, ch.ID)
}
func (c Discord) DelPrivChanCache(ch Channel) {
	delete(c.priCache, ch.Recipient.ID)
}

func (c Discord) FindNameUserCache(name string) (ret User, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.usrCache {
		if val.Username == name {
			ret, err = val, nil
		}
	}
	return
}
/* func (c Discord) FindNameMessageCache(name string) (ret Message, err error) {
	err = NameNotFoundError(name)
	for key, val := range c.msgCache {
		if val.Name == name {
			ret, err = val, nil
		}
	}
	return
} */
func (c Discord) FindNameGuildCache(name string) (ret Guild, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.gldCache {
		if val.Name == name {
			ret, err = val, nil
		}
	}
	return
}
func (c Discord) FindNameChanCache(name string) (ret Channel, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.chnCache {
		if val.Name == name {
			ret, err = val, nil
		}
	}
	return
}
func (c Discord) FindNamePrivChanCache(name string) (ret Channel, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.priCache {
		if val.Name == name {
			ret, err = val, nil
		}
	}
	return
}

//channels
func (c *Discord) AddChan(ch Channel) {
	//dont update guild channel list.  thats ready only
	c.SetChanCache(ch)
}
func (c *Discord) RemChan(ch Channel) {
	//dont update guild channel list.  thats ready only
	c.DelChanCache(ch)
}

func (c *Discord) ChannelParseWS(event string, ch Channel) {
	switch event {
	case "CHANNEL_CREATE","CHANNEL_UPDATE":
		c.AddChan(ch)
	case "CHANNEL_DELETE":
		c.RemChan(ch)
	}
}
//private channels
func (c *Discord) AddPrivChan(ch Channel) {
	c.SetPrivChanCache(ch)
}
func (c *Discord) RemPrivChan(ch Channel) {
	c.DelPrivChanCache(ch)
}

func (c *Discord) PrivateChannelParseWS(event string, ch Channel) {
	switch event {
	case "CHANNEL_CREATE","CHANNEL_UPDATE":
		c.AddPrivChan(ch)
	case "CHANNEL_DELETE":
		c.RemPrivChan(ch)
	}
}

//guilds
func (c *Discord) AddGuild(g Guild) {
	c.SetGuildCache(g)
}
func (c *Discord) RemGuild(g Guild) {
	c.DelGuildCache(g)
}

func (c *Discord) GuildParseWS(event string, g Guild) {
	switch event {
	case "GUILD_CREATE","GUILD_UPDATE":
		c.AddGuild(g)
	case "GUILD_DELETE":
		c.RemGuild(g)
	}
}
//////////////////////////////////////////////////

//TODO: Figure out a way to have this work for members
/* //guild members
func (c *Discord) AddMember(m Member) {
	c.SetMemberCache(m)
}
func (c *Discord) RemMember(m Member) {
	c.DelMemberCache(m)
}

func (c *Discord) GuildMemberParseWS(event string, m Member) {
	switch event {
	case "GUILD_MEMBER_ADD","GUILD_MEMBER_UPDATE":
		c.AddMember(m)
	case "GUILD_MEMBER_REMOVE":
		c.RemMember(m)
	}
} */
//////////////////////////////////////////////////

//guild members

func (g *Guild) AddMember(m Member) {
	g.Members = append(g.Members, m)
}
func (g *Guild) RemMemberIdx(idx int) {
	if idx == len(g.Members)-1 {
		g.Members = g.Members[:idx]
	} else {
		g.Members = append(g.Members[:idx], g.Members[idx+1:]...)
	}
}
func (g *Guild) RemMember(m Member) {
	idx, err := members(g.Members).FindIdxID(m.User.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	g.RemMemberIdx(idx)
}

func (c *Discord) GuildMemberParseWS(event string, m Member) {
	g, err := c.GuildCache(m.GuildID)
	if err != nil {
		fmt.Println("cache error:", err)
		return
	}
	if event != "GUILD_MEMBER_ADD" {
		g.RemMember(m)
	}
	if event != "GUILD_MEMBER_REMOVE" {
		g.AddMember(m)
	}
	c.SetGuildCache(g)
}

//guild bans

func (c *Discord) GuildBanParseWS(event string, b WSBan) {
	//for now do nothing
}

//guild roles
func (g *Guild) AddRole(r Role) {
	g.Roles = append(g.Roles, r)
}
func (g *Guild) RemRole(r Role) {
	idx, err := roles(g.Roles).FindIdxID(r.ID)
	if err != nil {
		fmt.Println("remove role warning:", err)
		return
	}

	if idx == len(g.Roles)-1 {
		g.Roles = g.Roles[:idx]
	} else {
		g.Roles = append(g.Roles[:idx], g.Roles[idx+1:]...)
	}
}
func (c *Discord) GuildRoleParseWS(event string, r WSRole) {
	//update guild roles
	//delete only gets role ID
	g, err := c.GuildCache(r.GetGuildID())
	if err != nil {
		fmt.Println("invalid guild:", err)
		return
	}
	if event != "GUILD_ROLE_CREATE" {
		g.RemRole(r.GetRole())
	}
	if event != "GUILD_ROLE_DELETE" {
		g.AddRole(r.GetRole())
	}
	c.SetGuildCache(g)
}

//WS

func (c *Discord) wsFillCaches(ws READY) {
	for _, channel := range ws.PrivateChannels {
		c.priCache[channel.Recipient.ID] = channel
	}
	for _, guild := range ws.Guilds {
		c.gldCache[guild.ID] = guild
		for _, member := range guild.Members {
			//TODO: member-guild link
			c.usrCache[member.User.ID] = member.User
		}
		for _, channel := range guild.Channels {
			channel.GuildID = guild.ID
			c.chnCache[channel.ID] = channel
			//fmt.Println("channel:", channel, ", guildID:", channel.GuildID)
		}
		if debug {
			for _, pres := range guild.Presences {
				if pres.Game != nil {
					fmt.Println("game:", pres.Game.Name)
					fmt.Println("userID:", pres.User.ID)
					fmt.Println("username:", pres.User.Username)
				}
			}
		}
	}
	fmt.Println("Caches filled!")
}
