package rtccamclient

import (
	"context"
	"github.com/gorilla/websocket"
	"rtccam/rtccamlog"
	"rtccam/rtccammessage"
	"sync"
)

func NewRTCCamClient(ws *websocket.Conn) *RTCCamClient {
	client := &RTCCamClient{
		ws: ws,
	}

	client.ctx, client.cancel = context.WithCancel(context.Background())
	client.channel = make(chan interface{})
	client.wg.Add(1)
	go client.sender()

	return client
}

type RTCCamClient struct {
	ClientId   int64 `json:"client_id"`
	JoinRoomId int64 `json:"-"`

	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	channel chan interface{}
	ws      *websocket.Conn
}

var sendMutex sync.Mutex

func (c *RTCCamClient) Send(message interface{}) {
	c.channel <- message
}

func (c *RTCCamClient) Recv() (*rtccammessage.RTCCamRequestMessage, error) {
	rtcCamRequestMessage := rtccammessage.NewRTCCamRequestMessage()
	err := c.ws.ReadJSON(rtcCamRequestMessage)
	if err != nil {
		return nil, err
	}

	return rtcCamRequestMessage, nil
}

func (c *RTCCamClient) Close() {
	c.cancel()
	c.wg.Wait()
	_ = c.ws.Close()
}

func (c *RTCCamClient) sender() {
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		case message := <-c.channel:
			err := c.ws.WriteJSON(message)
			if err != nil {
				rtccamlog.Error().Err(err).Any("ClientId", c.ClientId).Send()
				continue
			}
		}
	}
}
