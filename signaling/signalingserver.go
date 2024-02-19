package signaling

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type SignalingServer struct {
}

func NewSignalingServer() *SignalingServer {
	return &SignalingServer{}
}

func (s *SignalingServer) Start() {

}

func (s *SignalingServer) Stop() {

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *SignalingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	defer conn.Close()

	log.Println("[SignalingServer] Client Connect:", r.RemoteAddr)
}
