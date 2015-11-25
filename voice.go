package discord

import (
//	"fmt"
//	"net/http"
)

func (c Discord) VoiceRegions() (resp []Region, err error) {
	resp = make([]Region, 0)
	err = c.Get(VoiceRegionsURL, &resp)
	return
}

func (c Discord) VoiceIce() (resp Ice, err error) {
	resp = Ice{}
	err = c.Get(VoiceIceURL, &resp)
	return
}

