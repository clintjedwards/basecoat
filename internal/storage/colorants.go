package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	qb "github.com/Masterminds/squirrel"
)

type Colorant struct {
	Account      string
	ID           string
	Label        string
	Manufacturer string
	Created      int64
}

type UpdatableColorantFields struct {
	Label        *string
	Manufacturer *string
}

func (db *DB) ListColorants(conn Queryable, account string, offset, limit int) ([]Colorant, error) {
	if limit == 0 || limit > db.maxResultsLimit {
		limit = db.maxResultsLimit
	}

	query, args := qb.Select("account", "id", "label", "manufacturer", "created").
		From("colorants").Where(qb.Eq{"account": account}).
		OrderBy("id").Limit(uint64(limit)).Offset(uint64(offset)).MustSql()

	colorants := []Colorant{}
	err := conn.Select(&colorants, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return colorants, nil
}

func (db *DB) InsertColorant(conn Queryable, colorant *Colorant) error {
	_, err := qb.Insert("colorants").Columns("account", "id", "label", "manufacturer", "created").
		Values(colorant.Account, colorant.ID, colorant.Label, colorant.Manufacturer, colorant.Created).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) GetColorant(conn Queryable, account, id string) (Colorant, error) {
	query, args := qb.Select("account", "id", "label", "manufacturer", "created").From("colorants").
		Where(qb.Eq{"account": account, "id": id}).MustSql()

	colorant := Colorant{}
	err := conn.Get(&colorant, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Colorant{}, ErrEntityNotFound
		}

		return Colorant{}, fmt.Errorf("colorant error occurred: %v; %w", err, ErrInternal)
	}

	return colorant, nil
}

func (db *DB) UpdateColorant(conn Queryable, account, id string, fields UpdatableColorantFields) error {
	query := qb.Update("colorants")

	if fields.Label != nil {
		query = query.Set("label", fields.Label)
	}

	if fields.Manufacturer != nil {
		query = query.Set("manufacturer", fields.Manufacturer)
	}

	_, err := query.Where(qb.Eq{"account": account, "id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) DeleteColorant(conn Queryable, account, id string) error {
	_, err := qb.Delete("colorants").Where(qb.Eq{"account": account, "id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) AssociateColorantWithFormula(conn Queryable, formulaColorant *FormulaColorant) error {
	_, err := qb.Insert("formula_colorants").Columns("account", "formula", "colorant", "amount").Values(
		formulaColorant.Account, formulaColorant.Formula, formulaColorant.Colorant, formulaColorant.Amount,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) ListFormulaColorants(conn Queryable, account, formula string) ([]FormulaColorant, error) {
	query, args := qb.Select("account", "formula", "colorant", "amount").From("formula_colorants").
		Where(qb.Eq{"account": account, "formula": formula}).MustSql()

	formulaColorants := []FormulaColorant{}
	err := conn.Select(&formulaColorants, query, args...)
	if err != nil {
		return nil, fmt.Errorf("data error occurred: %v; %w", err, ErrInternal)
	}

	return formulaColorants, nil
}

func (db *DB) ListColorantFormulas(conn Queryable, account, colorant string) ([]FormulaColorant, error) {
	query, args := qb.Select("account", "formula", "colorant", "amount").From("formula_colorants").
		Where(qb.Eq{"account": account, "colorant": colorant}).MustSql()

	formulaColorants := []FormulaColorant{}
	err := conn.Select(&formulaColorants, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return formulaColorants, nil
}

func (db *DB) DeleteFormulaColorant(conn Queryable, account, formula, colorant string) error {
	_, err := qb.Delete("formula_colorants").Where(qb.Eq{"account": account, "formula": formula, "colorant": colorant}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}
