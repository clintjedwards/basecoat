package service

import (
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/search"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/toolkit/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// API represents a basecoat grpc backend service
type API struct {
	storage storage.Engine
	config  *config.Config
	search  *search.Search
}

// NewBasecoatAPI inits a grpc basecoat api service
func NewBasecoatAPI(config *config.Config) *API {
	basecoatAPI := API{}

	storage, err := storage.InitStorage()
	if err != nil {
		logger.Log().Fatalw("failed to initialize storage", "error", err)
	}

	searchIndex, err := search.InitSearch()
	if err != nil {
		logger.Log().Fatalw("failed to initialize search indexes", "error", err)
	}

	go searchIndex.BuildIndex(storage)

	basecoatAPI.config = config
	basecoatAPI.storage = storage
	basecoatAPI.search = searchIndex

	return &basecoatAPI
}

// CreateGRPCServer creates a grpc server with all the proper settings; TLS enabled
func CreateGRPCServer(basecoatAPI *API) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(basecoatAPI.authenticate),
			grpc_prometheus.UnaryServerInterceptor,
		)),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()

	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)
	api.RegisterBasecoatServer(grpcServer, basecoatAPI)

	return grpcServer
}
