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

func NewRTCCamRoomListMessage(roomManager interface{}) *RTCCamResultMessage {
	successMessage := NewRTCCamSuccessMessage()
	successMessage.RoomManager = roomManager
	return successMessage
}

type RTCCamResultMessage struct {
	ResultMessage string `json:"result_message"`

	ErrorMessage string `json:"error_message,omitempty"`

	RoomManager interface{} `json:"rooms,omitempty"`
}
