package storage

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCRUDFormulas(t *testing.T) {
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

	formula := Formula{
		Account:  "test_account",
		ID:       "test_formula",
		Name:     "Test Formula",
		Number:   "Formula number",
		Notes:    "Formula notes",
		Created:  0,
		Modified: 0,
	}

	err = db.InsertFormula(db, &formula)
	if err != nil {
		t.Fatal(err)
	}

	formulas, err := db.ListFormulas(db, account.ID, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(formulas) != 1 {
		t.Errorf("expected 1 element in list found %d", len(formulas))
	}

	if diff := cmp.Diff(formula, formulas[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	fetchedFormula, err := db.GetFormula(db, account.ID, formula.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(formula, fetchedFormula); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	formula.Name = "Updated Formula"
	formula.Modified = 1

	err = db.UpdateFormula(db, account.ID, formula.ID, UpdatableFormulaFields{
		Name:     &formula.Name,
		Modified: &formula.Modified,
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedFormula, err = db.GetFormula(db, account.ID, formula.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(formula, fetchedFormula); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	newColorant := Colorant{
		Account:      account.ID,
		ID:           "test_colorant",
		Manufacturer: "test_colorant_manu",
	}

	err = db.InsertColorant(db, &newColorant)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AssociateColorantWithFormula(db, &FormulaColorant{
		Account:  account.ID,
		Formula:  formula.ID,
		Colorant: newColorant.ID,
		Amount:   "test_amount",
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedFormulaColorants, err := db.ListFormulaColorants(db, account.ID, formula.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(FormulaColorant{
		Account:  account.ID,
		Formula:  formula.ID,
		Colorant: newColorant.ID,
		Amount:   "test_amount",
	}, fetchedFormulaColorants[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	newBase := Base{
		Account:      account.ID,
		ID:           "test_base",
		Manufacturer: "test_base_manu",
	}

	err = db.InsertBase(db, &newBase)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AssociateBaseWithFormula(db, &FormulaBase{
		Account: account.ID,
		Formula: formula.ID,
		Base:    newBase.ID,
		Amount:  "test_amount",
	})
	if err != nil {
		t.Fatal(err)
	}

	fetchedFormulaBases, err := db.ListFormulaBases(db, account.ID, formula.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(FormulaBase{
		Account: account.ID,
		Formula: formula.ID,
		Base:    newBase.ID,
		Amount:  "test_amount",
	}, fetchedFormulaBases[0]); diff != "" {
		t.Errorf("unexpected map values (-want +got):\n%s", diff)
	}

	err = db.DeleteFormula(db, account.ID, formula.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GetFormula(db, account.ID, formula.ID)
	if !errors.Is(err, ErrEntityNotFound) {
		t.Fatal("expected error Not Found; found alternate error")
	}
}
