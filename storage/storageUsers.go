package storage

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

// GetAllUsers returns all user accounts by map
func (db *BoltDB) GetAllUsers() (map[string]*api.User, error) {
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

// GetUser returns a single user by account id
func (db *BoltDB) GetUser(name string) (*api.User, error) {

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

// CreateUser adds a new user account and creates the appropriate buckets
func (db *BoltDB) CreateUser(id string, newUser *api.User) error {
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
