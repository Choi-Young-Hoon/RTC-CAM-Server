package rtccamserver

import (
	"errors"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccamlog"
	"rtccam/rtccammessage"
)

func SignalingRouteHander(client *rtccamclient.RTCCamClient, signalingRequestMessage *rtccammessage.SignalingMessage) {
	rtccamlog.Info().
		Any("ClientId", client.ClientId).
		Str("RequestType", signalingRequestMessage.RequestType).
		Send()

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		errorMessage := rtccammessage.NewRTCCamErrorMessage(err)
		rtccamlog.Error().
			Err(errors.New(err.Message)).
			Any("ClientId", client.ClientId).
			Int64("RoomId", client.JoinRoomId).
			Send()
		client.Send(errorMessage)
		return
	}

	responseClient, err := room.GetClient(signalingRequestMessage.ResponseClientId)
	if err != nil {
		rtccamlog.Error().
			Err(errors.New(err.Message)).
			Any("ClientId", client.ClientId).
			Any("ResponseClientId", signalingRequestMessage.ResponseClientId).
			Send()
		errorMessage := rtccammessage.NewRTCCamErrorMessage(err)
		client.Send(errorMessage)
		return
	}

	responseClient.Send(signalingRequestMessage)
}

func SignalingHandler(client *rtccamclient.RTCCamClient, signalingRequestMessage *rtccammessage.SignalingMessage) {
	switch signalingRequestMessage.RequestType {
	case rtccammessage.SignalingRequestTypeOffer:
		SignalingRouteHander(client, signalingRequestMessage)
	case rtccammessage.SignalingRequestTypeAnswer:
		SignalingRouteHander(client, signalingRequestMessage)
	case rtccammessage.SignalingRequestTypeCandidate:
		SignalingRouteHander(client, signalingRequestMessage)
	}
}
