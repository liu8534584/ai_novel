package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_novel/internal/model"
	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/prompt"
)

type CharacterAgent struct {
	llmProvider core.Provider
}

func NewCharacterAgent(provider core.Provider) *CharacterAgent {
	return &CharacterAgent{llmProvider: provider}
}

// GenerateCharacters 生成主要角色设定
func (a *CharacterAgent) GenerateCharacters(ctx context.Context, worldView, outline string) ([]model.Character, error) {
	data := map[string]interface{}{
		"WorldView": worldView,
		"Outline":   outline,
	}
	rendered, err := prompt.GetRegistry().Render("character", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render character prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		JSONMode:    true,
		MaxTokens:   4000, // Ensure enough tokens for multiple characters with descriptions
		Temperature: 0.7,
	}
	core.GetStrategy(core.TaskCharacterDesign).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	content := resp.Content
	if content == "" {
		return nil, fmt.Errorf("LLM returned empty content. Check your model configuration and prompt.")
	}

	var characters []model.Character
	// 1. Try unmarshaling as object with "characters" key (Preferred)
	var wrapper struct {
		Characters []model.Character `json:"characters"`
	}
	if err := json.Unmarshal([]byte(content), &wrapper); err == nil && len(wrapper.Characters) > 0 {
		characters = wrapper.Characters
	} else {
		// 2. Fallback: try unmarshaling as array
		if err2 := json.Unmarshal([]byte(content), &characters); err2 != nil {
			// If both fail, return detailed error
			return nil, fmt.Errorf("failed to unmarshal characters. Content might be truncated or invalid JSON. Error 1: %v, Error 2: %v, Content: %s", err, err2, content)
		}
	}

	if len(characters) == 0 {
		return nil, fmt.Errorf("no characters generated. content: %s", content)
	}

	// Validation: Check for empty descriptions
	for i, char := range characters {
		if char.Description == "" {
			fmt.Printf("Warning: Character %s (index %d) has empty description.\n", char.Name, i)
			// Optional: Set a default description if needed, or leave it for frontend to handle
			if char.Role != "" {
				characters[i].Description = fmt.Sprintf("Role: %s. (Description missing from generation)", char.Role)
			}
		}
	}

	return characters, nil
}
