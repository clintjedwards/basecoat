package service

import (
	"context"
	"time"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/utils"
	jwt "github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

var contextAccount = contextKey("account")

// CreateAPIToken returns a temporary api key that can be used on all subsequent requests
func (basecoat *API) CreateAPIToken(context context.Context, request *api.CreateAPITokenRequest) (*api.CreateAPITokenResponse, error) {

	if request.User == "" {
		return &api.CreateAPITokenResponse{}, status.Error(codes.FailedPrecondition, "user name required")
	}

	if request.Password == "" {
		return &api.CreateAPITokenResponse{}, status.Error(codes.FailedPrecondition, "password required")
	}

	// Limit length of duration requests; 946708560 = 30 years
	if request.Duration > 946708560 {
		return &api.CreateAPITokenResponse{}, status.Error(codes.FailedPrecondition, "duration request is too long; greater than 30 years")
	}

	user, err := basecoat.storage.GetUser(request.User)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return &api.CreateAPITokenResponse{}, status.Error(codes.NotFound, "could not authenticate user")
		}
		utils.StructuredLog(utils.LogLevelError, "could not authenticate user", err)
		return &api.CreateAPITokenResponse{}, status.Error(codes.Internal, "could not authenticate user; internal error")
	}

	if !utils.CheckPasswordHash([]byte(user.Hash), []byte(request.Password)) {
		return &api.CreateAPITokenResponse{}, status.Error(codes.NotFound, "could not authenticate user")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Name,
		"expiry":   int64(time.Now().Unix() + request.Duration),
	})

	tokenString, error := token.SignedString([]byte(basecoat.config.Backend.SecretKey))
	if error != nil {
		utils.StructuredLog(utils.LogLevelError, "could not sign jwt token", err)
		return &api.CreateAPITokenResponse{}, status.Error(codes.Internal, "could not authenticate user; internal error")
	}

	utils.StructuredLog(utils.LogLevelInfo, "api token created", request.User)
	return &api.CreateAPITokenResponse{Key: tokenString}, nil
}

func (basecoat *API) authenticate(ctx context.Context) (context.Context, error) {

	// Exclude the route to get the API token
	method, _ := grpc.Method(ctx)
	if method == "/api.Basecoat/CreateAPIToken" {
		return ctx, nil
	}

	// Exclude the route to get system information
	if method == "/api.Basecoat/GetSystemInfo" {
		return ctx, nil
	}

	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return ctx, err
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return ctx, err
		}

		return []byte(basecoat.config.Backend.SecretKey), nil
	})
	if err != nil {
		utils.StructuredLog(utils.LogLevelError, "could not decode jwt token", nil)
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token: %v", nil)
	}

	if _, present := jwtToken.Claims.(jwt.MapClaims); !present {
		utils.StructuredLog(utils.LogLevelError, "could not verify jwt token", nil)
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token: %v", nil)
	}

	if !jwtToken.Valid {
		utils.StructuredLog(utils.LogLevelError, "could not verify jwt token; not valid", nil)
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token: %v", nil)
	}

	claims := jwtToken.Claims.(jwt.MapClaims)

	if _, present := claims["username"]; !present {
		utils.StructuredLog(utils.LogLevelError, "misformatted jwt token; missing username", nil)
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token: %v", nil)
	}
	if _, present := claims["expiry"]; !present {
		utils.StructuredLog(utils.LogLevelError, "misformatted jwt token; missing expiry", nil)
		return ctx, grpc.Errorf(codes.Unauthenticated, "could not decode token: %v", nil)
	}

	expiry := int64(claims["expiry"].(float64))
	if time.Now().Unix() > expiry && expiry != 0 {
		utils.StructuredLog(utils.LogLevelInfo, "token has expired",
			map[string]int64{"current_time": time.Now().Unix(), "expiry_time": expiry})
		return ctx, grpc.Errorf(codes.Unauthenticated, "token has expired: %v", time.Unix(expiry, 0).UTC())
	}

	newCtx := context.WithValue(ctx, contextAccount, claims["username"].(string))
	return newCtx, nil
}

// getAccountFromContext gets the account name string from the context
func getAccountFromContext(ctx context.Context) (string, bool) {
	account, present := ctx.Value(contextAccount).(string)
	return account, present
}
