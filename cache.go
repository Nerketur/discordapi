package discord

import "fmt"

func (c Discord) UserCache(userID string) User {
	return c.usrCache[userID]
}
func (c Discord) MessageCache(msgID string) Message {
	return c.msgCache[msgID]
}
func (c Discord) GuildCache(guildID string) Guild {
	return c.gldCache[guildID]
}
func (c Discord) ChanCache(chanID string) Channel {
	return c.chnCache[chanID]
}

//channels
func (c *Discord) AddChan(gIdx int, ch Channel) {
	guild := c.cache.Guilds[gIdx]
	guild.Channels = append(guild.Channels, ch)
	c.cache.Guilds[gIdx] = guild
}
func (c *Discord) RemChanIdx(gIdx, idx int) {
	guild := c.cache.Guilds[gIdx]
	if idx == len(guild.Channels)-1 {
		guild.Channels = guild.Channels[:idx]
	} else {
		guild.Channels = append(guild.Channels[:idx], guild.Channels[idx+1:]...)
	}
	c.cache.Guilds[gIdx] = guild
}


func (c *Discord) ChannelParseWS(event string, ch Channel) {
	gIdx, err := guilds(c.cache.Guilds).FindIdxID(ch.GuildID)
	if err != nil {
		fmt.Println("chan parse err:", err)
		return
	}
	cIdx, err := _chan(c.cache.Guilds[gIdx].Channels).FindIdxID(ch.ID)
	if err != nil {
		fmt.Println(err)
	}
	if event != "CHANNEL_CREATE" && err == nil {
		c.RemChanIdx(gIdx, cIdx)
	}
	if event != "CHANNEL_DELETE" {
		c.AddChan(gIdx, ch)
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
