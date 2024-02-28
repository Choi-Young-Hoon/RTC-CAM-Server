package rtccamclient

import (
	"github.com/gorilla/websocket"
	"math"
	"rtccam/message"
	"sync"
)

var nextClientIdMutex sync.Mutex
var nextClientId int64 = 0

func NewRTCCamClient(ws *websocket.Conn) *RTCCamClient {
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

func (c *RTCCamClient) Recv() (*message.RTCCamRequestMessage, error) {
	rtcCamRequestMessage := message.NewRTCCamRequestMessage()
	err := c.ws.ReadJSON(rtcCamRequestMessage)
	if err != nil {
		return nil, err
	}

	return rtcCamRequestMessage, nil
}

func (c *RTCCamClient) Close() {
	c.ws.Close()
}

func (c *RTCCamClient) GenerateClientId() int64 {
	nextClientIdMutex.Lock()
	defer nextClientIdMutex.Unlock()

	if nextClientId == math.MaxInt64 {
		nextClientId = 0
	}

	nextClientId++
	c.ClientId = nextClientId

	return c.ClientId
}
