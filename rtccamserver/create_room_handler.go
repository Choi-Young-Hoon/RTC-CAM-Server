package rtccamserver

import (
	"log"
	"rtccam/message"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
)

func CreateRoomHandler(client *rtccamclient.RTCCamClient, createRoomRequestMessage *message.CreateRoomRequestMessage) {
	roomManager := roommanager.GetRoomManager()
	room := roommanager.NewRoom(createRoomRequestMessage.Title, createRoomRequestMessage.Password)
	authToken := room.GenerateAuthToken()
	roomManager.AddRoom(room)

	roomListMessage := message.NewRTCCamRoomListMessage(roomManager)
	rtccamclient.GetRTCCamClientManager().Broadcast(roomListMessage)

	err := client.Send(message.NewCreateRoomMessage(room.Id, authToken))
	if err != nil {
		log.Println("[CreateRoomHandler] ClientId:", client.ClientId, "Error:", err)
	}
}
