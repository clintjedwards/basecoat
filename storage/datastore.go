package storage

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"

	"github.com/clintjedwards/toolkit/listutil"
	"github.com/clintjedwards/toolkit/logger"
	"github.com/clintjedwards/toolkit/tkerrors"
)

type googleDatastore struct {
	client   *datastore.Client
	idLength int
	timeout  time.Duration
}

func newGoogleDatastore(config *config.Config) (googleDatastore, error) {

	db := googleDatastore{}

	if config.Database.GoogleDatastore.EmulatorHost != "" {
		os.Setenv("DATASTORE_EMULATOR_HOST", config.Database.GoogleDatastore.EmulatorHost)
		logger.Log().Infow("connecting to google datastore",
			"emulator_host", config.Database.GoogleDatastore.EmulatorHost)
	}

	client, err := datastore.NewClient(context.Background(), config.Database.GoogleDatastore.ProjectID)
	if err != nil {
		return googleDatastore{}, err
	}

	db.client = client
	db.timeout, err = time.ParseDuration(config.Database.GoogleDatastore.Timeout)
	if err != nil {
		return googleDatastore{}, err
	}

	db.idLength = config.Database.IDLength

	return db, nil
}

// CreateParentKeys creates the initial account string key in all buckets so that assets
// can be separated by account
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
			return nil, tkerrors.ErrEntityNotFound
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
			return tkerrors.ErrEntityExists
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
			return nil, tkerrors.ErrEntityNotFound
		}
		return nil, err
	}

	return &formula, nil
}

func (db *googleDatastore) AddFormula(account string, newFormula *api.Formula) (key string, err error) {

	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err = db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		key, err = db.getNewKey(tx, account, FormulasBucket)
		if err != nil {
			return err
		}

		newNameKey := datastore.NameKey(string(FormulasBucket), key, parentKey)

		newFormula.Id = key
		// If the user has not entered a formula number just make it the ID
		if newFormula.Number == "" {
			newFormula.Number = newFormula.Id
		}

		// insert new formula
		_, err = tx.Put(newNameKey, newFormula)
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

func (db *googleDatastore) UpdateFormula(account, key string, updatedFormula *api.Formula) error {
	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	updateKey := datastore.NameKey(string(FormulasBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {

		// make sure item exists
		var formula api.Formula
		err := tx.Get(updateKey, &formula)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return tkerrors.ErrEntityNotFound
			}
			return err
		}

		// persist formula changes
		_, err = tx.Put(updateKey, updatedFormula)
		if err != nil {
			return err
		}

		// figure out which changes to need to be made to other objects due to job additions/removals
		additions, removals := listutil.FindListUpdates(formula.Jobs, updatedFormula.Jobs)

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

func (db *googleDatastore) DeleteFormula(account, key string) error {
	parentKey := datastore.NameKey(string(FormulasBucket), account, nil)
	deleteKey := datastore.NameKey(string(FormulasBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {

		// make sure item exists
		var formula api.Formula
		err := tx.Get(deleteKey, &formula)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return tkerrors.ErrEntityNotFound
			}
			return err
		}

		// persist delete
		err = tx.Delete(deleteKey)
		if err != nil {
			return err
		}

		// Remove formula id from formulas list in jobs this was linked to
		for _, jobID := range formula.Jobs {
			err := db.unlinkFormulaFromJob(tx, account, formula.Id, jobID)
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
			return nil, tkerrors.ErrEntityNotFound
		}
		return nil, err
	}

	return &job, nil
}

func (db *googleDatastore) AddJob(account string, newJob *api.Job) (key string, err error) {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err = db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		key, err = db.getNewKey(tx, account, JobsBucket)
		if err != nil {
			return err
		}

		newNameKey := datastore.NameKey(string(JobsBucket), key, parentKey)
		newJob.Id = key

		// insert new job
		_, err = tx.Put(newNameKey, newJob)
		if err != nil {
			return err
		}

		// for all formulas included in newJob make sure to add new job ID to all
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
	return key, err
}

func (db *googleDatastore) UpdateJob(account, key string, updatedJob *api.Job) error {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	updateKey := datastore.NameKey(string(JobsBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {

		// make sure item exists
		var job api.Job
		err := tx.Get(updateKey, &job)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return tkerrors.ErrEntityNotFound
			}
			return err
		}

		// persist job changes
		_, err = tx.Put(updateKey, updatedJob)
		if err != nil {
			return err
		}

		// figure out which changes to need to be made to other objects due to job additions/removals
		additions, removals := listutil.FindListUpdates(job.Formulas, updatedJob.Formulas)

		// Append job id to job list in formula
		for _, formulaID := range additions {
			err := db.linkJobToFormula(tx, account, updatedJob.Id, formulaID)
			if err != nil {
				return err
			}
		}

		// Remove job id from job list in formulas removed
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

func (db *googleDatastore) DeleteJob(account, key string) error {
	parentKey := datastore.NameKey(string(JobsBucket), account, nil)
	deleteKey := datastore.NameKey(string(JobsBucket), key, parentKey)

	tctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()

	_, err := db.client.RunInTransaction(tctx, func(tx *datastore.Transaction) error {
		// make sure item exists
		var job api.Job
		err := tx.Get(deleteKey, &job)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				return tkerrors.ErrEntityNotFound
			}
			return err
		}

		// persist delete
		err = tx.Delete(deleteKey)
		if err != nil {
			return err
		}

		// Remove job id from jobs list in formulas this was linked to
		for _, formulaID := range job.Formulas {
			err := db.unlinkJobFromFormula(tx, account, job.Id, formulaID)
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
