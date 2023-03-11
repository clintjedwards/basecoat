// Package api controls the bulk of the Basecoat API logic.
package api

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/clintjedwards/basecoat/internal/config"
	"github.com/clintjedwards/basecoat/internal/metrics"
	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/search"
	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func ptr[T any](v T) *T {
	return &v
}

// API represents a basecoat grpc backend service
type API struct {
	// Storage represents the main backend storage implementation. Basecoat stores most of its critical state information
	// using this storage mechanism.
	db storage.DB

	// Config represents the relative configuration for the Basecoat API. This is a combination of envvars and config values
	// gleaned at startup time.
	config *config.API

	search *search.Search

	// We opt out of forward compatibility with this embedded interface. This is required by GRPC.
	//
	// We don't embed the "proto.UnimplementedBasecoatServer" as there should never(I assume this will come back to bite me)
	// be an instance where we add proto methods without also updating the server to support those methods.
	// There is the added benefit that without it embedded we get compile time errors when a function isn't correctly
	// implemented. Saving us from weird "Unimplemented" RPC bugs.
	proto.UnsafeBasecoatServer
}

// NewAPI creates a new instance of the main Basecoat API service.
func NewAPI(config *config.API, db storage.DB) (*API, error) {
	api := API{}

	searchIndex, err := search.InitSearch(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize search index: %w", err)
	}

	rebuildTime := time.Duration(config.SearchIndexRebuildTime) * time.Second
	go func() {
		searchIndex.BuildIndex(db)
		time.Sleep(rebuildTime)
	}()

	api.config = config
	api.db = db
	api.search = searchIndex

	// For dev mode we auto create an account which can be used for development purposes.
	if config.Development.AutoCreateAccount {
		err = api.db.InsertAccount(api.db, &storage.Account{
			ID:       "dev",
			Name:     "Development Account",
			Hash:     "$2a$14$erzdfpk.ZTGNQLoAwkpoFu8dN9dAFhB/9I9uuwcPOaRUmKJK4ZsBC",
			State:    string(models.AccountStateActive),
			Created:  time.Now().UnixMilli(),
			Modified: 0,
		})
		if err != nil {
			if !errors.Is(err, storage.ErrEntityExists) {
				return nil, err
			}
		}
		log.Warn().Str("id", "dev").Str("password", "test").
			Msg("development config auto_create_account activated; created testing account;")
	}

	return &api, nil
}

// StartAPIService starts the Basecoat API service and blocks until a SIGINT or SIGTERM is received.
func (api *API) StartAPIService() {
	grpcServer, err := api.createGRPCServer()
	if err != nil {
		log.Fatal().Err(err).Msg("could not create GRPC service")
	}

	tlsConfig, err := api.generateTLSConfig(api.config.Server.TLSCertPath, api.config.Server.TLSKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get proper TLS config")
	}

	httpServer := wrapGRPCServer(api.config, grpcServer)
	httpServer.TLSConfig = tlsConfig

	go metrics.InitPrometheusService(api.config.Metrics.Endpoint)

	if api.config.Development.UseLocalhostTLS {
		log.Warn().Msg("loaded localhost TLS due to development config use_localhost_tls")
	}

	// Run our server in a goroutine and listen for signals that indicate graceful shutdown
	go func() {
		if err := httpServer.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server exited abnormally")
		}
	}()
	log.Info().Str("url", api.config.Server.Host).Msg("started basecoat grpc/http service")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

	// Doesn't block if no connections, otherwise will wait until the timeout deadline or connections to finish,
	// whichever comes first.
	ctx, cancel := context.WithTimeout(context.Background(), api.config.Server.ShutdownTimeout) // shutdown gracefully
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not shutdown server in timeout specified")
		return
	}

	log.Info().Msg("grpc server exited gracefully")
}

// The logging middleware has to be run before the final call to return the request.
// This is because we wrap the responseWriter to gain information from it after it
// has been written to.
// To speed this process up we call Serve as soon as possible and log afterwards.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		log.Debug().Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status_code", ww.Status()).
			Int("response_size_bytes", ww.BytesWritten()).
			Dur("elapsed_ms", time.Since(start)).
			Msg("")
	})
}

// wrapGRPCServer returns a combined grpc/http (grpc-web compatible) service with all proper settings;
// Rather than going through the trouble of setting up a separate proxy and extra for the service in order to server http/grpc/grpc-web
// this keeps things simple by enabling the operator to deploy a single binary and serve them all from one endpoint.
// This reduces operational burden, configuration headache and overall just makes for a better time for both client and operator.
func wrapGRPCServer(config *config.API, grpcServer *grpc.Server) *http.Server {
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	router := mux.NewRouter()

	// Define GRPC/HTTP request detection middleware
	GRPCandHTTPHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Content-Type"), "application/grpc") || wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
		} else {
			router.ServeHTTP(resp, req)
		}
	})

	httpServer := http.Server{
		Addr:    config.Server.Host,
		Handler: loggingMiddleware(GRPCandHTTPHandler),
		// Timeouts set here unfortunately also apply to the backing GRPC server. Because GRPC might have long running calls
		// we have to set these to 0 or a very high number. This creates an issue where running the frontend in this configuration
		// could possibly open us up to DOS attacks where the client holds the request open for long periods of time. To mitigate
		// this we both implement timeouts for routes on both the GRPC side and the pure HTTP side.
		WriteTimeout: 0,
		ReadTimeout:  0,
	}

	return &httpServer
}

// createGRPCServer creates the basecoat grpc server with all the proper settings; TLS enabled.
func (api *API) createGRPCServer() (*grpc.Server, error) {
	tlsConfig, err := api.generateTLSConfig(api.config.Server.TLSCertPath, api.config.Server.TLSKeyPath)
	if err != nil {
		return nil, err
	}

	panicHandler := func(p interface{}) (err error) {
		log.Error().Err(err).Interface("panic", p).Msg("server has encountered a fatal error")
		log.Error().Msg(string(debug.Stack()))
		return status.Errorf(codes.Unknown, "server has encountered a fatal error and could not process request")
	}

	grpcServer := grpc.NewServer(
		// recovery should always be first
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler)),
				grpc_auth.UnaryServerInterceptor(api.authenticate),
				grpc_prometheus.UnaryServerInterceptor,
			),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler)),
				grpc_auth.StreamServerInterceptor(api.authenticate),
			),
		),

		// Handle TLS
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()

	reflection.Register(grpcServer)
	grpc_prometheus.Register(grpcServer)
	proto.RegisterBasecoatServer(grpcServer, api)

	return grpcServer, nil
}

// grpcDial establishes a connection with the request URL via GRPC.
func grpcDial(url string) (*grpc.ClientConn, error) {
	host, port, ok := strings.Cut(url, ":")
	if !ok {
		return nil, fmt.Errorf("could not parse url %q; format should be <host>:<port>", url)
	}

	var opt []grpc.DialOption
	var tlsConf *tls.Config

	// If we're testing in development bypass the cert checks.
	if host == "localhost" || host == "127.0.0.1" {
		tlsConf = &tls.Config{
			InsecureSkipVerify: true,
		}
		opt = append(opt, grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)))
	}

	opt = append(opt, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(time.Millisecond*100)))))

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), opt...)
	if err != nil {
		return nil, fmt.Errorf("could not connect to server: %w", err)
	}

	return conn, nil
}
