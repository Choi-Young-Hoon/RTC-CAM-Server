package message

func NewRTCCamSuccessMessage() *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "success",
	}
}

func NewRTCCamJoinSuccessMessage() *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "join_success",
	}
}

func NewRTCCamErrorMessage(errorMessage string) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "error",
		ErrorMessage:  errorMessage,
	}
}

func NewRTCCamLeaveMessage(clientId int64) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage: "leave_client",
		LeaveClientId: clientId,
	}
}

func NewRTCCamRoomListMessage(roomManager interface{}) *RTCCamResponseMessage {
	successMessage := NewRTCCamSuccessMessage()
	successMessage.RoomManager = roomManager
	return successMessage
}

func NewRTCCamConnectMessage(clientId int64) *RTCCamResponseMessage {
	return &RTCCamResponseMessage{
		ResultMessage:  "connect_result",
		ConnectMessage: NewConnectResponseMessage(clientId),
	}
}

type RTCCamResponseMessage struct {
	ResultMessage string `json:"result_message"`

	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`

	LeaveClientId int64 `json:"leave_client_id,omitempty"`

	RoomManager interface{} `json:"rooms,omitempty"` // RoomManager 구조체를 넣어줘야한다.

	ConnectMessage *ConnectReponseMessage `json:"connect_message,omitempty"`
}
