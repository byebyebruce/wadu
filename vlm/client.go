package vlm

import (
	env "github.com/caarlos0/env/v9"
	"github.com/sashabaranov/go-openai"
)

type Config struct {
	OpenAIAPIKey  string `env:"OPENAI_API_KEY"`
	OpenAIModel   string `env:"OPENAI_MODEL"`
	OpenAIBaseURL string `env:"OPENAI_BASE_URL"`
}

type Client struct {
	*openai.Client
	Model string
}

func NewClientFromEnv() (*Client, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return NewClient(cfg)
}

func NewClient(cfg Config) (*Client, error) {
	c := openai.DefaultConfig(cfg.OpenAIAPIKey)
	if len(cfg.OpenAIBaseURL) > 0 {
		c.BaseURL = cfg.OpenAIBaseURL
	}
	cli := openai.NewClientWithConfig(c)
	return &Client{
		Client: cli,
		Model:  cfg.OpenAIModel,
	}, nil
}
