package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	qb "github.com/Masterminds/squirrel"
)

type Base struct {
	Account      string
	ID           string
	Label        string
	Manufacturer string
	Created      int64
}

type UpdatableBaseFields struct {
	Label        *string
	Manufacturer *string
}

func (db *DB) ListBases(conn Queryable, account string, offset, limit int) ([]Base, error) {
	if limit == 0 || limit > db.maxResultsLimit {
		limit = db.maxResultsLimit
	}

	query, args := qb.Select("account", "id", "label", "manufacturer", "created").
		From("bases").Where(qb.Eq{"account": account}).
		OrderBy("id").Limit(uint64(limit)).Offset(uint64(offset)).MustSql()

	bases := []Base{}
	err := conn.Select(&bases, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return bases, nil
}

func (db *DB) InsertBase(conn Queryable, base *Base) error {
	_, err := qb.Insert("bases").Columns("account", "id", "label", "manufacturer", "created").
		Values(base.Account, base.ID, base.Label, base.Manufacturer, base.Created).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) GetBase(conn Queryable, account, id string) (Base, error) {
	query, args := qb.Select("account", "id", "label", "manufacturer", "created").From("bases").
		Where(qb.Eq{"account": account, "id": id}).MustSql()

	base := Base{}
	err := conn.Get(&base, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Base{}, ErrEntityNotFound
		}

		return Base{}, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return base, nil
}

func (db *DB) UpdateBase(conn Queryable, account, id string, fields UpdatableBaseFields) error {
	query := qb.Update("bases")

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

func (db *DB) DeleteBase(conn Queryable, account, id string) error {
	_, err := qb.Delete("bases").Where(qb.Eq{"account": account, "id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}
