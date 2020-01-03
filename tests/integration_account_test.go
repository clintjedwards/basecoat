package tests

import (
	"context"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func (info *testHarness) TestCreateAccount(t *testing.T) {
	t.Run("CreateAccount", func(t *testing.T) {

		createAccountRequest := &api.CreateAccountRequest{
			Id:       "test",
			Password: "test",
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createResponse, err := info.client.CreateAccount(ctx, createAccountRequest)
		require.NoError(t, err)
		require.NotNil(t, createResponse)
	})
}

func (info *testHarness) TestGetAccount(t *testing.T) {
	t.Run("GetAccount", func(t *testing.T) {

		expectedResponse := &api.GetAccountResponse{
			Account: &api.Account{
				Id: "test",
			},
		}

		getAccountRequest := &api.GetAccountRequest{
			Id: "test",
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		getResponse, err := info.client.GetAccount(ctx, getAccountRequest)
		require.NoError(t, err)
		require.NotNil(t, getResponse)
		require.Equal(t, getResponse.Account.Id, expectedResponse.Account.Id)
	})
}

func (info *testHarness) TestListAccounts(t *testing.T) {
	t.Run("ListAccounts", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		listResponse, err := info.client.ListAccounts(ctx, &api.ListAccountsRequest{})
		require.NoError(t, err)
		require.NotNil(t, listResponse)
		require.NotEmpty(t, listResponse.Accounts)
		require.Equal(t, listResponse.Accounts["test"].Id, "test")
	})
}

func (info *testHarness) TestDisableAccount(t *testing.T) {
	t.Run("DisableAccount", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		disableResponse, err := info.client.DisableAccount(ctx, &api.DisableAccountRequest{Id: "test"})
		require.NoError(t, err)
		require.NotNil(t, disableResponse)

		getResponse, err := info.client.GetAccount(ctx, &api.GetAccountRequest{Id: "test"})
		require.NoError(t, err)
		require.NotNil(t, disableResponse)
		require.Equal(t, getResponse.Account.State.String(), "DISABLED")
	})
}

func (info *testHarness) TestUpdateAccount(t *testing.T) {
	t.Run("UpdateAccount", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		getResponse, err := info.client.GetAccount(ctx, &api.GetAccountRequest{Id: "test"})
		require.NoError(t, err)
		require.NotNil(t, getResponse)

		updateAccountRequest := &api.UpdateAccountRequest{
			Id:    "test",
			Hash:  getResponse.Account.Hash,
			State: api.UpdateAccountRequest_ACTIVE,
		}

		updateResponse, err := info.client.UpdateAccount(ctx, updateAccountRequest)
		require.NoError(t, err)
		require.NotNil(t, updateResponse)

		getResponse, err = info.client.GetAccount(ctx, &api.GetAccountRequest{Id: "test"})
		require.NoError(t, err)
		require.NotNil(t, getResponse)
		require.Equal(t, "ACTIVE", getResponse.Account.State.String())
	})
}
