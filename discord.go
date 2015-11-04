package discord

import (
	"fmt"
	"net/http"
)

func Login(email, pass string) (Discord, error) {
	client := Discord{
		Client: &http.Client{ },
		LoggingIn: true,
	}
	req := Creds{
			Email: email,
			Pass: pass,
	}
	resp := CredsResp{}
	err := client.Post(LoginURL, req, &resp)
	if err != nil {
		return client, err
	}
	
	client.Token = resp.Token
	if client.Token == "" {
		return client, CredsError(resp)
	}
	fmt.Printf("User %s logged in successfully!\n", req.Email)
	return client, nil
}
func (c Discord) Logout() error {
	type Req struct{
		Token string `json:"token"`
	}
	request := Req{Token: c.Token}
	err := c.Post(LogoutURL, request, nil)
	if err != nil {
		return err
	}
	c.Token = ""
	fmt.Println("User was logged out successfully! (once implemented server-side)")
	return nil
}
