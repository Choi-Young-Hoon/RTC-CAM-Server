package rtccammessage

import (
	"github.com/pion/webrtc/v3"
)

var MessageOffer = "offer"
var MessageAnswer = "answer"
var MessageCandidate = "candidate"

var MessageRoomList = "roomList"
var MessageRoomJoin = "roomJoin"
var MessageRoomLeave = "roomLeave"

var ResponseMessageError = "error"
var ResponseMessageSuccess = "success"

func NewErrorMessage(clientId int, errorMessage string) Message {
	return Message{
		Type:         ResponseMessageError,
		ClientId:     clientId,
		ErrorMessage: errorMessage,
	}
}

func NewSuccessMessage(clientId int) Message {
	return Message{
		Type:     ResponseMessageSuccess,
		ClientId: clientId,
	}

}

type Message struct {
	// 주고받는 메시지의 타입.
	Type string `json:"type"`

	RoomId   int `json:"roomId,omitempty"`
	ClientId int `json:"clientId,omitempty"`

	ErrorMessage string `json:"errorMessage,omitempty"`

	Offer     *webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    *webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate *webrtc.ICECandidateInit   `json:"candidate,omitempty"`
}
