package volcanotts

import (
	"context"
	"testing"

	tts_ "github.com/byebyebruce/wadu/tts"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	godotenv.Overload()
	tts, err := NewTTSWithEnv()
	require.Nil(t, err)

	err = tts.SynthesisSaveFile(context.TODO(),
		"22种情感/风格】通用、愉悦、抱歉、嗔怪、开心、愤怒、惊讶、厌恶、悲伤、害怕、哭腔、客服、专业、严肃、傲娇、安慰鼓励、绿茶、娇媚、情感电台、撒娇、瑜伽、讲故事",
		"test1.wav",
		tts_.WithAudioEmotion("sad"))
	require.Nil(t, err)
}
