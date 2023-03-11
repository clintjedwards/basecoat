package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCRUDContacts(t *testing.T) {
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

	contact := Contact{
		Account:  "test_account",
		ID:       "test_contact",
		Name:     "Test Contact",
		Email:    "lolwut@gmail.com",
		Phone:    "0",
		Created:  0,
		Modified: 0,
	}

	err = db.InsertContact(db, &contact)
	if err != nil {
		t.Fatal(err)
	}

	contacts, err := db.ListContacts(db, account.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(contacts) != 1 {
		t.Errorf("expected 1 element in list found %d", len(contacts))
	}

	if diff := cmp.Diff(contact, contacts[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	fetchedContact, err := db.GetContact(db, account.ID, contact.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(contact, fetchedContact); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	contact.Name = "Updated Contact"
	contact.Modified = 1

	err = db.UpdateContact(db, account.ID, contact.ID, UpdatableContactFields{
		Name:     &contact.Name,
		Modified: &contact.Modified,
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedContact, err = db.GetContact(db, account.ID, contact.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(contact, fetchedContact); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.DeleteContact(db, account.ID, contact.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetContact(db, account.ID, contact.ID)
	if !errors.Is(err, ErrEntityNotFound) {
		t.Fatal("expected error Not Found; found alternate error")
	}
}
