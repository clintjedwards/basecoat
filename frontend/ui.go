package frontend

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/shurcooL/httpgzip"
	"go.uber.org/zap"
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

	// We bake frontend files directly into the binary
	// assets is an implementation of an http.filesystem created by
	// github.com/shurcooL/vfsgen that points to the "public" folder
	fileServerHandler := httpgzip.FileServer(assets, httpgzip.FileServerOptions{IndexHTML: true})

	file, err := assets.Open("index.html")
	if err != nil {
		zap.S().Fatalf("could not find index.html file: %v", err)
	}
	defer file.Close()

	indexContent, err := ioutil.ReadAll(file)
	if err != nil {
		zap.S().Fatalf("could not read index.html file: %v", err)
	}

	router.PathPrefix("/").Handler(historyModeHandler(fileServerHandler, indexContent))
}
