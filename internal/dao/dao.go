package dao

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var (
	BookBucket      = []byte("book")
	BookIndexBucket = []byte("book_index")
	BookIndexKey    = "book_index_key"
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

	d := &Dao{
		db: db,
	}
	return d, nil
}

func (d *Dao) InitDB() error {
	err := d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(BookBucket)
		return err
	})
	if err != nil {
		return err
	}

	err = upgradeDB_V1(d.db)
	if err != nil {
		return err
	}

	return err
}

// NextID 生成下一个ID
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

/*
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
*/

func listForward[T any](tx *bolt.Tx, bucket string, from int, count int) (as []T, total int, err error) {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, 0, fmt.Errorf("bucket not found")
	}
	total = b.Stats().KeyN
	var (
		idx = 0
		c   = b.Cursor()
	)

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if count > 0 {
			if len(as) >= count {
				break
			}
		}
		if idx < from {
			idx++
			continue
		}
		idx++

		var a T
		if err = json.Unmarshal(v, &a); err != nil {
			return
		}
		as = append(as, a)
	}
	return
}
func listBackward[T any](tx *bolt.Tx, bucket string, from int, count int) (as []T, total int, err error) {
	b := tx.Bucket([]byte(bucket))
	if b == nil {
		return nil, 0, fmt.Errorf("bucket not found")
	}
	total = b.Stats().KeyN
	var (
		idx = 0
		c   = b.Cursor()
	)

	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		if count > 0 {
			if len(as) >= count {
				break
			}
		}
		if idx < from {
			idx++
			continue
		}
		idx++

		var a T
		if err = json.Unmarshal(v, &a); err != nil {
			return
		}
		as = append(as, a)
	}
	return
}
