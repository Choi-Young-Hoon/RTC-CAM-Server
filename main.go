package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"rtccam/roommanager"
	"rtccam/rtccamserver"
	"rtccam/rtccamweb"
	"strconv"
	"syscall"
)

var httpServer *http.Server

func createDummyRoom() {
	roomManager := roommanager.GetRoomManager()
	for i := 1; i < 6; i++ {
		room := roommanager.NewRoom("Room - " + strconv.Itoa(i))
		roomManager.AddRoom(room)
	}
}

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
		Addr: ":50001",
	}

	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/js/", fs)

	http.HandleFunc("/", rtccamweb.HTTPIndexHandler)
	http.HandleFunc("/rtccam", rtccamserver.RTCCamWSHandler)

	infoLog(httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func main() {
	createDummyRoom()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		sigStopServer()
		os.Exit(0)
	}()

	startServer()
}
