package rtccammessage

type CreateRoomRequestMessage struct {
	// 방 생성
	Title    string `json:"title"`
	Password string `json:"password,omitempty"`

	MaxClientCount int `json:"max_client_count"`
}
