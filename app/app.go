package app

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/frontend"
	"github.com/clintjedwards/basecoat/metrics"
	"github.com/clintjedwards/basecoat/service"
	"github.com/clintjedwards/toolkit/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

// StartServices initializes a GRPC-web compatible webserver and a GRPC service
func StartServices() {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	api := service.NewBasecoatAPI(config)
	grpcServer := service.CreateGRPCServer(api)

	go metrics.InitPrometheusService(config)

	initCombinedService(config, grpcServer)
}

// initCombinedService starts a long running combined grpc/http (grpc-web compatible) service with all proper settings; TLS enabled
func initCombinedService(config *config.Config, server *grpc.Server) {
	wrappedGrpc := grpcweb.WrapServer(server)

	router := mux.NewRouter()

	if config.Frontend.Enable {
		frontend := frontend.NewFrontend()
		frontend.RegisterUIRoutes(router)
		logger.Log().Infow("basecoat frontend enabled",
			"enabled", config.Frontend.Enable)
	}

	combinedHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Content-Type"), "application/grpc") || wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		router.ServeHTTP(resp, req)
	})

	httpServer := http.Server{
		Addr:         config.URL,
		Handler:      combinedHandler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Set default http headers
	httpServer.Handler = defaultHeaders(httpServer.Handler)

	if config.Debug {
		httpServer.Handler = handlers.LoggingHandler(os.Stdout, httpServer.Handler)
	}

	logger.Log().Infow("starting basecoat grpc/http service", "url", config.URL)
	log.Fatal(httpServer.ListenAndServeTLS(config.TLSCertPath, config.TLSKeyPath))
}

// Wrapper function setting http headers
func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		next.ServeHTTP(w, r)
	})
}
