package storage

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// Bucket represents the name of a section of key/value pairs
// usually a grouping of some sort
// ex. A key/value pair of userid-userdata would belong in the users bucket
type Bucket string

const (
	// formulasBucket represents the container in which formulas are kept in the database
	formulasBucket Bucket = "formulas"

	// jobsBucket represents the container in which jobs are kept in the database
	jobsBucket Bucket = "jobs"

	// contractorsBucket represents the container in which contractors are kept in the database
	contractorsBucket Bucket = "contractors"
)

// BoltDB is a representation of the bolt datastore
type BoltDB struct {
	idLength int // length of generated IDs
	store    *bolt.DB
}

// NewBoltDB creates a new boltdb from settings in config file
func NewBoltDB(path string, idlength int) (BoltDB, error) {
	db := BoltDB{}

	store, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return BoltDB{}, err
	}

	db.store = store
	db.idLength = idlength

	return db, nil
}

// createBuckets creates buckets inside of another bucket
func (db *BoltDB) createBuckets(root *bolt.Bucket, buckets ...Bucket) error {

	for _, bucket := range buckets {
		_, err := root.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("could not create bucket: %s; %v", bucket, err)
		}
	}
	return nil
}
