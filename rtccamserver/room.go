package rtccamserver

import (
	"log"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccamerrors"
	"rtccam/rtccammessage"
)

var defaultRoomMessageDispatcher = NewRoomMessageDispatcher()

func GetRoomMessageDispatcher() *RoomMessageDispatcher {
	return defaultRoomMessageDispatcher
}

func NewRoomMessageDispatcher() *RoomMessageDispatcher {
	roomMessageDispatcher := &RoomMessageDispatcher{
		handles: make(map[string]RoomMessageHandle),
	}

	roomMessageDispatcher.AddHandleHandler(rtccammessage.RoomRequestTypeRoomList, roomListHandler)
	roomMessageDispatcher.AddHandleHandler(rtccammessage.RoomRequestTypeJoinRoom, roomJoinHandler)
	roomMessageDispatcher.AddHandleHandler(rtccammessage.RoomRequestTypeLeaveRoom, roomLeaveHandler)
	roomMessageDispatcher.AddHandleHandler(rtccammessage.RoomRequestAuthToken, roomAuthTokenHandler)

	return roomMessageDispatcher
}

type RoomMessageHandle func(*rtccamclient.RTCCamClient, *rtccammessage.RoomRequestMessage)

type RoomMessageDispatcher struct {
	handles map[string]RoomMessageHandle
}

func (r *RoomMessageDispatcher) AddHandleHandler(requestType string, handle RoomMessageHandle) {
	r.handles[requestType] = handle
}

func (r *RoomMessageDispatcher) RoomHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	log.Println("[RoomHandler] ClientId:", client.ClientId, "RequestType:", roomRequestMessage.RequestType)

	handle, ok := r.handles[roomRequestMessage.RequestType]
	if !ok {
		log.Println("[RoomHandler] ClientId:", client.ClientId, "Error: Not Found RequestType:", roomRequestMessage.RequestType)
		client.Send(rtccammessage.NewRTCCamErrorMessage(rtccamerrors.ErrorRequestTypeError.Error()))
		return
	}

	handle(client, roomRequestMessage)
}

func broadcastRoomList() {
	roomManager := roommanager.GetRoomManager()
	roomListMessage := rtccammessage.NewRTCCamRoomListMessage(roomManager)

	clientManager := rtccamclient.GetRTCCamClientManager()
	clientManager.Broadcast(roomListMessage)
}

func roomListHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	roomManager := roommanager.GetRoomManager()
	roomListMessage := rtccammessage.NewRTCCamRoomListMessage(roomManager)
	client.Send(roomListMessage)
}

func roomJoinHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	if client.JoinRoomId == roomRequestMessage.RoomId {
		broadcastRoomList()
		return
	}
	roomLeave(client)

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(roomRequestMessage.RoomId)
	if err != nil {
		log.Println("[roomJoinHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(rtccammessage.NewRTCCamErrorMessage(err.Error()))
		return
	}

	if room.Authenticate(roomRequestMessage.AuthToken) == false {
		log.Println("[roomJoinHandler] ClientId:", client.ClientId, "Error: Invalid AuthToken")
		client.Send(rtccammessage.NewRTCCamErrorMessage("Invalid AuthToken"))
		return
	}

	log.Println("[roomJoinHandler] ClientId:", client.ClientId, "RoomId:", room.Id)
	room.JoinClient(client)
	broadcastRoomList()
}

func roomLeaveHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	roomLeave(client)
	broadcastRoomList()
}

func roomAuthTokenHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	log.Println("[roomAuthTokenHandler] ClientId:", client.ClientId, "JoinRoomId:", roomRequestMessage.RoomId)

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(roomRequestMessage.RoomId)
	if err != nil {
		log.Println("[roomAuthTokenHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(rtccammessage.NewRTCCamErrorMessage(err.Error()))
		return
	}

	if room.IsPassword && room.Password != roomRequestMessage.Password {
		log.Println("[roomAuthTokenHandler] ClientId:", client.ClientId, "Error: Invalid Password")
		client.Send(rtccammessage.NewRTCCamErrorMessage("Invalid Password"))
		return
	}

	authToken := room.GenerateAuthToken()
	client.Send(rtccammessage.NewRTCCamAuthTokenMessage(authToken, room))
}

func roomLeave(client *rtccamclient.RTCCamClient) {
	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		if client.JoinRoomId == 0 {
			return
		}

		log.Println("[roomLeaveHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(rtccammessage.NewRTCCamErrorMessage(err.Error()))
		return
	}

	log.Println("[roomLeaveHandler] ClientId:", client.ClientId, "LeaveRoomId:", room.Id)
	room.LeaveClient(client)
}
