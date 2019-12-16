package storage

import (
	"cloud.google.com/go/datastore"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/listutil"
	"github.com/clintjedwards/toolkit/logger"
)

// linkFormulaToJob adds a formula's id to the formula list for a particular job
func (*googleDatastore) linkFormulaToJob(tx *datastore.Transaction, account, formulaID string, jobID string) error {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	getKey := datastore.NameKey(string(JobsBucket), jobID, parentKey)

	var updatedJob api.Job
	err := tx.Get(getKey, &updatedJob)
	if err != nil {
		return err
	}

	updatedJob.Formulas = append(updatedJob.Formulas, formulaID)

	_, err = tx.Put(getKey, &updatedJob)
	if err != nil {
		logger.Log().Errorf("could not link formula to job: %v", err)
		return err
	}

	return nil
}

// unlinkFormulaFromJob removes a formula's id from the formula list of a specific job
func (*googleDatastore) unlinkFormulaFromJob(tx *datastore.Transaction, account, formulaID string, jobID string) error {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	getKey := datastore.NameKey(string(JobsBucket), jobID, parentKey)

	var updatedJob api.Job
	err := tx.Get(getKey, &updatedJob)
	if err != nil {
		return err
	}

	updatedJob.Formulas = listutil.RemoveStringFromList(updatedJob.Formulas, formulaID)

	_, err = tx.Put(getKey, &updatedJob)
	if err != nil {
		logger.Log().Errorf("could not unlink formula from job: %v", err)
		return err
	}

	return nil
}

// linkJobToFormula adds a job's id to the job list for a particular formula
func (*googleDatastore) linkJobToFormula(tx *datastore.Transaction, account, jobID string, formulaID string) error {
	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	getKey := datastore.NameKey(string(FormulasBucket), formulaID, parentKey)

	var updatedFormula api.Formula
	err := tx.Get(getKey, &updatedFormula)
	if err != nil {
		return err
	}

	updatedFormula.Jobs = append(updatedFormula.Jobs, jobID)

	_, err = tx.Put(getKey, &updatedFormula)
	if err != nil {
		logger.Log().Errorf("could not link job to formula: %v", err)
		return err
	}

	return nil
}

// unlinkJobFromFormula removes a job's id from the job list of a specific formula
func (*googleDatastore) unlinkJobFromFormula(tx *datastore.Transaction, account, jobID string, formulaID string) error {
	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	getKey := datastore.NameKey(string(FormulasBucket), formulaID, parentKey)

	var updatedFormula api.Formula
	err := tx.Get(getKey, &updatedFormula)
	if err != nil {
		return err
	}

	updatedFormula.Jobs = listutil.RemoveStringFromList(updatedFormula.Jobs, jobID)

	_, err = tx.Put(getKey, &updatedFormula)
	if err != nil {
		logger.Log().Errorf("could not unlink job from formula: %v", err)
		return err
	}

	return nil
}
