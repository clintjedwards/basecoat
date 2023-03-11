package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	qb "github.com/Masterminds/squirrel"
)

type Account struct {
	ID       string
	Name     string
	Hash     string
	State    string
	Created  int64
	Modified int64
}

type UpdatableAccountFields struct {
	Name     *string
	Hash     *string
	State    *string
	Modified *int64
}

func (db *DB) ListAccounts(conn Queryable, offset, limit int) ([]Account, error) {
	if limit == 0 || limit > db.maxResultsLimit {
		limit = db.maxResultsLimit
	}

	query, args := qb.Select("id", "name", "hash", "state", "created", "modified").
		From("accounts").OrderBy("id").Limit(uint64(limit)).Offset(uint64(offset)).MustSql()

	accounts := []Account{}
	err := conn.Select(&accounts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return accounts, nil
}

func (db *DB) InsertAccount(conn Queryable, account *Account) error {
	_, err := qb.Insert("accounts").Columns("id", "name", "hash", "state", "created", "modified").Values(
		account.ID, account.Name, account.Hash, account.State, account.Created, account.Modified,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) GetAccount(conn Queryable, id string) (Account, error) {
	query, args := qb.Select("id", "name", "hash", "state", "created", "modified").
		From("accounts").Where(qb.Eq{"id": id}).MustSql()

	account := Account{}
	err := conn.Get(&account, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Account{}, ErrEntityNotFound
		}

		return Account{}, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return account, nil
}

func (db *DB) UpdateAccount(conn Queryable, id string, fields UpdatableAccountFields) error {
	query := qb.Update("accounts")

	if fields.Name != nil {
		query = query.Set("name", fields.Name)
	}

	if fields.State != nil {
		query = query.Set("state", fields.State)
	}

	if fields.Hash != nil {
		query = query.Set("hash", fields.Hash)
	}

	if fields.Modified != nil {
		query = query.Set("modified", fields.Modified)
	}

	_, err := query.Where(qb.Eq{"id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEntityNotFound
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) DeleteAccount(conn Queryable, id string) error {
	_, err := qb.Delete("accounts").Where(qb.Eq{"id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}
