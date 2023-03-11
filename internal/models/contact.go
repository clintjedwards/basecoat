package models

import (
	"time"

	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/lithammer/shortuuid/v4"
)

// A contact is the details for a specified person or point of contact.
type Contact struct {
	Account  string `json:"account"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

func NewContact(account, name string) *Contact {
	NewContact := &Contact{
		Account:  account,
		ID:       shortuuid.New()[0:7],
		Name:     name,
		Email:    "",
		Phone:    "",
		Created:  time.Now().UnixMilli(),
		Modified: 0,
	}

	return NewContact
}

func (c *Contact) ToProto() *proto.Contact {
	return &proto.Contact{
		Account:  c.Account,
		Id:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		Phone:    c.Phone,
		Created:  c.Created,
		Modified: c.Modified,
	}
}

func (c *Contact) ToStorage() *storage.Contact {
	return &storage.Contact{
		Account:  c.Account,
		ID:       c.ID,
		Name:     c.Name,
		Email:    c.Email,
		Phone:    c.Phone,
		Created:  c.Created,
		Modified: c.Modified,
	}
}

func (c *Contact) FromStorage(s *storage.Contact) {
	c.Account = s.Account
	c.ID = s.ID
	c.Name = s.Name
	c.Email = s.Email
	c.Phone = s.Phone
	c.Created = s.Created
	c.Modified = s.Modified
}
