package agent

import (
	"context"
	"testing"

	"ai_novel/internal/config"
	"ai_novel/internal/service/llm/core"
)

type MockProvider struct {
	response string
}

func (m *MockProvider) Chat(ctx context.Context, messages []core.Message, options core.Options) (core.Response, error) {
	return core.Response{
		Content: m.response,
		Role:    core.RoleAssistant,
	}, nil
}

func (m *MockProvider) StreamChat(ctx context.Context, messages []core.Message, options core.Options) (<-chan core.StreamResponse, error) {
	ch := make(chan core.StreamResponse)
	close(ch)
	return ch, nil
}

func (m *MockProvider) CreateEmbedding(ctx context.Context, input string, options core.Options) ([]float32, error) {
	return []float32{0.1, 0.2, 0.3}, nil
}

func TestDirectorAgent_InitWorld(t *testing.T) {
	config.GlobalConfig.LLM.Model = "mock"
	provider := &MockProvider{
		response: `{
			"book_title": "Test World",
			"background": "A dark testing realm",
			"power_system": {
				"tiers": ["F", "E"],
				"currency": "Bits",
				"upgrade_path": "Coding"
			},
			"protagonist": {
				"name": "Neo",
				"personality": "Stoic",
				"stats": {"HP": 100, "MP": 50},
				"inventory": ["Laptop"],
				"skills": ["Debug"],
				"cheat_ability": "Infinite Coffee"
			},
			"global_rules": ["No bugs"]
		}`,
	}

	agent := NewDirectorAgent(provider)
	world, err := agent.InitWorld(context.Background(), "test idea", "fantasy", 100)
	if err != nil {
		t.Fatalf("InitWorld failed: %v", err)
	}
	if world.Content == "" {
		t.Fatalf("unexpected empty content")
	}
}
