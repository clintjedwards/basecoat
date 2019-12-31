package storage

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/password"
	"github.com/clintjedwards/toolkit/tkerrors"
	go_proto "github.com/golang/protobuf/proto"
)

// GetAllAccounts returns all user accounts via map
// func (db *BoltDB) GetAllAccounts() (map[string]*api.Account, error) {
// 	results := map[string]*api.Account{}

// 	db.store.View(func(tx *bolt.Tx) error {
// 		err := tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
// 			storedAccount := &api.Account{}
// 			accountRaw := bucket.Get(name)
// 			if accountRaw == nil {
// 				return fmt.Errorf("could not get user information for account: %s; %w", name, tkerrors.ErrEntityNotFound)
// 			}

// 			err := go_proto.Unmarshal(accountRaw, storedAccount)
// 			if err != nil {
// 				return err
// 			}

// 			results[string(name)] = storedAccount

// 			return nil
// 		})
// 		return err
// 	})

// 	return results, nil
// }

// GetAccount returns a single account by id
// Accounts are both a bucket and an item inside that bucket so we must ask twice for id
func (db *BoltDB) GetAccount(id string) (*api.Account, error) {

	var storedAccount api.Account

	err := db.store.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(id))
		if bucket == nil {
			return fmt.Errorf("could not get bucket for account: %s; %w", id, tkerrors.ErrEntityNotFound)
		}

		accountRaw := bucket.Get([]byte(id))
		if accountRaw == nil {
			return fmt.Errorf("could not get account: %s;%w", id, tkerrors.ErrEntityNotFound)
		}

		err := go_proto.Unmarshal(accountRaw, &storedAccount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &storedAccount, nil
}

// GetAllAccounts returns all user accounts by map
func (db *BoltDB) GetAllAccounts() (map[string]*api.Account, error) {
	results := map[string]*api.Account{}

	db.store.View(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			storedAccount := &api.Account{}
			accountRaw := bucket.Get(name)
			if accountRaw == nil {
				return tkerrors.ErrEntityNotFound
			}

			err := go_proto.Unmarshal(accountRaw, storedAccount)
			if err != nil {
				return err
			}

			results[string(name)] = storedAccount

			return nil
		})
		return err
	})

	return results, nil
}

// CreateAccount adds a new user account and creates the required buckets
func (db *BoltDB) CreateAccount(id, pass string) error {
	err := db.store.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucketIfNotExists([]byte(id))
		if err != nil {
			return fmt.Errorf("could not create account: %s; %v", id, err)
		}

		err = db.createBuckets(bucket, formulasBucket, contractorsBucket)
		if err != nil {
			return fmt.Errorf("could not create account buckets: %s; %v", id, err)
		}

		hash, err := password.HashPassword([]byte(pass))
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}

		newAccount := api.Account{
			Id:       id,
			Hash:     hash,
			Created:  time.Now().Unix(),
			Modified: time.Now().Unix(),
		}

		accountRaw, err := go_proto.Marshal(&newAccount)
		if err != nil {
			return fmt.Errorf("could not marshal user account: %v", err)
		}

		err = bucket.Put([]byte(id), accountRaw)
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
