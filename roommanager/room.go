package roommanager

import (
	"errors"
	"rtccam/message"
	"rtccam/rtccamclient"
	"sync"
)

var nextRoomIdMutex sync.Mutex
var nextRoomId int64 = 0

func NewRoom(title, password string) *Room {
	nextRoomIdMutex.Lock()
	defer nextRoomIdMutex.Unlock()

	nextRoomId++
	return &Room{
		Id:         nextRoomId,
		Title:      title,
		IsPassword: password != "",
		Password:   password,
		Clients:    make(map[int64]*rtccamclient.RTCCamClient),
	}
}

type Room struct {
	Id int64 `json:"id"`

	Title      string `json:"title"`
	IsPassword bool   `json:"is_password"`
	Password   string `json:"-"`

	clientsMutex sync.Mutex
	Clients      map[int64]*rtccamclient.RTCCamClient `json:"clients"`
}

func (r *Room) JoinClient(client *rtccamclient.RTCCamClient) {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	client.JoinRoomId = r.Id

	r.Clients[client.ClientId] = client

	joinSuccessMessage := message.NewRTCCamJoinSuccessMessage()
	client.Send(joinSuccessMessage)
}

func (r *Room) LeaveClient(client *rtccamclient.RTCCamClient) {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	_, ok := r.Clients[client.ClientId]
	if !ok {
		return
	}

	client.JoinRoomId = 0
	delete(r.Clients, client.ClientId)

	if len(r.Clients) == 0 {
		roomManager := GetRoomManager()
		roomManager.RemoveRoom(r)
	}

	leaveMessage := message.NewRTCCamLeaveMessage(client.ClientId)
	for _, client := range r.Clients {
		client.Send(leaveMessage)
	}
}

func (r *Room) GetClient(clientId int64) (*rtccamclient.RTCCamClient, error) {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	client, ok := r.Clients[clientId]
	if !ok {
		return nil, errors.New("Client not found")
	}
	return client, nil
}
