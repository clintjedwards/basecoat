package api

import (
	"context"
	"strings"

	"github.com/clintjedwards/basecoat/proto"
)

var appVersion = "0.0.dev_000000"

func parseVersion(versionString string) (version, commit string) {
	version, commit, err := strings.Cut(versionString, "_")
	if !err {
		return "", ""
	}

	return
}

// GetSystemInfo returns system information and health
func (api *API) GetSystemInfo(_ context.Context, _ *proto.GetSystemInfoRequest) (*proto.GetSystemInfoResponse, error) {
	version, commit := parseVersion(appVersion)

	return &proto.GetSystemInfoResponse{
		Commit:          commit,
		FrontendEnabled: api.config.Frontend.Enable,
		Semver:          version,
	}, nil
}
