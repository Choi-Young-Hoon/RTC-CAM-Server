package rtccamserver

import (
	"log"
	"rtccam/message"
	"rtccam/rtccamclient"
)

func ConnectHandler(client *rtccamclient.RTCCamClient, connectRequestMessage *message.ConnectRequestMessage) {
	clientId := client.GenerateClientId()
	err := client.Send(message.NewConnectResponseMessage(clientId))
	if err != nil {
		log.Println("[ConnectHandler] ClientId:", clientId, "Error:", err)
	}
}
