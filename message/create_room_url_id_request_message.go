package message

type CreateRoomIdRequestMessage struct {
	// 방 생성
	Title    string `json:"title,omitempty"`
	Password string `json:"password,omitempty"`
}
