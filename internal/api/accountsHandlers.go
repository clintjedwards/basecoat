package api

import (
	"context"
	"fmt"
	"time"

	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAccount returns a single formula by key
func (api *API) GetAccount(_ context.Context, request *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	if request.Id == "" {
		return &proto.GetAccountResponse{}, status.Error(codes.FailedPrecondition, "account id required")
	}

	accountRaw, err := api.db.GetAccount(api.db, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.GetAccountResponse{}, status.Error(codes.NotFound, "account requested not found")
		}
		return &proto.GetAccountResponse{}, status.Error(codes.Internal, "failed to retrieve account from database")
	}

	account := models.Account{}
	account.FromStorage(&accountRaw)

	return &proto.GetAccountResponse{Account: account.ToProto()}, nil
}

// ListAccounts returns a list of all accounts on basecoat service
func (api *API) ListAccounts(_ context.Context, _ *proto.ListAccountsRequest) (*proto.ListAccountsResponse, error) {
	accounts, err := api.db.ListAccounts(api.db, 0, 0)
	if err != nil {
		return &proto.ListAccountsResponse{}, status.Error(codes.Internal, "failed to retrieve accounts from database")
	}

	protoAccounts := []*proto.Account{}
	for _, accountRaw := range accounts {
		var account models.Account
		account.FromStorage(&accountRaw)
		protoAccounts = append(protoAccounts, account.ToProto())
	}

	return &proto.ListAccountsResponse{Accounts: protoAccounts}, nil
}

// CreateAccount registers a new account
func (api *API) CreateAccount(_ context.Context, request *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	if request.Name == "" {
		return &proto.CreateAccountResponse{}, status.Error(codes.FailedPrecondition, "account name required")
	}

	if request.Password == "" {
		return &proto.CreateAccountResponse{}, status.Error(codes.FailedPrecondition, "account password required")
	}

	if len(request.Password) > 72 {
		return &proto.CreateAccountResponse{}, status.Error(codes.FailedPrecondition,
			"account password not allowed; password must be less than 72 chars")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return &proto.CreateAccountResponse{}, status.Error(codes.FailedPrecondition, "could not process password")
	}

	account := models.NewAccount(request.Name, string(hash))

	err = api.db.InsertAccount(api.db, account.ToStorage())
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.CreateAccountResponse{}, status.Error(codes.AlreadyExists, "could not save account; account already exists")
		}
		log.Error().Err(err).Msg("could not save account")
		return &proto.CreateAccountResponse{}, status.Error(codes.Internal, "could not save account")
	}

	log.Info().Str("id", account.ID).Str("name", account.Name).Msg("account created")
	return &proto.CreateAccountResponse{
		Account: account.ToProto(),
	}, nil
}

// UpdateAccount registers a new account
func (api *API) UpdateAccount(_ context.Context, request *proto.UpdateAccountRequest) (*proto.UpdateAccountResponse, error) {
	if request.Id == "" {
		return &proto.UpdateAccountResponse{}, status.Error(codes.FailedPrecondition, "account id required")
	}

	account, err := api.db.GetAccount(api.db, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.UpdateAccountResponse{}, status.Error(codes.NotFound, "account requested not found")
		}
		return &proto.UpdateAccountResponse{}, status.Error(codes.Internal, "failed to retrieve account from database")
	}

	var name *string

	if request.Name != "" {
		name = &request.Name
	}

	var hash *string

	if request.Password != "" {
		hashBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
		if err != nil {
			return &proto.UpdateAccountResponse{}, status.Error(codes.FailedPrecondition, "could not process password")
		}

		hash = ptr(string(hashBytes))
	}

	err = api.db.UpdateAccount(api.db, request.Id, storage.UpdatableAccountFields{
		Name:     name,
		Hash:     hash,
		Modified: ptr(time.Now().UnixMilli()),
	})
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.UpdateAccountResponse{}, status.Error(codes.NotFound, "account requested not found")
		}
		log.Error().Err(err).Msg("could not save account")
		return &proto.UpdateAccountResponse{}, status.Error(codes.Internal, "could not save account")
	}

	log.Debug().Str("id", account.ID).Str("name", account.Name).Msg("account created")
	return &proto.UpdateAccountResponse{}, nil
}

// ToggleAccountState enables or disables an account depending on what it's original state was.
func (api *API) ToggleAccountState(_ context.Context, request *proto.ToggleAccountStateRequest) (*proto.ToggleAccountStateResponse, error) {
	account := models.Account{}
	newState := models.AccountStateUnknown

	err := storage.InsideTx(api.db.DB, func(*sqlx.Tx) error {
		accountRaw, err := api.db.GetAccount(api.db, request.Id)
		if err != nil {
			return fmt.Errorf("could not get account: %w", err)
		}

		account.FromStorage(&accountRaw)

		if account.State == models.AccountStateDisabled {
			newState = models.AccountStateActive
		} else {
			newState = models.AccountStateDisabled
		}

		err = api.db.UpdateAccount(api.db, request.Id, storage.UpdatableAccountFields{
			State: ptr(string(newState)),
		})
		if err != nil {
			return fmt.Errorf("could not update account: %w", err)
		}

		return nil
	})
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return nil, status.Error(codes.NotFound, "account requested not found")
		}

		if err == storage.ErrEntityExists {
			return nil, status.Error(codes.AlreadyExists, "could not save account; account already exists")
		}

		return nil, status.Error(codes.Internal, "failed to retrieve account from database")
	}

	account.State = newState

	return &proto.ToggleAccountStateResponse{
		State: account.ToProto().State,
	}, nil
}
