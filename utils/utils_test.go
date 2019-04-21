package utils

import (
	"testing"
)

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func ListsAreEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestRemoveIntFromList(t *testing.T) {

	testList := []int{1, 2, 3, 4}
	removedInteger := 2
	returnedList := RemoveIntFromList(testList, removedInteger)

	for _, item := range returnedList {
		if item == 2 {
			t.Errorf("Item to be removed: %d was still found in slice: %v", removedInteger, returnedList)
		}
	}
}

func TestFindListUpdates(t *testing.T) {

	currentList := []int{1, 2, 3, 4}
	updateList := []int{1, 3, 5, 6}

	listAdditions, listRemovals := FindListUpdates(currentList, updateList)

	expectedAdditions := []int{5, 6}
	expectedRemovals := []int{2, 4}

	if !ListsAreEqual(expectedAdditions, listAdditions) {
		t.Errorf("List does not contain correct additions; expected %v; got %v", expectedAdditions, listAdditions)
	}

	if !ListsAreEqual(expectedRemovals, listRemovals) {
		t.Errorf("List does not contain correct removals; expected %v; got %v", expectedRemovals, listRemovals)
	}
}
