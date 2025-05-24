package config

type AnthropicConfig struct {
	APIKey    string `envconfig:"ANTHROPIC_API_KEY" required:"true"`
	Model     string `envconfig:"ANTHROPIC_MODEL" default:"claude-3-sonnet-20240229"`
	MaxTokens int64  `envconfig:"ANTHROPIC_MAX_TOKENS" default:"1024"`
}
