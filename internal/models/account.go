package models

import (
	"time"

	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/lithammer/shortuuid/v4"
)

type AccountState string

const (
	AccountStateUnknown  AccountState = "UNKNOWN"
	AccountStateActive   AccountState = "ACTIVE"
	AccountStateDisabled AccountState = "DISABLED"
)

// Accounts are used to divide users of Basecoat. It is the highest level unit and things like
// formulas and jobs belong to specific accounts.
type Account struct {
	ID       string       `json:"id"`       // Unique identifier
	Name     string       `json:"name"`     // Humanized name; great for reading from UIs.
	Hash     string       `json:"hash"`     // Password hash
	State    AccountState `json:"state"`    // Whether the account is disabled or not.
	Created  int64        `json:"created"`  // The creation time in epoch milli.
	Modified int64        `json:"modified"` // The modified time in epoch milli;
}

func NewAccount(name, hash string) *Account {
	newAccount := &Account{
		ID:       shortuuid.New()[0:7],
		Name:     name,
		Hash:     hash,
		State:    AccountStateActive,
		Created:  time.Now().UnixMilli(),
		Modified: 0,
	}

	return newAccount
}

func (a *Account) ToProto() *proto.Account {
	return &proto.Account{
		Id:       a.ID,
		Name:     a.Name,
		State:    proto.AccountState(proto.AccountState_value[string(a.State)]),
		Created:  a.Created,
		Modified: a.Modified,
	}
}

func (a *Account) ToStorage() *storage.Account {
	return &storage.Account{
		ID:       a.ID,
		Name:     a.Name,
		Hash:     a.Hash,
		State:    string(a.State),
		Created:  a.Created,
		Modified: a.Modified,
	}
}

func (a *Account) FromStorage(s *storage.Account) {
	a.ID = s.ID
	a.Name = s.Name
	a.Hash = s.Hash
	a.State = AccountState(s.State)
	a.Created = s.Created
	a.Modified = s.Modified
}
