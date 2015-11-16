package discord

import (
	"fmt"
	"net/http"
)

type UnknownError string
type TokenError string
type EncodingToError string
type EncodingFromError string
type PostError string
type ReadError string
type CredsError Creds
type HTTPError http.Response
type PermissionsError string



func (e PermissionsError) Error() string {
	return fmt.Sprintf("permission error: %s\n", string(e))
}
func (e UnknownError) Error() string {
	return fmt.Sprintf("error: %s", string(e))
}
func (e CredsError) Error() string {
	login := Creds(e)
	resp := "token error:\n"
	for _, line := range login.Email {
		resp += "\t" + line + "\n"
	}
	for _, line := range login.Pass {
		resp += "\t" + line + "\n"
	}
	return resp
}
func (e TokenError) Error() string {
	return fmt.Sprintf("token error: %s", string(e))
}
func (e EncodingToError) Error() string {
	return fmt.Sprintf("encoding to json error: %s", string(e))
}
func (e EncodingFromError) Error() string {
	return fmt.Sprintf("encoding from json error: %s", string(e))
}
func (e PostError) Error() string {
	return fmt.Sprintf("post error: %s", string(e))
}
func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %v", http.Response(e).StatusCode)
}
func (e ReadError) Error() string {
	return fmt.Sprintf("read error: %s", string(e))
}
