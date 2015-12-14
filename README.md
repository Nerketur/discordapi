# discordapi
A GoLang API wrapper for the Discord REST API (unofficial) (and eventually websockets)

v0.6 alpha

##About
This is a simple low-level REST API Wrapper for Discord.  I am not affiliated with or endorsed by Discrd at all, this is meant as a way to learn Go better as a language.  It may contain errors, bugs, and any other such oddities, but it seems to work for me :)

Later on, I will add websocket support, but if you want that right now, then please use [go-discord](https://github.com/gdraynz/go-discord) instead (or in addition to mine)!  What this means for you is you have to use polling if you want a functional bot.

Version system will be sporadic until 1.0, safe to ignore, other than alpha means preliminary, beta means set on paper, and stable means more or less set in stone.

##Usage
```go
package main

import (
	"fmt"
	"github.com/Nerketur/discordapi"
)

func main() {
	client, err := discord.Login("temp1@example.com", "12345")
	if err != nil {
		fmt.Println(err)
		return
	}
	if client.Client == nil {
		fmt.Println("Empty client!  Check your network connection.")
		return
	}
	defer func() {
		err = client.Logout()
		if err != nil {
			fmt.Println(err)
		}
	}()
	//other code
}
```

##Disclaimer
This is in ALPHA stage which means nothing is guaranteed to stay the same.  Things will likely change, and there will be absolutely NO attempt to maintain backwards compatibility until version 1.0.  This means things can and likely will break if this is used, upon upgrading to new versions.  Do not use for time- or safety-critical tasks.  You have been warned.

##Other libraries
- [discord.py](https://github.com/Rapptz/discord.py) (Python)
- [discord.js](https://github.com/discord-js/discord.js) (JS)
- [discord.io](https://github.com/izy521/discord.io) (JS)
- [Discord.NET](https://github.com/RogueException/Discord.Net) (C#)
- [DiscordSharp](https://github.com/Luigifan/DiscordSharp) (C#)
- [Discord4J](https://github.com/knobody/Discord4J) (Java)
- [discordrb](https://github.com/meew0/discordrb) (Ruby)
- [Discordgo](https://github.com/bwmarrin/Discordgo) (Go)
- [discord](https://github.com/Xackery/discord) (Go)
- [go-discord](https://github.com/gdraynz/go-discord) (Go)
