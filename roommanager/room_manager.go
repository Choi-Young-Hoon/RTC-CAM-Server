package roommanager

import (
	"rtccam/rtccamerrors"
	"sync"
)

var defaultRoomManager = &RoomManager{
	Rooms: make(map[int64]*Room),
}

func GetRoomManager() *RoomManager {
	return defaultRoomManager
}

type RoomManager struct {
	roomsMutex sync.Mutex
	Rooms      map[int64]*Room `json:"rooms"`
}

func (rm *RoomManager) AddRoom(room *Room) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	rm.Rooms[room.Id] = room
}

func (rm *RoomManager) RemoveRoom(room *Room) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	_, ok := rm.Rooms[room.Id]
	if !ok {
		return
	}
	delete(rm.Rooms, room.Id)
}

func (rm *RoomManager) GetRoom(roomId int64) (*Room, error) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	room, ok := rm.Rooms[roomId]
	if !ok {
		return nil, rtccamerrors.ErrorRoomNotFound
	}
	return room, nil
}
