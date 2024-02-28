package rtccamweb

import (
	"html/template"
	"log"
	"net/http"
)

const BasePath = "web/static"
const BaseHtmlPath = BasePath + "/html"
const BaseJsPath = BasePath + "/js"

func HTTPIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/index.html")
}

func HTTPDesignTestHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/designtest_home.html")
}

func HTTPDesignTestRoomHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/designtest_room.html")
}

func HTTPDesignTestRoomMobileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/designtest_room_mobile.html")
}

func HTTPRTCCamHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/rtccam_home.html")
}

func HTTPRTCCamRoomHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(BaseHtmlPath + "/rtccam_room.html")
	templateMap := make(map[string]string)

	createRoom := r.URL.Query().Get("create_room")
	if createRoom != "" {
		// room?create_room=1 이런식으로 오면 해당 페이지에서
		// create_room 요청을 보내서 만드는 식으로 진행
		templateMap["room_request_type"] = "create_room"
		templateMap["request_id"] = createRoom
	}

	joinRoom := r.URL.Query().Get("join_room")
	if joinRoom != "" {
		templateMap["room_request_type"] = "join_room"
		templateMap["request_id"] = joinRoom
	}

	if templateMap["room_request_type"] == "" {
		http.NotFound(w, r)
	} else {
		err := t.Execute(w, templateMap)
		if err != nil {
			log.Println("[HTTPRTCCamRoomHandler] Template Execute Error:", err)
		}
	}
}

func HTTPRTCCamRoomMobileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/rtccam_room_mobile.html")
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
