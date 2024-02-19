package rtccam

import (
	"github.com/gorilla/websocket"
	"rtccam/message"
)

var nextClientId int = 0

func NewRTCCamClient(ws *websocket.Conn) *RTCCamClient {
	nextClientId++
	return &RTCCamClient{
		Id:        nextClientId,
		Websocket: ws,
	}
}

type RTCCamClient struct {
	Id        int
	Websocket *websocket.Conn
}

func (r *RTCCamClient) Send(message message.Message) error {
	err := r.Websocket.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}

func (r *RTCCamClient) Recv(message *message.Message) error {
	err := r.Websocket.ReadJSON(&message)
	if err != nil {
		return err
	}
	return nil
}
