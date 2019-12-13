package service

import (
	"context"
	"strings"

	"github.com/clintjedwards/basecoat/api"
)

var appVersion = "v0.0.dev <commit>"

// GetSystemInfo returns system information and health
func (basecoat *API) GetSystemInfo(context context.Context, request *api.GetSystemInfoRequest) (*api.GetSystemInfoResponse, error) {

	versionTuple := strings.Split(appVersion, " ")

	return &api.GetSystemInfoResponse{
		DebugEnabled:    basecoat.config.Debug,
		FrontendEnabled: basecoat.config.Frontend.Enable,
		Version:         versionTuple[0],
		Commit:          versionTuple[1],
		DatabaseEngine:  basecoat.config.Database.Engine,
	}, nil
}
