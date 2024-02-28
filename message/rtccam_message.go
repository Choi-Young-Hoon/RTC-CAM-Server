package message

func NewRTCCamRequestMessage() *RTCCamRequestMessage {
	return &RTCCamRequestMessage{}
}

type RTCCamRequestMessage struct {
	Room         *RoomRequestMessage         `json:"room,omitempty"`
	Signaling    *SignalingMessage           `json:"signaling,omitempty"`
	Connect      *ConnectRequestMessage      `json:"connect,omitempty"`
	CreateRoomId *CreateRoomIdRequestMessage `json:"create_room_url,omitempty"`
}
