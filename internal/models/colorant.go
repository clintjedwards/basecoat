package models

import (
	"time"

	"github.com/clintjedwards/basecoat/internal/storage"
	proto "github.com/clintjedwards/basecoat/proto"
	"github.com/lithammer/shortuuid/v4"
)

// A colorant is a pigment added to a base to create an specific overall color.
type ColorantMetadata struct {
	Account      string `json:"account"`      // Account ID this colorant belongs to.
	ID           string `json:"id"`           // Unique identifier;
	Label        string `json:"label"`        // Humanized name; great for reading from UIs.
	Manufacturer string `json:"manufacturer"` // Company name who created the colorant.
	Created      int64  `json:"created"`      // The creation time in epoch milli.
}

func NewColorantMetadata(account, label, manufacturer string) *ColorantMetadata {
	newColorantMetadata := &ColorantMetadata{
		Account:      account,
		ID:           shortuuid.New()[0:7],
		Label:        label,
		Manufacturer: manufacturer,
		Created:      time.Now().UnixMilli(),
	}

	return newColorantMetadata
}

func (b *ColorantMetadata) ToProto() *proto.ColorantMetadata {
	return &proto.ColorantMetadata{
		Account:      b.Account,
		Id:           b.ID,
		Label:        b.Label,
		Manufacturer: b.Manufacturer,
		Created:      b.Created,
	}
}

func (b *ColorantMetadata) ToStorage() *storage.Colorant {
	return &storage.Colorant{
		Account:      b.Account,
		ID:           b.ID,
		Label:        b.Label,
		Manufacturer: b.Manufacturer,
		Created:      b.Created,
	}
}

func (b *ColorantMetadata) FromStorage(s *storage.Colorant) {
	b.Account = s.Account
	b.ID = s.ID
	b.Label = s.Label
	b.Manufacturer = s.Manufacturer
	b.Created = s.Created
}

// A colorant is the pigment which is mixed in to give a base a specific color.
// A FormulaColorant is metadata about a formula and colorant relationship.
type FormulaColorant struct {
	Formula  string `json:"account"`  // Unique ID for formula.
	Colorant string `json:"colorant"` // Unique ID for colorant.
	Amount   string `json:"amount"`
}

func NewFormulaColorant(formula, colorant, amount string) *FormulaColorant {
	newFormulaColorant := &FormulaColorant{
		Formula:  formula,
		Colorant: colorant,
		Amount:   amount,
	}

	return newFormulaColorant
}

func (fb *FormulaColorant) ToProto() *proto.FormulaColorant {
	return &proto.FormulaColorant{
		Formula:  fb.Formula,
		Colorant: fb.Colorant,
		Amount:   fb.Amount,
	}
}

func (fb *FormulaColorant) ToStorage() *storage.FormulaColorant {
	return &storage.FormulaColorant{
		Formula:  fb.Formula,
		Colorant: fb.Colorant,
		Amount:   fb.Amount,
	}
}

func (fb *FormulaColorant) FromStorage(s *storage.FormulaColorant) {
	fb.Formula = s.Formula
	fb.Colorant = s.Colorant
	fb.Amount = s.Amount
}
