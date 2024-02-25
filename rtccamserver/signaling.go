package rtccamserver

import (
	"rtccam/message"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
)

func SignalingRouteHander(client *rtccamclient.RTCCamClient, signalingRequestMessage *message.SignalingMessage) {
	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		errorMessage := message.NewRTCCamErrorMessage(err.Error())
		client.Send(errorMessage)
		return
	}

	responseClient, err := room.GetClient(signalingRequestMessage.ResponseClientId)
	if err != nil {
		errorMessage := message.NewRTCCamErrorMessage(err.Error())
		client.Send(errorMessage)
		return
	}

	responseClient.Send(signalingRequestMessage)
}

func SignalingHandler(client *rtccamclient.RTCCamClient, signalingRequestMessage *message.SignalingMessage) {
	switch signalingRequestMessage.RequestType {
	case message.SignalingRequestTypeOffer:
		SignalingRouteHander(client, signalingRequestMessage)
		break
	case message.SignalingRequestTypeAnswer:
		SignalingRouteHander(client, signalingRequestMessage)
		break
	case message.SignalingRequestTypeCandidate:
		SignalingRouteHander(client, signalingRequestMessage)
		break
	}
}
