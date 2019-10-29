package service

import (
	"log"
	"net"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/search"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/basecoat/utils"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
		utils.StructuredLog(utils.LogLevelFatal, "failed to initialize storage", err)
	}

	searchIndex, err := search.InitSearch()
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to initialize search functions", err)
	}

	go searchIndex.BuildIndex()

	basecoatAPI.config = config
	basecoatAPI.storage = storage
	basecoatAPI.search = searchIndex

	return &basecoatAPI
}

// CreateGRPCServer creates a grpc server with all the proper settings; TLS enabled
func CreateGRPCServer(basecoatAPI *API) *grpc.Server {

	creds, err := credentials.NewServerTLSFromFile(basecoatAPI.config.TLSCertPath, basecoatAPI.config.TLSKeyPath)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get certificates", err)
	}

	serverOption := grpc.Creds(creds)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(basecoatAPI.authenticate),
			grpc_prometheus.UnaryServerInterceptor,
		)),
		serverOption,
	)

	grpc_prometheus.EnableHandlingTimeHistogram()

	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)
	api.RegisterBasecoatServer(grpcServer, basecoatAPI)

	return grpcServer
}

// InitGRPCService starts a GPRC server
func InitGRPCService(config *config.Config, server *grpc.Server) {

	listen, err := net.Listen("tcp", config.Backend.GRPCURL)
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "could not initialize tcp listener", err)
	}

	utils.StructuredLog(utils.LogLevelInfo, "starting basecoat grpc service",
		map[string]string{"url": config.Backend.GRPCURL})

	log.Fatal(server.Serve(listen))
}
