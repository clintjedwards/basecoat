package storage

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/random"
)

// getNewKey generates a unique key for a given bucket
func (db *BoltDB) getNewKey(bucket *bolt.Bucket) (string, error) {

	const retryLimit int = 3
	var key string

	if bucket == nil {
		return "", fmt.Errorf("could generate key bucket is nil")
	}

	for i := 1; i <= retryLimit; i++ {
		key = string(random.GenerateRandString(db.idLength))

		existingKey := bucket.Get([]byte(key))
		if existingKey == nil {
			return key, nil
		}
		continue
	}

	// exceeded retries
	return "", fmt.Errorf("exceeded maximum retries(%d) for key generation", retryLimit)
}

// func (db *BoltDB) linkFormulaToJob(tx *bolt.Tx, account, formulaID, jobID string) error {

// 	var storedJob api.Job

// 	accountBucket := tx.Bucket([]byte(account))
// 	jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

// 	jobRaw := jobsBucket.Get([]byte(jobID))
// 	if jobRaw == nil {
// 		return tkerrors.ErrEntityNotFound
// 	}

// 	err := go_proto.Unmarshal(jobRaw, &storedJob)
// 	if err != nil {
// 		return err
// 	}

// 	storedJob.Formulas = append(storedJob.Formulas, formulaID)

// 	updatedJobRaw, err := go_proto.Marshal(&storedJob)
// 	if err != nil {
// 		return err
// 	}

// 	err = jobsBucket.Put([]byte(jobID), updatedJobRaw)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (db *BoltDB) unlinkFormulaFromJob(tx *bolt.Tx, account, formulaID, jobID string) error {
// 	var storedJob api.Job

// 	accountBucket := tx.Bucket([]byte(account))
// 	jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

// 	jobRaw := jobsBucket.Get([]byte(jobID))
// 	if jobRaw == nil {
// 		return tkerrors.ErrEntityNotFound
// 	}

// 	err := go_proto.Unmarshal(jobRaw, &storedJob)
// 	if err != nil {
// 		return err
// 	}

// 	storedJob.Formulas = listutil.RemoveStringFromList(storedJob.Formulas, formulaID)

// 	updatedJobRaw, err := go_proto.Marshal(&storedJob)
// 	if err != nil {
// 		return err
// 	}

// 	err = jobsBucket.Put([]byte(jobID), updatedJobRaw)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (db *BoltDB) linkJobToFormula(tx *bolt.Tx, account, jobID, formulaID string) error {

// 	var storedFormula api.Formula

// 	accountBucket := tx.Bucket([]byte(account))
// 	formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

// 	formulaRaw := formulasBucket.Get([]byte(formulaID))
// 	if formulaRaw == nil {
// 		return tkerrors.ErrEntityNotFound
// 	}

// 	err := go_proto.Unmarshal(formulaRaw, &storedFormula)
// 	if err != nil {
// 		return err
// 	}

// 	storedFormula.Jobs = append(storedFormula.Jobs, jobID)

// 	updatedFormulaRaw, err := go_proto.Marshal(&storedFormula)
// 	if err != nil {
// 		return err
// 	}

// 	err = formulasBucket.Put([]byte(formulaID), updatedFormulaRaw)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (db *BoltDB) unlinkJobFromFormula(tx *bolt.Tx, account, jobID, formulaID string) error {
// 	var storedFormula api.Formula

// 	accountBucket := tx.Bucket([]byte(account))
// 	formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

// 	formulaRaw := formulasBucket.Get([]byte(formulaID))
// 	if formulaRaw == nil {
// 		return tkerrors.ErrEntityNotFound
// 	}

// 	err := go_proto.Unmarshal(formulaRaw, &storedFormula)
// 	if err != nil {
// 		return err
// 	}

// 	storedFormula.Jobs = listutil.RemoveStringFromList(storedFormula.Jobs, jobID)

// 	updatedFormulaRaw, err := go_proto.Marshal(&storedFormula)
// 	if err != nil {
// 		return err
// 	}

// 	err = formulasBucket.Put([]byte(formulaID), updatedFormulaRaw)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

// handleNewJob generates a new key for a created job
func (db *BoltDB) handleNewJob(tx *bolt.Tx, account string, newJob *api.Job) error {
	newJob.Id = string(random.GenerateRandString(db.idLength))
	newJob.Created = time.Now().Unix()
	newJob.Modified = time.Now().Unix()

	// // all formulas within job add jodId to formula list
	// for _, formulaID := range newJob.Formulas {
	// 	err := db.linkJobToFormula(tx, account, newJob.Id, formulaID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

// func (db *BoltDB) handleUpdatedJob(tx *bolt.Tx, account string, newJob *api.Job) error {
// 	newJob.Id = string(random.GenerateRandString(db.idLength))
// 	newJob.Created = time.Now().Unix()
// 	newJob.Modified = time.Now().Unix()

// 	// all formulas within job add jodId to formula list
// 	for _, formulaID := range newJob.Formulas {
// 		err := db.linkJobToFormula(tx, account, newJob.Id, formulaID)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
