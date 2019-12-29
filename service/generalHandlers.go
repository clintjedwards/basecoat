package service

import (
	"context"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/version"
)

var appVersion = "v0.0.dev_<build_time>_<commit>"

// GetSystemInfo returns system information and health
func (basecoat *API) GetSystemInfo(context context.Context, request *api.GetSystemInfoRequest) (*api.GetSystemInfoResponse, error) {

	info, err := version.Parse(appVersion)
	if err != nil {
		return nil, err
	}

	return &api.GetSystemInfoResponse{
		BuildTime:       info.Epoch,
		Commit:          info.Hash,
		DebugEnabled:    basecoat.config.Debug,
		FrontendEnabled: basecoat.config.Frontend.Enable,
		Semver:          info.Semver,
	}, nil
}
