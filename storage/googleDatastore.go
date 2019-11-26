package storage

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/utils"
)

type googleDatastore struct {
	client  *datastore.Client
	timeout time.Duration
}

func (db *googleDatastore) Init(config *config.Config) error {

	if config.Database.GoogleDatastore.EmulatorHost != "" {
		os.Setenv("DATASTORE_EMULATOR_HOST", config.Database.GoogleDatastore.EmulatorHost)
		utils.StructuredLog(utils.LogLevelInfo, "connecting to google datastore emulator", config.Database.GoogleDatastore.EmulatorHost)
	}

	client, err := datastore.NewClient(context.Background(), config.Database.GoogleDatastore.ProjectID)
	if err != nil {
		return err
	}

	db.client = client
	db.timeout, err = time.ParseDuration(config.Database.GoogleDatastore.Timeout)
	if err != nil {
		return err
	}

	return nil
}

// CreateParentKeys creates the initial account string key in all buckets so that assets
// can be seperated by account
func (db *googleDatastore) CreateParentKeys(account string) error {
	newFormulaParentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	newJobParentKey := datastore.NameKey(string(JobsBucket), account, nil)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {

		// Create formula parent key
		empty := api.User{}
		err := tx.Get(newFormulaParentKey, &empty)
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err = tx.Put(newFormulaParentKey, &empty)
		if err != nil {
			return err
		}

		// Create Job parent key
		err = tx.Get(newJobParentKey, &empty)
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err = tx.Put(newJobParentKey, &empty)
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

func (db *googleDatastore) GetAllUsers() (map[string]*api.User, error) {

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var rawUsers []*api.User
	query := datastore.NewQuery(string(UsersBucket))
	keys, err := db.client.GetAll(tctx, query, &rawUsers)
	if err != nil {
		utils.StructuredLog(utils.LogLevelError, "could not retrieve users from database", err)
		return nil, err
	}

	users := map[string]*api.User{}
	for index, key := range keys {
		users[key.Name] = rawUsers[index]
	}

	return users, nil
}

func (db *googleDatastore) GetUser(name string) (*api.User, error) {

	getUser := datastore.NameKey(string(UsersBucket), name, nil)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var user api.User
	err := db.client.Get(tctx, getUser, &user)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, utils.ErrEntityNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (db *googleDatastore) CreateUser(name string, newUser *api.User) error {
	newKey := datastore.NameKey(string(UsersBucket), name, nil)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		var user api.User
		err := tx.Get(newKey, &user)
		if err != datastore.ErrNoSuchEntity {
			return utils.ErrEntityExists
		}

		_, err = tx.Put(newKey, newUser)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.CreateParentKeys(name)
	if err != nil {
		return err
	}

	return nil
}

func (db *googleDatastore) GetAllFormulas(account string) (map[string]*api.Formula, error) {

	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var rawFormulas []*api.Formula
	query := datastore.NewQuery(string(FormulasBucket)).Ancestor(parentKey).Filter("__key__ >", parentKey)
	keys, err := db.client.GetAll(tctx, query, &rawFormulas)
	if err != nil {
		utils.StructuredLog(utils.LogLevelError, "could not retrieve formulas from database", err)
		return nil, err
	}

	formulas := map[string]*api.Formula{}
	for index, key := range keys {
		formulas[key.Name] = rawFormulas[index]
	}

	return formulas, nil
}

func (db *googleDatastore) GetFormula(account, key string) (*api.Formula, error) {

	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	getKey := datastore.NameKey(string(FormulasBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var formula api.Formula
	err := db.client.Get(tctx, getKey, &formula)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, utils.ErrEntityNotFound
		}
		return nil, err
	}

	return &formula, nil
}

func (db *googleDatastore) AddFormula(account, key string, newFormula *api.Formula) error {

	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	newKey := datastore.NameKey(string(FormulasBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		var formula api.Formula
		err := tx.Get(newKey, &formula)
		if err != datastore.ErrNoSuchEntity {
			return utils.ErrEntityExists
		}

		_, err = tx.Put(newKey, newFormula)
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

func (db *googleDatastore) UpdateFormula(account, key string, updatedFormula *api.Formula) error {
	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	updateKey := datastore.NameKey(string(FormulasBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		var formula api.Formula
		err := tx.Get(updateKey, &formula)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return utils.ErrEntityNotFound
			}
			return err
		}

		_, err = tx.Put(updateKey, updatedFormula)
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

func (db *googleDatastore) DeleteFormula(account, key string) error {
	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	deleteKey := datastore.NameKey(string(FormulasBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		var formula api.Formula
		err := tx.Get(deleteKey, &formula)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return utils.ErrEntityNotFound
			}
			return err
		}

		err = tx.Delete(deleteKey)
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

func (db *googleDatastore) GetAllJobs(account string) (map[string]*api.Job, error) {

	parentKey := datastore.NameKey(string(JobsBucket), account, nil)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var rawJobs []*api.Job
	query := datastore.NewQuery(string(JobsBucket)).Ancestor(parentKey).Filter("__key__ >", parentKey)
	keys, err := db.client.GetAll(tctx, query, &rawJobs)
	if err != nil {
		return nil, err
	}

	jobs := map[string]*api.Job{}
	for index, key := range keys {
		jobs[key.Name] = rawJobs[index]
	}

	return jobs, nil
}

func (db *googleDatastore) GetJob(account, key string) (*api.Job, error) {

	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	getKey := datastore.NameKey(string(JobsBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	var job api.Job
	err := db.client.Get(tctx, getKey, &job)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, utils.ErrEntityNotFound
		}
		return nil, err
	}

	return &job, nil
}

// First check if the ID exits
func (db *googleDatastore) AddJob(account, key string, newJob *api.Job) error {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	newKey := datastore.NameKey(string(JobsBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		var job api.Job
		err := tx.Get(newKey, &job)
		if err != datastore.ErrNoSuchEntity {
			return utils.ErrEntityExists
		}

		_, err = tx.Put(newKey, newJob)
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

func (db *googleDatastore) UpdateJob(account, key string, updatedJob *api.Job) error {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	updateKey := datastore.NameKey(string(JobsBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		var job api.Job
		err := tx.Get(updateKey, &job)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return utils.ErrEntityNotFound
			}
			return err
		}

		_, err = tx.Put(updateKey, updatedJob)
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

func (db *googleDatastore) DeleteJob(account, key string) error {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	deleteKey := datastore.NameKey(string(JobsBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		var job api.Job
		err := tx.Get(deleteKey, &job)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return utils.ErrEntityNotFound
			}
			return err
		}

		err = tx.Delete(deleteKey)
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
