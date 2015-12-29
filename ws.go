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
		data := PRESENCE_UPDATE{}
		err = json.Unmarshal(rawData, &data)
		msg.Data = data
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
		fmt.Printf("object:\n\t%#v\n", msg)
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
+ means parsed
- means eventual new event
    +CHANNEL_CREATE
    +CHANNEL_CREATE (private)
    +CHANNEL_DELETE
    +CHANNEL_DELETE (private)
    +CHANNEL_UPDATE
    +GUILD_CREATE
    -GUILD_CREATE (unavailable)
    +GUILD_DELETE
    -GUILD_DELETE (unavailable)
    +GUILD_UPDATE
    +GUILD_BAN_ADD
    +GUILD_BAN_REMOVE
    +GUILD_MEMBER_ADD
    +GUILD_MEMBER_REMOVE
    +GUILD_MEMBER_UPDATE
    +GUILD_ROLE_CREATE
    +GUILD_ROLE_DELETE
    +GUILD_ROLE_UPDATE
    +GUILD_INTEGRATIONS_UPDATE
    +MESSAGE_CREATE
    +MESSAGE_DELETE
    +MESSAGE_UPDATE
    +MESSAGE_UPDATE (embeds only)
    +PRESENCE_UPDATE
    +READY
    +TYPING_START
    +USER_SETTINGS_UPDATE
    +VOICE_STATE_UPDATE
	+MESSAGE_ACK
    OP 7

*/
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
type GUILD_BAN_ADD struct{
	GuildID string `json:"guild_id"`
	User    User   `json:"user"`
}
type GUILD_BAN_REMOVE struct{
	GuildID string `json:"guild_id"`
	User    User   `json:"user"`
}
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
type PRESENCE_UPDATE WSPres
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
	Properties  *Properties `json:"properties,omitempty"`
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

func wsSend(con *websocket.Conn, msgSend, other chan WSMsg, stopWS, exit chan int) {
	seq := 1
	for {
		select {
		case <-stopWS:
			//send close message
			fmt.Println("got stop mssage")
			fmt.Println("sending close frame (send, before read)")
			err := con.WriteControl(websocket.CloseMessage, nil, time.Now().Add(3*time.Second))
			if err != nil { //if theres an error sending, assume corrupted, exit
				fmt.Println("control frame send err:", err)
				select {
					case _, ok := <-exit:
						if ok {
							close(exit)
						}
					default:
				}
				close(msgSend)// panic on trying to send more
			}
			return //end for, exit immediately
			//its okay if we continue recieving a bit before the close message is read.
		case nextMsg, ok := <-msgSend:
			if ok {
				//send the message on the channel to the connection
				
				nextMsg.Seq = seq
				seq++
				fmt.Println("sending msg", nextMsg.Type)
				j, _ := json.Marshal(nextMsg)
				fmt.Printf("msg sent: `%s`\n", j)
				if err := con.WriteJSON(&nextMsg); err != nil {
					fmt.Println("wsSend:",err)
				}
			} else {
				//if ok is false, channel was closed by read, so send close frame and exit
				fmt.Println("sending close frame (send, after read)")
				con.WriteControl(websocket.CloseMessage, nil, time.Now())
				close(other)
				close(exit)
				return
			}
		}
	}
}

func wsRead(con *websocket.Conn, other, msgRead chan WSMsg, stopWS, exit, timer chan int) {
	var nextMsg WSMsg
	for {
		//read the next message, put it on the channel
		err := con.ReadJSON(&nextMsg)
		if err != nil {
			
			if _, ok := err.(*websocket.CloseError); !ok {
				//act as if timer elapsed
				close(timer)
			}
			//print value from connection
			fmt.Println("wsRead:",err)
			select {
			case <-stopWS:
				//send already sent frame.  close and exit
				fmt.Println("wsRead: send already sent, so exiting")
				close(msgRead)
				close(exit) // exit
			default:
				//send has NOT sent frame.  close send and exit
				fmt.Println("wsRead: send not sent close, so closing send then exiting")
				close(other)
				//do NOT close exit until close frame sent)
			}
			//send will close exit channel
			return
		}
		//fmt.Println("Read from conn")
		//err = json.Unmarshal(msg, &nextMsg)
		fmt.Printf("Read %T from conn\n", nextMsg.Data)
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

type Callback func(event string, data interface{})

func (c Discord) WSProcess(con *websocket.Conn, msgSend, msgRead chan WSMsg, stopWS, exit chan int, CB *Callback) {
	if CB == nil {
		def := Callback(func(string, interface{}) {}) //the do nothing callback
		CB = &def
	}

	defer con.Close()

	//process events until a close message is encountered, or network error occurs.
	
	//Close frames must be sent and recieved.
	//A close frame sent requires waiting for recieving before closing
	//a close frame recieved requires sending, then closing
	//Gorrilla handles close frames by returning an error (along with the frame read)
	fmt.Println("starting sender")
	go wsRead(con, msgSend, msgRead, stopWS, exit, c.sigTime) // if we err on read, we have to send close frame then exit.
	fmt.Println("starting reader")
	go wsSend(con, msgSend, msgRead, stopWS, exit) // if we send close frame, we have to wait for a response
	fmt.Println("starting process")
	for msg := range msgRead {
		//process messages
		fmt.Println("processing message")
		switch msg.Op {
		case 0:
			//default, most
			switch msg.Type {
			//here, we only catch types that change internal state
			//as of yet, the oly one to do that is "READY"
			//  (because heartbeats)
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
				fmt.Println("filling guild and chan arrys...")
				c.MyGuilds = parsed.Guilds
				c.MyChans = parsed.PrivateChannels
				fmt.Println("Arrays filled!")

			//TODO: add code differentiating between unavailable guild
			// messages and normal messages
			// (make new events)
			default:
				d, ok := msg.Data.(*json.RawMessage)
				if ok {
					//json.rawmessage, so unexpected type
					fmt.Print("unexpected ")
				}
				fmt.Printf("type read '%v':\n", msg.Type)
				if ok {
					fmt.Printf("%s\n\n", d)
				} else {
					fmt.Printf("%#v\n\n", msg.Data)
				}
			}
			
			call := *CB
			call(msg.Type, msg.Data)
		default:
			fmt.Printf("unexpected op '%v':\n%#v\n\n", msg.Op, msg.Data)
		}
	}
	fmt.Println("Reads closed, exiting process")
}
func (c Discord) WSInit(con *websocket.Conn, msgChan chan WSMsg) {
	//send init on wire
	p := Properties{
					 OS: "DiscordBot",
				Browser: "discord.go",
				 Device: "console",
			   Referrer: "",
		ReferringDomain: "",
	}
	msgData := INIT{
		Token: c.Token,
		Version: 3, // hard-coded so changes will only happen when coded in
		Properties: &p,
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
	go c.WSProcess(con, msgSend, msgRead, c.sigStop, c.sigSafe, call)
	return
}