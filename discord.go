package discord

import (
	"fmt"
	"net/http"
	"crypto/sha1"
	"io/ioutil"
)

var debug bool = false

func Login(email, pass string) (Discord, error) {
	
	client := Discord{ // only created once per client.
		Client: &http.Client{ },
		MyGuilds: []Guild{},
		MyChans: []Channel{},
	}
	
	//start by trying to read the token from a file
	//files are named by sha1 of username+pass.
	
	//get sha1:
	s := email+pass
	////////////////////////////////////

//The pattern for generating a hash is sha1.New(), sha1.Write(bytes), then sha1.Sum([]byte{}). Here we start with a new hash.
	

    h := sha1.New()

//Write expects bytes. If you have a string s, use []byte(s) to coerce it to bytes.
	

    h.Write([]byte(s))

//This gets the finalized hash result as a byte slice. The argument to Sum can be used to append to an existing byte slice: it usually isn’t needed.
	

    bs := h.Sum(nil)

	/////////////////////////////
	
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
	
	return client, nil
}
func (c Discord) Logout() error {
	request := struct{Token string `json:"token"`}{Token: c.Token}
	err := c.Post(LogoutURL, request, nil)
	if err != nil {
		return err
	}
	c.Token = ""
	fmt.Println("User was logged out successfully! (once implemented server-side)")
	return nil
}

func Version() string {
	return "v0.5 alpha"
}
func VersionString() string {
	return fmt.Sprintf("Discord Go API %s", Version())
}

// GET https://discordapp.com/api/gateway
//
// {"url": "wss://gateway-nidhogg.discord.gg"}

func (c Discord) Gateway() (string, error) {
	req := struct{URL string `json"url"`}{}
	err := c.Get(GatewayURL, &req)
	return req.URL, err
}
func (c Discord) UserConnections() ([]Connection, error) {
	req := make([]Connection, 0)
	err := c.Get(MyConnectionsURL, &req)
	return req, err
}
func (c Discord) UserSettings() (Setting, error) {
	req := Setting{}
	err := c.Get(MySettingsURL, &req)
	return req, err
}
func (c Discord) SetUserSettings(s Setting) error {
	err := c.Put(MySettingsURL, &s)
	return err
}

func (c Discord) TutorialInfo() (Tutorial, error) {
	req := Tutorial{}
	err := c.Get(TutorialURL, &req)
	return req, err
}
/* func (c Discord) TutorialIndicatorsInfo() (struct{}, error) {
	req := struct{}{}
	err := c.Get(TutorialIndURL, &req)
	return req, err
} */