package storage

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
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
)

// BoltDB is a representation of the bolt datastore
type BoltDB struct {
	idLength int // length of generated IDs
	store    *bolt.DB
}

// Create a new boltdb from settings in config file
func newBoltDB(config *config.Config) (BoltDB, error) {
	db := BoltDB{}

	store, err := bolt.Open(config.Database.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return BoltDB{}, err
	}

	db.store = store
	db.idLength = config.Database.IDLength

	return db, nil
}

// createBuckets creates buckets inside of another bucket
func (db *BoltDB) createBuckets(tx *bolt.Tx, root *bolt.Bucket, buckets ...Bucket) error {

	for _, bucket := range buckets {
		_, err := root.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("could not create bucket: %s; %v", bucket, err)
		}
	}
	return nil
}

// InitStorage creates a storage object
func InitStorage(config *config.Config) (*BoltDB, error) {
	boltDBEngine, err := newBoltDB(config)
	if err != nil {
		return nil, err
	}

	return &boltDBEngine, nil
}
