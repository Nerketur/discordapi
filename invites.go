package discord

import (
	"fmt"
)

func (c Discord) InviteRevoke(code string) (err error) {
	url := fmt.Sprintf(InviteURL, code)
	if err = c.Delete(url); err != nil { // givs response
		return
	}
	
	fmt.Println("revoked invite!")
	return
}
//validation of code returns full invite
func (c Discord) InviteInfo(code string) (resp Invite, err error) {
	url := fmt.Sprintf(InviteURL, code)
	if err = c.Get(url, &resp); err != nil {
		return 
	}
	
	fmt.Println("got invite info!")
	return
}
func (c Discord) InviteAccept(code string) (resp Invite, err error) {
	url := fmt.Sprintf(InviteURL, code)
	if err = c.Post(url, nil, &resp); err != nil {
		return
	}
	
	fmt.Println("accepted invite info!")
	return
}
