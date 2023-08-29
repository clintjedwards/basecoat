package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	qb "github.com/Masterminds/squirrel"
)

type Formula struct {
	Account  string
	ID       string
	Name     string
	Number   string
	Notes    string
	Created  int64
	Modified int64
}

type FormulaColorant struct {
	Account  string
	Formula  string
	Colorant string
	Amount   string
}

type FormulaBase struct {
	Account string
	Formula string
	Base    string
	Amount  string
}

type UpdatableFormulaFields struct {
	Name     *string
	Number   *string
	Notes    *string
	Modified *int64
}

func (db *DB) ListFormulas(conn Queryable, account string, offset, limit int) ([]Formula, error) {
	if limit == 0 || limit > db.maxResultsLimit {
		limit = db.maxResultsLimit
	}

	query, args := qb.Select("account", "id", "name", "number", "notes", "created", "modified").
		From("formulas").
		Where(qb.Eq{"account": account}).
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		MustSql()

	formulas := []Formula{}
	err := conn.Select(&formulas, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return formulas, nil
}

func (db *DB) InsertFormula(conn Queryable, formula *Formula) error {
	_, err := qb.Insert("formulas").Columns("account", "id", "name", "number", "notes", "created", "modified").Values(
		formula.Account, formula.ID, formula.Name, formula.Number, formula.Notes, formula.Created, formula.Modified,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) AssociateBaseWithFormula(conn Queryable, formulaBase *FormulaBase) error {
	_, err := qb.Insert("formula_bases").Columns("account", "formula", "base", "amount").Values(
		formulaBase.Account, formulaBase.Formula, formulaBase.Base, formulaBase.Amount,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) ListFormulaBases(conn Queryable, account, formula string) ([]FormulaBase, error) {
	query, args := qb.Select("account", "formula", "base", "amount").From("formula_bases").
		Where(qb.Eq{"account": account, "formula": formula}).MustSql()

	formulaBases := []FormulaBase{}
	err := conn.Select(&formulaBases, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return formulaBases, nil
}

func (db *DB) DeleteFormulaBase(conn Queryable, account, formula, base string) error {
	_, err := qb.Delete("formula_bases").Where(qb.Eq{"account": account, "formula": formula, "base": base}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) ListFormulaJobs(conn Queryable, account, formula string) ([]FormulaJob, error) {
	query, args := qb.Select("account", "job", "formula").From("formula_jobs").
		Where(qb.Eq{"account": account, "formula": formula}).MustSql()

	jobFormulas := []FormulaJob{}
	err := conn.Select(&jobFormulas, query, args...)
	if err != nil {
		return nil, fmt.Errorf("data error occurred: %v; %w", err, ErrInternal)
	}

	return jobFormulas, nil
}

func (db *DB) GetFormula(conn Queryable, account, id string) (Formula, error) {
	query, args := qb.Select("account", "id", "name", "number", "notes", "created", "modified").From("formulas").
		Where(qb.Eq{"account": account, "id": id}).MustSql()

	formula := Formula{}
	err := conn.Get(&formula, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Formula{}, ErrEntityNotFound
		}

		return Formula{}, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return formula, nil
}

func (db *DB) UpdateFormula(conn Queryable, account, id string, fields UpdatableFormulaFields) error {
	query := qb.Update("formulas")

	if fields.Name != nil {
		query = query.Set("name", fields.Name)
	}

	if fields.Number != nil {
		query = query.Set("number", fields.Number)
	}

	if fields.Notes != nil {
		query = query.Set("notes", fields.Notes)
	}

	if fields.Modified != nil {
		query = query.Set("modified", fields.Modified)
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

func (db *DB) DeleteFormula(conn Queryable, account, id string) error {
	_, err := qb.Delete("formulas").Where(qb.Eq{"account": account, "id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}
