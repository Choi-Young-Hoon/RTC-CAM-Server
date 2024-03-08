package rtccammessage

import "rtccam/rtccamerrors"

func NewRTCCamSuccessMessage() *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "success",
	}
}

func NewRTCCamJoinSuccessMessage(room interface{}, joinClientId int64) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "join_success",
		RoomInfo:      room,
		ClientId:      joinClientId,
		ICEServers:    GetICEServers(),
	}
}

func NewRTCCamErrorMessage(rtccamError *rtccamerrors.RTCCamError) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "error",
		RTCCamError:   rtccamError,
	}
}

func NewRTCCamLeaveMessage(clientId int64) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "leave_client",
		ClientId:      clientId,
	}
}

func NewRTCCamRoomListMessage(roomManager interface{}) *RTCCamResponseMessage {
	successMessage := NewRTCCamSuccessMessage()
	successMessage.RoomManager = roomManager
	return successMessage
}

func NewRTCCamAuthTokenMessage(authToken string, room interface{}) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "auth_token",
		AutoToken:     authToken,
		RoomInfo:      room,
	}
}

func NewRTCCamPublicAuthTokenMessage(authToken string) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "public_auth_token",
		AutoToken:     authToken,
	}
}

type RTCCamResponseMessage struct {
	ResultMessage string `json:"result_message"`

	RTCCamError *rtccamerrors.RTCCamError `json:"error,omitempty"`

	ClientId int64 `json:"client_id,omitempty"`

	AutoToken string `json:"auth_token,omitempty"` // AuthToken 생성 요청시 생성되는 값

	RoomManager interface{} `json:"rooms,omitempty"`     // RoomManager 구조체를 넣어줘야한다.
	RoomInfo    interface{} `json:"room_info,omitempty"` // Room 구조체를 넣어줘야한다. join_success 시 생성

	ICEServers []ICEServer `json:"ice_servers,omitempty"`
}
