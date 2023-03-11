package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func ptr[T any](v T) *T {
	return &v
}

func TestCRUDJobs(t *testing.T) {
	path := tempFile()
	db, err := New(path, 200)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	account := Account{
		ID: "test_account",
	}

	err = db.InsertAccount(db, &account)
	if err != nil {
		t.Fatal(err)
	}

	err = db.InsertContact(db, &Contact{
		Account: "test_account",
		ID:      "test_contact",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = db.InsertContractor(db, &Contractor{
		Account: "test_account",
		ID:      "test_contractor",
		Contact: ptr("test_contact"),
	})
	if err != nil {
		t.Fatal(err)
	}

	job := Job{
		Account:    "test_account",
		Contractor: "test_contractor",
		ID:         "test_job",
		Contact:    ptr("test_contact"),
		Created:    0,
		Modified:   0,
	}

	err = db.InsertJob(db, &job)
	if err != nil {
		t.Fatal(err)
	}

	jobs, err := db.ListJobs(db, account.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(jobs) != 1 {
		t.Errorf("expected 1 element in list found %d", len(jobs))
	}

	if diff := cmp.Diff(job, jobs[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	fetchedJob, err := db.GetJob(db, account.ID, job.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(job, fetchedJob); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	job.Name = "Updated Job"
	job.Modified = 1

	err = db.UpdateJob(db, account.ID, job.ID, UpdatableJobFields{
		Name:     &job.Name,
		Modified: &job.Modified,
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedJob, err = db.GetJob(db, account.ID, job.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(job, fetchedJob); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.InsertFormula(db, &Formula{
		Account: "test_account",
		ID:      "test_formula",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AssociateFormulaWithJob(db, &FormulaJob{
		Account: account.ID,
		Job:     job.ID,
		Formula: "test_formula",
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedJobFormulas, err := db.ListJobFormulas(db, account.ID, job.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(FormulaJob{
		Account: account.ID,
		Job:     job.ID,
		Formula: "test_formula",
	}, fetchedJobFormulas[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.DeleteJob(db, account.ID, job.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetJob(db, account.ID, job.ID)
	if !errors.Is(err, ErrEntityNotFound) {
		t.Fatal("expected error Not Found; found alternate error")
	}
}
