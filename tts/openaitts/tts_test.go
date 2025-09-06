package openaitts

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
)

func TestTTS(t *testing.T) {
	godotenv.Overload()
	tts := NewTTSFromEnv()
	err := tts.SynthesisFile(context.Background(), "Hello, world!", "test.wav")
	if err != nil {
		t.Fatalf("Failed to synthesize speech: %v", err)
	}
}
