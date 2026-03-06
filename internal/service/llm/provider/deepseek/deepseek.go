package deepseek

import (
	"context"
	"errors"
	"io"

	"ai_novel/internal/service/llm/core"

	goopenai "github.com/sashabaranov/go-openai"
)

const (
	DeepSeekBaseURL = "https://api.deepseek.com"
)

type DeepSeekProvider struct {
	client *goopenai.Client
}

func NewDeepSeekProvider(apiKey string, baseURL string) *DeepSeekProvider {
	if baseURL == "" {
		baseURL = DeepSeekBaseURL
	}
	config := goopenai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := goopenai.NewClientWithConfig(config)
	return &DeepSeekProvider{client: client}
}

func (p *DeepSeekProvider) Chat(ctx context.Context, messages []core.Message, options core.Options) (core.Response, error) {
	reqMessages := make([]goopenai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		reqMessages[i] = goopenai.ChatCompletionMessage{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}

	req := goopenai.ChatCompletionRequest{
		Model:       options.Model,
		Messages:    reqMessages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
	}

	if options.JSONMode {
		// LM Studio / Local models often have issues with "json_object"
		// If the error is: 'response_format.type' must be 'json_schema' or 'text'
		// It means the provider doesn't support the standard 'json_object' yet.
		// For maximum compatibility with local providers, we fallback to 'text' 
		// and rely on our ParseJSON utility to extract the JSON.
		
		/*
		req.ResponseFormat = &goopenai.ChatCompletionResponseFormat{
			Type: goopenai.ChatCompletionResponseFormatTypeJSONObject,
		}
		*/
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return core.Response{}, err
	}

	if len(resp.Choices) == 0 {
		return core.Response{}, errors.New("no choices in response")
	}

	choice := resp.Choices[0]

	// Remove <think>...</think> content
	content := core.RemoveReasoningContent(choice.Message.Content)

	// Convert tool calls
	var toolCalls []core.ToolCall
	for _, tc := range choice.Message.ToolCalls {
		toolCalls = append(toolCalls, core.ToolCall{
			ID:   tc.ID,
			Type: string(tc.Type),
			Function: struct {
				Name      string `json:"name"`
				Arguments string `json:"arguments"`
			}{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		})
	}

	return core.Response{
		Content: content,
		Role:    core.Role(choice.Message.Role),
		Usage: core.Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		FinishReason: string(choice.FinishReason),
		ToolCalls:    toolCalls,
	}, nil
}

func (p *DeepSeekProvider) StreamChat(ctx context.Context, messages []core.Message, options core.Options) (<-chan core.StreamResponse, error) {
	reqMessages := make([]goopenai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		reqMessages[i] = goopenai.ChatCompletionMessage{
			Role:    string(m.Role),
			Content: m.Content,
		}
	}

	req := goopenai.ChatCompletionRequest{
		Model:       options.Model,
		Messages:    reqMessages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Stream:      true,
	}
	if options.JSONMode {
		// LM Studio / Local models often have issues with "json_object"
		// See comment in Chat() above.
		/*
		req.ResponseFormat = &goopenai.ChatCompletionResponseFormat{
			Type: goopenai.ChatCompletionResponseFormatTypeJSONObject,
		}
		*/
	}

	stream, err := p.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}

	outputChan := make(chan core.StreamResponse)

	go func() {
		defer close(outputChan)
		defer stream.Close()
		thinkFilter := core.NewThinkTagFilter()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				if rest := thinkFilter.Flush(); rest != "" {
					select {
					case outputChan <- core.StreamResponse{Content: rest}:
					case <-ctx.Done():
					}
				}
				return
			}
			if err != nil {
				select {
				case outputChan <- core.StreamResponse{Error: err.Error()}:
				case <-ctx.Done():
				}
				return
			}

			if len(response.Choices) > 0 {
				choice := response.Choices[0]
				content := thinkFilter.Process(choice.Delta.Content)
				if content != "" || choice.FinishReason != "" {
					select {
					case outputChan <- core.StreamResponse{
						Content:      content,
						FinishReason: string(choice.FinishReason),
					}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return outputChan, nil
}

func (p *DeepSeekProvider) CreateEmbedding(ctx context.Context, input string, options core.Options) ([]float32, error) {
	// DeepSeek might not support embeddings in the same way or at all via OpenAI-compatible endpoint.
	// We'll provide a placeholder or return error if not supported.
	return nil, errors.New("deepseek does not support embeddings via this adapter yet")
}
