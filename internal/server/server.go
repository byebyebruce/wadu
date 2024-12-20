package server

import (
	"net/http"

	"github.com/byebyebruce/wadu/internal/biz"

	"github.com/gin-gonic/gin"
)

const (
	APIPathCreateBook = "/api/book/create"
)

type Server struct {
	biz    *biz.Biz
	assets string
}

func NewServer(b *biz.Biz, assets string) *Server {
	return &Server{
		biz:    b,
		assets: assets,
	}
}

func (w *Server) Run(addr string) error {
	r := gin.Default()
	r.Static("/assets", w.assets)
	r.LoadHTMLGlob("internal/server/static/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})
	r.GET("/read", func(c *gin.Context) {
		c.HTML(200, "read.html", gin.H{})
	})
	r.GET("/upload", func(c *gin.Context) {
		c.HTML(200, "upload.html", gin.H{})
	})
	r.StaticFS("/audio", http.Dir("audio"))

	r.GET("/api/book/list", w.ListBook)
	r.GET("/api/book/:id", w.GetBook)
	r.DELETE("/api/book/delete/:id", w.DeleteBook)
	r.POST("/api/book/gen", w.GenBook)
	r.POST("/api/book/create", w.CreateFromRawBook)

	return r.Run(addr)
}
