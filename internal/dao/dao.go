package dao

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	BookBucket = []byte("book")
)
var (
	// ErrNotFound not found
	ErrNotFound = fmt.Errorf("not found")
)

type Dao struct {
	db *bolt.DB
}

func New(f string) (*Dao, error) {
	db, err := bolt.Open(f, 0644, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(BookBucket)
		return err
	})
	if err != nil {
		return nil, err
	}

	d := &Dao{
		db: db,
	}
	return d, nil
}

func (d *Dao) NextID() (uint64, error) {
	var (
		id  uint64
		err error
	)
	d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BookBucket))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		id, err = b.NextSequence()
		return nil
	})

	return id, err
}

func (d *Dao) Close() error {
	return d.db.Close()
}

func get[T any](tx *bolt.Tx, bucket string, key string) (*T, error) {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, fmt.Errorf("bucket not found")
	}
	v := b.Get([]byte(key))
	if v == nil {
		return nil, ErrNotFound
	}
	if len(v) == 0 {
		return nil, nil
	}
	var ret T
	if err := json.Unmarshal(v, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func set[T any](tx *bolt.Tx, bucket string, key string, data T) error {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return fmt.Errorf("bucket not found")
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return b.Put([]byte(key), bs)
}

func list[T any](tx *bolt.Tx, bucket string) ([]T, error) {
	var as []T
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, fmt.Errorf("bucket not found")
	}
	c := b.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var a T
		if err := json.Unmarshal(v, &a); err != nil {
			return nil, err
		}
		as = append(as, a)
	}
	return as, nil
}
