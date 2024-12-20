package volcanotts

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/byebyebruce/wadu/tts"

	"github.com/caarlos0/env/v9"
	"github.com/google/uuid"
)

// doc: https://www.volcengine.com/docs/6561/79820

// TTS_API_URL is the url for tts
var (
	TTS_API_URL = `https://openspeech.bytedance.com/api/v1/tts`
)

/*
https://www.volcengine.com/docs/6561/97465
	免费 声音类型

通用场景	3	BV700灿灿",
"BV001通用女声",
"BV002通用男声
有声阅读	5	BV701擎苍",
"BV119通用赘婿",
"BV102儒雅青年",
"BV113甜宠少御",
"BV115古风少御
智能助手/视频配音/特色/教育	6	BV007亲切女声",
"BV056阳光男声",
"BV005活泼女声",
"BV051奶气萌娃",
"BV034知性姐姐-双语",
"BV033温柔小哥
方言	3	BV021东北老铁",
" BV019重庆小伙 ",
"BV213广西表哥
英语	2	BV503活力女声-Ariana",
"BV504活力男声-Jackson
日语	2	BV522气质女生",
"BV524日语男声
*/

/*
wav / pcm / ogg_opus / mp3，默认为 pcm
注意：wav 不支持流式
*/

type EmotionType string

var EmotionTypes = []EmotionType{
	"pleased",
	"sorry",
	"annoyed",
	"customer_service",
	"professional",
	"serious",
	"happy",
	"sad",
	"angry",
	"scare",
	"hate",
	"surprise",
	"tear",
	"conniving",
	"comfort",
	"radio",
	"lovey-dovey",
	"tsundere",
	"charming",
	"yoga",
	"storytelling",
}

type Config struct {
	AppID     string `json:"app_id" env:"VOLCANO_TTS_APP_ID"`
	AppToken  string `json:"app_token" env:"VOLCANO_TTS_APP_TOKEN"`
	ClusterID string `json:"cluster_id" env:"VOLCANO_TTS_CLUSTER_ID"`
	VoiceType string `json:"voice_type" env:"VOLCANO_TTS_VOICE_TYPE"` // default BV213_streaming
}

var _ tts.TTS = (*TTS)(nil)

type TTS struct {
	cfg Config
}

func NewTTS(cfg Config) (*TTS, error) {
	if cfg.ClusterID == "" {
		cfg.ClusterID = "volcano_tts"
	}
	return &TTS{
		cfg: cfg,
	}, nil
}

func NewTTSWithEnv() (*TTS, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return NewTTS(cfg)
}

func (h *TTS) applyTTSConfig(opts ...tts.TTSOption) *tts.TTSConfig {
	cfg := &tts.TTSConfig{
		VoiceType:    "BV213_streaming",
		AudioRate:    16000,
		AudioType:    "wav", // default wav
		AudioSpeed:   1.0,   // 语速 0.3~3
		AudioEmotion: "",
	}

	if len(h.cfg.VoiceType) > 0 {
		cfg.VoiceType = h.cfg.VoiceType
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func (h *TTS) Synthesis(ctx context.Context, text string, option ...tts.TTSOption) ([]byte, error) {
	cfg := h.applyTTSConfig(option...)
	return h.SynthesisWithConfig(ctx, text, cfg)
}

func (h *TTS) SynthesisSaveFile(ctx context.Context, text string, out string, option ...tts.TTSOption) error {
	audio, err := h.Synthesis(ctx, text, option...)
	if err != nil {
		return err
	}
	return os.WriteFile(out, audio, 0644)
}

func (h *TTS) SynthesisWithConfig(
	ctx context.Context,
	text string,
	cfg *tts.TTSConfig,
) ([]byte, error) {

	reqID := uuid.NewString()
	params := make(map[string]map[string]interface{})
	params["app"] = make(map[string]interface{})
	//填写平台申请的appid
	params["app"]["appid"] = h.cfg.AppID
	//这部分的token不生效，填写下方的默认值就好
	//params["app"]["token"] = appToken
	//填写平台上显示的集群名称
	params["app"]["cluster"] = h.cfg.ClusterID
	params["user"] = make(map[string]interface{})
	//这部分如有需要，可以传递用户真实的ID，方便问题定位
	params["user"]["uid"] = "uid"
	params["audio"] = make(map[string]interface{})
	//填写选中的音色代号
	params["audio"]["voice_type"] = cfg.VoiceType
	params["audio"]["encoding"] = cfg.AudioType
	params["audio"]["speed_ratio"] = cfg.AudioSpeed
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["audio"]["rate"] = cfg.AudioRate
	if len(cfg.AudioEmotion) > 0 {
		params["audio"]["emotion"] = cfg.AudioEmotion
	}
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = reqID
	params["request"]["text"] = text
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = "query"

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	//bearerToken为saas平台对应的接入认证中的Token
	headers["Authorization"] = fmt.Sprintf("Bearer;%s", h.cfg.AppToken)

	// URL查看上方第四点: 4.并发合成接口(POST)
	url := TTS_API_URL
	bodyStr, _ := json.Marshal(params)
	synResp, err := httpPost(ctx, url, headers, bodyStr)
	if err != nil {
		return nil, err
	}
	var respJSON TTSServResponse
	err = json.Unmarshal(synResp, &respJSON)
	if err != nil {
		return nil, err
	}
	code := respJSON.Code
	if code != 3000 {
		return nil, fmt.Errorf("resp code fail: %d,message:%s", code, respJSON.Message)
	}

	audio, err := base64.StdEncoding.DecodeString(respJSON.Data)
	if err != nil {
		return nil, fmt.Errorf("base64 decode fail: %w", err)
	}
	return audio, nil
}

func httpPost(ctx context.Context, url string, headers map[string]string, body []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	retBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return retBody, err
}

// TTSServResponse response from backend srvs
type TTSServResponse struct {
	ReqID     string `json:"reqid"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
	Sequence  int    `json:"sequence"`
	Data      string `json:"data"`
}
