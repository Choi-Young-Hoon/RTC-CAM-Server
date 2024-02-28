package rtccamserver

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

	broadcastRoomList()
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

	log.Println("Client Connect Client Addr:", r.RemoteAddr)
	RTCCamServerRun(client)
}
