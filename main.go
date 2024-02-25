package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"rtccam/rtccamserver"
	"rtccam/rtccamweb"
	"syscall"
)

var httpServer *http.Server

func infoLog(servicePort string) {
	println("============================================")
	println("=           RTCCam Server Start            =")
	println("============================================")
	println("= ServicePort: " + servicePort + "                      =")
	println("============================================")

}

func sigStopServer() {
	if err := httpServer.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}

func startServer() {
	httpServer = &http.Server{
		Addr: ":40001",
	}

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/js/", fs)

	http.HandleFunc("/", rtccamweb.HTTPIndexHandler)
	http.HandleFunc("/rtccam", rtccamserver.RTCCamWSHandler)

	infoLog(httpServer.Addr)

	//if err := httpServer.ListenAndServeTLS("cert.pem", "privKey.pem"); err != nil && !errors.Is(err, http.ErrServerClosed) {
	//	panic(err)
	//}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		sigStopServer()
		os.Exit(0)
	}()

	startServer()
}
