package message

func NewRTCCamSuccessMessage() *RTCCamResultMessage {
	return &RTCCamResultMessage{
		IsSuccess: true,
	}
}

func NewRTCCamErrorMessage(errorMessage string) *RTCCamResultMessage {
	return &RTCCamResultMessage{
		IsSuccess:    false,
		ErrorMessage: errorMessage,
	}
}

func NewRTCCamRoomListMessage(roomManager interface{}) *RTCCamResultMessage {
	successMessage := NewRTCCamSuccessMessage()
	successMessage.RoomManager = roomManager
	return successMessage
}

type RTCCamResultMessage struct {
	IsSuccess bool `json:"is_success"`

	ErrorMessage string `json:"error_message,omitempty"`

	RoomManager interface{} `json:"rooms,omitempty"`
}
