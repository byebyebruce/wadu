package server

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/byebyebruce/wadu/model"

	"github.com/gin-gonic/gin"
)

func (w *Server) GenBook(c *gin.Context) {
	//title := c.PostForm("user")
	pdf, err := c.FormFile("doc")
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 读取文件
	docFile, err := pdf.Open()
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		slog.Error("open file error", "error", err)
		return
	}
	defer docFile.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	// audio
	var audioBytes []byte
	audio, err := c.FormFile("audio")
	if err == nil && audio != nil {

		// 读取文件
		audioFile, err := audio.Open()
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			slog.Error("open file error", "error", err)
			return
		}
		defer audioFile.Close()
		audioBytes, err = io.ReadAll(audioFile)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			slog.Error("open file error", "error", err)
			return
		}
	}

	a, err := w.biz.GenFromPDF(ctx, docFile, audioBytes)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, a)
}

func (w *Server) CreateFromRawBook(c *gin.Context) {
	var a model.RawBook
	if err := c.BindJSON(&a); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	book, err := w.biz.CreateFromRawBook(ctx, &a)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (w *Server) ListBook(c *gin.Context) {
	// TODO
	as, _, err := w.biz.DB.ListBook(0, 0)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := make([]BookInfo, 0, len(as))
	for _, a := range as {
		b := BookInfo{
			ID:        a.ID,
			Title:     a.Title,
			PublishAt: a.PublishAt,
			TotalPage: len(a.Pages),
		}
		for _, p := range a.Pages {
			if p.ImageURL != "" {
				b.CoverURL = p.ImageURL
				break
			}
		}
		resp = append(resp, b)
	}

	c.JSON(http.StatusOK, resp)
}

func (w *Server) GetBook(c *gin.Context) {
	a, err := w.biz.DB.GetBook(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (w *Server) DeleteBook(c *gin.Context) {
	err := w.biz.DB.DeleteBook(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
