package message

var RoomRequestTypeJoinRoom = "join_room"
var RoomRequestTypeLeaveRoom = "leave_room"
var RoomRequestTypeRoomList = "room_list"

type RoomRequestMessage struct {
	RequestType string `json:"request_type"`
	JoinRoomId  int64  `json:"join_room_id,omitempty"`
}
