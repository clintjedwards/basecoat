package storage

import (
	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/listutil"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

// GetAllJobs returns all jobs using a map
func (db *BoltDB) GetAllJobs(account string) (map[string]*api.Job, error) {
	results := map[string]*api.Job{}

	db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(jobsBucket))

		err := targetBucket.ForEach(func(key, value []byte) error {
			var job api.Job
			err := go_proto.Unmarshal(value, &job)
			if err != nil {
				return err
			}
			results[string(key)] = &job
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	return results, nil
}

// GetJob returns a single job object by key
func (db *BoltDB) GetJob(account, key string) (*api.Job, error) {
	var storedJob api.Job

	err := db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(jobsBucket))

		jobRaw := targetBucket.Get([]byte(key))
		if jobRaw == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(jobRaw, &storedJob)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &storedJob, nil
}

// AddJob creates a single job and returns a generated key
func (db *BoltDB) AddJob(account string, newJob *api.Job) (key string, err error) {
	err = db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(jobsBucket))

		key, err = db.getNewKey(targetBucket)

		newJob.Id = key

		jobRaw, err := go_proto.Marshal(newJob)
		if err != nil {
			return err
		}

		err = targetBucket.Put([]byte(key), jobRaw)
		if err != nil {
			return err
		}

		// for all jobs included in newjob make sure to add new job ID to all
		for _, formulaID := range newJob.Formulas {
			err := db.linkJobToFormula(accountBucket, key, formulaID)
			if err != nil {
				return err
			}
		}

		if newJob.ContractorId != "" {
			err := db.linkJobToContractor(accountBucket, key, newJob.ContractorId)
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

// UpdateJob modifies a single job by key
func (db *BoltDB) UpdateJob(account, key string, updatedJob *api.Job) error {
	var storedJob api.Job

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(jobsBucket))

		// First check if key exists
		currentJob := targetBucket.Get([]byte(key))
		if currentJob == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(currentJob, &storedJob)
		if err != nil {
			return err
		}

		jobRaw, err := go_proto.Marshal(updatedJob)
		if err != nil {
			return err
		}

		err = targetBucket.Put([]byte(key), jobRaw)
		if err != nil {
			return err
		}

		// figure out which changes to need to be made to other objects due to job additions/removals
		additions, removals := listutil.FindListUpdates(storedJob.Formulas, updatedJob.Formulas)

		// Append job id to job list in job
		for _, formulaID := range additions {
			err := db.linkJobToFormula(accountBucket, updatedJob.Id, formulaID)
			if err != nil {
				return err
			}
		}

		// Remove job id from jobs list in jobs removed
		for _, formulaID := range removals {
			err := db.unlinkJobFromFormula(accountBucket, updatedJob.Id, formulaID)
			if err != nil {
				return err
			}
		}

		// handle contractor linking
		if updatedJob.ContractorId != storedJob.ContractorId {
			err := db.unlinkJobFromContractor(accountBucket, updatedJob.Id, storedJob.ContractorId)
			if err != nil {
				return err
			}

			if updatedJob.ContractorId != "" {
				err := db.linkJobToContractor(accountBucket, updatedJob.Id, updatedJob.ContractorId)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteJob removes a job by key
func (db *BoltDB) DeleteJob(account, key string) error {
	var storedJob api.Job

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		targetBucket := accountBucket.Bucket([]byte(jobsBucket))

		// First check if key exists
		currentJob := targetBucket.Get([]byte(key))
		if currentJob == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(currentJob, &storedJob)
		if err != nil {
			return err
		}

		err = targetBucket.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove job id from jobs list in jobs this was linked to
		for _, formulaID := range storedJob.Formulas {
			err := db.unlinkJobFromFormula(accountBucket, key, formulaID)
			if err != nil {
				return err
			}
		}

		if storedJob.ContractorId != "" {
			err = db.unlinkJobFromContractor(accountBucket, storedJob.Id, storedJob.ContractorId)
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
