package storage

import (
	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/listutil"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

// GetAllFormulas returns all formulas as a map
func (db *BoltDB) GetAllFormulas(account string) (map[string]*api.Formula, error) {
	results := map[string]*api.Formula{}

	db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(formulasBucket))

		err := formulasBucket.ForEach(func(key, value []byte) error {
			var formula api.Formula
			err := go_proto.Unmarshal(value, &formula)
			if err != nil {
				return err
			}
			results[string(key)] = &formula
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	return results, nil
}

// GetFormula returns a single formula as a formula object
func (db *BoltDB) GetFormula(account, key string) (*api.Formula, error) {
	var storedFormula api.Formula

	err := db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(formulasBucket))

		formulaRaw := formulasBucket.Get([]byte(key))
		if formulaRaw == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(formulaRaw, &storedFormula)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &storedFormula, nil
}

// AddFormula creates a single formula, returns a generated ID
func (db *BoltDB) AddFormula(account string, newFormula *api.Formula) (key string, err error) {
	err = db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(formulasBucket))

		key, err = db.getNewKey(formulasBucket)

		newFormula.Id = key

		// If the user has not entered a formula number just make it the ID
		if newFormula.Number == "" {
			newFormula.Number = newFormula.Id
		}

		formulaRaw, err := go_proto.Marshal(newFormula)
		if err != nil {
			return err
		}

		err = formulasBucket.Put([]byte(key), formulaRaw)
		if err != nil {
			return err
		}

		// for all jobs included in newformula make sure to add new formula ID to all
		for _, jobID := range newFormula.Jobs {
			err := db.linkFormulaToJob(tx, account, key, jobID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return key, err
}

// UpdateFormula modifies a single formula
func (db *BoltDB) UpdateFormula(account, key string, updatedFormula *api.Formula) error {
	var storedFormula api.Formula

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(formulasBucket))

		// First check if key exists
		currentFormula := formulasBucket.Get([]byte(key))
		if currentFormula == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(currentFormula, &storedFormula)
		if err != nil {
			return err
		}

		formulaRaw, err := go_proto.Marshal(updatedFormula)
		if err != nil {
			return err
		}

		err = formulasBucket.Put([]byte(key), formulaRaw)
		if err != nil {
			return err
		}

		// figure out which changes to need to be made to other objects due to job additions/removals
		additions, removals := listutil.FindListUpdates(storedFormula.Jobs, updatedFormula.Jobs)

		// Append formula id to formula list in job
		for _, jobID := range additions {
			err := db.linkFormulaToJob(tx, account, updatedFormula.Id, jobID)
			if err != nil {
				return err
			}
		}

		// Remove formula id from formulas list in jobs removed
		for _, jobID := range removals {
			err := db.unlinkFormulaFromJob(tx, account, updatedFormula.Id, jobID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteFormula removes a formula by key
func (db *BoltDB) DeleteFormula(account, key string) error {
	var storedFormula api.Formula

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(formulasBucket))

		// First check if key exists
		currentFormula := formulasBucket.Get([]byte(key))
		if currentFormula == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(currentFormula, &storedFormula)
		if err != nil {
			return err
		}

		err = formulasBucket.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove formula id from formulas list in jobs this was linked to
		for _, jobID := range storedFormula.Jobs {
			err := db.unlinkFormulaFromJob(tx, account, key, jobID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
