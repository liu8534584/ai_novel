package llm

import (
	"ai_novel/internal/service/llm/core"
	"ai_novel/internal/service/llm/provider/openai"
	"ai_novel/internal/service/llm/provider/deepseek"
	"context"
	"errors"
)

// ErrorProvider returns a fixed error for all calls.
type ErrorProvider struct {
	Err error
}

func (p *ErrorProvider) Chat(ctx context.Context, messages []core.Message, options core.Options) (core.Response, error) {
	return core.Response{}, p.Err
}

func (p *ErrorProvider) StreamChat(ctx context.Context, messages []core.Message, options core.Options) (<-chan core.StreamResponse, error) {
	return nil, p.Err
}

func (p *ErrorProvider) CreateEmbedding(ctx context.Context, input string, options core.Options) ([]float32, error) {
	return nil, p.Err
}

// NewProvider creates a specific provider based on the given configuration.
func NewProvider(providerName, apiKey, baseURL string) core.Provider {
	if apiKey == "" {
		return &ErrorProvider{Err: errors.New("LLM API Key not configured")}
	}

	switch providerName {
	case "openai":
		return openai.NewOpenAIProvider(apiKey, baseURL)
	case "deepseek":
		return deepseek.NewDeepSeekProvider(apiKey, baseURL)
	case "glm", "zhipu", "aliyun", "siliconflow":
		// These providers are often OpenAI compatible
		return openai.NewOpenAIProvider(apiKey, baseURL)
	default:
		// If unknown but we have key/baseURL, try OpenAI compatible
		if baseURL != "" {
			return openai.NewOpenAIProvider(apiKey, baseURL)
		}
		return &ErrorProvider{Err: errors.New("Unsupported LLM provider: " + providerName)}
	}
}
