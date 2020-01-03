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

func (db *BoltDB) getNewKey(bucket *bolt.Bucket) (string, error) {

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

func (db *BoltDB) linkFormulaToJob(accountBucket *bolt.Bucket, formulaID, jobID string) error {

	var storedJob api.Job

	targetBucket := accountBucket.Bucket([]byte(jobsBucket))

	jobRaw := targetBucket.Get([]byte(jobID))
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

	err = targetBucket.Put([]byte(jobID), updatedJobRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *BoltDB) unlinkFormulaFromJob(accountBucket *bolt.Bucket, formulaID, jobID string) error {
	var storedJob api.Job

	targetBucket := accountBucket.Bucket([]byte(jobsBucket))

	jobRaw := targetBucket.Get([]byte(jobID))
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

	err = targetBucket.Put([]byte(jobID), updatedJobRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *BoltDB) linkJobToFormula(accountBucket *bolt.Bucket, jobID, formulaID string) error {
	var storedFormula api.Formula

	formulasBucket := accountBucket.Bucket([]byte(formulasBucket))

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

func (db *BoltDB) unlinkJobFromFormula(accountBucket *bolt.Bucket, jobID, formulaID string) error {
	var storedFormula api.Formula

	formulasBucket := accountBucket.Bucket([]byte(formulasBucket))

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

func (db *BoltDB) linkContractorToJob(accountBucket *bolt.Bucket, contractorID, jobID string) error {

	var storedJob api.Job

	targetBucket := accountBucket.Bucket([]byte(jobsBucket))

	jobRaw := targetBucket.Get([]byte(jobID))
	if jobRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(jobRaw, &storedJob)
	if err != nil {
		return err
	}

	storedJob.ContractorId = contractorID

	updatedJobRaw, err := go_proto.Marshal(&storedJob)
	if err != nil {
		return err
	}

	err = targetBucket.Put([]byte(jobID), updatedJobRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *BoltDB) unlinkContractorFromJob(accountBucket *bolt.Bucket, jobID string) error {

	var storedJob api.Job

	targetBucket := accountBucket.Bucket([]byte(jobsBucket))

	jobRaw := targetBucket.Get([]byte(jobID))
	if jobRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(jobRaw, &storedJob)
	if err != nil {
		return err
	}

	storedJob.ContractorId = ""

	updatedJobRaw, err := go_proto.Marshal(&storedJob)
	if err != nil {
		return err
	}

	err = targetBucket.Put([]byte(jobID), updatedJobRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *BoltDB) linkJobToContractor(accountBucket *bolt.Bucket, jobID, contractorID string) error {

	var storedContractor api.Contractor

	targetBucket := accountBucket.Bucket([]byte(contractorsBucket))

	contractorRaw := targetBucket.Get([]byte(contractorID))
	if contractorRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(contractorRaw, &storedContractor)
	if err != nil {
		return err
	}

	storedContractor.Jobs = append(storedContractor.Jobs, jobID)

	updatedContractorRaw, err := go_proto.Marshal(&storedContractor)
	if err != nil {
		return err
	}

	err = targetBucket.Put([]byte(contractorID), updatedContractorRaw)
	if err != nil {
		return err
	}

	return nil
}

func (db *BoltDB) unlinkJobFromContractor(accountBucket *bolt.Bucket, jobID, contractorID string) error {

	var storedContractor api.Contractor

	targetBucket := accountBucket.Bucket([]byte(contractorsBucket))

	contractorRaw := targetBucket.Get([]byte(contractorID))
	if contractorRaw == nil {
		return tkerrors.ErrEntityNotFound
	}

	err := go_proto.Unmarshal(contractorRaw, &storedContractor)
	if err != nil {
		return err
	}

	storedContractor.Jobs = listutil.RemoveStringFromList(storedContractor.Jobs, jobID)

	updatedContractorRaw, err := go_proto.Marshal(&storedContractor)
	if err != nil {
		return err
	}

	err = targetBucket.Put([]byte(contractorID), updatedContractorRaw)
	if err != nil {
		return err
	}

	return nil
}
