package rtccamweb

import (
	"context"
	"errors"
	"net/http"
	"rtccam/rtccamclient"
	"rtccam/rtccamserver"
)

var httpServer *http.Server = &http.Server{}
var HTTPProtocol string

func initHTTPHandler() {
	//http.Handle("/js/", fs)
	// rtccam.js 웹소켓 서버 주소 동적 삽입.
	http.HandleFunc("/js/rtccam_default.js", RTCCAMDefaultJavascriptHandler)
	http.Handle("/js/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/css/", http.FileServer(http.Dir("./web/static")))
	http.Handle("/img/", http.FileServer(http.Dir("./web/resource/")))

	http.HandleFunc("/", HTTPRTCCamHomeHandler)
	http.HandleFunc("/room", HTTPRTCCamRoomHandler)
	http.HandleFunc("/rtccam", rtccamserver.RTCCamWSHandler)
}

func StartHTTPSServer(servicePort string, certPem string, privKeyPem string) {
	initHTTPHandler()

	HTTPProtocol = "https"
	httpServer.Addr = ":" + servicePort
	if err := httpServer.ListenAndServeTLS(certPem, privKeyPem); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func StartHTTPServer(servicePort string) {
	initHTTPHandler()

	HTTPProtocol = "http"
	httpServer.Addr = ":" + servicePort
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func StopServer() {
	clientManager := rtccamclient.GetRTCCamClientManager()
	clientManager.CloseAll()

	if err := httpServer.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
