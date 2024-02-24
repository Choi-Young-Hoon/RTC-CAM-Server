package rtccamweb

import "net/http"

const BasePath = "web/static"
const BaseHtmlPath = BasePath + "/html"
const BaseJsPath = BasePath + "/js"

func HTTPIndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, BaseHtmlPath+"/index.html")
}
