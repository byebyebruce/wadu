package imagex

import (
	"bytes"
	"image/jpeg"
	"image/png"
)

// IsPNG 判断是否为PNG文件
func IsPNG(data []byte) bool {
	return len(data) > 4 && bytes.Equal(data[:4], []byte{0x89, 0x50, 0x4E, 0x47})
}

// IsJPEG 判断是否为JPEG文件
func IsJPEG(data []byte) bool {
	return len(data) > 2 && bytes.Equal(data[:2], []byte{0xFF, 0xD8})
}

// ConvertPNGtoJPEG 将PNG转换为JPEG并返回字节数组
func ConvertPNGtoJPEG(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var jpegBuffer bytes.Buffer
	err = jpeg.Encode(&jpegBuffer, img, nil)
	if err != nil {
		return nil, err
	}

	return jpegBuffer.Bytes(), nil
}
