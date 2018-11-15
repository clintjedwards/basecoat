package models

import (
	"strings"

	"github.com/go-pg/pg/orm"
)

//Formula describes a single formula
type Formula struct {
	TableName struct{}          `sql:"formulas" json:"-"`
	ID        int               `json:"id"`
	Name      string            `json:"name" sql:",unique"`
	Number    string            `json:"number"`
	Notes     string            `json:"notes"`
	Jobs      []int             `json:"jobs"` // Job ids
	Colorants map[string]string `json:"colorants"`
	Base      map[string]string `json:"base"`
	Created   int64             `json:"created"`
	Modified  int64             `json:"modified"`
}

//BeforeInsert defines data actions to occur before SQL inserts
func (formula *Formula) BeforeInsert(db orm.DB) error {

	formula.Name = strings.ToLower(formula.Name)
	formula.Number = strings.ToLower(formula.Number)

	return nil
}

//BeforeUpdate defines data actions to occur before SQL updates
func (formula *Formula) BeforeUpdate(db orm.DB) error {

	formula.Name = strings.ToLower(formula.Name)
	formula.Number = strings.ToLower(formula.Number)

	return nil
}
