package rtccam

import (
	"fmt"
	"net/http"
)

type RTCCamServer struct {
}

func NewRTCCamServer() *RTCCamServer {
	return &RTCCamServer{}
}

func (s *RTCCamServer) Start() {

}

func (s *RTCCamServer) Stop() {

}

func (s *RTCCamServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RTCCamServer ServeHTTP()")
}
