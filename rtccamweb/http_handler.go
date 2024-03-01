package rtccamweb

import (
	"html/template"
	"log"
	"net/http"
)

const BasePath = "web/static"
const BaseHtmlPath = BasePath + "/html"
const BaseJsPath = BasePath + "/js"

func HTTPNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func CreateTemplate() *template.Template {
	// 항상 요청할때마다 html 파일들을 읽어올 수 없으니 나중에 한번만 읽어오도록 바꿔야한다.
	t := template.Must(template.ParseGlob(BaseHtmlPath + "/*.html"))
	return t
}

func HTTPRTCCamHomeHandler(w http.ResponseWriter, r *http.Request) {
	t := CreateTemplate()

	pageData := NewPageData("Home")
	err := t.ExecuteTemplate(w, "rtccam_home.html", pageData)
	if err != nil {
		log.Println("[HTTPRTCCamHomeHandler] Template Execute Error:", err)
	}
}

func RoomPageHandler(w http.ResponseWriter, r *http.Request, htmlFile string) {
	log.Println("[RoomPageHandler] Start")

	t := CreateTemplate()
	pageData := NewPageData("Room")
	joinRoom := r.URL.Query().Get("join_room")
	authToken := r.URL.Query().Get("auth_token")

	if (joinRoom == "") || (authToken == "") {
		HTTPNotFoundHandler(w, r)
		return
	}

	pageData.RoomRequestType = "join_room"
	pageData.RequestId = joinRoom
	pageData.AuthToken = authToken

	err := t.ExecuteTemplate(w, htmlFile, pageData)
	if err != nil {
		log.Println("[HTTPRTCCamRoomHandler] Template Execute Error:", err)
	}
}

func HTTPRTCCamRoomHandler(w http.ResponseWriter, r *http.Request) {
	RoomPageHandler(w, r, "rtccam_room.html")
}

func JavascriptHandler(w http.ResponseWriter, r *http.Request, jsFile string) {
	t, _ := template.ParseFiles(BaseJsPath + "/" + jsFile)
	webSocketUrl := make(map[string]string)

	if HTTPProtocol == "https" {
		webSocketUrl["WebSocketURL"] = "wss://" + r.Host + "/rtccam"
	} else {
		webSocketUrl["WebSocketURL"] = "ws://" + r.Host + "/rtccam"
	}
	t.Execute(w, webSocketUrl)
}

func RTCCAMDefaultJavascriptHandler(w http.ResponseWriter, r *http.Request) {
	JavascriptHandler(w, r, "rtccam_default.js")
}
