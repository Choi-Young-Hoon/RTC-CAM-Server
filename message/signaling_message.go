package message

import "github.com/pion/webrtc/v3"

type SignalingMessage struct {
	RequestClientId  int `json:"request_client_id"`
	ResponseClientId int `json:"response_client_id,omitempty"`

	Offer     *webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    *webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate *webrtc.ICECandidateInit   `json:"candidate,omitempty"`
}
