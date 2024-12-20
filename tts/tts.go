package tts

import (
	"context"
)

type TTSConfig struct {
	VoiceType    string  `json:"voice_type"`    // 语音类型
	AudioRate    int     `json:"audio_rate"`    // 声音频率
	AudioType    string  `json:"audio_type"`    // 音频类型
	AudioSpeed   float64 `json:"audio_speed"`   // 语速
	AudioEmotion string  `json:"audio_emotion"` // 情绪
}

type TTSOption func(*TTSConfig)

func WithVoiceType(voiceType string) TTSOption {
	return func(cfg *TTSConfig) {
		cfg.VoiceType = voiceType
	}
}
func WithAudioRate(audioRate int) TTSOption {
	return func(cfg *TTSConfig) {
		cfg.AudioRate = audioRate
	}
}
func WithAudioType(audioType string) TTSOption {
	return func(cfg *TTSConfig) {
		cfg.AudioType = audioType
	}
}
func WithAudioSpeed(audioSpeed float64) TTSOption {
	return func(cfg *TTSConfig) {
		cfg.AudioSpeed = audioSpeed
	}
}
func WithAudioEmotion(audioEmotion string) TTSOption {
	return func(cfg *TTSConfig) {
		cfg.AudioEmotion = audioEmotion
	}
}

type TTS interface {
	Synthesis(ctx context.Context, text string, option ...TTSOption) ([]byte, error)
}
