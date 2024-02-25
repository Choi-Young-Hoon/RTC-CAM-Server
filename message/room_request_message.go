package message

var RoomRequestTypeCreateRoom = "create_room"
var RoomRequestTypeJoinRoom = "join_room"
var RoomRequestTypeLeaveRoom = "leave_room"
var RoomRequestTypeRoomList = "room_list"

type RoomRequestMessage struct {
	RequestType string `json:"request_type"`

	// 방 생성
	Title    string `json:"title,omitempty"`
	Password string `json:"password,omitempty"`

	// 방 접속
	JoinRoomId int64 `json:"join_room_id,omitempty"`
}
