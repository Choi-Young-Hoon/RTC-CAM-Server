package message

import "github.com/pion/webrtc/v3"

const SignalingRequestTypeOffer = "offer"
const SignalingRequestTypeAnswer = "answer"
const SignalingRequestTypeCandidate = "candidate"

type SignalingMessage struct {
	RequestType string `json:"request_type"`

	RequestClientId  int64 `json:"request_client_id"`
	ResponseClientId int64 `json:"response_client_id"`

	Offer     *webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    *webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate *webrtc.ICECandidateInit   `json:"candidate,omitempty"`
}
