package roommanager

import (
	"math/rand"
	"rtccam/message"
	"rtccam/rtccamclient"
	"rtccam/rtccamerrors"
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
		AuthTokens: make(map[string]int),
	}
}

type Room struct {
	Id int64 `json:"id"`

	Title      string `json:"title"`
	IsPassword bool   `json:"is_password"`
	Password   string `json:"-"`

	clientsMutex sync.Mutex
	Clients      map[int64]*rtccamclient.RTCCamClient `json:"clients"`

	AuthTokens map[string]int `json:"-"`
}

func (r *Room) GenerateAuthToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	authToken := make([]byte, 20)
	for i := range authToken {
		authToken[i] = charset[rand.Intn(len(charset))]
	}

	r.AuthTokens[string(authToken)]++

	return string(authToken)
}

func (r *Room) Authenticate(authToken string) bool {
	value, ok := r.AuthTokens[authToken]
	if ok && value > 0 {
		r.AuthTokens[authToken]--
		if r.AuthTokens[authToken] <= 0 {
			delete(r.AuthTokens, authToken)
		}
		return true
	}
	return false
}

func (r *Room) JoinClient(client *rtccamclient.RTCCamClient) error {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	client.JoinRoomId = r.Id

	r.Clients[client.ClientId] = client

	joinSuccessMessage := message.NewRTCCamJoinSuccessMessage(r, client.ClientId)
	err := client.Send(joinSuccessMessage)
	if err != nil {
		return err
	}

	return nil
}

func (r *Room) LeaveClient(client *rtccamclient.RTCCamClient) error {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	_, ok := r.Clients[client.ClientId]
	if !ok {
		return nil
	}

	client.JoinRoomId = 0
	delete(r.Clients, client.ClientId)

	if len(r.Clients) == 0 {
		roomManager := GetRoomManager()
		roomManager.RemoveRoom(r)
	}

	leaveMessage := message.NewRTCCamLeaveMessage(client.ClientId)
	for _, client := range r.Clients {
		err := client.Send(leaveMessage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Room) GetClient(clientId int64) (*rtccamclient.RTCCamClient, error) {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	client, ok := r.Clients[clientId]
	if !ok {
		return nil, rtccamerrors.ErrorClientNotFound
	}
	return client, nil
}
