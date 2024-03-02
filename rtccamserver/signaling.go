package rtccamserver

import (
	"log"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccammessage"
)

func SignalingRouteHander(client *rtccamclient.RTCCamClient, signalingRequestMessage *rtccammessage.SignalingMessage) {
	log.Println("[SignalingRouteHander] ClientId:", client.ClientId, "RequestType:", signalingRequestMessage.RequestType)

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		errorMessage := rtccammessage.NewRTCCamErrorMessage(err.Error())
		log.Println("[SignalingRouteHander] GetRoom Failed ClientId:", client.ClientId, "RoomId:", client.JoinRoomId, "Error:", err)
		client.Send(errorMessage)
		return
	}

	responseClient, err := room.GetClient(signalingRequestMessage.ResponseClientId)
	if err != nil {
		log.Println("[SignalingRouteHander] GetClient Failed ClientId:", client.ClientId, "ResponseClientId:", signalingRequestMessage.ResponseClientId, "Error:", err)
		errorMessage := rtccammessage.NewRTCCamErrorMessage(err.Error())
		client.Send(errorMessage)
		return
	}

	responseClient.Send(signalingRequestMessage)
}

func SignalingHandler(client *rtccamclient.RTCCamClient, signalingRequestMessage *rtccammessage.SignalingMessage) {
	switch signalingRequestMessage.RequestType {
	case rtccammessage.SignalingRequestTypeOffer:
		SignalingRouteHander(client, signalingRequestMessage)
		break
	case rtccammessage.SignalingRequestTypeAnswer:
		SignalingRouteHander(client, signalingRequestMessage)
		break
	case rtccammessage.SignalingRequestTypeCandidate:
		SignalingRouteHander(client, signalingRequestMessage)
		break
	}
}
