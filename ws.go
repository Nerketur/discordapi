package discord

import (
	"fmt"
	"time"
	"github.com/gorilla/websocket"
	"encoding/json"
	//"strconv"
)

type State struct{
	MentionCount  int    `json:"mention_count"`
	LastMessageId string `json:"last_message_id,omitempty"`
	ID            string `json:"id"`
}
type WSVoiceState struct{
	UserID    string  `json:"user_id"`
	Token     string  `json:"token"`
	Suppress  bool    `json:"suppress"`
	SessionID string  `json:"session_id"`
	SelfMute  *bool   `json:"self_mute"`
	SelfDeaf  *bool   `json:"self_deaf"`
	Mute      bool    `json:"mute"`
	Deaf      bool    `json:"deaf"`
	ChannelID *string `json:"channel_id"`
	GuildID   string  `json:"guild_id"`
}
type WSPres struct{
	User      User     `json:"user"`
	Status    string   `json:"status"`
	Roles     []string `json:"roles,omitempty"`
	GuildID   string   `json:"guild_id,omitempty"`
	Game      *Game    `json:"game"`
}

type _pres []WSPres

type Game struct{
	Name string `json:"name"`
} 
func (m *Game) UnmarshalJSON(raw []byte) (err error) {
	var rawMess json.RawMessage
	rawDat := struct{
		Name *json.RawMessage `json:"name"`
	}{
		Name: &rawMess,
	}
	err = json.Unmarshal(raw, &rawDat)
	if err != nil {
		return // sould be no error.  if there is, something is very wrong
	}
	msg := Game{
		Name: fmt.Sprint(rawMess), // convert to string
	}
	m = &msg
	return
}
/* func (m *WSPres) UnmarshalJSON(raw []byte) (err error) {
	type wsPres WSPres
	tmp := wsPres{}
	err = json.Unmarshal(raw, &tmp)
	if err != nil {
		fmt.Println("err in presence:")
		return
	}
	if tmp.GameIDStr != "" {
		tmpint, err := strconv.Atoi(tmp.GameIDStr)
		if err != nil {
			return err
		}
		tmp.GameID = &tmpint
		tmp.GameIDStr = ""
	}
	*m = WSPres(tmp)
	return
} */

type WSMsg struct{
    Type string      `json:"t,omitempty"`
    Seq  int         `json:"s,omitempty"`
    Op   int         `json:"op"`
	Data interface{} `json:"d"`  //deferred until we know what it is.
	time time.Time
}

