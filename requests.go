package discord

import (
	"fmt"
	"io"
	"bytes"
	"net/http"
	"encoding/json"
	//"os"
)

const (
	APIURL = "https://discordapp.com/api"
		AuthURL = APIURL  + "/auth"
			LoginURL = AuthURL  + "/login" //post
			LogoutURL = AuthURL + "/logout" //post
			RegisterURL = AuthURL + "/register"
			VerifyURL = AuthURL + "/verify"
				VerifyResendURL = VerifyURL + "/resend"
			ForgotURL = AuthURL + "/forgot"
			ResetURL = AuthURL + "/reset"
	
		InviteURL = APIURL + "/invite/%s" //invite ID (get forinfo, post to accept)
		
		GuildsURL = APIURL + "/guilds" //get for list, post for new, patch for edit)
			GuildIDURL = GuildsURL + "/%s" //guild (server) ID
				GuildBansURL = GuildIDURL + "/bans"
				GuildMembersURL = GuildIDURL + "/members" // get
				GuildRolesURL = GuildIDURL + "/roles"     // get, patch
					GuildRoleIDURL = GuildRolesURL + "/%s"     // get?, put?, patch
				GuildChansURL = GuildIDURL + "/channels"  // get
		
		ChansURL                  = APIURL  + "/channels" // chanID
			ChanIDURL             = ChansURL + "/%s" // chanID
				ChanMsgsURL       = ChanIDURL + "/messages" // chanID (get)
					MsgURL        = ChanMsgsURL + "/%s" //channel ID, message ID (Patch edit, delete del)
						MsgAckURL = MsgURL + "/ack" //channel ID, message ID (Post)
				ChanInviteURL     = ChanIDURL + "/invites" // chanID
				ChanTypingURL     = ChanIDURL + "/typing" // chanID (post only)
				ChanPermsURL      = ChanIDURL + "/permissions" // chanID
					ChanPermIDURL = ChanPermsURL + "/%s" // chanID, permID
		
		UsersURL = APIURL + "/users" //invite ID
			UserIDURL = UsersURL + "/%s" // user ID
				UserChansURL = UserIDURL + "/channels" //(get chans)
				UserGuildsURL = UserIDURL + "/guilds" // get
				UserAvatarsURL = UserIDURL + "/avatars" // get
					UserAvatarIDURL = UserAvatarsURL + "/%s" // avatar ID
			MyURL = UsersURL + "/@me"
				MySettingsURL = MyURL + "/settings" // get
				MyDevicesURL = MyURL + "/devices" // get
				MyConnectionsURL = MyURL + "/connections" // get
				MyChansURL = MyURL + "/channels"
				MyGuildsURL = MyURL + "/guilds"
)

func (c Discord) Send(method, url string, data, want interface{}) error {
	//note data and want are interfaces, so theoreticaly any object can be sent via JSON.
	//We don't have to worry about it, as an error will propogate if we send something unmarshalable
	if c.Token == "" && !c.LoggingIn { // not logged in or logging in
		return TokenError("Not logged in!")
	}
	
	b, err := json.Marshal(data)
	if err != nil {
		return EncodingToError(fmt.Sprintf("%s", err))
	}
	var send bytes.Buffer
	
	if data != nil { // really for speed. can be removed with no issue
		send.Write(b)
		//os.Stdout.Write(b)
		//os.Stdout.Write([]byte{10})
	}
	req, err := http.NewRequest(method, url, &send)
	if err != nil {
		return UnknownError(fmt.Sprintf("%s", err))
	}
	if data == nil {
		req.Body = nil
	}
	if c.Token != "" {
		req.Header.Add("Authorization", c.Token)
	}
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	resp, err := c.Client.Do(req) // at this point, bytes buffer is closed if needed
	if err != nil { //if theres an err, body couldbe nil
		return PostError(fmt.Sprintf("%s", err)) // body is nil here
	}
	defer resp.Body.Close() //after so panic wont occur for nil body
	
	var buff bytes.Buffer
	io.Copy(&buff, resp.Body)
	
	if resp.StatusCode == http.StatusForbidden {
		message := struct{Message string}{}
		err = json.Unmarshal(buff.Bytes(), &message)
		return PermissionsError(message.Message)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		message := struct{Message string}{}
		err = json.Unmarshal(buff.Bytes(), &message)
		return PermissionsError(fmt.Sprintf("%s -- %v", message.Message, resp.StatusCode))
	}
	
	if want != nil {
		err = json.Unmarshal(buff.Bytes(), &want)
		//os.Stdout.Write(buff.Bytes())
		//os.Stdout.Write([]byte{10})
		
		if err != nil {
			return EncodingFromError(fmt.Sprintf("%s", err))
		}
	}
	return nil
}
func (c Discord) Post(url string, req, resp interface{}) error {
	return c.Send("POST", url, req, resp)
}
func (c Discord) Put(url string, req interface{}) error {
	return c.Send("PUT", url, req, nil)
}
func (c Discord) Get(url string, resp interface{}) error {
	return c.Send("GET", url, nil, resp)
}
func (c Discord) Delete(url string) error {
	return c.Send("DELETE", url, nil, nil)
}
