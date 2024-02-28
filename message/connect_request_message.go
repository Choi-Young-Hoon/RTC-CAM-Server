package message

func NewConnectRequestMessage() *ConnectRequestMessage {
	return &ConnectRequestMessage{}
}

type ConnectRequestMessage struct {
	RequestType string `json:"request_type"`
}
