package rtccammessage

// 'iceServers': [
//
//	{'urls': stunServerUrl},
//	{'urls': turnServerUrl, 'username': 'test', 'credential': 'test'},
//	{'urls': 'turn:kyj9447.iptime.org:50001', 'username': 'test', 'credential': 'test'},
//
// ]

func GetICEServers() []ICEServer {
	iceServers := []ICEServer{
		{
			Urls: "stun:stun.l.google.com:19302",
		},
		{
			Urls: "stun:stun1.l.google.com:19302",
		},
		{
			Urls: "stun:stun2.l.google.com:19302",
		},
		{
			Urls:       "turn:choiyh.synology.me:50001",
			UserName:   "test",
			Credential: "test",
		},
	}
	return iceServers
}

type ICEServer struct {
	Urls       string `json:"urls"`
	UserName   string `json:"username,omitempty"`
	Credential string `json:"credential,omitempty"`
}
