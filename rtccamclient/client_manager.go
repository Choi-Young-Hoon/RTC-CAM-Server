package rtccamclient

import "sync"

var defaultRTCCamClientManager = &RTCCamClientManager{
	clients: make(map[int64]*RTCCamClient),
}

func GetRTCCamClientManager() *RTCCamClientManager {
	return defaultRTCCamClientManager
}

type RTCCamClientManager struct {
	clientsMutex sync.Mutex
	clients      map[int64]*RTCCamClient
}

func (cm *RTCCamClientManager) AddClient(client *RTCCamClient) {
	cm.clientsMutex.Lock()
	defer cm.clientsMutex.Unlock()

	client.GenerateClientId()
	cm.clients[client.ClientId] = client
}

func (cm *RTCCamClientManager) RemoveClient(client *RTCCamClient) {
	cm.clientsMutex.Lock()
	defer cm.clientsMutex.Unlock()

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
