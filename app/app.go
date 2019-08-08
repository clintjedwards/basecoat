package app

import (
	"log"

	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/frontend"
	"github.com/clintjedwards/basecoat/metrics"
	"github.com/clintjedwards/basecoat/service"
)

// StartServices initializes a GRPC-web compatible webserver and a GPRC service
func StartServices() {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	api := service.NewBasecoatAPI(config)
	grpcServer := service.CreateGRPCServer(api)

	go frontend.InitHTTPService(config, grpcServer)
	go metrics.InitPrometheusService(config)
	service.InitGRPCService(config, grpcServer)
}
