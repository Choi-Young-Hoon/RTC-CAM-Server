package rtccamclient

import (
	"github.com/gorilla/websocket"
	"rtccam/message"
	"sync"
)

var nextClientIdMutex sync.Mutex
var nextClientId int64 = 0

func NewRTCCamClient(ws *websocket.Conn) *RTCCamClient {
	nextClientIdMutex.Lock()
	defer nextClientIdMutex.Unlock()

	nextClientId++
	return &RTCCamClient{
		ClientId: nextClientId,
		ws:       ws,
	}
}

type RTCCamClient struct {
	ClientId   int64 `json:"client_id"`
	JoinRoomId int64 `json:"-"`
	ws         *websocket.Conn
}

func (c *RTCCamClient) Send(message interface{}) error {
	return c.ws.WriteJSON(message)
}

func (c *RTCCamClient) Recv() (*message.RTCCamMessage, error) {
	rtcCamRequestMessage := message.NewRTCCamMessage()
	err := c.ws.ReadJSON(rtcCamRequestMessage)
	if err != nil {
		return nil, err
	}

	return rtcCamRequestMessage, nil
}

func (c *RTCCamClient) Close() {
	c.ws.Close()
}
