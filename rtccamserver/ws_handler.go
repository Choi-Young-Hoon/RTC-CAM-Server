package rtccamserver

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"rtccam/message"
	"rtccam/roommanager"
	"rtccam/rtccamclient"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RTCCamWSCientClose(client *rtccamclient.RTCCamClient) {
	defer client.Close()

	clientManager := rtccamclient.GetRTCCamClientManager()
	clientManager.RemoveClient(client)

	roomManager := roommanager.GetRoomManager()
	room, err := roomManager.GetRoom(client.JoinRoomId)
	if err != nil {
		return
	}

	room.LeaveClient(client)

	BroadcastRoomList()
}

func RTCCamWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	clientManager := rtccamclient.GetRTCCamClientManager()
	client := rtccamclient.NewRTCCamClient(conn)
	defer RTCCamWSCientClose(client)
	clientManager.AddClient(client)

	connecMessage := message.NewConnectMessage(client.ClientId, "stun:stun.l.google.com:19302", "turn:choiyh.synology.me:50001")
	err = client.Send(connecMessage)
	if err != nil {
		log.Println("[RTCCamWSHandler] ConnectMessage ClientId:", client.ClientId, "Send Error:", err)
		return
	}

	log.Println("Client Connect Client Addr:", r.RemoteAddr)
	RTCCamServerRun(client)
}
