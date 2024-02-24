package rtccamserver

import (
	"log"
	"rtccam/rtccamclient"
)

func RTCCamServerRun(client *rtccamclient.RTCCamClient) {
	for {
		message, err := client.Recv()
		if err != nil {
			log.Println("[RTCCamServer] ClientId:", client.ClientId, "Recv Error:", err)
			return
		}

		if message.Room != nil {
			RoomHandler(client, message.Room)
		} else if message.Signaling != nil {
			SignalingHandler(client, message.Signaling)
		}
	}
}
