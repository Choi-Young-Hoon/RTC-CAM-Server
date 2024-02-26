package rtccamweb

import (
	"context"
	"errors"
	"net/http"
	"rtccam/rtccamserver"
)

var httpServer *http.Server = &http.Server{}

func initHTTPHandler() {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/js/", fs)
	http.Handle("/css/", fs)

	http.HandleFunc("/", HTTPIndexHandler)
	http.HandleFunc("/rtccam", rtccamserver.RTCCamWSHandler)
}

func StartHTTPSServer(servicePort string, certPem string, privKeyPem string) {
	initHTTPHandler()

	httpServer.Addr = ":" + servicePort
	if err := httpServer.ListenAndServeTLS(certPem, privKeyPem); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func StartHTTPServer(servicePort string) {
	initHTTPHandler()

	httpServer.Addr = ":" + servicePort
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func StopServer() {
	if err := httpServer.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
