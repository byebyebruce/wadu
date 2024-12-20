package biz

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/byebyebruce/wadu/internal/dao"
	"github.com/byebyebruce/wadu/model"
	"github.com/byebyebruce/wadu/pdfbook"
	"github.com/byebyebruce/wadu/tts"
	"github.com/byebyebruce/wadu/vlm"
)

type Biz struct {
	openaiCli *vlm.Client
	tts       tts.TTS
	assetsDir string
	DB        *dao.Dao
}

func NewBiz(DB *dao.Dao, openaiCli *vlm.Client, tts tts.TTS, assetsDir string) *Biz {
	return &Biz{
		DB:        DB,
		openaiCli: openaiCli,
		tts:       tts,
		assetsDir: assetsDir,
	}
}

func (b *Biz) CreateFromRawBook(ctx context.Context, rawBook *model.RawBook) (*model.Book, error) {
	if len(rawBook.Pages) == 0 {
		return nil, fmt.Errorf("book has no pages")
	}
	id, err := b.DB.NextID()
	if err != nil {
		return nil, err
	}
	idStr := fmt.Sprintf("%d", id)

	dir := b.assetsDir
	book := &model.Book{
		ID:        idStr,
		Title:     rawBook.Title,
		PublishAt: time.Now().Unix(),
		Pages:     make([]model.Page, 0, len(rawBook.Pages)),
	}
	for i, rawPage := range rawBook.Pages {
		var sentences = make([]model.Sentence, 0, len(rawPage.Sentences))
		for j, sentence := range rawPage.Sentences {
			audioFile := fmt.Sprintf("%s_%d_%d.mp3", idStr, i, j)
			if _, err := b.doTTS(ctx, sentence, audioFile); err != nil {
				slog.Error("tts", "error", err)
				return nil, err
			}
			sentences = append(sentences, model.Sentence{
				Content:  sentence,
				AudioURL: audioFile,
			})
		}

		imageFile := fmt.Sprintf("%s_%d.jpeg", idStr, i)
		if err := os.WriteFile(filepath.Join(dir, imageFile), rawPage.RawImage, 0644); err != nil {
			slog.Error("write image", "error", err)
			return nil, err
		}
		/*
			imageFile := fmt.Sprintf("%s_%d.webp", idStr, i)
			if err := imagex.JPEGWriteWebp(rawPage.RawImage, filepath.Join(dir, imageFile)); err != nil {
				slog.Error("write image", "error", err)
				return nil, err
			}
		*/

		p := model.Page{
			ID:        i,
			ImageURL:  imageFile,
			Sentences: sentences,
		}

		book.Pages = append(book.Pages, p)
	}
	err = b.DB.CreateBook(book)
	return book, err
}

func (b *Biz) doTTS(ctx context.Context, text string, file string) (string, error) {
	buf, err := b.tts.Synthesis(ctx,
		text,
		tts.WithAudioType("mp3"),
		tts.WithAudioSpeed(0.8),
	)
	if err != nil {
		slog.Error("tts", "text", text, "error", err)
		return "", err
	}

	dest := filepath.Join(b.assetsDir, file)
	if err := os.WriteFile(dest, buf, 0644); err != nil {
		slog.Error("write audio", "error", err)
		return "", err
	}

	slog.Info("tts", "text", text, "file", dest)
	return dest, err
}
func (b *Biz) GenFromPDF(ctx context.Context, pdf io.Reader) (*model.RawBook, error) {
	return pdfbook.GenFromPDF(ctx, b.openaiCli.Client, b.openaiCli.Model, pdf)
}
