package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"rtccam/rtccamlog"
	"rtccam/rtccamweb"
	"syscall"
)

func infoLog(httpProtocol, servicePort, certPem, privKeyPem string) {
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
		break
	case "https":
		rtccamweb.StartHTTPSServer(servicePort, certPem, privKeyPem)
		break
	default:
		rtccamlog.Error().Msg("Invalid http protocol: " + httpProtocol)
	}

}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
	infoLog(*httpProtocol, *servicePort, *certPem, *privKeyPem)
	go startServer(*httpProtocol, *servicePort, *certPem, *privKeyPem)

	// Stop Server
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	rtccamweb.StopServer()
	os.Exit(0)
}
