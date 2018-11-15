package frontend

import (
	"net/http"

	"github.com/clintjedwards/basecoat/config"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
)

//Frontend represents an instance of the frontend application
type Frontend struct {
	config *config.Config
}

//NewFrontend initializes a new UI application
func NewFrontend(config *config.Config) *Frontend {

	return &Frontend{
		config: config,
	}
}

//RegisterUIRoutes registers the endpoints needed for the frontend
// with an already established router
func (ui *Frontend) RegisterUIRoutes(router *mux.Router) {

	router.HandleFunc("/logout", nil)

	box := packr.NewBox("./public")
	router.PathPrefix("/").Handler(handleUI(http.FileServer(box)))
}

//Allows proper headers for frontend while maintaining api default headers
func handleUI(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		header := w.Header()
		header.Del("Access-Control-Allow-Origin")
		header.Del("Content-Type")
		h.ServeHTTP(w, req)
		return
	})
}
