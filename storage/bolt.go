package storage

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/toolkit/listutil"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

type boltDB struct {
	idLength int // length of generated IDs
	store    *bolt.DB
}

// Create a new boltdb from settings in config file
func newBoltDB(config *config.Config) (boltDB, error) {
	db := boltDB{}

	store, err := bolt.Open(config.Database.Bolt.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return boltDB{}, err
	}

	db.store = store
	db.idLength = config.Database.IDLength

	return db, nil
}

// createBuckets creates buckets inside of another bucket
func (db *boltDB) createBuckets(tx *bolt.Tx, root *bolt.Bucket, buckets ...Bucket) error {

	for _, bucket := range buckets {
		_, err := root.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("could not create bucket: %s; %v", bucket, err)
		}
	}
	return nil
}

func (db *boltDB) GetAllUsers() (map[string]*api.User, error) {
	results := map[string]*api.User{}

	db.store.View(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			storedUser := &api.User{}
			userRaw := bucket.Get(name)
			if userRaw == nil {
				return tkerrors.ErrEntityNotFound
			}

			err := go_proto.Unmarshal(userRaw, storedUser)
			if err != nil {
				return err
			}

			results[string(name)] = storedUser

			return nil
		})
		return err
	})

	return results, nil
}

func (db *boltDB) GetUser(name string) (*api.User, error) {

	var storedUser api.User

	err := db.store.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))

		userRaw := bucket.Get([]byte(name))
		if userRaw == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(userRaw, &storedUser)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &storedUser, nil
}

func (db *boltDB) CreateUser(id string, newUser *api.User) error {
	err := db.store.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(id))
		if err != nil {
			return fmt.Errorf("could not create user: %s; %v", id, err)
		}

		err = db.createBuckets(tx, bucket, FormulasBucket, JobsBucket)
		if err != nil {
			return fmt.Errorf("could not create user buckets: %s; %v", id, err)
		}

		userRaw, err := go_proto.Marshal(newUser)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(id), userRaw)
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

func (db *boltDB) GetAllFormulas(account string) (map[string]*api.Formula, error) {
	results := map[string]*api.Formula{}

	db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

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

func (db *boltDB) GetFormula(account, key string) (*api.Formula, error) {
	var storedFormula api.Formula

	err := db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

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

func (db *boltDB) AddFormula(account string, newFormula *api.Formula) (key string, err error) {
	err = db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

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

func (db *boltDB) UpdateFormula(account, key string, updatedFormula *api.Formula) error {
	var storedFormula api.Formula

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

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

func (db *boltDB) DeleteFormula(account, key string) error {
	var storedFormula api.Formula

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		formulasBucket := accountBucket.Bucket([]byte(FormulasBucket))

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

func (db *boltDB) GetAllJobs(account string) (map[string]*api.Job, error) {
	results := map[string]*api.Job{}

	db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

		err := jobsBucket.ForEach(func(key, value []byte) error {
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

func (db *boltDB) GetJob(account, key string) (*api.Job, error) {
	var storedJob api.Job

	err := db.store.View(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

		jobRaw := jobsBucket.Get([]byte(key))
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

func (db *boltDB) AddJob(account string, newJob *api.Job) (key string, err error) {
	err = db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

		key, err = db.getNewKey(jobsBucket)

		newJob.Id = key

		jobRaw, err := go_proto.Marshal(newJob)
		if err != nil {
			return err
		}

		err = jobsBucket.Put([]byte(key), jobRaw)
		if err != nil {
			return err
		}

		// for all jobs included in newjob make sure to add new job ID to all
		for _, formulaID := range newJob.Formulas {
			err := db.linkJobToFormula(tx, account, key, formulaID)
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

func (db *boltDB) UpdateJob(account, key string, updatedJob *api.Job) error {
	var storedJob api.Job

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

		// First check if key exists
		currentJob := jobsBucket.Get([]byte(key))
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

		err = jobsBucket.Put([]byte(key), jobRaw)
		if err != nil {
			return err
		}

		// figure out which changes to need to be made to other objects due to job additions/removals
		additions, removals := listutil.FindListUpdates(storedJob.Formulas, updatedJob.Formulas)

		// Append job id to job list in job
		for _, formulaID := range additions {
			err := db.linkJobToFormula(tx, account, updatedJob.Id, formulaID)
			if err != nil {
				return err
			}
		}

		// Remove job id from jobs list in jobs removed
		for _, formulaID := range removals {
			err := db.unlinkJobFromFormula(tx, account, updatedJob.Id, formulaID)
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

func (db *boltDB) DeleteJob(account, key string) error {
	var storedJob api.Job

	err := db.store.Update(func(tx *bolt.Tx) error {
		accountBucket := tx.Bucket([]byte(account))
		if accountBucket == nil {
			return tkerrors.ErrEntityNotFound
		}
		jobsBucket := accountBucket.Bucket([]byte(JobsBucket))

		// First check if key exists
		currentJob := jobsBucket.Get([]byte(key))
		if currentJob == nil {
			return tkerrors.ErrEntityNotFound
		}

		err := go_proto.Unmarshal(currentJob, &storedJob)
		if err != nil {
			return err
		}

		err = jobsBucket.Delete([]byte(key))
		if err != nil {
			return err
		}

		// Remove job id from jobs list in jobs this was linked to
		for _, formulaID := range storedJob.Formulas {
			err := db.unlinkJobFromFormula(tx, account, key, formulaID)
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
