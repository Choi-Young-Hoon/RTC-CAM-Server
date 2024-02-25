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

func startServer(isRunningHttps bool) {
	httpServer = &http.Server{
		Addr: ":40001",
	}

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/js/", fs)

	http.HandleFunc("/", rtccamweb.HTTPIndexHandler)
	http.HandleFunc("/rtccam", rtccamserver.RTCCamWSHandler)

	infoLog(httpServer.Addr)

	if isRunningHttps {
		if err := httpServer.ListenAndServeTLS("cert.pem", "privKey.pem"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	} else {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
}

func main() {
	isRunningHttps := false
	for _, arg := range os.Args[1:] {
		if arg == "https" {
			isRunningHttps = true
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		sigStopServer()
		os.Exit(0)
	}()

	startServer(isRunningHttps)
}
