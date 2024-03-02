package rtccamserver

import (
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccammessage"
)

func CreateRoomHandler(client *rtccamclient.RTCCamClient, createRoomRequestMessage *rtccammessage.CreateRoomRequestMessage) {
	roomManager := roommanager.GetRoomManager()
	room := roommanager.NewRoom(createRoomRequestMessage.Title, createRoomRequestMessage.Password)
	authToken := room.GenerateAuthToken()
	roomManager.AddRoom(room)

	roomListMessage := rtccammessage.NewRTCCamRoomListMessage(roomManager)
	rtccamclient.GetRTCCamClientManager().Broadcast(roomListMessage)

	client.Send(rtccammessage.NewCreateRoomMessage(room.Id, authToken))
}
