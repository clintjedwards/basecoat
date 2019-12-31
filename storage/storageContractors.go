package storage

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

// GetAllContractors returns all contractors using a map
func (db *BoltDB) GetAllContractors(account string) (map[string]*api.Contractor, error) {
	results := map[string]*api.Contractor{}

	db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return fmt.Errorf("could not find account bucket for %s: %w", account, tkerrors.ErrEntityNotFound)
		}
		contractorBucket := accountBucket.Bucket([]byte(contractorsBucket))
		if contractorBucket == nil {
			return fmt.Errorf("could not find contractor bucket for account %s: %w", account, tkerrors.ErrEntityNotFound)
		}

		err := contractorBucket.ForEach(func(key, value []byte) error {
			var contractor api.Contractor
			err := go_proto.Unmarshal(value, &contractor)
			if err != nil {
				return err
			}
			results[string(key)] = &contractor
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	return results, nil
}

// GetContractor returns a single contractor object by key
func (db *BoltDB) GetContractor(account, key string) (*api.Contractor, error) {
	var storedContractor api.Contractor

	err := db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return fmt.Errorf("could not find account bucket for %s: %w", account, tkerrors.ErrEntityNotFound)
		}
		contractorBucket := accountBucket.Bucket([]byte(contractorsBucket))
		if contractorBucket == nil {
			return fmt.Errorf("could not find contractor bucket for account %s: %w", account, tkerrors.ErrEntityNotFound)
		}

		contractorRaw := contractorBucket.Get([]byte(key))
		if contractorRaw == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(contractorRaw, &storedContractor)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &storedContractor, nil
}

// CreateContractor creates a single contractor and returns a generated key
func (db *BoltDB) CreateContractor(account string, newContractor *api.Contractor) error {
	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return fmt.Errorf("could not find account bucket for %s: %w", account, tkerrors.ErrEntityNotFound)
		}
		contractorBucket := accountBucket.Bucket([]byte(contractorsBucket))
		if contractorBucket == nil {
			return fmt.Errorf("could not find contractor bucket for account %s: %w", account, tkerrors.ErrEntityNotFound)
		}

		// Generate new key for entity
		key, err := db.getNewKey(contractorBucket)
		if err != nil {
			return err
		}
		newContractor.Id = key

		for _, job := range newContractor.Jobs {
			err := db.handleNewJob(tx, account, job)
			if err != nil {
				return err
			}
		}

		contractorRaw, err := go_proto.Marshal(newContractor)
		if err != nil {
			return err
		}

		err = contractorBucket.Put([]byte(key), contractorRaw)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateContractor modifies a single contractor by key
func (db *BoltDB) UpdateContractor(account, key string, updatedContractor *api.Contractor) error {
	var storedContractor api.Contractor

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return fmt.Errorf("could not find account bucket for %s: %w", account, tkerrors.ErrEntityNotFound)
		}
		contractorBucket := accountBucket.Bucket([]byte(contractorsBucket))
		if contractorBucket == nil {
			return fmt.Errorf("could not find contractor bucket for account %s: %w", account, tkerrors.ErrEntityNotFound)
		}

		currentContractor := contractorBucket.Get([]byte(key))
		if currentContractor == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(currentContractor, &storedContractor)
		if err != nil {
			return err
		}

		contractorRaw, err := go_proto.Marshal(updatedContractor)
		if err != nil {
			return err
		}

		err = contractorBucket.Put([]byte(key), contractorRaw)
		if err != nil {
			return err
		}

		// // figure out which changes to need to be made to other objects due to contractor additions/removals
		// additions, removals := listutil.FindListUpdates(storedContractor.Formulas, updatedContractor.Formulas)

		// // Append contractor id to contractor list in contractor
		// for _, formulaID := range additions {
		// 	err := db.linkContractorToFormula(tx, account, updatedContractor.Id, formulaID)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// // Remove contractor id from contractors list in contractors removed
		// for _, formulaID := range removals {
		// 	err := db.unlinkContractorFromFormula(tx, account, updatedContractor.Id, formulaID)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// // // DeleteContractor removes a contractor by key
// // func (db *BoltDB) DeleteContractor(account, key string) error {
// // 	var storedContractor api.Contractor

// // 	err := db.store.Update(func(tx *bolt.Tx) error {
// // 		accountBucket := tx.Bucket([]byte(account))
// // 		if accountBucket == nil {
// // 			return tkerrors.ErrEntityNotFound
// // 		}
// // 		contractorsBucket := accountBucket.Bucket([]byte(ContractorsBucket))

// // 		// First check if key exists
// // 		currentContractor := contractorsBucket.Get([]byte(key))
// // 		if currentContractor == nil {
// // 			return tkerrors.ErrEntityNotFound
// // 		}

// // 		err := go_proto.Unmarshal(currentContractor, &storedContractor)
// // 		if err != nil {
// // 			return err
// // 		}

// // 		err = contractorsBucket.Delete([]byte(key))
// // 		if err != nil {
// // 			return err
// // 		}

// // 		// Remove contractor id from contractors list in contractors this was linked to
// // 		for _, formulaID := range storedContractor.Formulas {
// // 			err := db.unlinkContractorFromFormula(tx, account, key, formulaID)
// // 			if err != nil {
// // 				return err
// // 			}
// // 		}
// // 		return nil
// // 	})
// // 	if err != nil {
// // 		return err
// // 	}
// // 	return nil
// // }
