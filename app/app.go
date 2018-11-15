package app

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/clintjedwards/basecoat/utils"

	"github.com/gorilla/handlers"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/frontend"
	"github.com/gorilla/mux"
)

//APP represents an instance of the basecoat app as a whole
// the starting point for the entire application
type APP struct {
	config *config.Config
	router *mux.Router
}

//NewAPP creates a new instance of the basecoat app
func newAPP() *APP {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	return &APP{
		router: router,
		config: config,
	}
}

//InitializeApplication sets up the instance of the basecoat app to run
func InitializeApplication() *APP {
	app := newAPP()

	api := *api.NewAPI(app.config)
	api.RegisterFormulaRoutes(app.router)
	api.RegisterJobRoutes(app.router)
	api.RegisterAuthRoutes(app.router)

	if app.config.Frontend {
		frontend := *frontend.NewFrontend(app.config)
		frontend.RegisterUIRoutes(app.router)
		utils.StructuredLog("info", "frontend enabled", nil)
	}

	return app
}

//RunServer starts an instance of the basecoat app web server
func RunServer(app *APP) {
	server := http.Server{
		Addr:         app.config.ServerURL,
		Handler:      app.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	server.Handler = defaultHeaders(app.router)
	if app.config.Debug {
		server.Handler = handlers.LoggingHandler(os.Stdout, server.Handler)
	}

	utils.StructuredLog("info", "starting basecoat application", app.config.ServerURL)
	log.Fatal(server.ListenAndServe())
}

//Wrapper function setting the default headers
func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
