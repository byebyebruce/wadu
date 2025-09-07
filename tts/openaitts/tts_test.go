package openaitts

import (
	"context"
	"testing"

	atts "github.com/byebyebruce/wadu/tts"
	"github.com/joho/godotenv"
)

func TestTTS(t *testing.T) {
	godotenv.Overload()
	tts := NewTTSFromEnv()
	err := tts.SynthesisFile(context.Background(), "Hello, world!", "test.wav", atts.WithAudioSpeed(0.5))
	if err != nil {
		t.Fatalf("Failed to synthesize speech: %v", err)
	}
}
