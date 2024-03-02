package rtccamserver

import (
	"rtccam/rtccamclient"
	"rtccam/rtccamlog"
)

func RTCCamServerRun(client *rtccamclient.RTCCamClient) {
	roomMessageDispatcher := GetRoomMessageDispatcher()

	for {
		message, err := client.Recv()
		if err != nil {
			rtccamlog.Error().Err(err).Any("ClientId", client.ClientId).Send()
			return
		}

		if message.Room != nil {
			roomMessageDispatcher.RoomHandler(client, message.Room)
		} else if message.Signaling != nil {
			SignalingHandler(client, message.Signaling)
		} else if message.CreateRoomIdRequest != nil {
			CreateRoomHandler(client, message.CreateRoomIdRequest)
		}
	}
}
