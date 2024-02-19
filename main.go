package main

import (
	"log"
	"net/http"
	"rtccam/rtccamroom"
	"rtccam/rtccamserver"
	"rtccam/signaling"
)

func createRoom() {
	roomManager := rtccamroom.GetRoomManager()
	for i := 0; i < 5; i++ {
		err := roomManager.CreateRoom("Room - 1" + string(i))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	createRoom()

	rtcCamServer := rtccamserver.NewRTCCamServer()
	signalingServer := signaling.NewSignalingServer()

	http.Handle("/rtccam", rtcCamServer)
	http.Handle("/signaling", signalingServer)

	log.Println("Server Start Service Port: 50002")

	err := http.ListenAndServeTLS(":50002", "cert.pem", "privkey.pem", nil)
	if err != nil {
		panic(err)
	}
}
