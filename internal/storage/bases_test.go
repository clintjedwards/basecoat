package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCRUDBases(t *testing.T) {
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

	base := Base{
		Account:      account.ID,
		ID:           "test_base",
		Label:        "label",
		Manufacturer: "test_base_manu",
		Created:      2,
	}

	err = db.InsertBase(db, &base)
	if err != nil {
		t.Fatal(err)
	}

	bases, err := db.ListBases(db, account.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(bases) != 1 {
		t.Errorf("expected 1 element in list found %d", len(bases))
	}

	if diff := cmp.Diff(base, bases[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	fetchedBase, err := db.GetBase(db, account.ID, base.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(base, fetchedBase); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.DeleteBase(db, account.ID, base.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetBase(db, account.ID, base.ID)
	if !errors.Is(err, ErrEntityNotFound) {
		t.Fatal("expected error Not Found; found alternate error")
	}
}
