package utils

import (
	"github.com/clintjedwards/basecoat/api"
	"golang.org/x/crypto/bcrypt"
)

// RemoveStringFromList removes an element form an array of ints
// does not preserve list order
func RemoveStringFromList(list []string, value string) []string {
	for index, item := range list {
		if item == value {
			list[index] = list[len(list)-1]
			return list[:len(list)-1]
		}
	}

	return list
}

// FindListDifference returns list elements that are in list A
// but not found in B
func FindListDifference(a, b []string) []string {
	m := make(map[string]bool)
	diff := []string{}

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// FindListUpdates is used to compare a new and old version of lists
// it will compare the old version to the new version and return
// which elements have been added or removed from the old list
func FindListUpdates(oldList, newList []string) (additions, removals []string) {

	removals = FindListDifference(oldList, newList)
	additions = FindListDifference(newList, oldList)

	return additions, removals
}

// MergeFormulaStruct aids in updating formulas by overwriting old formula
// objects with only non-nil fields from a provided new formula object
func MergeFormulaStruct(oldFormula, newFormula *api.Formula) *api.Formula {

	if newFormula.Name != "" {
		oldFormula.Name = newFormula.Name
	}

	if newFormula.Number != "" {
		oldFormula.Number = newFormula.Number
	}

	if newFormula.Notes != "" {
		oldFormula.Notes = newFormula.Notes
	}

	if newFormula.Modified != 0 {
		oldFormula.Modified = newFormula.Modified
	}

	if newFormula.Jobs != nil {
		oldFormula.Jobs = newFormula.Jobs
	}

	if newFormula.Bases != nil {
		oldFormula.Bases = newFormula.Bases
	}

	if newFormula.Colorants != nil {
		oldFormula.Colorants = newFormula.Colorants
	}

	return oldFormula
}

// MergeJobStruct aids in updating formulas by overwriting old formula
// objects with only non-nil fields from a provided new formula object
func MergeJobStruct(oldJob, newJob *api.Job) *api.Job {

	if newJob.Name != "" {
		oldJob.Name = newJob.Name
	}

	if newJob.Street != "" {
		oldJob.Street = newJob.Street
	}

	if newJob.Street2 != "" {
		oldJob.Street2 = newJob.Street2
	}

	if newJob.City != "" {
		oldJob.City = newJob.City
	}

	if newJob.State != "" {
		oldJob.State = newJob.State
	}

	if newJob.Zipcode != "" {
		oldJob.Zipcode = newJob.Zipcode
	}

	if newJob.Notes != "" {
		oldJob.Notes = newJob.Notes
	}

	if newJob.Formulas != nil {
		oldJob.Formulas = newJob.Formulas
	}

	if newJob.Modified != 0 {
		oldJob.Modified = newJob.Modified
	}

	if newJob.Contact != nil {
		oldJob.Contact = newJob.Contact
	}

	return oldJob
}

// CheckPasswordHash validates a password against the stored hash
// to verify the user is authorized
func CheckPasswordHash(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

// HashPassword converts a byte string password into a bcrypt hash
// which is then stored as the only form of password
func HashPassword(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 14)
	return hash, err
}
