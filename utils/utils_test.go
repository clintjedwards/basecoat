package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindListDifference(t *testing.T) {
	testListA := []string{"1", "2", "3", "4"}
	testListB := []string{"1", "3"}
	returnedList := FindListDifference(testListA, testListB)

	expectedDifference := []string{"2", "4"}

	if !cmp.Equal(returnedList, expectedDifference) {
		t.Errorf("List does not contain correct difference: Diff below: \n%v", cmp.Diff(returnedList, expectedDifference))
	}
}

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
