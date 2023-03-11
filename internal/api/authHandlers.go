package api

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/clintjedwards/basecoat/proto"
	jwt "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

var contextAccount = contextKey("account")

var adminMethods = []string{
	"/api.Basecoat/GetAccount",
	"/api.Basecoat/CreateAccount",
	"/api.Basecoat/ListAccounts",
	"/api.Basecoat/UpdateAccount",
	"/api.Basecoat/ToggleAccountState",
}

var authlessMethods = []string{
	"/api.Basecoat/CreateAPIToken",
	"/api.Basecoat/GetSystemInfo",
}

// CreateAPIToken returns a temporary api key that can be used on all subsequent requests
func (api *API) CreateAPIToken(_ context.Context, request *proto.CreateAPITokenRequest) (*proto.CreateAPITokenResponse, error) {
	if request.Account == "" || request.Password == "" {
		return &proto.CreateAPITokenResponse{}, status.Error(codes.FailedPrecondition, "id and password required")
	}

	// Limit length of duration requests
	if request.Duration > api.config.TokenDurationLimit {
		return &proto.CreateAPITokenResponse{}, status.Errorf(codes.FailedPrecondition,
			"duration request is too long; greater than %d seconds", api.config.TokenDurationLimit)
	}

	accountRaw, err := api.db.GetAccount(api.db, request.Account)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			return &proto.CreateAPITokenResponse{}, status.Error(codes.NotFound, "could not authenticate account")
		}
		log.Error().Err(err).Str("account", request.Account).Msg("could not authenticate account")
		return &proto.CreateAPITokenResponse{}, status.Error(codes.Internal, "could not authenticate account; internal error")
	}

	account := models.Account{}
	account.FromStorage(&accountRaw)

	if account.State == models.AccountStateDisabled {
		return &proto.CreateAPITokenResponse{}, status.Error(codes.FailedPrecondition, "account is disabled")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Hash), []byte(request.Password))
	if err != nil {
		return &proto.CreateAPITokenResponse{}, status.Error(codes.NotFound, "could not authenticate account")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account.ID,
		"expiry":  int64(time.Now().Unix() + request.Duration),
	})

	tokenString, err := token.SignedString([]byte(api.config.EncryptionKey))
	if err != nil {
		log.Error().Err(err).Str("account", account.ID).Msg("could not sign jwt token")
		return &proto.CreateAPITokenResponse{}, status.Error(codes.Internal, "could not authenticate account; internal error")
	}

	log.Info().Str("account", account.ID).Msg("api token created")
	return &proto.CreateAPITokenResponse{Key: tokenString}, nil
}

// authenticate is run on every call to verify if the user is allowed to access a given rpc
func (api *API) authenticate(ctx context.Context) (context.Context, error) {
	method, _ := grpc.Method(ctx)

	// Exclude routes that don't need authentication
	for _, route := range authlessMethods {
		if method == route {
			return ctx, nil
		}
	}

	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return ctx, err
	}

	// Specially handle admin routes
	for _, route := range adminMethods {
		if method == route {
			if api.config.Development.BypassAuth {
				log.Debug().Msg("admin route accessed due to bypass_auth config set to true")
				return ctx, nil
			}

			admin := handleAdminRoutes(token, api.config.AdminToken)
			if admin {
				log.Info().Str("method", method).Msg("admin route accessed")
				return ctx, err
			}
			log.Debug().Str("method", method).Msg("could not verify admin token")
			return ctx, status.Errorf(codes.Unauthenticated, "could not verify admin token")
		}
	}

	if api.config.Development.BypassAuth {
		newCtx := context.WithValue(ctx, contextAccount, "dev")
		log.Debug().Msg("automatically authed due to bypass_auth config set to true; authed as 'dev' account")
		return newCtx, nil
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return ctx, err
		}

		return []byte(api.config.EncryptionKey), nil
	})
	if err != nil {
		log.Error().Err(err).Msg("could not decode jwt token")
		return ctx, status.Errorf(codes.Unauthenticated, "could not decode token")
	}

	if _, present := jwtToken.Claims.(jwt.MapClaims); !present {
		log.Error().Msg("could not verify jwt token")
		return ctx, status.Errorf(codes.Unauthenticated, "could not decode token")
	}

	if !jwtToken.Valid {
		log.Error().Msg("could not verify jwt token; not valid")
		return ctx, status.Errorf(codes.Unauthenticated, "could not decode token")
	}

	claims := jwtToken.Claims.(jwt.MapClaims)

	if _, present := claims["account"]; !present {
		log.Error().Msg("misformatted jwt token; missing account")
		return ctx, status.Errorf(codes.Unauthenticated, "could not decode token")
	}
	if _, present := claims["expiry"]; !present {
		log.Error().Msg("misformatted jwt token; missing expiry")
		return ctx, status.Errorf(codes.Unauthenticated, "could not decode token")
	}

	expiry := int64(claims["expiry"].(float64))
	if time.Now().Unix() > expiry && expiry != 0 {
		log.Debug().Interface("user", claims["account"]).Int64("expiry_time", expiry).Msg("token has expired")

		return ctx, status.Errorf(codes.Unauthenticated, "token has expired: %v", time.Unix(expiry, 0).UTC())
	}

	newCtx := context.WithValue(ctx, contextAccount, claims["account"].(string))
	return newCtx, nil
}

// getAccountFromContext gets the account name string from the context
func getAccountFromContext(ctx context.Context) (string, bool) {
	account, present := ctx.Value(contextAccount).(string)
	return account, present
}

func handleAdminRoutes(token, adminKey string) bool {
	return token == adminKey
}
