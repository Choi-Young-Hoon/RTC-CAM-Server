package rtccamserver

import (
	"errors"
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
	roomMessageDispatcher.AddHandleHandler(rtccammessage.RoomRequestPublicAuthToken, roomPublicAuthTokenHandler)

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
		resultTypeError := rtccamerrors.NewRequestTypeError()
		rtccamlog.Error().
			Err(errors.New(resultTypeError.Message)).
			Any("ClientId", client.ClientId).
			Str("RequestType", roomRequestMessage.RequestType).
			Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(resultTypeError))
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
		rtccamlog.Error().Err(errors.New(err.Message)).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(err))
		return
	}

	// 공개용 URL 을 사용해서 들어온 경우
	if room.PublicAuthToken != roomRequestMessage.AuthToken {
		// 일반 방 생성 및 방 입장일 경우
		if room.Authenticate(roomRequestMessage.AuthToken) == false {
			authTokenError := rtccamerrors.NewInvalidAuthToken()
			rtccamlog.Error().Err(errors.New(authTokenError.Message)).Any("ClientId", client.ClientId).Send()
			client.Send(rtccammessage.NewRTCCamErrorMessage(authTokenError))
			return
		}
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
		rtccamlog.Error().Err(errors.New(err.Message)).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(err))
		return
	}

	if room.MaxClientCount <= room.GetClientCount() {
		roomFullError := rtccamerrors.NewRoomIsFull()
		rtccamlog.Error().Err(errors.New(roomFullError.Message)).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(roomFullError))
		return
	}

	if room.IsPassword && room.Password != roomRequestMessage.Password {
		invalidPasswordError := rtccamerrors.NewInvalidPassword()
		rtccamlog.Error().Err(errors.New(invalidPasswordError.Message)).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(invalidPasswordError))
		return
	}

	authToken := room.GenerateAuthToken()
	client.Send(rtccammessage.NewRTCCamAuthTokenMessage(authToken, room))
}

func roomPublicAuthTokenHandler(client *rtccamclient.RTCCamClient, roomRequestMessage *rtccammessage.RoomRequestMessage) {
	rtccamlog.Info().
		Any("ClientId", client.ClientId).
		Int64("Client JoinRoomId", client.JoinRoomId).
		Send()

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		rtccamlog.Error().Err(errors.New(err.Message)).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(err))
		return
	}

	if room.PublicAuthToken == "" {
		room.GeneratePublicAuthToken()
	}
	client.Send(rtccammessage.NewRTCCamPublicAuthTokenMessage(room.PublicAuthToken))
}

func roomLeave(client *rtccamclient.RTCCamClient) {
	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		// 접속한 방이 없으면.
		if client.JoinRoomId == 0 {
			return
		}

		rtccamlog.Error().Err(errors.New(err.Message)).Any("ClientId", client.ClientId).Send()
		client.Send(rtccammessage.NewRTCCamErrorMessage(err))
		return
	}

	rtccamlog.Info().Any("ClientId", client.ClientId).Int64("LeaveRoomId", room.Id).Send()
	room.LeaveClient(client)
}
