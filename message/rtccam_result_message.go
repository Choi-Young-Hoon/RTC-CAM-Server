package message

func NewRTCCamSuccessMessage() *RTCCamResultMessage {
	return &RTCCamResultMessage{
		ResultMessage: "success",
	}
}

func NewRTCCamJoinSuccessMessage() *RTCCamResultMessage {
	return &RTCCamResultMessage{
		ResultMessage: "join_success",
	}
}

func NewRTCCamErrorMessage(errorMessage string) *RTCCamResultMessage {
	return &RTCCamResultMessage{
		ResultMessage: "error",
		ErrorMessage:  errorMessage,
	}
}

func NewRTCCamLeaveMessage(clientId int64) *RTCCamResultMessage {
	return &RTCCamResultMessage{
		ResultMessage: "leave_client",
		LeaveClientId: clientId,
	}
}

func NewRTCCamRoomListMessage(roomManager interface{}) *RTCCamResultMessage {
	successMessage := NewRTCCamSuccessMessage()
	successMessage.RoomManager = roomManager
	return successMessage
}

type RTCCamResultMessage struct {
	ResultMessage string `json:"result_message"`

	LeaveClientId int64 `json:"leave_client_id,omitempty"`

	ErrorMessage string `json:"error_message,omitempty"`

	RoomManager interface{} `json:"rooms,omitempty"`
}
