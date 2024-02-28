package message

func NewCreateRoomIdMessage(id int64) *CreateRoomIdMessage {
	return &CreateRoomIdMessage{
		Id: id,
	}
}

type CreateRoomIdMessage struct {
	Id int64 `json:"room_id"`
}
