package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Create a test which tests the admin key
func (info *testHarness) TestCreateAPIToken(t *testing.T) {
	t.Run("CreateAPIToken", func(t *testing.T) {

		var opts []grpc.DialOption

		creds, err := credentials.NewClientTLSFromFile("../localhost.crt", "")
		require.NoError(t, err)

		opts = append(opts, grpc.WithTransportCredentials(creds))

		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
		require.NoError(t, err)
		defer conn.Close()

		basecoatClient := api.NewBasecoatClient(conn)

		createAPITokenRequest := &api.CreateAPITokenRequest{
			User:     "test",
			Password: "test",
			Duration: 1200,
		}

		createResponse, err := basecoatClient.CreateAPIToken(context.Background(), createAPITokenRequest)
		require.NoError(t, err)

		info.apikey = createResponse.Key
	})
}
