package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"rtccam/rtccamlog"
	"rtccam/rtccamweb"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func info(httpProtocol, servicePort, certPem, privKeyPem string) {
	fmt.Println("============================================")
	fmt.Println("=           RTCCam Server Start            =")
	fmt.Println("=================  Info  ===================")
	fmt.Println("= HTTP Protocol: ", httpProtocol)
	fmt.Println("= Service Port: ", servicePort)
	if httpProtocol == "https" {
		fmt.Println("= Cert Pem: ", certPem)
		fmt.Println("= PrivKey Pem: ", privKeyPem)
	}
	fmt.Println("============================================")
}

func startServer(httpProtocol, servicePort, certPem, privKeyPem string) {
	switch httpProtocol {
	case "http":
		rtccamweb.StartHTTPServer(servicePort)
	case "https":
		rtccamweb.StartHTTPSServer(servicePort, certPem, privKeyPem)
	default:
		rtccamlog.Error().Msg("Invalid http protocol: " + httpProtocol)
	}

}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})

	httpProtocol := flag.String("protocol", "http", "http 프로토콜(http, https)")
	servicePort := flag.String("p", "40001", "포트번호")
	certPem := flag.String("c", "cert.pem", "인증서 파일")
	privKeyPem := flag.String("k", "privKey.pem", "개인키 파일")
	help := flag.Bool("h", false, "도와주세요")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Start Server
	info(*httpProtocol, *servicePort, *certPem, *privKeyPem)
	go startServer(*httpProtocol, *servicePort, *certPem, *privKeyPem)

	// Stop Server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	rtccamweb.StopServer()
	os.Exit(0)
}
