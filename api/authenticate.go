package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// RegisterAuthRoutes adds /authenticate to routing list
func (restAPI *API) RegisterAuthRoutes(router *mux.Router) {
	router.Handle("/auth/api", handlers.MethodHandler{
		"GET": restAPI.requireAuthentication(http.HandlerFunc(restAPI.authenticateViaAPIHandler)),
	})
}

func (restAPI *API) validateAuthToken(token string) bool {

	config, err := config.FromEnv()
	if err != nil {
		return false
	}

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	loginInfo := strings.Split(string(decodedToken), ":")

	if loginInfo[1] != config.SecretKey {
		return false
	}

	return true
}

func (restAPI *API) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token := request.Header.Get("Authorization")

		if token == "" {
			utils.StructuredLog("info", "token or auth header missing", nil)
			utils.SendResponse(writer, http.StatusUnauthorized, "token or auth header missing", true)
			return
		}

		authenticated := restAPI.validateAuthToken(token)

		if !authenticated {
			utils.StructuredLog("info", "could not authenticate user", nil)
			utils.SendResponse(writer, http.StatusUnauthorized, "could not authenticate user", true)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (restAPI *API) authenticateViaAPIHandler(w http.ResponseWriter, req *http.Request) {
	utils.SendResponse(w, http.StatusOK, "authenticated", false)
	return
}