func (m *WSMsg) UnmarshalJSON(raw []byte) (err error) {
	type wsMsg WSMsg
	var rawData json.RawMessage
	msg := wsMsg{
		Data: &rawData,
	}
	err = json.Unmarshal(raw, &msg)
	switch msg.Type {
	//code duplication because of Go restrictions
	//unsure how to make it shorter :(
	case "READY":
		data := READY{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "PRESENCE_UPDATE":
		var data json.RawMessage
		err = json.Unmarshal(rawData, &data)
		msg.Data = PRESENCE_UPDATE(data)
		//keep it a jsonRawMessage t parse later
	case "GUILD_CREATE":
		data := GUILD_CREATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_UPDATE":
		data := GUILD_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_DELETE":
		data := GUILD_DELETE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_ROLE_CREATE":
		data := GUILD_ROLE_CREATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_ROLE_UPDATE":
		data := GUILD_ROLE_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_ROLE_DELETE":
		data := GUILD_ROLE_DELETE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_MEMBER_ADD":
		data := GUILD_MEMBER_ADD{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_MEMBER_UPDATE":
		data := GUILD_MEMBER_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_MEMBER_REMOVE":
		data := GUILD_MEMBER_REMOVE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_BAN_ADD":
		data := GUILD_BAN_ADD{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_BAN_REMOVE":
		data := GUILD_BAN_REMOVE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "GUILD_INTEGRATIONS_UPDATE":
		data := GUILD_INTEGRATIONS_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "CHANNEL_CREATE":
		data := CHANNEL_CREATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "CHANNEL_UPDATE":
		data := CHANNEL_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "CHANNEL_DELETE":
		data := CHANNEL_DELETE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "MESSAGE_CREATE":
		data := MESSAGE_CREATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "MESSAGE_UPDATE":
		data := MESSAGE_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "MESSAGE_DELETE":
		data := MESSAGE_DELETE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "MESSAGE_ACK":
		data := MESSAGE_ACK{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "TYPING_START":
		data := TYPING_START{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "USER_SETTINGS_UPDATE":
		data := USER_SETTINGS_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	case "VOICE_STATE_UPDATE":
		data := VOICE_STATE_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
	default:
		fmt.Printf("unknown message type: %q\n", msg.Type)
	}
	if msg.Type != "READY" {
		tmp := msg
		if _, ok := msg.Data.(PRESENCE_UPDATE); ok {
			tmp.Data = fmt.Sprintf("%s", msg.Data)
		}
		fmt.Printf("object:\n\t%#v\n", tmp)
		fmt.Printf("raw:\n\t%s\n", rawData)
	}
	*m = WSMsg(msg)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("type: %s\n", msg.Type)
		if msg.Type != "READY" {
			fmt.Printf("rawdata:\n\t%s\n", rawData)
			fmt.Printf("struct:\n\t%#v\n\n", msg.Data)
		} else {
			if perr, ok := err.(*json.UnmarshalTypeError); ok {
				fmt.Printf("here: %s\n\n", rawData[perr.Offset-100:perr.Offset+100])
			}
		}
	}
	return
}

func (m *MESSAGE_CREATE) UnmarshalJSON(raw []byte) (err error) {
	msg := Message{}
	non := struct{
		Nonce int64
	}{}
	err = json.Unmarshal(raw, &msg)
	if err != nil {
		err = json.Unmarshal(raw, &non)
		if err != nil {
			return
		}
		if non.Nonce != 0 {
			msg.Nonce = non.Nonce
		}
	}
	
	*m = MESSAGE_CREATE(msg)
	return nil
}


//TODO add compression
type READY struct{ // op from server (0)
	Version           int        `json:"v"`
	User              User       `json:"user"`
	SessionId         string     `json:"session_id"`
	//SessionId         int        `json:"session_id"`//testing error
	ReadState         []State    `json:"read_state"`
	PrivateChannels   []Channel  `json:"private_channels"`
	HeartbeatInterval uint64     `json:"heartbeat_interval"`
	Guilds            []Guild    `json:"guilds"`
}
/*
* means unparsed
+ means implemented with cache
x means not planning to implement
- means eventual new event
= means not fully impemented
? means possible new event
    +CHANNEL_CREATE
    ?CHANNEL_CREATE (private)
    +CHANNEL_DELETE
    ?CHANNEL_DELETE (private)
    +CHANNEL_UPDATE
    +GUILD_CREATE
    -GUILD_CREATE (unavailable)
    +GUILD_DELETE
    -GUILD_DELETE (unavailable)
    +GUILD_UPDATE
    =GUILD_BAN_ADD
    =GUILD_BAN_REMOVE
    +GUILD_MEMBER_ADD
    +GUILD_MEMBER_REMOVE
    +GUILD_MEMBER_UPDATE
    +GUILD_ROLE_CREATE
    +GUILD_ROLE_DELETE
    +GUILD_ROLE_UPDATE
    GUILD_INTEGRATIONS_UPDATE
    MESSAGE_CREATE
    MESSAGE_DELETE
    MESSAGE_UPDATE
    MESSAGE_UPDATE (embeds only)
    +PRESENCE_UPDATE
    +READY
    xTYPING_START //pointless for bots?
    USER_SETTINGS_UPDATE //unsure if this is really needed
    VOICE_STATE_UPDATE
	xMESSAGE_ACK //unimportant,clients only
    OP 7

*/
type WSBan struct{
	GuildID string `json:"guild_id"`
	User    User   `json:"user"`
}
type WSRole interface{
	GetGuildID() string
	GetRoleID() string
	GetRole() Role
}
func (m GUILD_ROLE_CREATE) GetGuildID() string {return m.GuildID}
func (m GUILD_ROLE_UPDATE) GetGuildID() string {return m.GuildID}
func (m GUILD_ROLE_DELETE) GetGuildID() string {return m.GuildID}

func (m GUILD_ROLE_CREATE) GetRoleID() string {return m.Role.ID}
func (m GUILD_ROLE_UPDATE) GetRoleID() string {return m.Role.ID}
func (m GUILD_ROLE_DELETE) GetRoleID() string {return m.RoleID}

func (m GUILD_ROLE_CREATE) GetRole() Role {return m.Role}
func (m GUILD_ROLE_UPDATE) GetRole() Role {return m.Role}
func (m GUILD_ROLE_DELETE) GetRole() Role {return Role{ID: m.RoleID,}}

type MESSAGE_ACK struct{
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
}
type USER_SETTINGS_UPDATE Setting
type VOICE_STATE_UPDATE WSVoiceState
type GUILD_CREATE Guild
type GUILD_UPDATE Guild
type GUILD_DELETE Guild
type GUILD_MEMBER_ADD Member
type GUILD_MEMBER_UPDATE Member
type GUILD_MEMBER_REMOVE Member
type GUILD_BAN_ADD WSBan
type GUILD_BAN_REMOVE WSBan
type GUILD_ROLE_CREATE struct{
	GuildID string `json:"guild_id"`
	Role    Role   `json:"role"`
}
type GUILD_ROLE_UPDATE struct{
	GuildID string `json:"guild_id"`
	Role    Role   `json:"role"`
}
type GUILD_ROLE_DELETE struct{
	GuildID string `json:"guild_id"`
	RoleID  string `json:"role_id"`
}
type GUILD_INTEGRATIONS_UPDATE struct{
	GuildID string `json:"guild_id"`
}
type CHANNEL_CREATE Channel
type CHANNEL_UPDATE Channel
type CHANNEL_DELETE Channel
type MESSAGE_CREATE Message
type MESSAGE_UPDATE Message
type MESSAGE_DELETE Message
type PRESENCE_UPDATE json.RawMessage
type TYPING_START struct{
	ChanID    string `json:"channel_id"`
	//ChanID    int    `json:"channel_id"`  //check error
	Timestamp uint64 `json:"timestamp"`
	UserID    string `json:"user_id"`
}

type Properties struct{
	OS              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referrer        string `json:"$referrer"`
	ReferringDomain string `json:"$referring_domain"`
}
type INIT struct{ //op 2
	Token       string      `json:"token"`
	Version     int         `json:"v"`
	Properties  Properties `json:"properties"`
	LargeThresh int         `json:"large_threshold,omitempty"`
	Compress    bool        `json:"compress,omitempty"`
}
type Heartbeat time.Time //op 1
/*
{
        "op": 2,
        "d": {
                "token": "aaaaaaaabbbbbbbbccccccccddddddddeeeeeeeeffffffffgggggggghh",
                "v": 3,
                "properties": {
                        "$os": "Windows",
                        "$browser": "Chrome",
                        "$device": "",
                        "$referrer":" https://discordapp.com/@me",
                        "$referring_domain":"discordapp.com"
                },
                "large_threshold":100,
                "compress":true
        }
}
*/

/*
connect:
	1.) get the WS URL from thegateway endpoint
	2.) open a websocket connection to that URL.
	3.) send the correct first message (op 2) with v and current token in the data
persist:
	4.) read the ready event to see what the haeartbeat is
	5.) set up a way to send said heartbeat every heartbeat_interval milliseconds.  it can't miss, but it doesn't have to be exactly that many every time.
*/

func stopSend(con *websocket.Conn, done chan int) {
	//send close message
	fmt.Println("wsSend: got stop mssage")
	//fmt.Println("wsSend: closing send")
	//its okay if we continue recieving a bit before the close message is read.
	//channel was closed, send close frame and exit
	fmt.Println("wsSend: sending close frame")
	err := con.WriteControl(websocket.CloseMessage, nil, time.Now().Add(3*time.Second))
	if err != nil { //if theres an error sending, print err, exit
		fmt.Println("wsSend: control frame send err:", err)
	}
	close(done)
}

func wsSend(con *websocket.Conn, msgSend chan WSMsg, done, readDone, stopWS chan int) {
	defer fmt.Println("wsSend: goroutine stopped")
	seq := 1
	for {
		select {
		case <-stopWS:
			stopSend(con, done)
			return
		case <-readDone:
			stopSend(con, done)
			return
		case nextMsg, ok := <- msgSend:
			if ok {
				//send the message on the channel to the connection
				nextMsg.Seq = seq
				seq++
				fmt.Println("wsSend: sending msg", nextMsg.Type)
				j, _ := json.Marshal(nextMsg)
				fmt.Printf("wsSend: msg sent: `%s`\n", j)
				if err := con.WriteJSON(&nextMsg); err != nil {
					fmt.Println("wsSend:",err)
				}
			}
		}
	}
}

func wsRead(con *websocket.Conn, msgRead chan WSMsg, done chan int) {
	defer fmt.Println("wsRead: goroutine stopped")
	var nextMsg WSMsg
	for {
		//read the next message, put it on the channel
		err := con.ReadJSON(&nextMsg)
		if err != nil {
			fmt.Println("wsRead:", err)
			fmt.Println("wsRead: stopping")
			//we either need to send close, or we already sent close.
			//either way, close readdone
			close(done)  // this signals to wssend to start close if needed
			close(msgRead)  // this signals to process to stop
			//and exit
			return
		}
		//fmt.Println("Read from conn")
		//err = json.Unmarshal(msg, &nextMsg)
		fmt.Printf("wsRead: Read %T from conn\n", nextMsg.Data)
		nextMsg.time = time.Now()
		msgRead <- nextMsg
	}
}

func wsSendBeat(con *websocket.Conn, now time.Time) {
	beat := WSMsg{
		Op: 1,
		Data: now.Unix(),
	}
	fmt.Println("sent beat")
	con.WriteJSON(beat)
	j, _ := json.Marshal(beat)
	fmt.Printf("beat: %s\n", j)
}

func wsHeartbeat(con *websocket.Conn, msInterval uint64) {
	//send a heartbeat message now, and every msInterval mlliseconds after
	//don't use the msgSend channel because these MUST be sent at he requested time,
	//no matter how many messages are in the queue
	fmt.Println("started hearbeats.  interval (ms):", msInterval)
	
	t := time.NewTicker(time.Duration(msInterval) * time.Millisecond)
	
	//note that we currently never close the ticker.
	//later when we support resuming and more than one time use,
	//we will have to be able to close these tickers no longer in use
	
	wsSendBeat(con, time.Now()) //send a beat immediately
	for now := range t.C {
		wsSendBeat(con, now) //send a beat every tick
		fmt.Println("Tick:", now)
	}
}

func (c *Discord) wsFillCaches(ws READY) {
	for _, guild := range ws.Guilds {
		c.gldCache[guild.ID] = guild
		for _, member := range guild.Members {
			//TODO: member-guild link
			c.usrCache[member.User.ID] = member.User
		}
		for _, channel := range guild.Channels {
			c.chnCache[channel.ID] = channel
		}
	}
}

type Callback func(c Discord, event string, data interface{})

func (c *Discord) WSProcess(con *websocket.Conn, msgSend, msgRead chan WSMsg, CB *Callback) {
	if CB == nil {
		def := Callback(func(Discord, string, interface{}) {}) //the do nothing callback
		CB = &def
	}

	defer con.Close()

	//process events until a close message is encountered, or network error occurs.
	
	//Close frames must be sent and recieved.
	//A close frame sent requires waiting for recieving before closing
	//a close frame recieved requires sending, then closing
	//Gorrilla handles close frames by returning an error (along with the frame read)
	
	//accomplish this with read channel and sendsdone channels, which we wait for
	readDone := make(chan int)
	sendDone := make(chan int)
	fmt.Println("starting sender")
	go wsRead(con, msgRead, readDone)
	fmt.Println("starting reader")
	go wsSend(con, msgSend, sendDone, readDone, c.sigStop)
	fmt.Println("starting process")
	for msg := range msgRead {
		//process messages
		fmt.Println("processing message")
		switch msg.Op {
		case 0:
			//default, most
			d, ok := msg.Data.(*json.RawMessage)
			if ok {
				//json.rawmessage, so unexpected type
				fmt.Print("unexpected ")
			}
			fmt.Printf("type read '%v':\n", msg.Type)
			switch msg.Type {
			//here, we only catch types that change internal state
			//All the rest of the coding is done by the handeler of the callback.
			case "READY":
				parsed, ok := msg.Data.(READY)
				if !ok {
					fmt.Printf("Expected READY, got %T\n", parsed)
					break
				}
				
				//(current time - start time) - interval
				start := msg.time
				totalDur := time.Duration(parsed.HeartbeatInterval) * time.Millisecond
				//HBstartTimer := time.AfterFunc(totalDur-time.Since(start), func() {
				time.AfterFunc(totalDur-time.Since(start), func() {
					wsHeartbeat(con, parsed.HeartbeatInterval)
				})
				//fill arrays
				fmt.Println("filling cache...")
				c.cache = &parsed
				c.Me = &parsed.User
				go c.wsFillCaches(parsed)
				if len(c.cache.Guilds) == 0 {
					panic("cache should not be empty! -- just after assign --")
				}
				fmt.Println("cache filled!")

			//TODO: add code differentiating between unavailable guild
			// messages and normal messages
			// (make new events)
			case "GUILD_CREATE","GUILD_UPDATE","GUILD_DELETE":
				//parse guild stuff
				var parsed Guild
				switch msg.Data.(type) {
				case GUILD_CREATE:
					tmp, _ := msg.Data.(GUILD_CREATE)
					parsed = Guild(tmp)
				case GUILD_UPDATE:
					tmp, _ := msg.Data.(GUILD_UPDATE)
					parsed = Guild(tmp)
				case GUILD_DELETE:
					tmp, _ := msg.Data.(GUILD_DELETE)
					parsed = Guild(tmp)
				default:
					fmt.Printf("Expected discord.%s, got %T\n", msg.Type, msg.Data)
					close(c.sigStop)
					continue
				}
				//guild cache needs updating
				if msg.Type == "GUILD_DELETE" {
					delete(c.gldCache, parsed.ID)
				} else {
					c.gldCache[parsed.ID] = parsed
				}
				c.GuildParseWS(msg.Type, parsed)
			case "CHANNEL_CREATE","CHANNEL_UPDATE","CHANNEL_DELETE":
				//parse guild stuff
				var parsed Channel
				switch msg.Data.(type) {
				case CHANNEL_CREATE:
					tmp, _ := msg.Data.(CHANNEL_CREATE)
					parsed = Channel(tmp)
				case CHANNEL_UPDATE:
					tmp, _ := msg.Data.(CHANNEL_UPDATE)
					parsed = Channel(tmp)
				case CHANNEL_DELETE:
					tmp, _ := msg.Data.(CHANNEL_DELETE)
					parsed = Channel(tmp)
				default:
					fmt.Printf("Expected discord.%s, got %T\n", msg.Type, msg.Data)
					close(c.sigStop)
					continue
				}
				if parsed.Private {
				
					c.PrivateChannelParseWS(msg.Type, parsed)
				} else {
					if msg.Type == "CHANNEL_DELETE" {
						delete(c.chnCache, parsed.ID)
					} else {
						c.chnCache[parsed.ID] = parsed
					}
					c.ChannelParseWS(msg.Type, parsed)
				}
			case "GUILD_MEMBER_ADD","GUILD_MEMBER_UPDATE","GUILD_MEMBER_REMOVE":
				//parse guild member stuff
				var parsed Member
				switch msg.Data.(type) {
				case GUILD_MEMBER_ADD:
					tmp, _ := msg.Data.(GUILD_MEMBER_ADD)
					parsed = Member(tmp)
				case GUILD_MEMBER_UPDATE:
					tmp, _ := msg.Data.(GUILD_MEMBER_UPDATE)
					parsed = Member(tmp)
				case GUILD_MEMBER_REMOVE:
					tmp, _ := msg.Data.(GUILD_MEMBER_REMOVE)
					parsed = Member(tmp)
				default:
					fmt.Printf("Expected discord.%s, got %T\n", msg.Type, msg.Data)
					close(c.sigStop)
					continue
				}
				c.GuildMemberParseWS(msg.Type, parsed)
			case "GUILD_BAN_ADD", "GUILD_BAN_REMOVE":
				var parsed WSBan
				switch msg.Data.(type) {
				case GUILD_BAN_ADD:
					tmp, _ := msg.Data.(GUILD_BAN_ADD)
					parsed = WSBan(tmp)
				case GUILD_BAN_REMOVE:
					tmp, _ := msg.Data.(GUILD_BAN_REMOVE)
					parsed = WSBan(tmp)
				default:
					fmt.Printf("Expected discord.%s, got %T\n", msg.Type, msg.Data)
					close(c.sigStop)
					continue
				}
				c.GuildBanParseWS(msg.Type, parsed)
			case "GUILD_ROLE_CREATE", "GUILD_ROLE_UPDATE", "GUILD_ROLE_DELETE":
				var parsed WSRole
				switch msg.Data.(type) {
				case GUILD_ROLE_CREATE:
					parsed, _ = msg.Data.(GUILD_ROLE_CREATE)
				case GUILD_ROLE_UPDATE:
					parsed, _ = msg.Data.(GUILD_ROLE_UPDATE)
				case GUILD_ROLE_DELETE:
					parsed, _ = msg.Data.(GUILD_ROLE_DELETE)
				default:
					fmt.Printf("Expected discord.%s, got %T\n", msg.Type, msg.Data)
					close(c.sigStop)
					continue
				}
				c.GuildRoleParseWS(msg.Type, parsed)
			case "PRESENCE_UPDATE":
				parsed, ok := msg.Data.(PRESENCE_UPDATE)
				if ok {
					//parse presences
					var err error
					msg.Data, err = c.wsUpdatePres2(parsed)
					if err != nil {
						fmt.Println(err)
						fmt.Println("ignoring")
						
					}
				} else {
					fmt.Printf("Expected json.RawMessage, got %T\n", msg.Data)
					fmt.Println("Ignoring...")
				}
			default:
				if debug {
					if c.cache == nil {
						panic("cache should not be nil! -- default case --")
					} else if len(c.cache.Guilds) == 0 {
						panic("guild cache should not be empty! -- default case --")
					}
				}
				if ok {
					fmt.Printf("%s\n(needs adding)\n\n", d)
				} else {
					fmt.Printf("%#v\n\n", msg.Data)
				}
			}
			
			call := *CB
			call(*c, msg.Type, msg.Data)
		default:
			fmt.Printf("unexpected op '%v':\n%#v\n\n", msg.Op, msg.Data)
		}
	}
	fmt.Println("Reads closed, waiting for send to finish")
	<-sendDone
	fmt.Println("waiting for read to finish")
	<-readDone //just in case
	
	fmt.Println("closing safe/exit")
	close(c.sigSafe)
	select { // used to lose sigtime if not closed
	case <-c.sigTime: //dont close if closed
	default: //close if not closed
		close(c.sigTime)
	}
	fmt.Println("all finished, exiting process")
}
func (c *Discord) wsUpdatePres2(rep PRESENCE_UPDATE) (p WSPres, err error) {
	need := struct{
		GuildID string `json:"guild_id"`
		User struct{
			ID string
		}
	}{}
	var pres int
	
	raw := json.RawMessage(rep)
	
	if err = json.Unmarshal(raw, &need); err != nil {
		return
	}
	guild, err := guilds(c.cache.Guilds).FindIdxID(need.GuildID)
	if err != nil {
		fmt.Println("pres update guildID err:", err)
		return
	}
	g := c.cache.Guilds[guild]
	//after this point we update anyway, even if info is missing
	pres, err = _pres(g.Presences).FindIdx(need.User.ID)
	if err != nil {
		fmt.Println("pres update warning:", err)
		pres = len(g.Presences)
		g.Presences = append(g.Presences, WSPres{})
	}
	p = g.Presences[pres]
	if debug {
		fmt.Printf("\t\told:\n\t\t%#v\n", p)
	}
	err = json.Unmarshal(raw, &p)
	g.Presences[pres] = p
	c.cache.Guilds[guild] = g
	if debug {
		fmt.Printf("\t\tnew:\n\t\t%#v\n\n", c.cache.Guilds[guild].Presences[pres])
	}
	return
}

func (p _pres) FindIdx(ID string) (int, error) {
	for idx, ele := range p {
		if ele.User.ID == ID {
			return idx, nil
		}
	}
	return -1, IDNotFoundError(ID)
}

/*
func (c guilds) FindIdxID(ID string) (int, error) {
	for idx, ele := range c {
		if ele.ID == ID {
			return idx, nil
		}
	}
	return -1, IDNotFoundError(ID)
}
*/

func (c Discord) WSInit(con *websocket.Conn, msgChan chan WSMsg) {
	//send init on wire
	msgData := INIT{
		Token: c.Token,
		Version: 3, // hard-coded so changes will only happen when coded in
		Properties: Properties{
						 OS: "DiscordBot",
					Browser: "discord.go",
					 Device: "console",
				   Referrer: "",
			ReferringDomain: "",
		},
	}
	msg := WSMsg{
		Op: 2,
		Data: msgData,
	}
	j, _ := json.Marshal(msg)
	fmt.Printf("msgSentInit: `%s`\n", j)
	/* if err := con.WriteMessage(1, j); err != nil {
		fmt.Println("wsInit:",err)
	} */
	err := con.WriteJSON(msg)
	if err != nil {
		fmt.Println("wsInit:",err)
	}
}

func (c Discord) WSConnect(call *Callback) (err error) {
	
	gateway, err := c.Gateway()
	if err != nil {
		return
	}
	dialer := websocket.Dialer{}
	con, resp, err := dialer.Dial(gateway, nil)
	if err != nil {
		fmt.Printf("resp:\n%#v\n", resp)
		return
	}
	msgSend := make(chan WSMsg)
	msgRead := make(chan WSMsg)
	c.WSInit(con, msgSend)//ensure this is FIRST
	fmt.Println("init sent")
	go c.WSProcess(con, msgSend, msgRead, call)
	return
}