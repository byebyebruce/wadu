package dao

import (
	"encoding/json"
	"sort"

	"github.com/byebyebruce/wadu/model"

	bolt "go.etcd.io/bbolt"
)

func upgradeDB_V1(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BookIndexBucket))
		if b != nil {
			return nil
		}
		_, err := tx.CreateBucket(BookIndexBucket)
		if err != nil {
			return err
		}
		b = tx.Bucket([]byte(BookBucket))
		if b == nil {
			return nil
		}
		bookIndex := make([]model.BookInfo, 0)
		b.ForEach(func(k, v []byte) error {
			var a model.Book
			if err := json.Unmarshal(v, &a); err != nil {
				return err
			}
			bi := model.BookInfo{
				ID:        a.ID,
				Title:     a.Title,
				PublishAt: a.PublishAt,
				TotalPage: len(a.Pages),
			}

			for _, v := range a.Pages {
				if v.ImageURL != "" {
					bi.CoverURL = v.ImageURL
					break
				}
			}
			bookIndex = append(bookIndex, bi)
			return nil
		})

		bi := model.BookIndex{
			Total: len(bookIndex),
			Books: bookIndex,
		}
		sort.Slice(bookIndex, func(i, j int) bool {
			return bookIndex[i].PublishAt > bookIndex[j].PublishAt
		})
		return set(tx, string(BookIndexBucket), BookIndexKey, bi)
	})
}
