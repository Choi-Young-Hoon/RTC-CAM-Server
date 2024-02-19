package rtccamroom

import (
	"errors"
	"rtccam/rtccamclient"
	"sync"
)

// 싱글턴
var roomInstance = NewRoomManager()

func GetRoomManager() *RoomManager {
	return roomInstance
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		roomNumber: 0,
		Rooms:      make(map[int]*Room),
	}
}

type RoomManager struct {
	roomNumber int

	roomsMutex sync.Mutex
	Rooms      map[int]*Room `json:"rooms"`
}

func (rm *RoomManager) CreateRoom(title string) error {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	if _, ok := rm.Rooms[rm.roomNumber]; ok {
		return errors.New("Room already exists")
	}

	room := NewRoom(rm.roomNumber, title)
	rm.Rooms[room.Id] = room
	rm.roomNumber++

	return nil
}

func (rm *RoomManager) DeleteRoom(id int) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	if _, ok := rm.Rooms[id]; !ok {
		return
	}
	delete(rm.Rooms, id)
}

func (rm *RoomManager) RemoveClientFromRoom(client *rtccamclient.RTCCamClient) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	if room, ok := rm.Rooms[client.JoinRoomId]; !ok {
		room.RemoveClient(client)
	}
}

func (rm *RoomManager) GetRoom(id int) (*Room, error) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	if value, ok := rm.Rooms[id]; ok {
		return value, nil
	}

	return nil, errors.New("Room not found")
}
