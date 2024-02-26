package rtccamweb

import (
	"html/template"
	"net/http"
)

const BasePath = "web/static"
const BaseHtmlPath = BasePath + "/html"
const BaseJsPath = BasePath + "/js"

func HTTPIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/index.html")
}

func RTCCAMJavascriptHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(BaseJsPath + "/rtccam.js")
	webSocketUrl := make(map[string]string)

	if HTTPProtocol == "https" {
		webSocketUrl["WebSocketURL"] = "wss://" + r.Host + "/rtccam"
	} else {
		webSocketUrl["WebSocketURL"] = "ws://" + r.Host + "/rtccam"
	}
	t.Execute(w, webSocketUrl)
}
