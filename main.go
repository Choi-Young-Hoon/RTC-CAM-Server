package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"rtccam/rtccametc"
	"rtccam/rtccamlog"
	"rtccam/rtccamserver"
	"rtccam/rtccamweb"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func info(config *rtccametc.RTCCamConfig) {
	serverConifg := config.ServerConfig
	httpsCert := config.HTTPSCert

	fmt.Println("============================================")
	fmt.Println("=           RTCCam Server Start            =")
	fmt.Println("=================  Info  ===================")
	fmt.Println("= HTTP Protocol: ", serverConifg.Protocol)
	fmt.Println("= Service Port: ", serverConifg.Port)
	if serverConifg.Protocol == "https" {
		fmt.Println("= Cert Pem: ", httpsCert.CertFile)
		fmt.Println("= PrivKey Pem: ", httpsCert.PrivKeyFile)
	}
	fmt.Println("============================================")
}

func startServer(config *rtccametc.RTCCamConfig) {
	serverConfig := config.ServerConfig
	httpsCert := config.HTTPSCert

	rtccamweb.ImageServerUrl = config.ImageServerUrl
	rtccamserver.ICEServers = config.GetIceServersToJson()

	switch serverConfig.Protocol {
	case "http":
		rtccamweb.StartHTTPServer(serverConfig)
	case "https":
		rtccamweb.StartHTTPSServer(serverConfig, httpsCert)
	default:
		rtccamlog.Error().Msg("Invalid http protocol: " + serverConfig.Protocol)
	}

}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})

	defaultConifgGen := flag.Bool("g", false, "기본 설정 파일 생성")
	help := flag.Bool("h", false, "도와주세요")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *defaultConifgGen {
		rtccametc.NewDefaultConfig().WriteConfig()
		println("Generate 'config.yaml' file")
		os.Exit(0)
	}

	config := rtccametc.NewConifg()
	err := config.ReadConfig()
	if err != nil {
		rtccamlog.Error().Err(err).Send()
		println("please check 'config.yaml' file")
		os.Exit(1)
	}

	// Start Server
	info(config)
	go startServer(config)

	// Stop Server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	println("Stop RTC-CAM server..........")
	rtccamweb.StopServer()
	os.Exit(0)
}
