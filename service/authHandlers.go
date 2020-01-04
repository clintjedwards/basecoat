package service

import (
	"context"
	"errors"
	"time"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/password"
	"github.com/clintjedwards/toolkit/tkerrors"
	"go.uber.org/zap"

	jwt "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
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
	"/api.Basecoat/DisableAccount",
}
var authlessMethods = []string{
	"/api.Basecoat/CreateAPIToken",
	"/api.Basecoat/GetSystemInfo",
}

// CreateAPIToken returns a temporary api key that can be used on all subsequent requests
func (bc *API) CreateAPIToken(ctx context.Context, request *api.CreateAPITokenRequest) (*api.CreateAPITokenResponse, error) {

	if request.User == "" || request.Password == "" {
		return &api.CreateAPITokenResponse{}, status.Error(codes.FailedPrecondition, "id and password required")
	}

	// Limit length of duration requests
	if request.Duration > bc.config.Backend.APITokenDurationLimit {
		return &api.CreateAPITokenResponse{}, status.Errorf(codes.FailedPrecondition,
			"duration request is too long; greater than %d seconds", bc.config.Backend.APITokenDurationLimit)
	}

	account, err := bc.storage.GetAccount(request.User)
	if err != nil {
		if errors.Is(err, tkerrors.ErrEntityNotFound) {
			return &api.CreateAPITokenResponse{}, status.Error(codes.NotFound, "could not authenticate account")
		}
		zap.S().Errorw("could not authenticate account", "error", err, "account", request.User)
		return &api.CreateAPITokenResponse{}, status.Error(codes.Internal, "could not authenticate account; internal error")
	}

	if account.State == api.Account_DISABLED {
		return &api.CreateAPITokenResponse{}, status.Error(codes.FailedPrecondition, "account is disabled")
	}

	if !password.CheckPasswordHash([]byte(account.Hash), []byte(request.Password)) {
		return &api.CreateAPITokenResponse{}, status.Error(codes.NotFound, "could not authenticate account")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account.Id,
		"expiry":  int64(time.Now().Unix() + request.Duration),
	})

	tokenString, err := token.SignedString([]byte(bc.config.Backend.SecretKey))
	if err != nil {
		zap.S().Errorw("could not sign jwt token", "error", err)
		return &api.CreateAPITokenResponse{}, status.Error(codes.Internal, "could not authenticate account; internal error")
	}

	zap.S().Infow("api token created", "account", request.User)
	return &api.CreateAPITokenResponse{Key: tokenString}, nil
}

// authenticate is run on every call to verify if the user is allowed to access a given rpc
func (bc *API) authenticate(ctx context.Context) (context.Context, error) {

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
			admin := handleAdminRoutes(token, bc.config.Backend.AdminToken)
			if admin {
				zap.S().Infow("admin route accessed", "method", method)
				return ctx, err
			}
			zap.S().Warnw("could not verify admin token", "method", method)
			return ctx, grpc.Errorf(codes.Unauthenticated, "could not verify admin token")
		}
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return ctx, err
		}

		return []byte(bc.config.Backend.SecretKey), nil
	})
	if err != nil {
		zap.S().Errorw("could not decode jwt token", "error", err)
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token")
	}

	if _, present := jwtToken.Claims.(jwt.MapClaims); !present {
		zap.S().Error("could not verify jwt token")
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token")
	}

	if !jwtToken.Valid {
		zap.S().Error("could not verify jwt token; not valid")
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token")
	}

	claims := jwtToken.Claims.(jwt.MapClaims)

	if _, present := claims["account"]; !present {
		zap.S().Error("misformatted jwt token; missing account")
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token")
	}
	if _, present := claims["expiry"]; !present {
		zap.S().Error("misformatted jwt token; missing expiry")
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token")
	}

	expiry := int64(claims["expiry"].(float64))
	if time.Now().Unix() > expiry && expiry != 0 {
		zap.S().Infow("token has expired",
			"user", claims["account"],
			"current_time", time.Now().Unix(),
			"expiry_time", expiry)

		return ctx, grpc.Errorf(codes.Unauthenticated, "token has expired: %v", time.Unix(expiry, 0).UTC())
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
	if token != adminKey {
		return false
	}

	return true
}
