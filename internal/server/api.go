package server

type BookInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	PublishAt int64  `json:"publish_at"`
	CoverURL  string `json:"cover_url"`
	TotalPage int    `json:"total_page"`
}
