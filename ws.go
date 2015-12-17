package discord

import (
	"fmt"
	"time"
	"github.com/gorilla/websocket"
	"encoding/json"
)

type State struct{
	MentionCount  int    `json:"mention_count"`
	LastMessageId string `json:"last_message_id,omitempty"`
	ID            string `json:"id"`
}
type WSVoiceStates struct{
	UserID    string `json:"user_id"`
	Token     string `json:"token"`
	Suppress  bool   `json:"suppress"`
	SessionID string `json:"session_id"`
	SelfMute  bool   `json:"self_mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	Mute      bool   `json:"mute"`
	Deaf      bool   `json:"deaf"`
	ChannelID string `json:"channel_id"`
}
type WSPres struct{
	User   struct{
		ID string `json:"id"`
	} `json:"user"`
	Status string  `json:"status"`
	GameID *int    `json:"game_id"`
}
type WSGuilds struct{
	VoiceStates  []WSVoiceStates `json:"voice_states"`
	Roles        []Role          `json:"roles"`
	Region       string          `json:"region"`
	Presences    []WSPres        `json:"presences"`
	OwnerID      string          `json:"owner_id"`
	Name         string          `json:"name"`
	//Large        bool            `json:"large"`
	Members      []Member        `json:"members"`
	JoinedAt     time.Time       `json:"joined_at"`
	ID           string          `json:"id"`
	Icon         *string         `json:"icon"`
	Channels     []Channel       `json:"channels"`
	AfkTimeout   uint64          `json:"afk_timeout"`
	AfkChannelID *string         `json:"afk_channel_id"`
}

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
	if err = json.Unmarshal(raw, &msg); err != nil {
		fmt.Println(err)
	}
	switch msg.Type {
	case "READY":
		data := READY{}
		if err = json.Unmarshal(rawData, &data); err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("READY message: %s\n", msg.Data)
		msg.Data = data
		//fmt.Printf("READY message: %#v\n", msg)
	default:
		fmt.Printf("unknown message type: %q", msg.Type)
	}
	*m = WSMsg(msg)
	return

/*////////////////////

	type wsMsg WSMsg
	msg := wsMsg(m)
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return
	}
	switch msg.Type {
	case "READY":
		dat := READY{}
		err := json.Unmarshal(data, &dat)
		msg.Data = dat
	
	}
////////////////////*/
}

type READY struct{ // op from server (0)
	Version           int        `json:"v"`
	User              User       `json:"user"`
	SessionId         string     `json:"session_id"`
	ReadState         []State    `json:"read_state"`
	PrivateChannels   []Channel  `json:"private_channels"`
	HeartbeatInterval uint64     `json:"heartbeat_interval"`
	Guilds            []WSGuilds `json:"guilds"`
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

func wsSend(con *websocket.Conn, msgSend chan WSMsg, stopWS, exit chan int) {
	var nextMsg WSMsg
	for {
		select {
		case <-stopWS:
			//send close message
			fmt.Println("got stop mssage")
			fmt.Println("sending close frame (send)")
			err := con.WriteControl(websocket.CloseMessage, nil, time.Now().Add(3*time.Second))
			if err != nil {
				fmt.Println(err)
				exit <- 0
			}
			break
		case nextMsg = <-msgSend:
			//send the message on the channel to the connection
			
			
			fmt.Println("sending msg", nextMsg.Type)
			j, _ := json.Marshal(nextMsg)
			fmt.Printf("msg sent: `%s`\n", j)
			if err := con.WriteJSON(&nextMsg); err != nil {
				fmt.Println("wsSend:",err)
			}
		}
	}
}

func wsRead(con *websocket.Conn, msgRead chan WSMsg, exit chan<- int) {
	var nextMsg WSMsg
	for {
		//read the next message, put it on the channel
		err := con.ReadJSON(&nextMsg)
		if err != nil {
			fmt.Println("wsRead:",err)
			//close frame.  send then exit.
			
			fmt.Println("sending close frame (read)")
			con.WriteControl(websocket.CloseMessage, nil, time.Now())
			close(msgRead)
			exit <- 0
			break
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

func (c Discord) WSProcess(con *websocket.Conn, msgSend, msgRead chan WSMsg, stopWS, exit chan int) {
	defer con.Close()

	//process events until a close message is encountered, or network error occurs.
	
	//Close frames must be sent and recieved.
	//A close frame sent requires waiting for recieving before closing
	//a close frame recieved requires sending, then closing
	//Gorrilla handles close frames by returning an error (along with the frame read)
	fmt.Println("starting sender")
	go wsRead(con, msgRead, exit) // if we err on read, we have to send close frame then exit.
	fmt.Println("starting reader")
	go wsSend(con, msgSend, stopWS, exit) // if we send close frame, we have to wait for a response
	fmt.Println("starting process")
	for msg := range msgRead {
		//process messages
		fmt.Println("processing message")
		switch msg.Op {
		case 0:
			//default, most
			switch msg.Type {
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
				
			default:
				fmt.Printf("unexpected type '%v':\n%s\n", msg.Type, msg.Data)
			}
		default:
			fmt.Printf("unexpected op '%v':\n%s\n", msg.Op, msg.Data)
		}
	}
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

func (c Discord) WSConnect(stopWS, safe chan int) (err error) {
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
	go c.WSProcess(con, msgSend, msgRead, stopWS, safe)
	return
}