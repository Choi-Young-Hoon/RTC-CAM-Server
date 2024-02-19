package room

import "errors"

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
	Rooms      map[int]*Room `json:"rooms"`
}

func (rm *RoomManager) CreateRoom(title string) error {
	if _, ok := rm.Rooms[rm.roomNumber]; ok {
		return errors.New("Room already exists")
	}

	room := NewRoom(rm.roomNumber, title)
	rm.Rooms[room.Id] = room
	rm.roomNumber++

	return nil
}

func (rm *RoomManager) DeleteRoom(id int) {
	if _, ok := rm.Rooms[id]; !ok {
		return
	}
	delete(rm.Rooms, id)
}

func (rm *RoomManager) GetRoom(id int) (*Room, error) {
	if value, ok := rm.Rooms[id]; ok {
		return value, nil
	}

	return nil, errors.New("Room not found")
}
