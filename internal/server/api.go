package server

import "github.com/byebyebruce/wadu/model"

type BookListResp struct {
	Total int              `json:"total"`
	Books []model.BookInfo `json:"books"`
}
