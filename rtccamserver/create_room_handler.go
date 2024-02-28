package rtccamserver

import (
	"log"
	"rtccam/message"
	"rtccam/rtccamclient"
)

func CreateRoomUrlHandler(client *rtccamclient.RTCCamClient, createRoomUrlRequestMessage *message.CreateRoomIdRequestMessage) {
	waitList := GetCreateRoomWaitList()
	id := waitList.GenerateId()
	waitList.Add(id, createRoomUrlRequestMessage)

	err := client.Send(message.NewCreateRoomIdMessage(id))
	if err != nil {
		log.Println("[CreateRoomUrlHandler] ClientId:", client.ClientId, "Error:", err)
	}
}
