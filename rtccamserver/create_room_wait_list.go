package rtccamserver

import (
	"rtccam/message"
	"rtccam/rtccamerrors"
	"sync"
)

func NewCreateRoomWaitItem(id int64, createRoomIdRequestMessage *message.CreateRoomIdRequestMessage) CreateRoomWaitItem {
	return CreateRoomWaitItem{
		Id:       id,
		RoomInfo: createRoomIdRequestMessage,
	}
}

type CreateRoomWaitItem struct {
	Id       int64
	RoomInfo *message.CreateRoomIdRequestMessage
}

func NewCreateRoomWaitList() CreateRoomWaitList {
	return CreateRoomWaitList{
		items: []CreateRoomWaitItem{},
	}
}

var defaultCreateRoomWaitList = NewCreateRoomWaitList()

func GetCreateRoomWaitList() *CreateRoomWaitList {
	return &defaultCreateRoomWaitList
}

type CreateRoomWaitList struct {
	itemsMutex sync.Mutex
	items      []CreateRoomWaitItem
}

func (l *CreateRoomWaitList) Add(id int64, createRoomIdRequestMessage *message.CreateRoomIdRequestMessage) {
	l.itemsMutex.Lock()
	defer l.itemsMutex.Unlock()

	l.items = append(l.items, NewCreateRoomWaitItem(id, createRoomIdRequestMessage))
}

func (l *CreateRoomWaitList) remove(id int64) {
	for i, item := range l.items {
		if item.Id == id {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}
}

func (l *CreateRoomWaitList) Get(id int64) (*CreateRoomWaitItem, error) {
	l.itemsMutex.Lock()
	defer l.itemsMutex.Unlock()

	for _, item := range l.items {
		if item.Id == id {
			l.remove(id)
			return &item, nil
		}
	}

	return nil, rtccamerrors.ErrorNotFoundCreateRoomWaitItem
}

var createRoomWaitListIdMutex sync.Mutex
var createRoomWaitListId int64 = 0

func (l *CreateRoomWaitList) GenerateId() int64 {
	createRoomWaitListIdMutex.Lock()
	defer createRoomWaitListIdMutex.Unlock()

	createRoomWaitListId++
	return createRoomWaitListId
}
