package discord

import (
	"fmt"
)

type guilds []Guild

func (c Discord) Guild(name string) (Guild, error) {
	return c.FindNameGuildCache(name)
}
func (c Discord) GuildID(name string) (string, error) {
	resp, err := c.Guild(name)
	return resp.ID, err //resp.id will be "" if invalid
}
func (c guilds) Find(name string) (ret Guild, err error) {
	var idx int // to avoid shadowing
	if idx, err = c.FindIdx(name); err == nil {
		ret = c[idx]
	}
	return
}
func (c guilds) FindIdx(name string) (int, error) {
	for idx, ele := range c {
		if ele.Name == name {
			return idx, nil
		}
	}
	return -1, NameNotFoundError(name)
}
func (c guilds) FindID(ID string) (ret Guild, err error) {
	var idx int // to avoid shadowing
	if idx, err = c.FindIdxID(ID); err == nil {
		ret = c[idx]
	}
	return
}
func (c guilds) FindIdxID(ID string) (int, error) {
	for idx, ele := range c {
		if ele.ID == ID {
			return idx, nil
		}
	}
	return -1, IDNotFoundError(ID)
}

//now use cache instead
func (c Discord) GuildMembers(guildID string) (resp []Member, err error) {
	var g Guild // to prevent shadowing
	g, err = c.GuildCache(guildID)
	resp = g.Members
	if err != nil {
		return
	}
	
	fmt.Println("Got members successfully!")
	return
}

func (c Discord) GuildChannels(guildID string) (resp []Channel, err error) {
	resp = make([]Channel, 0)
	err = c.Get(fmt.Sprintf(GuildChansURL, guildID), &resp)
	if err == nil {
		fmt.Println("Got channels successfully!")
	}
	return
}

func (c Discord) GuildChanCreate(guildID, name, kind string) (resp Channel, err error) {
	if kind != "text" && kind != "voice" {
		return resp, InvalidTypeError(kind)
	}
	url := fmt.Sprintf(GuildChansURL, guildID)
	req := struct{
		Name string `json:"name"`
		Type string `json:"type"`
	}{
		Name: name,
		Type: kind,
	}
	if err = c.Post(url, req, &resp); err == nil {
		fmt.Println("created channel!")
	}
	return
}

func (c Discord) Guilds() []Guild {
	return c.GuildCacheGuilds()
}

func (c Discord) GuildsRest() (resp []Guild, err error) {
	resp, err = make([]Guild, 0), c.Get(MyGuildsURL, &resp)
	if err == nil {
		fmt.Println("Got guilds successfully!")
	}
	return
}

func (c Discord) GuildRoles(guildID string) (resp []Role, err error) {
	resp, err = make([]Role, 0), c.Get(fmt.Sprintf(GuildRolesURL, guildID), &resp)
	if err == nil {
		fmt.Println("Got roles successfully!")
	}
	return
}

func (c Discord) GuildAddRole(guildID string) (resp Role, err error) {
	err = c.Post(fmt.Sprintf(GuildRolesURL, guildID), nil, &resp)
	if err == nil {
		fmt.Println("added role successfully!")
	}
	return
}
func (c Discord) GuildEditRole(guildID string, r Role) (resp Role, err error) {
	err = c.Patch(fmt.Sprintf(GuildRoleIDURL, guildID, r.ID), &r, &resp)
	if err != nil {
		return
	}
	if resp.Name != r.Name {
		fmt.Println("edit unsuccessful")
	} else {
		fmt.Println("edited role successfully!")
	}
	return
}
func (c Discord) GuildAddNamedRole(guildID, name string) (resp Role, err error) {
	resp, err = c.GuildAddRole(guildID)
	if err != nil {
		return
	}
	resp.Name = name
	resp, err = c.GuildEditRole(guildID, resp)
	if err == nil {
		fmt.Println("added named role successfully!")
	}
	return
}
func (c Discord) GuildDeleteRole(guildID, roleID string) (err error) {
	err = c.Delete(fmt.Sprintf(GuildRoleIDURL, guildID, roleID))
	if err == nil {
		fmt.Println("deleted role successfully!")
	}
	return
}

type members []Member

func (ms members) Find(name string) (ret []Member, err error) {
	err = NameNotFoundError("member: " + name)
	for _, m := range ms {
		if m.User.Username == name {
			ret, err = append(ret, m), nil
		}
	}
	return
}
func (ms members) FindIdxID(ID string) (ret int, err error) {
	ret, err = -1, IDNotFoundError("member:" + ID)
	for idx, m := range ms {
		if m.User.ID == ID {
			ret, err = idx, nil
		}
	}
	return
}
func (ms members) FindID(ID string) (ret Member, err error) {
	idx, err := ms.FindIdxID(ID)
	if err == nil {
		ret = ms[idx]
	}
	return
}

