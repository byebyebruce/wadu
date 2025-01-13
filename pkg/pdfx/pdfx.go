package pdfx

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"

	"github.com/gen2brain/go-fitz"
)

// ConvertPDFFileToJPEG 将PDF文件转换为JPEG图像
func ConvertPDFFileToJPEG(f string) ([][]byte, error) {
	in, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer in.Close()
	return ConvertPDFToJPEGWithQuality(in, 100)
}

// ConvertPDFToJPEG 将PDF reader转换为JPEG图像
func ConvertPDFToJPEG(in io.Reader) ([][]byte, error) {
	return ConvertPDFToJPEGWithQuality(in, 100)
}

func ConvertPDFToJPEGWithQuality(in io.Reader, quality int) ([][]byte, error) {
	doc, err := fitz.NewFromReader(in)
	if err != nil {
		return nil, fmt.Errorf("打开PDF失败: %v", err)
	}
	defer doc.Close()

	var images [][]byte
	// 遍历PDF的每一页
	for i := 0; i < doc.NumPage(); i++ {
		// 将PDF页面转换为图像
		imgBytes, err := doc.Image(i)
		if err != nil {
			return nil, fmt.Errorf("转换第 %d 页失败: %v", i+1, err)
		}

		bf := &bytes.Buffer{}

		/*

			//"github.com/chai2010/webp"
			if err := webp.Encode(bf, imgBytes, &webp.Options{Lossless: true}); err != nil {
				return nil, fmt.Errorf("编码第 %d 页为WebP: %v", i, err)
			}
		*/

		// 将页面保存为JPEG图像
		err = jpeg.Encode(bf, imgBytes, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, fmt.Errorf("encode 第 %d 页为图像: %v", i, err)
		}
		images = append(images, bf.Bytes())
	}

	return images, nil
}
