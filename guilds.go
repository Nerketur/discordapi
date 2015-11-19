package discord

import (
	"fmt"
)

type NotFoundError string
type guild []Guild

func (c NotFoundError) Error() string {
	return fmt.Sprintf("name not found: %s.  returning \"\"", string(c))
}

func (c Discord) GuildID(name string) (string, error) {
	resp, err := guild(c.MyGuilds).Find(name)
	if err != nil {
		return "", err
	}
	return resp, nil
}
func (c guild) Find(name string) (string, error) {
	for _, ele := range c {
		if ele.Name == name {
			return ele.ID, nil
		}
	}
	return "", NotFoundError(name)
}

func (c Discord) GuildMembers(guildID string) ([]Member, error) {
	
	resp := make([]Member, 0)
	err := c.Get(fmt.Sprintf(GuildMembersURL, guildID), &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Got members successfully!")
	return resp, nil
}

func (c Discord) GuildChannels(guildID string) ([]Channel, error) {
	
	resp := make([]Channel, 0)
	err := c.Get(fmt.Sprintf(GuildChansURL, guildID), &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Got channels successfully!")
	return resp, nil
}

func (c Discord) GetMyGuilds() ([]Guild, error) {
	
	resp := make([]Guild, 0)
	err := c.Get(MyGuildsURL, &resp)
	if err != nil {
		//fmt.Println(err)
		return resp, err
	}
	
	fmt.Println("Got guilds successfully!")
	return resp, nil
}

func (c Discord) GuildRoles(guildID string) ([]Role, error) {
	
	resp := make([]Role, 0)
	err := c.Get(fmt.Sprintf(GuildRolesURL, guildID), &resp)
	if err != nil {
		return resp, err
	}
	
	fmt.Println("Got roles successfully!")
	return resp, nil
}

func (c Discord) GuildAddRole(guildID string) (Role, error) {
	
	resp := Role{}
	err := c.Post(fmt.Sprintf(GuildRolesURL, guildID), nil, &resp)
	if err != nil {
		return resp, err
	}
	
	fmt.Println("added role successfully!")
	return resp, nil
}
func (c Discord) GuildEditRole(guildID string, r Role) (Role, error) {
	
	resp := Role{}
	err := c.Send("PATCH", fmt.Sprintf(GuildRoleIDURL, guildID, r.ID), &r, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Name != r.Name {
		fmt.Println("edit unsuccessful")
	} else {
		fmt.Println("edited role successfully!")
	}
	return resp, nil
}
func (c Discord) GuildAddNamedRole(guildID, name string) (Role, error) {
	
	resp, err := c.GuildAddRole(guildID)
	if err != nil {
		return resp, err
	}
	resp.Name = name
	fmt.Printf("%+v\n", resp)
	resp, err = c.GuildEditRole(guildID, resp)
	if err != nil {
		return resp, err
	}
	
	fmt.Println("added named role successfully!")
	return resp, nil
}
func (c Discord) GuildDeleteRole(guildID string, r Role) error {
	
	err := c.Send("DELETE", fmt.Sprintf(GuildRoleIDURL, guildID, r.ID), nil, nil)
	if err != nil {
		return err
	}
	fmt.Println("deleted role successfully!")
	return nil
}

type Members []Member

func (ms Members) Find(name string) []Member {
	ret := make([]Member, 0)
	for _, m := range ms {
		if m.User.Username == name {
			ret = append(ret, m)
		}
	}
	return ret
}

func (c Discord) GuildFindMember(guildID string, n string) ([]Member, error) {
	
	membs, err := c.GuildMembers(guildID)
	if err != nil {
		return []Member{}, err
	}
	return Members(membs).Find(n), nil
}

func (c Discord) GuildBans(guildID string) ([]Member, error) {
	
	resp := make([]Member, 0)
	err := c.Get(fmt.Sprintf(GuildBansURL, guildID), &resp)
	if err != nil {
		return resp, err
	}
	
	fmt.Println("Got bans successfully!")
	return resp, nil
}
func (c Discord) GuildAddBan(guildID, userID string, days int) error {
	
	url := fmt.Sprintf(GuildBanIDURL, guildID, userID)
	if days >= 0 {
		url += fmt.Sprintf("?delete-message-days=%v", days)
	}
	
	if err := c.Put(url, nil); err != nil {
		return err
	}
	
	fmt.Println("added ban successfully!")
	return nil
}
//PUT https://discordapp.com/api/guilds/:guild_id/bans/:user_id?delete-message-days=0

func (c Discord) GuildRemoveBan(guildID, userID string) error {
	
	url := fmt.Sprintf(GuildBanIDURL, guildID, userID)
	
	if err := c.Delete(url); err != nil {
		return err
	}
	
	fmt.Println("added ban successfully!")
	return nil
}

func (c Discord) GuildCreate(name, region string) (Guild, error) {
	
/* 
	reigons := make([]struct{
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
	}
 */	
	req := make(map[string]string)
	req["name"] = name
	req["region"] = region
	resp := Guild{}
	if err := c.Post(GuildsURL, req, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("created guild successfully!")
	return resp, nil
}
func (c Discord) GuildEdit(guildID, name, region string) (Guild, error) {
	req := make(map[string]string)
	req["name"] = name
	if region != "" {req["region"] = region}
	resp := Guild{}
	if err := c.Patch(fmt.Sprintf(GuildIDURL, guildID), req, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("edited guild successfully!")
	return resp, nil
}
func (c Discord) GuildLeave(guildID string) error {
	//resp := Guild{}
	if err := c.Delete(fmt.Sprintf(GuildIDURL, guildID)); err != nil {
		return err
	}
	
	fmt.Println("left guild successfully!")
	return nil
}
func (c Discord) GuildMemberEdit(guildID, userID string, roleIDs []string) error {
	req := struct{
		Roles []string `json:"roles,omitempty"`
	}{
		Roles: roleIDs,
	}
	if err := c.Patch(fmt.Sprintf(GuildMemberIDURL, guildID, userID), req, nil); err != nil {
		return err
	}
	
	fmt.Println("edited member successfully!")
	return nil
}
func (c Discord) GuildMemberKick(guildID, userID string) error {
	//resp := Guild{}
	if err := c.Delete(fmt.Sprintf(GuildMemberIDURL, guildID, userID)); err != nil {
		return err
	}
	
	fmt.Println("kicked member successfully!")
	return nil
}

