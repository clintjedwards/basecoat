package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Create a test which tests the admin key
func (info *testHarness) TestCreateAPIToken(t *testing.T) {
	t.Run("CreateAPIToken", func(t *testing.T) {

		var opts []grpc.DialOption

		creds, err := credentials.NewClientTLSFromFile("../localhost.crt", "")
		if err != nil {
			log.Fatalf("failed to get certificates: %v", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))

		conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
		if err != nil {
			log.Fatalf("could not connect to basecoat: %v", err)
		}
		defer conn.Close()

		basecoatClient := api.NewBasecoatClient(conn)

		createAPITokenRequest := &api.CreateAPITokenRequest{
			User:     "test",
			Password: "test",
			Duration: 1200,
		}

		createResponse, err := basecoatClient.CreateAPIToken(context.Background(), createAPITokenRequest)
		if err != nil {
			log.Fatalf("could not create token: %v", err)
		}

		info.apikey = createResponse.Key
	})
}
