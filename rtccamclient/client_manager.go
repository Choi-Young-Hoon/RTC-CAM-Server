package rtccamclient

import (
	"math"
	"sync"
	"sync/atomic"
)

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

	client.ClientId = cm.generateClientId()
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

var nextClientId int64 = 0

func (c *RTCCamClientManager) generateClientId() int64 {
	if nextClientId == math.MaxInt64 {
		atomic.StoreInt64(&nextClientId, 0)
	}
	atomic.AddInt64(&nextClientId, 1)

	return nextClientId
}
