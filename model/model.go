package model

type BookIndex struct {
	Total int        `json:"total"`
	Books []BookInfo `json:"books"`
}

type BookInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	PublishAt int64  `json:"publish_at"`
	CoverURL  string `json:"cover_url"`
	TotalPage int    `json:"total_page"`
}

type Book struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	PublishAt int64  `json:"publish_at"`
	Pages     Pages  `json:"pages"`
	FullAudio string `json:"full_audio"`
}

type Pages []Page

type Sentence struct {
	Content  string `json:"content"`
	AudioURL string `json:"audio_url"`
}

type Page struct {
	ID        int        `json:"id"`
	ImageURL  string     `json:"image_url"`
	Sentences []Sentence `json:"sentences"`
}

type RawPage struct {
	RawImage  []byte   `json:"raw_image"` // jpeg base64
	Sentences []string `json:"sentences"`
}

type RawBook struct {
	Title string    `json:"title"`
	Pages []RawPage `json:"pages"`
	MP3   []byte    `json:"mp3"`
}
