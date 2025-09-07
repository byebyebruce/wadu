package openaitts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/byebyebruce/wadu/tts"
)

var _ tts.TTS = (*TTS)(nil)

/*
type TTS struct {
	cli   *openai.Client
	model string
}

// SynthesisFile implements tts.TTS.
func (t *TTS) SynthesisFile(ctx context.Context, text string, file string, option ...tts.Option) error {
	opt := tts.Config{
		Voice:     "FunAudioLLM/CosyVoice2-0.5B:alex",
		AudioRate: 41000,
	}
	tts.ApplyOption(&opt, option...)
	req := openai.CreateSpeechRequest{
		Model:          openai.SpeechModel(t.model),
		Input:          text,
		ResponseFormat: openai.SpeechResponseFormatWav,
		Voice:          openai.SpeechVoice(opt.Voice),
	}
	resp, err := t.cli.CreateSpeech(ctx, req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(resp.ReadCloser)
	if err != nil {
		return err
	}
	defer resp.ReadCloser.Close()

	return os.WriteFile(file, b, 0644)
}
*/

type TTS struct {
	apiKey  string
	baseURL string
	model   string
}

func NewTTS(apiKey, baseURL, model string) *TTS {
	return &TTS{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
	}
}

func NewTTSFromEnv() *TTS {
	apiKey := os.Getenv("TTS_API_KEY")
	baseURL := os.Getenv("TTS_BASE_URL")
	model := os.Getenv("TTS_MODEL")
	return NewTTS(apiKey, baseURL, model)
}

func (s *TTS) SynthesisFile(ctx context.Context, text string, file string, option ...tts.Option) error {
	b, err := s.Synthesis(ctx, text, option...)
	if err != nil {
		return err
	}
	return os.WriteFile(file, b, 0644)
}
func (s *TTS) Synthesis(ctx context.Context, text string, option ...tts.Option) ([]byte, error) {
	//func (s *TTS) SynthesisFile(ctx context.Context, text string, file string, option ...tts.Option) error {
	opt := tts.Config{
		//Voice:     "FunAudioLLM/CosyVoice2-0.5B:alex",
		Voice:     "FunAudioLLM/CosyVoice2-0.5B:david",
		AudioRate: 44100,
		AudioType: "wav",
	}
	tts.ApplyOption(&opt, option...)

	reqBody := map[string]interface{}{
		"model":           s.model,
		"input":           text,
		"voice":           opt.Voice,
		"sample_rate":     opt.AudioRate,
		"response_format": opt.AudioType, // Default to wav
		"speed":           opt.AudioSpeed,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/audio/speech", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("siliconflow tts failed: %s", string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
