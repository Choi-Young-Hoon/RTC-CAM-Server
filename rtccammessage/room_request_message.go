package rtccammessage

var RoomRequestTypeJoinRoom = "join_room"
var RoomRequestTypeLeaveRoom = "leave_room"
var RoomRequestTypeRoomList = "room_list"
var RoomRequestAuthToken = "auth_token"
var RoomRequestPublicAuthToken = "public_auth_token"

type RoomRequestMessage struct {
	RequestType string `json:"request_type"`

	// 요청한 방 Id, 패스워드
	RoomId   int64  `json:"room_id,omitempty"`
	Password string `json:"password,omitempty"`

	// 방 접속 요청시 사용하는 인증 토큰
	AuthToken string `json:"auth_token,omitempty"`
}
