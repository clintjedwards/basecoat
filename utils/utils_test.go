package utils

import (
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/google/go-cmp/cmp"
)

func TestRemoveStringFromList(t *testing.T) {

	testList := []string{"1", "2", "3", "4"}
	removedString := "2"
	returnedList := RemoveStringFromList(testList, removedString)

	for _, item := range returnedList {
		if item == "2" {
			t.Errorf("Item to be removed: %s was still found in slice: %v", removedString, returnedList)
		}
	}
}

func TestFindListUpdates(t *testing.T) {

	currentList := []string{"1", "2", "3", "4"}
	updateList := []string{"1", "3", "5", "6"}

	listAdditions, listRemovals := FindListUpdates(currentList, updateList)

	expectedAdditions := []string{"5", "6"}
	expectedRemovals := []string{"2", "4"}

	if !cmp.Equal(expectedAdditions, listAdditions) {
		t.Errorf("List does not contain correct additions; Diff below: \n%v", cmp.Diff(expectedAdditions, listAdditions))
	}

	if !cmp.Equal(expectedRemovals, listRemovals) {
		t.Errorf("List does not contain correct removals; Diff below: \n%v", cmp.Diff(expectedRemovals, listRemovals))
	}
}

func TestMergeFormulaStruct(t *testing.T) {

	oldFormula := api.Formula{
		Name:      "testFormula",
		Number:    "1",
		Jobs:      []string{"1", "2", "3"},
		Bases:     []*api.Base{&api.Base{Name: "testBase"}},
		Colorants: []*api.Colorant{&api.Colorant{Name: "testColorant"}},
	}

	newFormula := api.Formula{
		Name:      "changedFormula",
		Jobs:      []string{"1", "2", "3", "4"},
		Colorants: []*api.Colorant{},
	}

	updatedFormula := MergeFormulaStruct(&oldFormula, &newFormula)

	intendedFormula := api.Formula{
		Name:      "changedFormula",
		Number:    "1",
		Jobs:      []string{"1", "2", "3", "4"},
		Bases:     []*api.Base{&api.Base{Name: "testBase"}},
		Colorants: []*api.Colorant{},
	}

	if !cmp.Equal(*updatedFormula, intendedFormula) {
		t.Errorf("updated formula does not match intended formula; Diff below: \n%v", cmp.Diff(*updatedFormula, intendedFormula))
	}
}

func TestMergeJobStruct(t *testing.T) {

	oldJob := api.Job{
		Name:   "testJob",
		Street: "1280 E46th",
		Contact: &api.Contact{
			Name: "Bob",
			Info: "Bob@gmail.com",
		},
	}

	newJob := api.Job{
		Street2: "Apt 5X",
	}

	updatedJob := MergeJobStruct(&oldJob, &newJob)

	intendedJob := api.Job{
		Name:    "testJob",
		Street:  "1280 E46th",
		Street2: "Apt 5X",
		Contact: &api.Contact{
			Name: "Bob",
			Info: "Bob@gmail.com",
		},
	}

	if !cmp.Equal(*updatedJob, intendedJob) {
		t.Errorf("updated Job does not match intended Job; Diff included below. \n%v", cmp.Diff(*updatedJob, intendedJob))
	}
}
