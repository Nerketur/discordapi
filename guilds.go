package discord

import (
	"fmt"
)

type guild []Guild

func (c Discord) Guild(name string) (Guild, error) {
	return guild(c.MyGuilds).Find(name)
}
func (c Discord) GuildID(name string) (string, error) {
	resp, err := c.Guild(name)
	return resp.ID, err //resp.id will be "" if invalid
}
func (c guild) Find(name string) (Guild, error) {
	for _, ele := range c {
		if ele.Name == name {
			return ele, nil
		}
	}
	return Guild{}, NameNotFoundError(name)
}

//depricated.  may be removed because only in WS
func (c Discord) GuildMembers(guildID string) (resp []Member, err error) {
	resp = make([]Member, 0)
	err = c.Get(fmt.Sprintf(GuildMembersURL, guildID), &resp)
	if err != nil {
		return
	}
	
	fmt.Println("Got members successfully!")
	return
}

func (c Discord) GuildChannels(guildID string) (resp []Channel, err error) {
	resp = make([]Channel, 0)
	err = c.Get(fmt.Sprintf(GuildChansURL, guildID), &resp)
	if err != nil {
		return
	}
	
	fmt.Println("Got channels successfully!")
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
	if err = c.Post(url, req, &resp); err != nil {
		return
	}
	
	fmt.Println("created channel!")
	return
}

func (c Discord) GetMyGuilds() (resp []Guild, err error) {
	resp = make([]Guild, 0)
	err = c.Get(MyGuildsURL, &resp)
	if err != nil {
		return
	}
	
	fmt.Println("Got guilds successfully!")
	return
}

func (c Discord) GuildRoles(guildID string) (resp []Role, err error) {
	resp = make([]Role, 0)
	err = c.Get(fmt.Sprintf(GuildRolesURL, guildID), &resp)
	if err != nil {
		return
	}
	
	fmt.Println("Got roles successfully!")
	return
}

func (c Discord) GuildAddRole(guildID string) (resp Role, err error) {
	err = c.Post(fmt.Sprintf(GuildRolesURL, guildID), nil, &resp)
	if err != nil {
		return
	}
	
	fmt.Println("added role successfully!")
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
	if err != nil {
		return
	}
	
	fmt.Println("added named role successfully!")
	return
}
func (c Discord) GuildDeleteRole(guildID, roleID string) (err error) {
	err = c.Delete(fmt.Sprintf(GuildRoleIDURL, guildID, roleID))
	if err != nil {
		return
	}
	fmt.Println("deleted role successfully!")
	return
}

type Members []Member

func (ms Members) Find(name string) (ret []Member, err error) {
	ret = make([]Member, 0)
	for _, m := range ms {
		if m.User.Username == name {
			ret = append(ret, m)
		}
	}
	if len(ret) == 0 {
		err = NameNotFoundError("member: " + name)
	}
	return
}

func (c Discord) GuildFindMember(guildID, n string) ([]Member, error) {
	membs, err := c.GuildMembers(guildID)
	if err != nil {
		return []Member{}, err
	}
	return Members(membs).Find(n)
}

func (c Discord) GuildBans(guildID string) (resp []Member, err error) {
	resp = make([]Member, 0)
	err = c.Get(fmt.Sprintf(GuildBansURL, guildID), &resp)
	if err != nil {
		return
	}
	fmt.Println("Got bans successfully!")
	return
}
func (c Discord) GuildAddBan(guildID, userID string, days int) (err error) {
	url := fmt.Sprintf(GuildBanIDURL, guildID, userID)
	if days >= 0 {
		url += fmt.Sprintf("?delete-message-days=%v", days)
	}
	if err = c.Put(url, nil); err != nil {
		return
	}
	fmt.Println("added ban successfully!")
	return
}
//PUT https://discordapp.com/api/guilds/:guild_id/bans/:user_id?delete-message-days=0

func (c Discord) GuildRemoveBan(guildID, userID string) (err error) {
	url := fmt.Sprintf(GuildBanIDURL, guildID, userID)
	if err = c.Delete(url); err != nil {
		return
	}
	fmt.Println("added ban successfully!")
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
	if err = c.Post(GuildsURL, req, &resp); err != nil {
		return
	}
	fmt.Println("created guild successfully!")
	return
}
func (c Discord) GuildEdit(guildID, name, region string) (resp Guild, err error) {
	req := make(map[string]string)
	req["name"] = name
	if region != "" {req["region"] = region}
	if err = c.Patch(fmt.Sprintf(GuildIDURL, guildID), req, &resp); err != nil {
		return
	}
	fmt.Println("edited guild successfully!")
	return
}
func (c Discord) GuildLeave(guildID string) (err error) {
	if err = c.Delete(fmt.Sprintf(GuildIDURL, guildID)); err != nil {
		return
	}
	fmt.Println("left guild successfully!")
	return
}
func (c Discord) GuildMemberEdit(guildID, userID string, roleIDs []string) (err error) {
	req := struct{
		Roles []string `json:"roles,omitempty"`
	}{
		Roles: roleIDs,
	}
	if err = c.Patch(fmt.Sprintf(GuildMemberIDURL, guildID, userID), req, nil); err != nil {
		return
	}
	
	fmt.Println("edited member successfully!")
	return
}
func (c Discord) GuildMemberKick(guildID, userID string) (err error) {
	if err = c.Delete(fmt.Sprintf(GuildMemberIDURL, guildID, userID)); err != nil {
		return
	}
	fmt.Println("kicked member successfully!")
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
