package models

import (
	"time"

	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/lithammer/shortuuid/v4"
)

// A formula is a combination of paint bases and colorants to make a specific color for a particular customer.
type Formula struct {
	Metadata  FormulaMetadata   `json:"metadata"`
	Bases     []FormulaBase     `json:"bases"`
	Colorants []FormulaColorant `json:"colorants"`
}

// A formula is a combination of paint bases and colorants to make a specific color for a particular customer.
// Formula metadata is information about that specific formula without the core data like bases, colorants, and jobs.
type FormulaMetadata struct {
	Account  string `json:"account"` // Account identifier
	ID       string `json:"id"`      // Unique identifier
	Name     string `json:"name"`    // Humanized name; great for reading from UIs.
	Number   string `json:"number"`  // Some formulas have specific reference numbers from their manufacturers
	Notes    string `json:"notes"`
	Created  int64  `json:"created"`  // The creation time in epoch milli.
	Modified int64  `json:"modified"` // The modified time in epoch milli;
}

func NewFormulaMetadata(account, name string) *FormulaMetadata {
	NewFormulaMetadata := &FormulaMetadata{
		Account:  account,
		ID:       shortuuid.New()[0:7],
		Name:     name,
		Number:   "",
		Notes:    "",
		Created:  time.Now().UnixMilli(),
		Modified: 0,
	}

	return NewFormulaMetadata
}

func (f *FormulaMetadata) ToProto() *proto.FormulaMetadata {
	return &proto.FormulaMetadata{
		Account:  f.Account,
		Id:       f.ID,
		Name:     f.Name,
		Number:   f.Number,
		Notes:    f.Notes,
		Created:  f.Created,
		Modified: f.Modified,
	}
}

func (f *FormulaMetadata) ToStorage() *storage.Formula {
	return &storage.Formula{
		Account:  f.Account,
		ID:       f.ID,
		Name:     f.Name,
		Number:   f.Number,
		Notes:    f.Notes,
		Created:  f.Created,
		Modified: f.Modified,
	}
}

func (f *FormulaMetadata) FromStorage(s *storage.Formula) {
	f.Account = s.Account
	f.ID = s.ID
	f.Name = s.Name
	f.Number = s.Number
	f.Notes = s.Notes
	f.Created = s.Created
	f.Modified = s.Modified
}
