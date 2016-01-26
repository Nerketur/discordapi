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
func (g Guild) GuildCacheMembers() (ret []Member) {
	for _, val := range g.memCache {
		ret = append(ret, val)
	}
	return
}
func (c Discord) GuildCacheChans(guildID string) (ret []Channel) {
	for _, val := range c.chnCache {
		if val.GuildID == guildID {
			ret = append(ret, val)
		}
	}
	return
}
func (c Discord) ChanCacheAll() (ret []Channel) {
	for _, val := range c.chnCache {
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
func (c Discord) MemberCache(userID, guildID string) (ret Member, err error) {
	var g Guild
	if g, err = c.GuildCache(guildID); err != nil {
		return
	}
	var ok bool
	if ret, ok = g.memCache[userID]; !ok {
		err = IDNotFoundError(userID + " in " + guildID)
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
	} else if debug {
		fmt.Println("id:", guildID)
		fmt.Println("name:", c.gldCache[guildID].Name)
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

func (c *Discord) SetUserCache(u User) {
	c.usrCache[u.ID] = u
}
func (c *Discord) SetMemberCaches(m Member){
	g, err := c.GuildCache(m.GuildID)
	
	if err != nil {
		//fmt.Println("error adding to member cache:", err)
		//fmt.Println("ignoring...")
		return
	}
	u, err := c.UserCache(m.User.ID)
	if err != nil {
		u = m.User // if user doesn't exist, use the one we have
	}
	mb, ok := g.memCache[m.User.ID]
	if !ok {
		mb = m // if member doesn't exist, use the one we have
	}
	u.guildSet[g.ID] = struct{}{} //save guildID to user id slice
	c.SetUserCache(u) //replace user
	//user is updated.  now update user in member
	mb.User = u
	//member is updated.  now update guild cache
	g.memCache[mb.User.ID] = mb //save member to guild
	c.SetGuildCache(g) //replace guild
	//guild done, complete
	//we have users cache to get all guild ids a user is a member of
	//from this we go to the respective guild to get the member info
	//said member info now contains the same user info as before
	
	//user   -> guilds (id) -- non ID requires pointers
	//guild  -> members
	//member -> user //cant get member info without user and guild
	//               //but member itself gives that to us
}
func (c *Discord) SetMessageCache(m Message) {
	c.msgCache[m.ID] = m
}
func (c *Discord) SetGuildCache(g Guild) {
	c.gldCache[g.ID] = g
}
func (c *Discord) SetChanCache(ch Channel) {
	c.chnCache[ch.ID] = ch
}
func (c *Discord) SetPrivChanCache(ch Channel) {
	c.priCache[ch.Recipient.ID] = ch
}

func (c *Discord) DelUserCache(u User) {
	delete(c.usrCache, u.ID)
}
func (c *Discord) DelMemberCaches(m Member) {
	//here we do the opposite of set
	//we first check guild
	g, err := c.GuildCache(m.GuildID)
	if err != nil { //ignore remove if doesn't exist
		fmt.Println("error removing from member cache:", err)
		fmt.Println("ignoring...")
		return
	}
	u, err := c.UserCache(m.User.ID)
	if err != nil { //ignore remove if doesnt exist
		fmt.Println("error removing from member cache:", err)
		fmt.Println("ignoring...")
		return
	}
	m, ok := g.memCache[m.User.ID]
	if !ok { //ignore remove if doesn't exist
		fmt.Println("error removing from member cache:", err)
		fmt.Println("ignoring...")
		return
	}
	//we remove the guild data from our user
	delete(u.guildSet, g.ID)
	//we save updated user
	c.SetUserCache(u) //replace user
	//we delete member from guild
	delete(g.memCache, m.User.ID)
	//member is updated.  now update guild cache
	c.SetGuildCache(g) //replace guild
	
}
func (c *Discord) DelMessageCache(m Message) {
	delete(c.msgCache, m.ID)
}
func (c *Discord) DelGuildCache(g Guild) {
	delete(c.gldCache, g.ID)
}
func (c *Discord) DelChanCache(ch Channel) {
	delete(c.chnCache, ch.ID)
}
func (c *Discord) DelPrivChanCache(ch Channel) {
	delete(c.priCache, ch.Recipient.ID)
}

func (c Discord) FindNameUserCache(name string) (ret []User, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.usrCache {
		if val.Username == name {
			ret = append(ret, val)
			err = nil
		}
	}
	return
}
/* func (c Discord) FindNameMessageCache(name string) (ret []Message, err error) {
	err = NameNotFoundError(name)
	for key, val := range c.msgCache {
		if val.Name == name {
			ret = append(ret, val)
			err = nil
		}
	}
	return
} */
func (c Discord) FindNameGuildCache(name string) (ret []Guild, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.gldCache {
		if val.Name == name {
			ret = append(ret, val)
			err = nil
		}
	}
	return
}
func (c Discord) FindNameChanCache(name string) (ret []Channel, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.chnCache {
		if val.Name == name {
			ret = append(ret, val)
			err = nil
		}
	}
	return
}
func (c Discord) FindNamePrivChanCache(name string) (ret []Channel, err error) {
	err = NameNotFoundError(name)
	for _, val := range c.priCache {
		if val.Name == name {
			ret = append(ret, val)
			err = nil
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
	for _, ch := range g.Channels {
		ch.GuildID = g.ID //needed to let our cache lookup work
		c.SetChanCache(ch)
	}
	for _, m := range g.Members {
		u := m.User
		u.guildSet = make(map[string]struct{})
		c.SetUserCache(u)
		c.SetMemberCaches(m)
	}
	fmt.Printf("%#v", g.memCache)
}
func (c *Discord) RemGuild(g Guild) {
	for _, m := range g.Members {
		c.DelMemberCaches(m)
	}
	for _, ch := range g.Channels {
		c.DelChanCache(ch)
	}

	c.DelGuildCache(g)
}

func (c *Discord) GuildParseWS(event string, g Guild) {
	switch event {
	case "GUILD_CREATE":
		//init the map here
		g.memCache = make(map[string]Member)
		fallthrough
	case "GUILD_UPDATE":
		c.AddGuild(g)
	case "GUILD_DELETE":
		c.RemGuild(g)
	}
}

//guild members
func (c *Discord) AddMember(m Member) {
	c.SetMemberCaches(m)
}
func (c *Discord) RemMember(m Member) {
	c.DelMemberCaches(m)
}

func (c *Discord) GuildMemberParseWS(event string, m Member) {
	switch event {
	case "GUILD_MEMBER_ADD":
		m.User.guildSet = make(map[string]struct{})
		fallthrough
	case "GUILD_MEMBER_UPDATE":
		c.AddMember(m)
	case "GUILD_MEMBER_REMOVE":
		c.RemMember(m)
	}
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
	//c.Version   = ws.Version
	//c.SessionId = ws.SessionId
	//c.ReadState = ws.ReadState
	for _, channel := range ws.PrivateChannels {
		c.SetPrivChanCache(channel)
	}
	for _, guild := range ws.Guilds {
		//if debug {
			fmt.Println("guild:", guild.Name)
			fmt.Println("   id:", guild.ID)
		//}
		guild.memCache = make(map[string]Member) //make here to avoid errors
		c.SetGuildCache(guild)
		
		for _, ch := range guild.Channels {
			ch.GuildID = guild.ID //needed to let our cache lookup work
			c.SetChanCache(ch)
		}
		for _, m := range guild.Members {
			u := m.User
			m.GuildID = guild.ID //fixes lookup issues
			u.guildSet = make(map[string]struct{})
			c.SetUserCache(u)
			c.SetMemberCaches(m)
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
		g, err := c.GuildCache(guild.ID)
		if err != nil {
			fmt.Println("Error!")
		}
		//if debug {
			fmt.Printf("%v\t%v\n", len(g.memCache), len(guild.Members))
		//}
	}
	fmt.Println("Caches filled!")
	//c.Stop()
}
