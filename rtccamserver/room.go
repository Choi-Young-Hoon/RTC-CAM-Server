package rtccamserver

import (
	"rtccam/roommanager"
	"rtccam/rtccamclient"
	"rtccam/rtccamerrors"
	"rtccam/rtccamlog"
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
	rtccamlog.Info().
		Any("ClientId", client.ClientId).
		Str("RequestType", roomRequestMessage.RequestType).
		Send()

	handle, ok := r.handles[roomRequestMessage.RequestType]
	if !ok {
		rtccamlog.Error().
			Err(rtccamerrors.ErrorRequestTypeError).
			Any("ClientId", client.ClientId).
			Str("RequestType", roomRequestMessage.RequestType).
			Send()
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
		rtccamlog.Error().Err(err).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(err.Error()))
		return
	}

	if room.Authenticate(roomRequestMessage.AuthToken) == false {
		rtccamlog.Error().Err(rtccamerrors.ErrorInvalidAuthToken).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(rtccamerrors.ErrorInvalidAuthToken.Error()))
		return
	}

	rtccamlog.Info().Any("ClientId", client.ClientId).Int64("Join RoomId", room.Id).Send()
	room.JoinClient(client)
	broadcastRoomList()
}

func roomLeaveHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	roomLeave(client)
	broadcastRoomList()
}

func roomAuthTokenHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	rtccamlog.Info().
		Any("ClientId", client.ClientId).
		Int64("JoinRoomId", roomRequestMessage.RoomId).
		Send()

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(roomRequestMessage.RoomId)
	if err != nil {
		rtccamlog.Error().Err(err).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(err.Error()))
		return
	}

	if room.MaxClientCount <= room.GetClientCount() {
		rtccamlog.Error().Err(rtccamerrors.ErrorRoomIsFull).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(rtccamerrors.ErrorRoomIsFull.Error()))
		return
	}

	if room.IsPassword && room.Password != roomRequestMessage.Password {
		rtccamlog.Error().Err(rtccamerrors.ErrorInvalidPassword).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(rtccamerrors.ErrorInvalidPassword.Error()))
		return
	}

	authToken := room.GenerateAuthToken()
	client.Send(rtccammessage.NewRTCCamAuthTokenMessage(authToken, room))
}

func roomLeave(client *rtccamclient.RTCCamClient) {
	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		// 접속한 방이 없으면.
		if client.JoinRoomId == 0 {
			return
		}

		rtccamlog.Error().Err(err).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(err.Error()))
		return
	}

	rtccamlog.Info().Any("ClientId", client.ClientId).Int64("LeaveRoomId", room.Id).Send()
	room.LeaveClient(client)
}
