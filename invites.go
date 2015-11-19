package discord

import (
	"fmt"
)

func (c Discord) InviteRevoke(code string) error {
	url := fmt.Sprintf(InviteURL, code)
	if err := c.Delete(url); err != nil { // givs response
		return err
	}
	
	fmt.Println("revoked invite!")
	return nil
}
//validation of code returns full invite
func (c Discord) InviteInfo(code string) (Invite, error) {
	url := fmt.Sprintf(InviteURL, code)
	resp := Invite{}
	if err := c.Get(url, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("got invite info!")
	return resp, nil
}
func (c Discord) InviteAccept(code string) (Invite, error) {
	url := fmt.Sprintf(InviteURL, code)
	resp := Invite{}
	if err := c.Post(url, nil, &resp); err != nil {
		return resp, err
	}
	
	fmt.Println("accepted invite info!")
	return resp, nil
}
