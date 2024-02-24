package message

func NewRTCCamMessage() *RTCCamMessage {
	return &RTCCamMessage{}
}

type RTCCamMessage struct {
	Room      *RoomRequestMessage `json:"room,omitempty"`
	Signaling *SignalingMessage   `json:"signaling,omitempty"`
}
