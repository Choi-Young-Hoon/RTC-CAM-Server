package rtccamroom

import (
	"rtccam/rtccamclient"
	"rtccam/rtccammessage"
	"sync"
)

func NewRoom(id int, title string) *Room {
	return &Room{
		Id:      id,
		Title:   title,
		clients: make(map[int]*rtccamclient.RTCCamClient),
	}
}

type Room struct {
	Id    int    `json:"id"`
	Title string `json:"title"`

	clientMutex sync.Mutex
	clients     map[int]*rtccamclient.RTCCamClient
}

func (r *Room) AddClient(client *rtccamclient.RTCCamClient) {
	r.clientMutex.Lock()
	defer r.clientMutex.Unlock()

	client.JoinRoomId = r.Id
	r.clients[client.Id] = client
}

func (r *Room) RemoveClient(client *rtccamclient.RTCCamClient) {
	r.clientMutex.Lock()
	defer r.clientMutex.Unlock()

	if _, ok := r.clients[client.Id]; !ok {
		return
	}
	delete(r.clients, client.Id)
	client.JoinRoomId = -1
}

func (r *Room) Broadcast(requestClient *rtccamclient.RTCCamClient, message rtccammessage.Message) {
	for _, client := range r.clients {
		if client.Id == requestClient.Id {
			continue
		}

		go func(client *rtccamclient.RTCCamClient) {
			err := client.Send(message)
			if err != nil {
				r.RemoveClient(client)
			}
		}(client)
	}
}
