package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	qb "github.com/Masterminds/squirrel"
)

type Job struct {
	Account    string
	ID         string
	Contractor string
	Name       string
	Address    string
	Notes      string
	Contact    *string
	Created    int64
	Modified   int64
}

type FormulaJob struct {
	Account string
	Job     string
	Formula string
}

type UpdatableJobFields struct {
	Name     *string
	Address  *string
	Notes    *string
	Contact  *string
	Modified *int64
}

func (db *DB) ListJobs(conn Queryable, account string, offset, limit int) ([]Job, error) {
	if limit == 0 || limit > db.maxResultsLimit {
		limit = db.maxResultsLimit
	}

	query, args := qb.Select("account", "id", "contractor", "name", "address", "notes", "contact", "created", "modified").
		From("jobs").
		Where(qb.Eq{"account": account}).
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		MustSql()

	jobs := []Job{}
	err := conn.Select(&jobs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return jobs, nil
}

func (db *DB) InsertJob(conn Queryable, job *Job) error {
	_, err := qb.Insert("jobs").Columns("account", "id", "contractor", "name", "address", "notes", "contact", "created", "modified").Values(
		job.Account, job.ID, job.Contractor, job.Name, job.Address, job.Notes, job.Contact, job.Created, job.Modified,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) GetJob(conn Queryable, account, id string) (Job, error) {
	query, args := qb.Select("account", "id", "contractor", "name", "address", "notes", "contact", "created", "modified").From("jobs").
		Where(qb.Eq{"account": account, "id": id}).MustSql()

	job := Job{}
	err := conn.Get(&job, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, ErrEntityNotFound
		}

		return Job{}, fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return job, nil
}

func (db *DB) UpdateJob(conn Queryable, account, id string, fields UpdatableJobFields) error {
	query := qb.Update("jobs")

	if fields.Name != nil {
		query = query.Set("name", fields.Name)
	}

	if fields.Address != nil {
		query = query.Set("address", fields.Address)
	}

	if fields.Notes != nil {
		query = query.Set("notes", fields.Notes)
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

func (db *DB) DeleteJob(conn Queryable, account, id string) error {
	_, err := qb.Delete("jobs").Where(qb.Eq{"account": account, "id": id}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) AssociateFormulaWithJob(conn Queryable, formulaJob *FormulaJob) error {
	_, err := qb.Insert("formula_jobs").Columns("account", "job", "formula").Values(
		formulaJob.Account, formulaJob.Job, formulaJob.Formula,
	).RunWith(conn).Exec()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrEntityExists
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}

func (db *DB) ListJobFormulas(conn Queryable, account, job string) ([]FormulaJob, error) {
	query, args := qb.Select("account", "job", "formula").From("formula_jobs").
		Where(qb.Eq{"account": account, "job": job}).MustSql()

	jobFormulas := []FormulaJob{}
	err := conn.Select(&jobFormulas, query, args...)
	if err != nil {
		return nil, fmt.Errorf("data error occurred: %v; %w", err, ErrInternal)
	}

	return jobFormulas, nil
}

func (db *DB) DeleteJobFormula(conn Queryable, account, job, formula string) error {
	_, err := qb.Delete("formula_jobs").Where(qb.Eq{"account": account, "job": job, "formula": formula}).RunWith(conn).Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("database error occurred: %v; %w", err, ErrInternal)
	}

	return nil
}
