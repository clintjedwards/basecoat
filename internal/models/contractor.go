package models

import (
	"time"

	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/lithammer/shortuuid/v4"
)

// A contractor is a company or individual that does a specific job.
// Contractor metadata is information about that specific contractor.
type Contractor struct {
	Account  string  `json:"account"`  // Account identifier
	ID       string  `json:"id"`       // Unique identifier
	Company  string  `json:"company"`  // Humanized name; great for reading from UIs.
	Contact  *string `json:"contact"`  // contact ID
	Created  int64   `json:"created"`  // The creation time in epoch milli.
	Modified int64   `json:"modified"` // The modified time in epoch milli;
}

func NewContractor(account, company string) *Contractor {
	newContractor := &Contractor{
		Account:  account,
		ID:       shortuuid.New()[0:7],
		Company:  company,
		Contact:  nil,
		Created:  time.Now().UnixMilli(),
		Modified: 0,
	}

	return newContractor
}

func (f *Contractor) ToProto() *proto.Contractor {
	return &proto.Contractor{
		Account:  f.Account,
		Id:       f.ID,
		Company:  f.Company,
		Contact:  f.Contact,
		Created:  f.Created,
		Modified: f.Modified,
	}
}

func (f *Contractor) ToStorage() *storage.Contractor {
	return &storage.Contractor{
		Account:  f.Account,
		ID:       f.ID,
		Company:  f.Company,
		Contact:  f.Contact,
		Created:  f.Created,
		Modified: f.Modified,
	}
}

func (f *Contractor) FromStorage(s *storage.Contractor) {
	f.Account = s.Account
	f.ID = s.ID
	f.Company = s.Company
	f.Contact = s.Contact
	f.Created = s.Created
	f.Modified = s.Modified
}
