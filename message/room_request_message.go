package message

var RoomRequestTypeCreateRoom = "create_room"
var RoomRequestTypeJoinRoom = "join_room"
var RoomRequestTypeLeaveRoom = "leave_room"
var RoomRequestTypeRoomList = "room_list"

type RoomRequestMessage struct {
	RequestType string `json:"request_type, omitempty"`

	// 방 생성 클라이언트가 room?create_id=1234 이런식으로 방을 생성하면
	// 해당 방의 아이디를 클라이언트가 id 값으로 요청을 보내온다
	CreateRoomId int64 `json:"create_id,omitempty"`

	// join_room 요청시 패스워드
	JoinPassword string `json:"password,omitempty"`

	// 방 접속
	JoinRoomId int64 `json:"join_room_id,omitempty"`
}
