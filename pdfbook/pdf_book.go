package pdfbook

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"

	"github.com/byebyebruce/wadu/model"
	"github.com/byebyebruce/wadu/pkg/pdfx"
	"github.com/byebyebruce/wadu/vlm"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

//go:embed prompt.txt
var Prompt string

// BookPage 书页
type BookPage struct {
	Title     string   `json:"title"`
	Page      int      `json:"page"`
	Sentences []string `json:"sentences"`
}

func genPageInfo(ctx context.Context, openaiCli *openai.Client, vlmModel string, img []byte) (*BookPage, error) {
	return vlm.ChatImageJSON[BookPage](ctx, openaiCli, vlmModel, Prompt, img)
}

// GenFromPDF 从PDF生成 model.RawBook
func GenFromPDF(ctx context.Context, openaiCli *openai.Client, vlmModel string, pdf io.Reader) (*model.RawBook, error) {
	imgs, err := pdfx.ConvertPDFToJPEGWithQuality(pdf, 10)
	if err != nil {
		return nil, err
	}
	slog.Info("pdf", "pages", len(imgs))
	if len(imgs) == 0 {
		return nil, fmt.Errorf("PDF没有内容")
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(3) // 限制并发数

	book := &model.RawBook{
		Pages: make([]model.RawPage, len(imgs)),
	}
	for _i, _img := range imgs {
		i, img := _i, _img
		eg.Go(func() error {
			page, err := genPageInfo(egCtx, openaiCli, vlmModel, img)
			if err != nil {
				slog.Error("genPageInfo", "error", err)
				return err
			}
			slog.Info("page", "page", page)

			rawPage := model.RawPage{
				RawImage: img,
			}

			rawPage.Sentences = append(rawPage.Sentences, page.Sentences...)

			book.Pages[i] = rawPage

			if book.Title == "" && len(page.Title) > 0 {
				book.Title = page.Title
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return book, nil
}
