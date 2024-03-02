package rtccamclient

import (
	"github.com/gorilla/websocket"
	"rtccam/rtccamgen"
	"sync"
)

var defaultRTCCamClientManager = &RTCCamClientManager{
	clients:     make(map[int64]*RTCCamClient),
	idGenerator: rtccamgen.NewIDGenerator(),
}

func GetRTCCamClientManager() *RTCCamClientManager {
	return defaultRTCCamClientManager
}

type RTCCamClientManager struct {
	idGenerator rtccamgen.Generator

	clientsMutex sync.Mutex
	clients      map[int64]*RTCCamClient
}

func (cm *RTCCamClientManager) CreateClient(ws *websocket.Conn) *RTCCamClient {
	cm.clientsMutex.Lock()
	defer cm.clientsMutex.Unlock()

	// ID 생성
	client := NewRTCCamClient(ws)
	client.ClientId = cm.idGenerator.GenerateID()

	cm.clients[client.ClientId] = client

	return client
}

func (cm *RTCCamClientManager) RemoveClient(client *RTCCamClient) {
	cm.clientsMutex.Lock()
	defer cm.clientsMutex.Unlock()

	// 반납
	cm.idGenerator.ReturnID(client.ClientId)

	_, ok := cm.clients[client.ClientId]
	if !ok {
		return
	}
	delete(cm.clients, client.ClientId)
}

func (cm *RTCCamClientManager) Broadcast(message interface{}) {
	cm.clientsMutex.Lock()
	defer cm.clientsMutex.Unlock()

	for _, client := range cm.clients {
		client.Send(message)
	}
}

func (cm *RTCCamClientManager) CloseAll() {
	cm.clientsMutex.Lock()
	defer cm.clientsMutex.Unlock()

	wg := sync.WaitGroup{}
	wg.Add(len(cm.clients))

	for _, client := range cm.clients {
		go func(client *RTCCamClient) {
			defer wg.Done()
			client.Close()
		}(client)
	}

	wg.Wait()
}
