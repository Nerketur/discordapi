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
				GuildRolesURL = GuildIDURL + "/roles"     // get
				GuildChansURL = GuildIDURL + "/channels"  // get
		
		ChansURL                  = APIURL  + "/channels" // chanID
			ChanIDURL             = ChansURL + "/%s" // chanID
				ChanMsgsURL       = ChanIDURL + "/messages" // chanID (get)
					MsgURL        = ChanMsgsURL + "/%s" //channel ID, message ID (Patch edit, delete del)
						MsgAckURL = MsgURL + "/ack" //channel ID, message ID (Post)
				ChanInviteURL     = ChanIDURL + "/invites" // chanID
				ChanTypingURL     = ChanIDURL + "/typing" // chanID (post only)
				ChanPermsURL      = ChanIDURL + "/permissions" // chanID
		
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

				
	
	
)

func (c Discord) Send(method, url string, data, want interface{}) error {
	//note data and want are interfaces, so theoreticaly any object can be sent via JSON.
	//We don't have to worry about it, as an error will propogate if we send something unmarshalable
	if c.Token == "" && !c.LoggingIn { // not logged in or logging in
		return TokenError("Not logged in!")
	}
	var req *http.Request // important to define it before using
	var err error // define before use so we can use it later
	if data != nil {

		b, err := json.Marshal(data)
		if err != nil {
			return EncodingToError(fmt.Sprintf("%s", err))
		}
		//os.Stdout.Write(b)
		//os.Stdout.Write([]byte{10})
		req, err = http.NewRequest(method, url, bytes.NewBuffer(b))
		if err != nil {
			return UnknownError(fmt.Sprintf("%s", err))
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return UnknownError(fmt.Sprintf("%s", err))
		}
	}
	//err outside if would cause weird thing where err isn't populated
	if req == nil {
		return UnknownError("Why is req nil???") // should never happen
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
	
	//if resp.StatusCode < 200 || resp.StatusCode > 299 {
	//	return HTTPError(*resp)
	//}
	
	//on second thought, let the caller handle HTTP stuff
	
	var buff bytes.Buffer
	io.Copy(&buff, resp.Body)
	
	if want != nil {
		err = json.Unmarshal(buff.Bytes(), want)
		//os.Stdout.Write(buff.Bytes())
		//os.Stdout.Write([]byte{10})
		
		if err != nil {
			return EncodingFromError(fmt.Sprintf("%s", err))
		}
	}
	
	c.LoggingIn = false
	return nil
}
func (c Discord) Post(url string, req, resp interface{}) error {
	return c.Send("POST", url, req, resp)
}
func (c Discord) Get(url string, resp interface{}) error {
	return c.Send("GET", url, nil, resp)
}
