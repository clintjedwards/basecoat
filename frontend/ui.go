package frontend

import (
	"log"
	"net/http"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
)

//Frontend represents an instance of the frontend application
type Frontend struct{}

//NewFrontend initializes a new UI application
func NewFrontend() *Frontend {
	return &Frontend{}
}

//historyModeHandler is a hack so that our frontend can use history mode
// We essentially answer all requests for files normally but
// any other paths we just return the normal index.html file.
// https://router.vuejs.org/guide/essentials/history-mode.html
func historyModeHandler(fileServerHandler http.Handler, indexFile []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		// serve static files as normal
		if strings.Contains(path, ".") || path == "/" {
			fileServerHandler.ServeHTTP(w, req)
			return
		}

		// return index.html for any 404s
		w.Write(indexFile)
		return
	})
}

//RegisterUIRoutes registers the endpoints needed for the frontend
// with an already established router
func (ui *Frontend) RegisterUIRoutes(router *mux.Router) {

	box := packr.NewBox("./public")
	fileServerHandler := http.FileServer(box)
	indexHTMLfile, err := box.Find("./index.html")
	if err != nil {
		log.Println(err)
	}

	router.PathPrefix("/").Handler(historyModeHandler(fileServerHandler, indexHTMLfile))
}
