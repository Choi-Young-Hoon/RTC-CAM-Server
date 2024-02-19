package rtccamclient

import (
	"github.com/gorilla/websocket"
	"rtccam/rtccammessage"
)

var nextClientId int = 0

func NewRTCCamClient(ws *websocket.Conn) *RTCCamClient {
	nextClientId++
	return &RTCCamClient{
		Id:         nextClientId,
		Websocket:  ws,
		JoinRoomId: -1,
	}
}

type RTCCamClient struct {
	Id         int
	JoinRoomId int
	Websocket  *websocket.Conn
}

func (r *RTCCamClient) SendSuccessMessage() error {
	message := rtccammessage.NewSuccessMessage(r.Id)
	err := r.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func (r *RTCCamClient) SendErrorMessage(errorMessage string) error {
	message := rtccammessage.NewErrorMessage(r.Id, errorMessage)
	err := r.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func (r *RTCCamClient) Send(message rtccammessage.Message) error {
	message.ClientId = r.Id
	err := r.Websocket.WriteJSON(message)
	if err != nil {
		return err
	}
	return nil
}

func (r *RTCCamClient) Recv(message *rtccammessage.Message) error {
	message.ClientId = r.Id
	err := r.Websocket.ReadJSON(&message)
	if err != nil {
		return err
	}
	return nil
}
