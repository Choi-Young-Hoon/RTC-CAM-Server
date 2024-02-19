package room

import (
	"rtccam/message"
	"rtccam/rtccam"
	"sync"
)

func NewRoom(id int, title string) *Room {
	return &Room{
		Id:      id,
		Title:   title,
		clients: make(map[int]*rtccam.RTCCamClient),
	}
}

type Room struct {
	Id    int    `json:"id"`
	Title string `json:"title"`

	clientMutex sync.Mutex
	clients     map[int]*rtccam.RTCCamClient
}

func (r *Room) AddClient(client *rtccam.RTCCamClient) {
	r.clients[client.Id] = client
}

func (r *Room) RemoveClient(client *rtccam.RTCCamClient) {
	r.clientMutex.Lock()
	defer r.clientMutex.Unlock()

	if _, ok := r.clients[client.Id]; !ok {
		return
	}
	delete(r.clients, client.Id)
}

func (r *Room) Broadcast(requestClient *rtccam.RTCCamClient, message message.Message) {
	for _, client := range r.clients {
		if client.Id == requestClient.Id {
			continue
		}

		go func(client *rtccam.RTCCamClient) {
			err := client.Send(message)
			if err != nil {
				r.RemoveClient(client)
			}
		}(client)
	}
}
