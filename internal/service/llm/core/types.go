package core

import "context"

type ContextKey string

const (
	ContextKeyBookID ContextKey = "book_id"
)

func WithBookID(ctx context.Context, bookID uint) context.Context {
	return context.WithValue(ctx, ContextKeyBookID, bookID)
}

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role      Role       `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Options struct {
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	Timeout     int     `json:"timeout"` // in seconds
	JSONMode    bool    `json:"json_mode"`
	MaxRetries  int     `json:"max_retries"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Response struct {
	Content      string     `json:"content"`
	Role         Role       `json:"role"`
	Usage        Usage      `json:"usage"`
	FinishReason string     `json:"finish_reason"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
}

type StreamResponse struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason"`
	Error        string `json:"error,omitempty"`
}
