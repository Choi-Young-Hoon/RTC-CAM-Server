package rtccamserver

import (
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccamerrors"
	"rtccam/rtccamlog"
	"rtccam/rtccammessage"
)

func CreateRoomHandler(client *rtccamclient.RTCCamClient, createRoomRequestMessage *rtccammessage.CreateRoomRequestMessage) {
	if createRoomRequestMessage.MaxClientCount <= 0 && createRoomRequestMessage.MaxClientCount > 10 {
		rtccamlog.Error().
			Err(rtccamerrors.ErrorInvalidMaxClientCount).
			Any("ClientId", client.ClientId).
			Any("MaxClientCount", createRoomRequestMessage.MaxClientCount).
			Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(rtccamerrors.ErrorInvalidMaxClientCount.Error()))
		return
	}

	roomManager := roommanager.GetRoomManager()
	room := roommanager.NewRoom(createRoomRequestMessage.Title, createRoomRequestMessage.Password, createRoomRequestMessage.MaxClientCount)
	authToken := room.GenerateAuthToken()
	roomManager.AddRoom(room)

	roomListMessage := rtccammessage.NewRTCCamRoomListMessage(roomManager)
	rtccamclient.GetRTCCamClientManager().Broadcast(roomListMessage)

	client.Send(rtccammessage.NewCreateRoomMessage(room.Id, authToken))
}
