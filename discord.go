package discord

import (
	"fmt"
	"net/http"
)

func Login(email, pass string) (Discord, error) {
	client := Discord{ // only created once per client.
		Client: &http.Client{ },
		LoggingIn: true,
	}
	req := struct{ // no need to save this struct as its really only ever used ONCE per client.
		Email string `json:"email"`
		Pass  string `json:"password"`
	}{
		Email: email,
		Pass: pass,
	}
	resp := Creds{}
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
	request := struct{Token string `json:"token"`}{Token: c.Token}
	err := c.Post(LogoutURL, request, nil)
	if err != nil {
		return err
	}
	c.Token = ""
	fmt.Println("User was logged out successfully! (once implemented server-side)")
	return nil
}
