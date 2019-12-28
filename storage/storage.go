package storage

import (
	"fmt"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
)

// Bucket represents the name of a section of key/value pairs
// usually a grouping of some sort
// ex. A key/value pair of userid-userdata would belong in the users bucket
type Bucket string

const (
	// FormulasBucket represents the container in which formulas are kept in the database
	FormulasBucket Bucket = "formulas"

	// JobsBucket represents the container in which jobs are kept in the database
	JobsBucket Bucket = "jobs"

	// UsersBucket represents the container in which users are managed
	UsersBucket Bucket = "users"
)

// EngineType type represents the different possible storage engines available
type EngineType string

const (
	// EngineGoogleDatastore represents a google datastore
	// a distributed key-value store
	// https://cloud.google.com/datastore/docs/concepts/overview
	EngineGoogleDatastore EngineType = "googleDatastore"
)

// Engine represents backend storage implementations where items can be persisted
type Engine interface {
	Init(config *config.Config) error
	GetAllUsers() (map[string]*api.User, error)
	GetUser(name string) (*api.User, error)
	CreateUser(name string, user *api.User) error
	GetAllFormulas(account string) (map[string]*api.Formula, error)
	GetFormula(account, key string) (*api.Formula, error)
	AddFormula(account string, formula *api.Formula) (key string, err error)
	UpdateFormula(account, key string, formula *api.Formula) error
	DeleteFormula(account, key string) error
	GetAllJobs(account string) (map[string]*api.Job, error)
	GetJob(account, key string) (*api.Job, error)
	AddJob(account string, job *api.Job) (key string, err error)
	UpdateJob(account, key string, job *api.Job) error
	DeleteJob(account, key string) error
}

// InitStorage creates a storage object with the appropriate engine
func InitStorage() (Engine, error) {
	config, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	engineType := EngineType(config.Database.Engine)

	switch engineType {
	case EngineGoogleDatastore:

		datastoreEngine := googleDatastore{}
		err := datastoreEngine.Init(config)
		if err != nil {
			return nil, err
		}

		return &datastoreEngine, nil
	default:
		return nil, fmt.Errorf("storage backend not implemented: %s", engineType)
	}
}
