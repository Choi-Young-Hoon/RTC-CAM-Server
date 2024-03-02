package rtccammessage

func NewRTCCamRequestMessage() *RTCCamRequestMessage {
	return &RTCCamRequestMessage{}
}

type RTCCamRequestMessage struct {
	Room                *RoomRequestMessage       `json:"room,omitempty"`
	Signaling           *SignalingMessage         `json:"signaling,omitempty"`
	CreateRoomIdRequest *CreateRoomRequestMessage `json:"create_room,omitempty"`
}
