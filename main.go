package main

import (
	"net/http"
	"rtccam/rtccam"
	"rtccam/signaling"
)

func main() {
	rtcCamServer := rtccam.NewRTCCamServer()
	signalingServer := signaling.NewSignalingServer()

	http.Handle("/rtccam", rtcCamServer)
	http.Handle("/signaling", signalingServer)

	err := http.ListenAndServeTLS(":50002", "cert.pem", "privkey.pem", nil)
	if err != nil {
		panic(err)
	}
}
