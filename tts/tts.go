package tts

import (
	"context"
)

type Config struct {
	Voice        string  `json:"voice_type"`    // 语音类型
	AudioRate    int     `json:"audio_rate"`    // 声音频率
	AudioType    string  `json:"audio_type"`    // 音频类型 mp3, wav
	AudioSpeed   float64 `json:"audio_speed"`   // 语速
	AudioEmotion string  `json:"audio_emotion"` // 情绪
	Volume       int     `json:"volume"`        // 音量
}

type Option func(*Config)

func WithVoiceType(voiceType string) Option {
	return func(cfg *Config) {
		cfg.Voice = voiceType
	}
}
func WithAudioRate(audioRate int) Option {
	return func(cfg *Config) {
		cfg.AudioRate = audioRate
	}
}
func WithAudioType(audioType string) Option {
	return func(cfg *Config) {
		cfg.AudioType = audioType
	}
}
func WithAudioSpeed(audioSpeed float64) Option {
	return func(cfg *Config) {
		cfg.AudioSpeed = audioSpeed
	}
}
func WithAudioEmotion(audioEmotion string) Option {
	return func(cfg *Config) {
		cfg.AudioEmotion = audioEmotion
	}
}
func WithVolume(volume int) Option {
	return func(cfg *Config) {
		cfg.Volume = volume
	}
}

func (cfg *Config) Apply(o ...Option) {
	for _, opt := range o {
		opt(cfg)
	}
}

type TTS interface {
	Synthesis(ctx context.Context, text string, option ...Option) ([]byte, error)
}
