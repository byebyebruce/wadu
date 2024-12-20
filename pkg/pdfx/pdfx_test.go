package pdfx

import "testing"

func TestConvert(t *testing.T) {
	pdfPath := "../testdata/a.pdf"
	_, err := ConvertPDFFileToJPEG(pdfPath)
	if err != nil {
		t.Errorf("转换失败: %v", err)
	}
}
