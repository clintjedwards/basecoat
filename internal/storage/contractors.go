package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	qb "github.com/Masterminds/squirrel"
)

type Contractor struct {
	Account  string
	ID       string
	Company  string
	Contact  *string
	Created  int64
	Modified int64
}

type UpdatableContractorFields struct {
	Company  *string
	Contact  *string
	Modified *int64
}

func (db *DB) ListContractors(conn Queryable, account string, offset, limit int) ([]Contractor, error) {
	if limit == 0 || limit > db.maxResultsLimit {
		limit = db.maxResultsLimit
	}

	query, args := qb.Select("account", "id", "company", "contact", "created", "modified").
		From("contractors").
		Where(qb.Eq{"account": account}).
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		MustSql()

	contractors := []Contractor{}
	err := conn.Select(&contractors, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return contractors, nil
}

func (db *DB) InsertContractor(conn Queryable, contractor *Contractor) error {
	_, err := qb.Insert("contractors").Columns("account", "id", "company", "contact", "created", "modified").Values(
		contractor.Account, contractor.ID, contractor.Company, contractor.Contact, contractor.Created, contractor.Modified,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) GetContractor(conn Queryable, account, id string) (Contractor, error) {
	query, args := qb.Select("account", "id", "company", "contact", "created", "modified").From("contractors").
		Where(qb.Eq{"account": account, "id": id}).MustSql()

	contractor := Contractor{}
	err := conn.Get(&contractor, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Contractor{}, ErrEntityNotFound
		}

		return Contractor{}, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return contractor, nil
}

func (db *DB) UpdateContractor(conn Queryable, account, id string, fields UpdatableContractorFields) error {
	query := qb.Update("contractors")

	if fields.Company != nil {
		query = query.Set("company", fields.Company)
	}

	if fields.Contact != nil {
		query = query.Set("contact", fields.Contact)
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

func (db *DB) DeleteContractor(conn Queryable, account, id string) error {
	_, err := qb.Delete("contractors").Where(qb.Eq{"account": account, "id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}
