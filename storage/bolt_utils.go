package storage

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/listutil"
	"github.com/clintjedwards/toolkit/random"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

func (db *boltDB) getNewKey(bucket *bolt.Bucket) (string, error) {

	const retryLimit int = 3
	var key string

	for i := 1; i <= retryLimit; i++ {
		key = string(random.GenerateRandString(db.idLength))

		keyRaw := bucket.Get([]byte(key))
		if keyRaw == nil {
			return key, nil
		}

		continue
	}

	// exceeded retries
	return "", fmt.Errorf("exceeded maximum retries(%d) for key generation", retryLimit)
}

func (db *boltDB) linkFormulaToJob(tx *bolt.Tx, account, formulaID, jobID string) error {

	var storedJob api.Job

	accountBucket := tx.Bucket([]byte(account))
	jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

	jobRaw := jobsBucket.Get([]byte(jobID))
	if jobRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(jobRaw, &storedJob)
	if err != nil {
		return err
	}

	storedJob.Formulas = append(storedJob.Formulas, formulaID)

	updatedJobRaw, err := go_proto.Marshal(&storedJob)
	if err != nil {
		return err
	}

	err = jobsBucket.Put([]byte(jobID), updatedJobRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *boltDB) unlinkFormulaFromJob(tx *bolt.Tx, account, formulaID, jobID string) error {
	var storedJob api.Job

	accountBucket := tx.Bucket([]byte(account))
	jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

	jobRaw := jobsBucket.Get([]byte(jobID))
	if jobRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(jobRaw, &storedJob)
	if err != nil {
		return err
	}

	storedJob.Formulas = listutil.RemoveStringFromList(storedJob.Formulas, formulaID)

	updatedJobRaw, err := go_proto.Marshal(&storedJob)
	if err != nil {
		return err
	}

	err = jobsBucket.Put([]byte(jobID), updatedJobRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *boltDB) linkJobToFormula(tx *bolt.Tx, account, jobID, formulaID string) error {

	var storedFormula api.Formula

	accountBucket := tx.Bucket([]byte(account))
	formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

	formulaRaw := formulasBucket.Get([]byte(formulaID))
	if formulaRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(formulaRaw, &storedFormula)
	if err != nil {
		return err
	}

	storedFormula.Jobs = append(storedFormula.Jobs, jobID)

	updatedFormulaRaw, err := go_proto.Marshal(&storedFormula)
	if err != nil {
		return err
	}

	err = formulasBucket.Put([]byte(formulaID), updatedFormulaRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *boltDB) unlinkJobFromFormula(tx *bolt.Tx, account, jobID, formulaID string) error {
	var storedFormula api.Formula

	accountBucket := tx.Bucket([]byte(account))
	formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

	formulaRaw := formulasBucket.Get([]byte(formulaID))
	if formulaRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(formulaRaw, &storedFormula)
	if err != nil {
		return err
	}

	storedFormula.Jobs = listutil.RemoveStringFromList(storedFormula.Jobs, jobID)

	updatedFormulaRaw, err := go_proto.Marshal(&storedFormula)
	if err != nil {
		return err
	}

	err = formulasBucket.Put([]byte(formulaID), updatedFormulaRaw)
	if err != nil {
		return err
	}

	return nil

}
