package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCRUDColorants(t *testing.T) {
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

	colorant := Colorant{
		Account:      account.ID,
		ID:           "test_colorant",
		Label:        "label",
		Manufacturer: "test_colorant_manu",
		Created:      2,
	}

	err = db.InsertColorant(db, &colorant)
	if err != nil {
		t.Fatal(err)
	}

	colorants, err := db.ListColorants(db, account.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(colorants) != 1 {
		t.Errorf("expected 1 element in list found %d", len(colorants))
	}

	if diff := cmp.Diff(colorant, colorants[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	fetchedColorant, err := db.GetColorant(db, account.ID, colorant.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(colorant, fetchedColorant); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.DeleteColorant(db, account.ID, colorant.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetColorant(db, account.ID, colorant.ID)
	if !errors.Is(err, ErrEntityNotFound) {
		t.Fatal("expected error Not Found; found alternate error")
	}
}
