package rtccamserver

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"rtccam/rtccamclient"
	"rtccam/rtccammessage"
	"rtccam/rtccamroom"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewRTCCamServer() *RTCCamServer {
	return &RTCCamServer{}
}

type RTCCamServer struct {
}

func (s *RTCCamServer) run(client *rtccamclient.RTCCamClient) {
	for {
		message := &rtccammessage.Message{}
		err := client.Recv(message)
		if err != nil {
			log.Println("[RTCCamServer] Recv Error:", err)
			break
		}
		defer s.leaveRoom(client, client.JoinRoomId)

		switch message.Type {
		case rtccammessage.MessageRoomList:
			err = s.sendRoomList(client)
			break
		case rtccammessage.MessageRoomJoin:
			err = s.joinRoom(client, message.RoomId)
			break
		case rtccammessage.MessageRoomLeave:
			err = s.leaveRoom(client, message.RoomId)
			break
		}

		if err != nil {

			return
		}

		err = client.SendSuccessMessage()
		if err != nil {
			return
		}

	}
}

func (r *RTCCamServer) sendRoomList(client *rtccamclient.RTCCamClient) error {
	roomManager := rtccamroom.GetRoomManager()
	err := client.Websocket.WriteJSON(roomManager)
	if err != nil {
		log.Println("[RTCCamServer] SendRoomList Error:", err)
		return err
	}

	log.Println("[RTCCamServer] SendRoomList", "clientId:", client.Id)

	return nil
}

func (r *RTCCamServer) leaveRoom(client *rtccamclient.RTCCamClient, roomId int) error {
	roomManager := rtccamroom.GetRoomManager()
	room, err := roomManager.GetRoom(roomId)
	if err != nil {
		if client.JoinRoomId != -1 {
			log.Println("[RTCCamServer] leaveRoom Error:", err)
		}
		return err
	}

	log.Println("[RTCCamServer] LeaveRoom", "clientId:", client.Id, "leaveRoom:", room.Id)
	room.RemoveClient(client)

	return nil
}

func (r *RTCCamServer) joinRoom(client *rtccamclient.RTCCamClient, roomId int) error {
	r.leaveRoom(client, client.JoinRoomId)

	roomManager := rtccamroom.GetRoomManager()
	room, err := roomManager.GetRoom(roomId)
	if err != nil {
		log.Println("[RTCCamServer] joinRoom Error:", err)
		return err
	}
	log.Println("[RTCCamServer] JoinRoom", "clientId:", client.Id, "joinRoom:", room.Id)
	room.AddClient(client)

	return nil
}

func (s *RTCCamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	defer conn.Close()

	client := rtccamclient.NewRTCCamClient(conn)
	log.Println("[RTCCamServer] Client Connect:", r.RemoteAddr)
	log.Println("[RTCCamServer] ClientId:", client.Id)

	s.run(client)
}
