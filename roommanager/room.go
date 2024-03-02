package roommanager

import (
	"math/rand"
	"rtccam/rtccamclient"
	"rtccam/rtccamerrors"
	"rtccam/rtccammessage"
	"sync"
)

func NewRoom(title, password string, maxClientCount int) *Room {
	return &Room{
		Title:          title,
		IsPassword:     password != "",
		Password:       password,
		MaxClientCount: maxClientCount,
		Clients:        make(map[int64]*rtccamclient.RTCCamClient),
		AuthTokens:     make(map[string]int),
	}
}

type Room struct {
	Id int64 `json:"id"`

	Title          string `json:"title"`
	IsPassword     bool   `json:"is_password"`
	Password       string `json:"-"`
	MaxClientCount int    `json:"max_client_count"`

	clientsMutex sync.Mutex
	Clients      map[int64]*rtccamclient.RTCCamClient `json:"clients"`

	deleteAuthTokenMutex sync.Mutex
	AuthTokens           map[string]int `json:"-"`
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
			r.deleteAuthToken(authToken)
		}
		return true
	}
	return false
}

func (r *Room) deleteAuthToken(authToken string) {
	r.deleteAuthTokenMutex.Lock()
	defer r.deleteAuthTokenMutex.Unlock()
	delete(r.AuthTokens, authToken)
}

func (r *Room) JoinClient(client *rtccamclient.RTCCamClient) error {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	client.JoinRoomId = r.Id

	r.Clients[client.ClientId] = client

	joinSuccessMessage := rtccammessage.NewRTCCamJoinSuccessMessage(r, client.ClientId)
	client.Send(joinSuccessMessage)

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

	leaveMessage := rtccammessage.NewRTCCamLeaveMessage(client.ClientId)
	for _, client := range r.Clients {
		client.Send(leaveMessage)
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

func (r *Room) GetClientCount() int {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	return len(r.Clients)
}

func (r *Room) Broadcast(message interface{}) {
	r.clientsMutex.Lock()
	defer r.clientsMutex.Unlock()

	for _, client := range r.Clients {
		client.Send(message)
	}
}
