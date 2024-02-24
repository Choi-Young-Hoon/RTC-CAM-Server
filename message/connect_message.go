package message

func NewConnectMessage(clientId int64) *ConnectMessage {
	return &ConnectMessage{
		ClientId: clientId,
	}
}

type ConnectMessage struct {
	ClientId int64 `json:"client_id"`
}
