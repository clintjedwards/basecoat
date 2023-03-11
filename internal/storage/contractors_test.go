package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCRUDContractors(t *testing.T) {
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

	contractor := Contractor{
		Account:  "test_account",
		ID:       "test_contractor",
		Company:  "Test Contractor",
		Contact:  ptr("test_contact"),
		Created:  0,
		Modified: 0,
	}

	err = db.InsertContractor(db, &contractor)
	if err != nil {
		t.Fatal(err)
	}

	contractors, err := db.ListContractors(db, account.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(contractors) != 1 {
		t.Errorf("expected 1 element in list found %d", len(contractors))
	}

	if diff := cmp.Diff(contractor, contractors[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	fetchedContractor, err := db.GetContractor(db, account.ID, contractor.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(contractor, fetchedContractor); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	contractor.Company = "Updated Contractor"
	contractor.Modified = 1

	err = db.UpdateContractor(db, account.ID, contractor.ID, UpdatableContractorFields{
		Company:  &contractor.Company,
		Modified: &contractor.Modified,
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedContractor, err = db.GetContractor(db, account.ID, contractor.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(contractor, fetchedContractor); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.DeleteContractor(db, account.ID, contractor.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetContractor(db, account.ID, contractor.ID)
	if !errors.Is(err, ErrEntityNotFound) {
		t.Fatal("expected error Not Found; found alternate error")
	}
}
