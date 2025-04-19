package edgetts

const (
	TRUSTED_CLIENT_TOKEN = "6A5AA1D4EAFF4E9FB37E23D68491D6F4"
	WSS_URL              = "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?TrustedClientToken=" + TRUSTED_CLIENT_TOKEN
	VOICE_LIST           = "https://speech.platform.bing.com/consumer/speech/synthesize/readaloud/voices/list?trustedclienttoken=" + TRUSTED_CLIENT_TOKEN
)

// Locale
const (
	ZhCN = "zh-CN"
	EnUS = "en-US"
)

const (
	ChunkTypeAudio        = "Audio"
	ChunkTypeWordBoundary = "WordBoundary"
	ChunkTypeSessionEnd   = "SessionEnd"
	ChunkTypeEnd          = "ChunkEnd"
)

const (
	HiuGaaiNeural   = `{Microsoft Server Speech Text to Speech Voice (zh-HK, HiuGaaiNeural)`
	HiuMaanNeural   = `Microsoft Server Speech Text to Speech Voice (zh-HK, HiuMaanNeural)`
	WanLungNeural   = `Microsoft Server Speech Text to Speech Voice (zh-HK, WanLungNeural)`
	XiaoxiaoNeural  = `Microsoft Server Speech Text to Speech Voice (zh-CN, XiaoxiaoNeural)`
	XiaoyiNeural    = `Microsoft Server Speech Text to Speech Voice (zh-CN, XiaoyiNeural)`
	YunjianNeural   = `Microsoft Server Speech Text to Speech Voice (zh-CN, YunjianNeural)`
	YunxiNeural     = `Microsoft Server Speech Text to Speech Voice (zh-CN, YunxiNeural)`
	YunxiaNeural    = `Microsoft Server Speech Text to Speech Voice (zh-CN, YunxiaNeural)`
	YunyangNeural   = `Microsoft Server Speech Text to Speech Voice (zh-CN, YunyangNeural)`
	XiaobeiNeural   = `Microsoft Server Speech Text to Speech Voice (zh-CN-liaoning, XiaobeiNeural)`
	HsiaoChenNeural = `Microsoft Server Speech Text to Speech Voice (zh-TW, HsiaoChenNeural)`
	YunJheNeural    = `Microsoft Server Speech Text to Speech Voice (zh-TW, YunJheNeural)`
	HsiaoYuNeural   = `Microsoft Server Speech Text to Speech Voice (zh-TW, HsiaoYuNeural)`
	XiaoniNeural    = `Microsoft Server Speech Text to Speech Voice (zh-CN-shaanxi, XiaoniNeural)`
)
