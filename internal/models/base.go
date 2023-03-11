package models

import (
	"time"

	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/lithammer/shortuuid/v4"
)

// A base is the starting paint mix before significant color is added.
type BaseMetadata struct {
	Account      string `json:"account"`      // Account ID this base belongs to.
	ID           string `json:"id"`           // Unique identifier;
	Label        string `json:"label"`        // Humanized name; great for reading from UIs.
	Manufacturer string `json:"manufacturer"` // Company name who created the base.
	Created      int64  `json:"created"`      // The creation time in epoch milli.
}

func NewBaseMetadata(account, label, manufacturer string) *BaseMetadata {
	newBaseMetadata := &BaseMetadata{
		Account:      account,
		ID:           shortuuid.New()[0:7],
		Label:        label,
		Manufacturer: manufacturer,
		Created:      time.Now().UnixMilli(),
	}

	return newBaseMetadata
}

func (b *BaseMetadata) ToProto() *proto.BaseMetadata {
	return &proto.BaseMetadata{
		Account:      b.Account,
		Id:           b.ID,
		Label:        b.Label,
		Manufacturer: b.Manufacturer,
		Created:      b.Created,
	}
}

func (b *BaseMetadata) ToStorage() *storage.Base {
	return &storage.Base{
		Account:      b.Account,
		ID:           b.ID,
		Label:        b.Label,
		Manufacturer: b.Manufacturer,
		Created:      b.Created,
	}
}

func (b *BaseMetadata) FromStorage(s *storage.Base) {
	b.Account = s.Account
	b.ID = s.ID
	b.Label = s.Label
	b.Manufacturer = s.Manufacturer
	b.Created = s.Created
}

// A base is the starting paint mix before significant color is added.
// A FormulaBase is metadata about a formula and base relationship.
type FormulaBase struct {
	Formula string `json:"account"` // Unique ID for formula.
	Base    string `json:"base"`    // Unique ID for base.
	Amount  string `json:"amount"`
}

func NewFormulaBase(formula, base, amount string) *FormulaBase {
	newFormulaBase := &FormulaBase{
		Formula: formula,
		Base:    base,
		Amount:  amount,
	}

	return newFormulaBase
}

func (fb *FormulaBase) ToProto() *proto.FormulaBase {
	return &proto.FormulaBase{
		Formula: fb.Formula,
		Base:    fb.Base,
		Amount:  fb.Amount,
	}
}

func (fb *FormulaBase) ToStorage() *storage.FormulaBase {
	return &storage.FormulaBase{
		Formula: fb.Formula,
		Base:    fb.Base,
		Amount:  fb.Amount,
	}
}

func (fb *FormulaBase) FromStorage(s *storage.FormulaBase) {
	fb.Formula = s.Formula
	fb.Base = s.Base
	fb.Amount = s.Amount
}
