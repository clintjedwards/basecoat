package storage

import (
	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/listutil"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

// GetAllContractors returns all contractors using a map
func (db *BoltDB) GetAllContractors(account string) (map[string]*api.Contractor, error) {
	results := map[string]*api.Contractor{}

	db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(contractorsBucket))

		err := targetBucket.ForEach(func(key, value []byte) error {
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
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(contractorsBucket))

		contractorRaw := targetBucket.Get([]byte(key))
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

// AddContractor creates a single contractor and returns a generated key
func (db *BoltDB) AddContractor(account string, newContractor *api.Contractor) (key string, err error) {
	err = db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(contractorsBucket))

		key, err = db.getNewKey(targetBucket)

		newContractor.Id = key

		contractorRaw, err := go_proto.Marshal(newContractor)
		if err != nil {
			return err
		}

		err = targetBucket.Put([]byte(key), contractorRaw)
		if err != nil {
			return err
		}

		// for all jobs included in contractors make sure to link the contractor id
		for _, jobID := range newContractor.Jobs {
			err := db.linkContractorToJob(accountBucket, key, jobID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

// UpdateContractor modifies a single contractor by key
func (db *BoltDB) UpdateContractor(account, key string, updatedContractor *api.Contractor) error {
	var storedContractor api.Contractor

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(contractorsBucket))

		currentContractor := targetBucket.Get([]byte(key))
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

		err = targetBucket.Put([]byte(key), contractorRaw)
		if err != nil {
			return err
		}

		// figure out which changes to need to be made to other objects due to job additions/removals
		additions, removals := listutil.FindListUpdates(storedContractor.Jobs, updatedContractor.Jobs)

		// Append contractor id to list in job
		for _, jobID := range additions {
			err := db.linkContractorToJob(accountBucket, updatedContractor.Id, jobID)
			if err != nil {
				return err
			}
		}

		// Remove contractor id from list in jobs removed
		for _, jobID := range removals {
			err := db.unlinkContractorFromJob(accountBucket, jobID)
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

// DeleteContractor removes a contractor by key
func (db *BoltDB) DeleteContractor(account, key string) error {
	var storedContractor api.Contractor
	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(contractorsBucket))

		currentContractor := targetBucket.Get([]byte(key))
		if currentContractor == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(currentContractor, &storedContractor)
		if err != nil {
			return err
		}

		err = targetBucket.Delete([]byte(key))
		if err != nil {
			return err
		}

		// for all jobs included in contractors make sure to unlink the contractor id
		for _, jobID := range storedContractor.Jobs {
			err := db.unlinkContractorFromJob(accountBucket, jobID)
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
