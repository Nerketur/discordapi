package discord

import (
	"fmt"
	"net/http"
	"crypto/sha1"
	"io/ioutil"
	"time"
)

var debug bool = false

func Login(email, pass string, call *Callback) (Discord, error) {
	
	client := Discord{ // only created once per client.
		Client: &http.Client{ },
		MyGuilds: []Guild{},
		MyChans: []Channel{},
		sigStop: make(chan int),
		sigSafe: make(chan int),
	}
	
	//start by trying to read the token from a file
	//files are named by sha1 of username+pass.
	
	//get sha1:
	s := email+pass
    h := sha1.New()
    h.Write([]byte(s))
    bs := h.Sum(nil)
	
	filename := fmt.Sprintf("%x.json", bs)
	dat, err := ioutil.ReadFile(filename)
	if err == nil {
		//file exists, dat has token
		fmt.Println("Using previous Token...")
		client.Token = string(dat)
	} else {
		
		// if we get here, ask for a new token:
		fmt.Println("Making new Token...")
		
		client.LoggingIn = true
		req := struct{ // no need to save this struct as its really only ever used ONCE per client.
			Email string `json:"email"`
			Pass  string `json:"password"`
		}{
			Email: email,
			Pass: pass,
		}
		resp := Creds{}
		err = client.Post(LoginURL, req, &resp) // does not error as long as request succeeds
		if err != nil {
			return client, err
		}
		
		client.Token = resp.Token
		if client.Token == "" {
			return client, CredsError(resp)
		}
		client.LoggingIn = false
		//cache result in file
		err = ioutil.WriteFile(filename, []byte(client.Token), 0644)
		if err != nil {
			fmt.Printf("WARNING: token info not saved - %T:%s", err, err)
		}
	}
	fmt.Printf("User %s (token %s) logged in successfully!\n", email, client.Token)
	fmt.Println("filling guild and chan arrys...")
	client.MyGuilds, err = client.GetMyGuilds()
	if err != nil {
		return client, err
	}
	client.MyChans, err = client.GetMyPrivateChans()
	if err != nil {
		return client, err
	}
	fmt.Println("Arrays filled!")
	fmt.Println("Attempting websocket connection...")
	err = client.WSConnect(client.sigStop, client.sigSafe, call)
	if err != nil {
		fmt.Println("websocket error:", err)
	} else {
		fmt.Println("websocket connected!")
	}
	return client, err
}
func (c Discord) Logout() (err error) {
	//first wait for websocket'
	fmt.Println("Waiting for websocket...")
	select {
	case <-c.sigSafe:
		fmt.Println(c.quitMsg)
		return
	}

	request := struct{Token string `json:"token"`}{Token: c.Token}
	err = c.Post(LogoutURL, request, nil)
	if err != nil {
		return
	}
	c.Token = ""
	fmt.Println("User was logged out successfully! (once implemented server-side)")
	return
}

func (c Discord) SetMaxRuntime(amt time.Duration, expireMsg, quitMsg string) {
	fmt.Println("setting timer")
	endTimer := time.NewTimer(amt)
	c.quitMsg = quitMsg
	for {
		select {
		case <-endTimer.C:
			fmt.Println(expireMsg)
			close(c.sigStop)
		case <-c.sigSafe:
			endTimer.Stop()
		}
		return
	}
}

func Version() string {
	return "v0.6 alpha"
}
func VersionString() string {
	return fmt.Sprintf("Discord Go API %s", Version())
}


func (c Discord) Gateway() (string, error) {
	req := struct{URL string `json"url"`}{}
	err := c.Get(GatewayURL, &req)
	return req.URL, err
}
func (c Discord) UserConnections() (req []Connection, err error) {
	req = make([]Connection, 0)
	err = c.Get(MyConnectionsURL, &req)
	return
}
func (c Discord) UserSettings() (req Setting, err error) {
	err = c.Get(MySettingsURL, &req)
	return
}
func (c Discord) SetUserSettings(s Setting) (err error) {
	err = c.Put(MySettingsURL, &s)
	return
}

func (c Discord) TutorialInfo() (req Tutorial, err error) {
	err = c.Get(TutorialURL, &req)
	return
}
/* func (c Discord) TutorialIndicatorsInfo() (req struct{}, err error) {
	err = c.Get(TutorialIndURL, &req)
	return
} */