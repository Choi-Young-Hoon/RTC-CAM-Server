package rtccammessage

func NewCreateRoomMessage(id int64, authToken string) *CreateRoomMessage {
	return &CreateRoomMessage{
		Id:        id,
		AuthToken: authToken,
	}
}

type CreateRoomMessage struct {
	Id        int64  `json:"room_id"`
	AuthToken string `json:"auth_token"`
}
