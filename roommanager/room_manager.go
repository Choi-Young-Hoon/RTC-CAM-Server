package roommanager

import (
	"rtccam/rtccamerrors"
	"rtccam/rtccamgen"
	"sync"
)

var defaultRoomManager = &RoomManager{
	idGenerator: rtccamgen.NewIDGenerator(),
	Rooms:       make(map[int64]*Room),
}

func GetRoomManager() *RoomManager {
	return defaultRoomManager
}

type RoomManager struct {
	idGenerator rtccamgen.Generator

	roomsMutex sync.Mutex
	Rooms      map[int64]*Room `json:"rooms"`
}

func (rm *RoomManager) CreatRoom(title, password string, maxClientCount int) *Room {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	room := NewRoom(title, password, maxClientCount)
	room.Id = rm.idGenerator.GenerateID()

	rm.Rooms[room.Id] = room

	return room
}

func (rm *RoomManager) RemoveRoom(room *Room) {
	rm.roomsMutex.Lock()
	defer rm.roomsMutex.Unlock()

	_, ok := rm.Rooms[room.Id]
	if !ok {
		return
	}
	delete(rm.Rooms, room.Id)

	rm.idGenerator.ReturnID(room.Id)
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
