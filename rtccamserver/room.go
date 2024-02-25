package rtccamserver

import (
	"log"
	"rtccam/message"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
)

func BroadcastRoomList() {
	roomManager := roommanager.GetRoomManager()
	roomListMessage := message.NewRTCCamRoomListMessage(roomManager)

	clientManager := rtccamclient.GetRTCCamClientManager()
	clientManager.Broadcast(roomListMessage)
}

func RoomListHandler(client *rtccamclient.RTCCamClient) {
	roomManager := roommanager.GetRoomManager()
	roomListMessage := message.NewRTCCamRoomListMessage(roomManager)
	err := client.Send(roomListMessage)
	if err != nil {
		log.Println("[RoomListHandler] ClientId:", client.ClientId, "Error:", err)
		return
	}
}

func RoomJoinHandler(client *rtccamclient.RTCCamClient, roomId int64) {
	roomLeave(client)

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(roomId)
	if err != nil {
		log.Println("[RoomJoinHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(message.NewRTCCamErrorMessage(err.Error()))
		return
	}

	log.Println("[RoomJoinHandler] ClientId:", client.ClientId, "JoinRoomId:", room.Id)
	room.JoinClient(client)
	BroadcastRoomList()
}

func RoomLeaveHandler(client *rtccamclient.RTCCamClient) {
	roomLeave(client)
	BroadcastRoomList()
}

func roomLeave(client *rtccamclient.RTCCamClient) {
	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		if client.JoinRoomId == 0 {
			return
		}

		log.Println("[RoomLeaveHandler] ClientId:", client.ClientId, "Error:", err)
		client.Send(message.NewRTCCamErrorMessage(err.Error()))
		return
	}

	log.Println("[RoomLeaveHandler] ClientId:", client.ClientId, "LeaveRoomId:", room.Id)
	room.LeaveClient(client)
}

func RoomHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *message.RoomRequestMessage) {
	log.Println("[RoomHandler] ClientId:", client.ClientId, "RequestType:", roomRequestMessage.RequestType)

	switch roomRequestMessage.RequestType {
	case message.RoomRequestTypeRoomList:
		RoomListHandler(client)
		break

	case message.RoomRequestTypeJoinRoom:
		RoomJoinHandler(client, roomRequestMessage.JoinRoomId)
		break

	case message.RoomRequestTypeLeaveRoom:
		RoomLeaveHandler(client)
		break
	}
}
