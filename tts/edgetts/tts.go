package edgetts

import (
	"bytes"
	"context"
	"os"

	"github.com/byebyebruce/wadu/tts"
)

var _ tts.TTS = (*EdgeTTS)(nil)

type EdgeTTS struct {
}

func New() *EdgeTTS {
	return &EdgeTTS{}
}

// Synthesis implements tts.TTS.
func (e *EdgeTTS) SynthesisFile(ctx context.Context, text string, file string, option ...tts.Option) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := e.Synthesis(ctx, text, option...)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// Synthesis implements tts.TTS.
func (e *EdgeTTS) Synthesis(ctx context.Context, text string, option ...tts.Option) ([]byte, error) {
	cfg := &tts.Config{
		Voice:     XiaoxiaoNeural,
		AudioType: "mp3",
	}
	cfg.Apply(option...)

	o := Option{
		voice:  cfg.Voice,
		rate:   "+0%",
		volume: "+0%",
	}

	bytesBuffer := bytes.NewBuffer(nil)

	err := TTS(ctx, text, bytesBuffer, o)
	if err != nil {
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}
