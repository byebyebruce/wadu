package server

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/byebyebruce/wadu/internal/biz"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var htmlFS embed.FS

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

func (w *Server) Run(addr string, debug bool) error {
	r := gin.Default()
	r.Static("/assets", w.assets)

	// static
	if debug {
		r.LoadHTMLGlob("internal/server/static/*.html")
	} else {
		tmpl := template.Must(template.ParseFS(htmlFS, "static/*.html"))
		r.SetHTMLTemplate(tmpl)
	}

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
	r.POST(APIPathCreateBook, w.CreateFromRawBook)

	return r.Run(addr)
}
