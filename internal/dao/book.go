package dao

import (
	"encoding/json"
	"fmt"
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

		b = tx.Bucket([]byte(BookIndexBucket))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		bi, err := get[model.BookIndex](tx, string(BookIndexBucket), BookIndexKey)
		if err != nil {
			return err
		}

		bio := model.BookInfo{
			ID:        a.ID,
			Title:     a.Title,
			PublishAt: a.PublishAt,
			TotalPage: len(a.Pages),
		}
		for _, v := range a.Pages {
			if v.ImageURL != "" {
				bio.CoverURL = v.ImageURL
				break
			}
		}
		bi.Books = append([]model.BookInfo{bio}, bi.Books...)
		bi.Total = len(bi.Books)
		return set(tx, string(BookIndexBucket), BookIndexKey, bi)
	})
	return nil
}

// ListBook 列出书
func (d *Dao) ListBook(from, count int) ([]model.BookInfo, int, error) {
	var (
		bi    *model.BookIndex
		err   error
		total int
	)
	d.db.View(func(tx *bolt.Tx) error {
		bi, err = get[model.BookIndex](tx, string(BookIndexBucket), BookIndexKey)
		return err
	})
	if err != nil {
		return nil, total, err
	}
	total = bi.Total
	if from >= total {
		return nil, total, nil
	}
	to := from + count
	if count <= 0 {
		to = total
	}
	if to > total {
		to = total
	}
	return bi.Books[from:to], total, nil
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
		bi, err := get[model.BookIndex](tx, string(BookIndexBucket), BookIndexKey)
		if err != nil {
			return err
		}
		for i, v := range bi.Books {
			if v.ID == id {
				bi.Books = append(bi.Books[:i], bi.Books[i+1:]...)
				break
			}
		}
		bi.Total = len(bi.Books)
		if err := set(tx, string(BookIndexBucket), BookIndexKey, bi); err != nil {
			return err
		}
		bs := tx.Bucket([]byte(BookBucket))
		bs.Delete([]byte(id))
		return nil
	})
}
