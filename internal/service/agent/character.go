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
	messages, options, err := a.buildCharacterRequest(worldView, outline)
	if err != nil {
		return nil, err
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// Filter think tags
	thinkFilter := core.NewThinkTagFilter()
	cleanContent := thinkFilter.Process(resp.Content)
	cleanContent += thinkFilter.Flush()
	
	content := core.ParseJSON(cleanContent)
	if content == "" {
		return nil, fmt.Errorf("failed to extract valid JSON from LLM output. Output was: %s", resp.Content)
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

func (a *CharacterAgent) GenerateCharactersStream(ctx context.Context, worldView, outline string) (<-chan core.StreamResponse, error) {
	messages, options, err := a.buildCharacterRequest(worldView, outline)
	if err != nil {
		return nil, err
	}
	
	streamResp, err := a.llmProvider.StreamChat(ctx, messages, options)
	if err != nil {
		return nil, err
	}

	// Add filter logic
	outputChan := make(chan core.StreamResponse)
	thinkFilter := core.NewThinkTagFilter()

	go func() {
		defer close(outputChan)
		for r := range streamResp {
			if r.Error != "" {
				select {
				case outputChan <- core.StreamResponse{Error: r.Error}:
				case <-ctx.Done():
				}
				return
			}

			// Process content through filter
			filteredContent := thinkFilter.Process(r.Content)

			if filteredContent != "" || r.FinishReason != "" {
				newResp := r
				newResp.Content = filteredContent
				select {
				case outputChan <- newResp:
				case <-ctx.Done():
					return
				}
			}
		}
		if rest := thinkFilter.Flush(); rest != "" {
			select {
			case outputChan <- core.StreamResponse{Content: rest}:
			case <-ctx.Done():
			}
		}
	}()

	return outputChan, nil
}

func (a *CharacterAgent) buildCharacterRequest(worldView, outline string) ([]core.Message, core.Options, error) {
	data := map[string]interface{}{
		"WorldView": worldView,
		"Outline":   outline,
	}
	rendered, err := prompt.GetRegistry().Render("character", data)
	if err != nil {
		return nil, core.Options{}, fmt.Errorf("failed to render character prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		JSONMode:    true,
		MaxTokens:   8192, // Increased from 4000 to handle potential long thinking process + JSON
		Temperature: 0.7,
	}
	core.GetStrategy(core.TaskCharacterDesign).ApplyToOptions(&options)
	return messages, options, nil
}
