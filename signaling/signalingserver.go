package signaling

import (
	"fmt"
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

func (s *SignalingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SignalingServer ServeHTTP()")
}
