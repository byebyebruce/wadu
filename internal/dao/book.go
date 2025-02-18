package dao

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/byebyebruce/wadu/model"

	bolt "go.etcd.io/bbolt"
)

// CreateBook 创建书
func (d *Dao) CreateBook(a *model.Book) error {
	d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BookBucket))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		id, _ := b.NextSequence()
		a.ID = fmt.Sprintf("%d", id)
		a.PublishAt = time.Now().Unix()
		data, err := json.Marshal(a)
		if err != nil {
			return err
		}

		if err := b.Put([]byte(a.ID), data); err != nil {
			return err
		}
		return nil
	})
	return nil
}

// ListBook 列出书
func (d *Dao) ListBook(from, count int) ([]model.Book, int, error) {
	var (
		as    []model.Book
		err   error
		total int
	)
	d.db.View(func(tx *bolt.Tx) error {
		as, total, err = listBackward[model.Book](tx, string(BookBucket), from, count)
		b := tx.Bucket([]byte(BookBucket))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		return err
	})
	if err != nil {
		return nil, total, err
	}
	sort.Slice(as, func(i, j int) bool {
		return as[i].PublishAt > as[j].PublishAt
	})
	return as, total, nil
}

// UpdateBook 更新书
func (d *Dao) UpdateBook(a *model.Book) error {
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BookBucket))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		if err := b.Put([]byte(a.ID), data); err != nil {
			return err
		}
		return err
	})
}

// GetBook 获取书
func (d *Dao) GetBook(id string) (*model.Book, error) {
	var (
		a   *model.Book
		err error
	)
	d.db.View(func(tx *bolt.Tx) error {
		a, err = get[model.Book](tx, string(BookBucket), id)
		return err
	})

	return a, err
}

// DeleteBook 删除书
func (d *Dao) DeleteBook(id string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BookBucket))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		v := b.Get([]byte(id))
		if len(v) == 0 {
			return ErrNotFound
		}
		return b.Delete([]byte(id))
	})
}
