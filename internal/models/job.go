package models

import (
	"encoding/json"
	"time"

	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/lithammer/shortuuid/v4"
	"github.com/rs/zerolog/log"
)

type Address struct {
	Street  string `json:"street"`
	Street2 string `json:"street2"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zipcode string `json:"zipcode"`
}

func (a *Address) ToProto() *proto.Address {
	return &proto.Address{
		Street:  a.Street,
		Street2: a.Street2,
		City:    a.City,
		State:   a.State,
		Zipcode: a.Zipcode,
	}
}

func (a *Address) FromProto(pa *proto.Address) {
	a.Street = pa.Street
	a.Street2 = pa.Street2
	a.City = pa.City
	a.State = pa.State
	a.Zipcode = pa.Zipcode
}

func (a *Address) ToJSON() string {
	jsonAddress, err := json.Marshal(a)
	if err != nil {
		log.Error().Err(err).Msg("could not encode address json")
	}

	return string(jsonAddress)
}

// A job is a specific instance of work where formulas might have been used.
// Job metadata is information about that specific job.
type Job struct {
	Account    string  `json:"account"`    // Account identifier
	Contractor string  `json:"contractor"` // Contractor identifier
	ID         string  `json:"id"`         // Unique identifier
	Name       string  `json:"name"`
	Address    Address `json:"address"`
	Notes      string  `json:"notes"`
	Contact    *string `json:"contact"`
	Created    int64   `json:"created"`  // The creation time in epoch milli.
	Modified   int64   `json:"modified"` // The modified time in epoch milli;
}

func NewJob(account, contractor, name string) *Job {
	NewJob := &Job{
		Account:    account,
		Contractor: contractor,
		ID:         shortuuid.New()[0:7],
		Name:       name,
		Address:    Address{},
		Notes:      "",
		Contact:    nil,
		Created:    time.Now().UnixMilli(),
		Modified:   0,
	}

	return NewJob
}

func (f *Job) ToProto() *proto.Job {
	return &proto.Job{
		Account:    f.Account,
		Id:         f.ID,
		Name:       f.Name,
		Address:    f.Address.ToProto(),
		Notes:      f.Notes,
		Contractor: f.Contractor,
		Contact:    f.Contact,
		Created:    f.Created,
		Modified:   f.Modified,
	}
}

func (f *Job) ToStorage() *storage.Job {
	return &storage.Job{
		Account:    f.Account,
		ID:         f.ID,
		Contractor: f.Contractor,
		Name:       f.Name,
		Address:    f.Address.ToJSON(),
		Notes:      f.Notes,
		Contact:    f.Contact,
		Created:    f.Created,
		Modified:   f.Modified,
	}
}

func (f *Job) FromStorage(s *storage.Job) {
	address := Address{}
	err := json.Unmarshal([]byte(s.Address), &address)
	if err != nil {
		log.Error().Err(err).Msg("could not decode address json")
	}

	f.Account = s.Account
	f.ID = s.ID
	f.Name = s.Name
	f.Address = address
	f.Notes = s.Notes
	f.Contractor = s.Contractor
	f.Contact = s.Contact
	f.Created = s.Created
	f.Modified = s.Modified
}
