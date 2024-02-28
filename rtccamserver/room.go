package rtccamserver

import (
	"fmt"
	"log"
	"rtccam/message"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccamerrors"
)

func NewRoomMessageDispatcher() *RoomMessageDispatcher {
	roomMessageDispatcher := &RoomMessageDispatcher{
		handles: make(map[string]RoomMessageHandle),
	}

	roomMessageDispatcher.AddHandleHandler(message.RoomRequestTypeCreateRoom, roomCreateHandler)
	roomMessageDispatcher.AddHandleHandler(message.RoomRequestTypeRoomList, roomListHandler)
	roomMessageDispatcher.AddHandleHandler(message.RoomRequestTypeJoinRoom, roomJoinHandler)
	roomMessageDispatcher.AddHandleHandler(message.RoomRequestTypeLeaveRoom, roomLeaveHandler)

	return roomMessageDispatcher
}

type RoomMessageHandle func(*rtccamclient.RTCCamClient, *message.RoomRequestMessage)

type RoomMessageDispatcher struct {
	handles map[string]RoomMessageHandle
}

func (r *RoomMessageDispatcher) AddHandleHandler(requestType string, handle RoomMessageHandle) {
	r.handles[requestType] = handle
}

func (r *RoomMessageDispatcher) RoomHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *message.RoomRequestMessage) {
	log.Println("[RoomHandler] ClientId:", client.ClientId, "RequestType:", roomRequestMessage.RequestType)

	handle, ok := r.handles[roomRequestMessage.RequestType]
	if !ok {
		log.Println("[RoomHandler] ClientId:", client.ClientId, "Error: Not Found RequestType:", roomRequestMessage.RequestType)
		client.Send(message.NewRTCCamErrorMessage(rtccamerrors.ErrorRequestTypeError.Error()))
		return
	}

	handle(client, roomRequestMessage)
}

func broadcastRoomList() {
	roomManager := roommanager.GetRoomManager()
	roomListMessage := message.NewRTCCamRoomListMessage(roomManager)

	clientManager := rtccamclient.GetRTCCamClientManager()
	clientManager.Broadcast(roomListMessage)
}

func roomCreateHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *message.RoomRequestMessage) {
	log.Println("[roomCreateHandler] Start ClientId:", client.ClientId, "CreateRoomId:", roomRequestMessage.CreateRoomId)
	waitList := GetCreateRoomWaitList()
	waitId, err := waitList.Get(roomRequestMessage.CreateRoomId)
	if err != nil {
		fmt.Println("[roomCreateHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(message.NewRTCCamErrorMessage(err.Error()))
		return
	}

	roomManager := roommanager.GetRoomManager()
	room := roommanager.NewRoom(waitId.RoomInfo.Title, roomRequestMessage.JoinPassword)
	roomManager.AddRoom(room)

	roomRequestMessage.JoinRoomId = room.Id
	roomJoinHandler(client, roomRequestMessage)
}

func roomListHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *message.RoomRequestMessage) {
	roomManager := roommanager.GetRoomManager()
	roomListMessage := message.NewRTCCamRoomListMessage(roomManager)
	err := client.Send(roomListMessage)
	if err != nil {
		log.Println("[roomListHandler] ClientId:", client.ClientId, "Error:", err)
		return
	}
}

func roomJoinHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *message.RoomRequestMessage) {
	if client.JoinRoomId == roomRequestMessage.JoinRoomId {
		broadcastRoomList()
		return
	}
	roomLeave(client)

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(roomRequestMessage.JoinRoomId)
	if err != nil {
		log.Println("[roomJoinHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(message.NewRTCCamErrorMessage(err.Error()))
		return
	}

	if room.IsPassword && room.Password != roomRequestMessage.JoinPassword {
		log.Println("[roomJoinHandler] ClientId:", client.ClientId, "Error: JoinPassword is incorrect")
		client.Send(message.NewRTCCamErrorMessage("JoinPassword is incorrect"))
		return
	}

	log.Println("[roomJoinHandler] ClientId:", client.ClientId, "JoinRoomId:", room.Id)
	room.JoinClient(client)
	broadcastRoomList()
}

func roomLeaveHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *message.RoomRequestMessage) {
	roomLeave(client)
	broadcastRoomList()
}

func roomLeave(client *rtccamclient.RTCCamClient) {
	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		if client.JoinRoomId == 0 {
			return
		}

		log.Println("[roomLeaveHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(message.NewRTCCamErrorMessage(err.Error()))
		return
	}

	log.Println("[roomLeaveHandler] ClientId:", client.ClientId, "LeaveRoomId:", room.Id)
	room.LeaveClient(client)
}
