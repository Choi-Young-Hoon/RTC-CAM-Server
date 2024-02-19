package message

import (
	"github.com/pion/webrtc/v3"
)

var MessageOffer = "offer"
var MessageAnswer = "answer"
var MessageCandidate = "candidate"

var MessageRoomJoin = "roomJoin"
var MessageRoomLeave = "roomLeave"

type Message struct {
	// 주고받는 메시지의 타입.
	MessageType string `json:"messageType"`

	RoomId   int `json:"roomId"`
	ClientId int `json:"clientId"`

	Offer     *webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    *webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate *webrtc.ICECandidateInit   `json:"candidate,omitempty"`
}
