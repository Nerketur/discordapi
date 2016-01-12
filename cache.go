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
	/* g := c.cache.Guilds[gIdx]
	guild.Channels = append(guild.Channels, ch)
	c.cache.Guilds[gIdx] = guild */
	//g.Channels = append(g.Channels, ch)
	//c.SetGuildCache(g)
	c.SetChanCache(ch)
}
func (c *Discord) RemChan(ch Channel) {
	//guild := c.cache.Guilds[gIdx]
	//dont update guild channel list.  thats ready only
	/* if idx == len(g.Channels)-1 {
		g.Channels = g.Channels[:idx]
	} else {
		g.Channels = append(g.Channels[:idx], g.Channels[idx+1:]...)
	} */
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

//private chan
func (c *_chan) AddChan(ch Channel) {
	*c = append(*c, ch)
}
func (ch *_chan) RemChanIdx(idx int) {
	c := *ch
	if idx == len(c)-1 {
		c = c[:idx]
	} else {
		c = append(c[:idx], c[idx+1:]...)
	}
	*ch = c
}

func (c *Discord) PrivateChannelParseWS(event string, ch Channel) {
	chans := _chan(c.cache.PrivateChannels)
	cIdx, err := chans.FindIdxID(ch.ID)
	if err != nil {
		fmt.Println(err)
	}
	if event != "CHANNEL_CREATE" && err == nil {
		chans.RemChanIdx(cIdx)
	}
	if event != "CHANNEL_DELETE" {
		chans.AddChan(ch)
	}
}

//guild

func (c *Discord) AddGuild(g Guild) {
	c.cache.Guilds = append(c.cache.Guilds, g)
}
func (c *Discord) RemGuildIdx(idx int) {
	if idx == len(c.cache.Guilds)-1 {
		c.cache.Guilds = c.cache.Guilds[:idx]
	} else {
		c.cache.Guilds = append(c.cache.Guilds[:idx], c.cache.Guilds[idx+1:]...)
	}
}

func (c *Discord) GuildParseWS(event string, g Guild) {
	if g.Unavailable != nil {
		return // ignore these messages for now
	}
	oldIdx, err := guilds(c.cache.Guilds).FindIdxID(g.ID)
	if err != nil {
		fmt.Println(err)
	}
	if event != "GUILD_CREATE" && err == nil {
		c.RemGuildIdx(oldIdx)
	}
	if event != "GUILD_DELETE" {
		c.AddGuild(g)
	}
}

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
	idx, err := guilds(c.cache.Guilds).FindIdxID(m.GuildID)
	g := c.cache.Guilds[idx]
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
	c.cache.Guilds[idx] = g
}

//guild bans

func (c *Discord) GuildBanParseWS(event string, b WSBan) {
	//for now do nothing
}

//guld roles
func (c *Discord) AddGuildRole(gIdx int, r Role) {
	guild := c.cache.Guilds[gIdx]
	guild.Roles = append(guild.Roles, r)
	c.cache.Guilds[gIdx] = guild
}
func (c *Discord) RemGuildRoleIdx(gIdx, idx int) {
	guild := c.cache.Guilds[gIdx]
	if idx == len(guild.Roles)-1 {
		guild.Roles = guild.Roles[:idx]
	} else {
		guild.Roles = append(guild.Roles[:idx], guild.Roles[idx+1:]...)
	}
	c.cache.Guilds[gIdx] = guild
}
func (c *Discord) GuildRoleParseWS(event string, r WSRole) {
	//update guild roles
	//delete only gets role ID
	gIdx, err := guilds(c.cache.Guilds).FindIdxID(r.GetGuildID())
	if err != nil {
		fmt.Println("invalid guild:", err)
		return
	}
	g := c.cache.Guilds[gIdx]
	rIdx, err := roles(g.Roles).FindIdxID(r.GetRoleID())
	if err != nil {
		fmt.Println("invalid guild:", err)
		return
	}
	if event != "GUILD_ROLE_CREATE" && err == nil {
		c.RemGuildRoleIdx(gIdx, rIdx)
	}
	if event != "GUILD_ROLE_DELETE" {
		c.AddGuildRole(gIdx, r.GetRole())
	}
	c.cache.Guilds[gIdx] = g
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
}
