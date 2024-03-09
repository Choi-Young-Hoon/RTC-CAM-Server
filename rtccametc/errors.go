package rtccametc

func NewClientNotFound() *RTCCamError {
	return &RTCCamError{
		Code:    1001,
		Message: "클라이언트를 찾을 수 없습니다.",
	}
}
func NewRoomNotFound() *RTCCamError {
	return &RTCCamError{
		Code:    1002,
		Message: "방을 찾을 수 없습니다.",
	}
}

func NewRoomIsFull() *RTCCamError {
	return &RTCCamError{
		Code:    1003,
		Message: "방이 꽉 찼습니다.",
	}
}

func NewRequestTypeError() *RTCCamError {
	return &RTCCamError{
		Code:    1004,
		Message: "요청 타입이 올바르지 않습니다.",
	}
}

func NewTitleIsEmpty() *RTCCamError {
	return &RTCCamError{
		Code:    1005,
		Message: "방 제목이 비어있습니다.",
	}
}

func NewInvalidMaxClientCount() *RTCCamError {
	return &RTCCamError{
		Code:    1006,
		Message: "최대 인원 수가 올바르지 않습니다.",
	}
}

func NewInvalidPassword() *RTCCamError {
	return &RTCCamError{
		Code:    1007,
		Message: "비밀번호가 올바르지 않습니다.",
	}
}

func NewInvalidAuthToken() *RTCCamError {
	return &RTCCamError{
		Code:    1008,
		Message: "잘못된 접근입니다.",
	}
}

type RTCCamError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}
