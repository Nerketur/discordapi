package discord

import (
	"fmt"
	"net/http"
)

type UnknownError string
type TokenError string
type EncodingError string
type PostError string
type ReadError string
type CredsError Creds
type HTTPError http.Response


func (e UnknownError) Error() string {
	return fmt.Sprintf("error: %s", string(e))
}
func (e CredsError) Error() string {
	login := Creds(e)
	if len(login.Email) > 0 {
		return fmt.Sprintf("token error: %s", login.Email[0])
	}
	return fmt.Sprintf("token error: %s", login.Pass[0])
}
func (e TokenError) Error() string {
	return fmt.Sprintf("token error: %s", string(e))
}
func (e EncodingError) Error() string {
	return fmt.Sprintf("encoding to json error: %s", string(e))
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
