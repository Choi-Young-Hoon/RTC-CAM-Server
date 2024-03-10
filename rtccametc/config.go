package rtccametc

import (
	"gopkg.in/yaml.v3"
	"os"
)

func NewConifg() *RTCCamConfig {
	return &RTCCamConfig{}
}

func NewDefaultConfig() *RTCCamConfig {
	return &RTCCamConfig{
		ImageServerUrl: "http://localhost:40001",

		ServerConfig: ServerConfig{
			Protocol: "http",
			Port:     "40001",
			HTTPSCert: HTTPSCert{
				CertFile:    "cert.pem",
				PrivKeyFile: "privKey.pem",
			},
		},

		ICEServers: ICEServers{
			StunServers: []StunServer{
				{
					URL: "stun:stun.l.google.com:19302",
				},
				{
					URL: "stun:stun1.l.google.com:19302",
				},
				{
					URL: "stun:stun2.l.google.com:19302",
				},
			},

			TurnServers: []TurnServer{
				{
					URL:        "turn:turn.example.com:3478",
					Username:   "username",
					Credential: "credential",
				},
			},
		},
	}
}

type RTCCamConfig struct {
	ImageServerUrl string `yaml:"image_server_url"`

	ServerConfig `yaml:"server"`

	ICEServers `yaml:"ice_servers"`
}

func (c *RTCCamConfig) WriteConfig() error {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile("config.yaml", yamlData, 0755)
}

func (c *RTCCamConfig) ReadConfig() error {
	yamlData, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlData, c)
}

func (c *RTCCamConfig) GetIceServersToJson() []interface{} {
	iceServers := make([]interface{}, 0, len(c.StunServers)+len(c.TurnServers))
	for _, stunServer := range c.StunServers {
		iceServers = append(iceServers, struct {
			Urls string `yaml:"urls"`
		}{
			Urls: stunServer.URL,
		})
	}

	for _, turnServer := range c.TurnServers {
		iceServers = append(iceServers, struct {
			Urls       string `json:"urls"`
			UserName   string `json:"username"`
			Credential string `json:"credential"`
		}{
			Urls:       turnServer.URL,
			UserName:   turnServer.Username,
			Credential: turnServer.Credential,
		})
	}

	return iceServers
}

type HTTPSCert struct {
	CertFile    string `yaml:"cert"`
	PrivKeyFile string `yaml:"priv_key"`
}

type ServerConfig struct {
	Protocol  string    `yaml:"protocol"`
	HTTPSCert HTTPSCert `yaml:"https_cert,omitempty"`

	Port string `yaml:"port"`
}

type StunServer struct {
	URL string `yaml:"url"`
}

type TurnServer struct {
	URL        string `yaml:"url"`
	Username   string `yaml:"username"`
	Credential string `yaml:"credential"`
}

type ICEServers struct {
	StunServers []StunServer `yaml:"stun"`

	TurnServers []TurnServer `yaml:"turn"`
}
