package message

func NewConnectMessage(clientId int64, stunAddr string, turnAddr string) *ConnectMessage {
	return &ConnectMessage{
		ClientId: clientId,
		StunAddr: stunAddr,
		TurnAddr: turnAddr,
	}
}

type ConnectMessage struct {
	ClientId int64  `json:"client_id"`
	StunAddr string `json:"stun_addr"`
	TurnAddr string `json:"turn_addr"`
}
