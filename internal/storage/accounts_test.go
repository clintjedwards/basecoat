package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCRUDAccounts(t *testing.T) {
	path := tempFile()
	db, err := New(path, 200)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	account := Account{
		ID:       "test_account",
		Name:     "Test Account",
		Hash:     "This is a test account",
		State:    "SOME_STATE",
		Created:  0,
		Modified: 0,
	}

	err = db.InsertAccount(db, &account)
	if err != nil {
		t.Fatal(err)
	}

	accounts, err := db.ListAccounts(db, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(accounts) != 1 {
		t.Errorf("expected 1 element in list found %d", len(accounts))
	}

	if diff := cmp.Diff(account, accounts[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	fetchedAccount, err := db.GetAccount(db, account.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(account, fetchedAccount); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	account.Name = "Updated Account"
	account.State = "updated account"
	account.Hash = "updated password"
	account.Modified = 1

	err = db.UpdateAccount(db, account.ID, UpdatableAccountFields{
		Name:     &account.Name,
		State:    &account.State,
		Hash:     &account.Hash,
		Modified: &account.Modified,
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedAccount, err = db.GetAccount(db, account.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(account, fetchedAccount); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.DeleteAccount(db, account.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetAccount(db, account.ID)
	if !errors.Is(err, ErrEntityNotFound) {
		t.Fatal("expected error Not Found; found alternate error")
	}
}
