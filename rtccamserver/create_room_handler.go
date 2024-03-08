package rtccamserver

import (
	"errors"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccamerrors"
	"rtccam/rtccamlog"
	"rtccam/rtccammessage"
)

func CreateRoomHandler(client *rtccamclient.RTCCamClient, createRoomRequestMessage *rtccammessage.CreateRoomRequestMessage) {
	rtccamlog.Info().Msg("Creat Room Start")
	if createRoomRequestMessage.MaxClientCount <= 0 && createRoomRequestMessage.MaxClientCount > 10 {
		countError := rtccamerrors.NewInvalidMaxClientCount()
		rtccamlog.Error().
			Err(errors.New(countError.Message)).
			Any("ClientId", client.ClientId).
			Any("MaxClientCount", createRoomRequestMessage.MaxClientCount).
			Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(countError))
		return
	}

	roomManager := roommanager.GetRoomManager()
	room := roomManager.CreatRoom(createRoomRequestMessage.Title, createRoomRequestMessage.Password, createRoomRequestMessage.MaxClientCount)
	authToken := room.GenerateAuthToken()

	roomListMessage := rtccammessage.NewRTCCamRoomListMessage(roomManager)
	rtccamclient.GetRTCCamClientManager().Broadcast(roomListMessage)

	client.Send(rtccammessage.NewCreateRoomMessage(room.Id, authToken))
}