func (c Discord) GuildFindMember(guildID, n string) ([]Member, error) {
	membs, err := c.GuildMembers(guildID)
	if err != nil {
		return []Member{}, err
	}
	return members(membs).Find(n)
}

func (c Discord) GuildBans(guildID string) (resp []User, err error) {
	resp, err = make([]User, 0), c.Get(fmt.Sprintf(GuildBansURL, guildID), &resp)
	if err == nil {
		fmt.Println("Got bans successfully!")
	}
	return
}
func (c Discord) GuildAddBan(guildID, userID string, days int) (err error) {
	url := fmt.Sprintf(GuildBanIDURL, guildID, userID)
	if days >= 0 {
		url += fmt.Sprintf("?delete-message-days=%v", days)
	}
	if err = c.Put(url, nil); err == nil {
		fmt.Println("added ban successfully!")
	}
	return
}
//PUT https://discordapp.com/api/guilds/:guild_id/bans/:user_id?delete-message-days=0

func (c Discord) GuildRemoveBan(guildID, userID string) (err error) {
	url := fmt.Sprintf(GuildBanIDURL, guildID, userID)
	if err = c.Delete(url); err == nil {
		fmt.Println("added ban successfully!")
	}
	return
}

func (c Discord) GuildCreate(name, region string) (resp Guild, err error) {
/* 	reigons := make([]struct{
		//{"sample_hostname": "us-west19.discord.gg", "sample_port": 80, "id": "us-west", "name": "US West"}
		Hostname string `json:"sample_hostname"`
		Port     uint64 `json:"sample_port"`
		ID       string `json:"id"`
		Name     string `json:"name"`
	}{})
	err := c.VoiceRegions(&regions)
	if err != nil {
		return err
	}
	find := false
	for _, reg := range regions {
		if reg.ID == region {
			find = true
		}
	}
	if !find {
		region = ""
	} */	
	req := make(map[string]string)
	req["name"] = name
	req["region"] = region
	if err = c.Post(GuildsURL, req, &resp); err == nil {
		fmt.Println("created guild successfully!")
	}
	return
}
func (c Discord) GuildEdit(guildID, name, region string) (resp Guild, err error) {
	req := make(map[string]string)
	req["name"] = name
	if region != "" {req["region"] = region}
	if err = c.Patch(fmt.Sprintf(GuildIDURL, guildID), req, &resp); err == nil {
		fmt.Println("edited guild successfully!")
	}
	return
}
func (c Discord) GuildLeave(guildID string) (err error) {
	if err = c.Delete(fmt.Sprintf(GuildIDURL, guildID)); err == nil {
		fmt.Println("left guild successfully!")
	}
	return
}

func (c Discord) GuildMemberEdit(guildID, userID string, roleIDs []string) (err error) {
	req := struct{
		Roles []string `json:"roles,omitempty"`
	}{
		Roles: roleIDs,
	}
	if err = c.Patch(fmt.Sprintf(GuildMemberIDURL, guildID, userID), req, nil); err == nil {
		fmt.Println("edited member successfully!")
	}
	return
}
func (c Discord) GuildMemberKick(guildID, userID string) (err error) {
	if err = c.Delete(fmt.Sprintf(GuildMemberIDURL, guildID, userID)); err == nil {
		fmt.Println("kicked member successfully!")
	}
	return
}

func (c Discord) GuildInvitesList(guildID string) (resp []Invite, err error) {
	resp = make([]Invite, 0)
	err = c.Get(fmt.Sprintf(GuildInvitesURL, guildID), &resp)
	return
}
func (c Discord) GuildPruneInfo(guildID string, days int) (int, error) {
	resp := struct{Pruned int `json:"pruned"`}{}
	url := fmt.Sprintf(GuildPruneURL, guildID)
	if days > 0 {
		url += fmt.Sprintf("?days=%v", days)
	}
	err := c.Get(url, &resp)
	return resp.Pruned, err
}
func (c Discord) guildPrune(guildID string, days int) (int, error) { // private so people won't call it
	resp := struct{Pruned int `json:"pruned"`}{}
	url := fmt.Sprintf(GuildPruneURL, guildID)
	if days > 0 {
		url += fmt.Sprintf("?days=%v", days)
	}
	err := c.Post(url, nil, &resp)
	return resp.Pruned, err
}

func (c Discord) GuildEmbed(guildID string) (resp EmbedInfo, err error) {
	err = c.Get(fmt.Sprintf(GuildEmbedURL, guildID), &resp)
	return
}
func (c Discord) GuildIntegrations(guildID string) (resp []Integration, err error) {
	resp = make([]Integration, 0)
	err = c.Get(fmt.Sprintf(GuildIntegrationsURL, guildID), &resp)
	return
}

type roles []Role
func (r roles) FindIdxID(ID string) (int, error) {
	for idx, ele := range r {
		if ele.ID == ID {
			return idx, nil
		}
	}
	return -1, IDNotFoundError(ID)
}
