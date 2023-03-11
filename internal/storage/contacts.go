package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	qb "github.com/Masterminds/squirrel"
)

type Contact struct {
	Account  string
	ID       string
	Name     string
	Email    string
	Phone    string
	Created  int64
	Modified int64
}

type UpdatableContactFields struct {
	Name     *string
	Email    *string
	Phone    *string
	Modified *int64
}

func (db *DB) ListContacts(conn Queryable, account string, offset, limit int) ([]Contact, error) {
	if limit == 0 || limit > db.maxResultsLimit {
		limit = db.maxResultsLimit
	}

	query, args := qb.Select("account", "id", "name", "email", "phone", "created", "modified").
		From("contacts").
		Where(qb.Eq{"account": account}).
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		MustSql()

	contacts := []Contact{}
	err := conn.Select(&contacts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return contacts, nil
}

func (db *DB) InsertContact(conn Queryable, contact *Contact) error {
	_, err := qb.Insert("contacts").Columns("account", "id", "name", "email", "phone", "created", "modified").Values(
		contact.Account, contact.ID, contact.Name, contact.Email, contact.Phone, contact.Created, contact.Modified,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) GetContact(conn Queryable, account, id string) (Contact, error) {
	query, args := qb.Select("account", "id", "name", "email", "phone", "created", "modified").From("contacts").
		Where(qb.Eq{"account": account, "id": id}).MustSql()

	contact := Contact{}
	err := conn.Get(&contact, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Contact{}, ErrEntityNotFound
		}

		return Contact{}, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return contact, nil
}

func (db *DB) UpdateContact(conn Queryable, account, id string, fields UpdatableContactFields) error {
	query := qb.Update("contacts")

	if fields.Name != nil {
		query = query.Set("name", fields.Name)
	}

	if fields.Email != nil {
		query = query.Set("email", fields.Email)
	}

	if fields.Phone != nil {
		query = query.Set("phone", fields.Phone)
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

func (db *DB) DeleteContact(conn Queryable, account, id string) error {
	_, err := qb.Delete("contacts").Where(qb.Eq{"account": account, "id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}
