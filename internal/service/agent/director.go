package agent

import (
	"context"
	"fmt"

	"ai_novel/internal/model"
	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/prompt"
)

type DirectorAgent struct {
	llmProvider core.Provider
}

func NewDirectorAgent(provider core.Provider) *DirectorAgent {
	return &DirectorAgent{llmProvider: provider}
}

// InitWorld 根据灵感生成世界观文档
func (a *DirectorAgent) InitWorld(ctx context.Context, description, genre string, chapters int) (*model.WorldConfig, error) {
	// 1. 使用 Registry 渲染动态 Prompt
	data := map[string]interface{}{
		"Description": description,
		"Genre":       genre,
		"Chapters":    chapters,
	}
	rendered, err := prompt.GetRegistry().Render("director", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render director prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:    "",
		JSONMode: false,
	}
	core.GetStrategy(core.TaskWorldBuilding).ApplyToOptions(&options)

	// 2. 调用 LLM
	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// 3. 返回结构化文本内容
	return &model.WorldConfig{
		Content: resp.Content,
	}, nil
}

// ChatForInspiration 与 LLM 进行灵感头脑风暴
func (a *DirectorAgent) ChatForInspiration(ctx context.Context, history []core.Message) (string, error) {
	systemMsg, err := prompt.GetRegistry().Render("inspiration_chat", nil)
	if err != nil {
		return "", fmt.Errorf("failed to render inspiration_chat prompt: %w", err)
	}

	messages := append([]core.Message{
		{Role: core.RoleSystem, Content: systemMsg},
	}, history...)

	options := core.Options{}
	core.GetStrategy(core.TaskWorldBuilding).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Content, nil
}

// FinalizeInspiration 将对话记录加工成结构化的小说方案
func (a *DirectorAgent) FinalizeInspiration(ctx context.Context, conversation string) (string, error) {
	data := map[string]interface{}{
		"Conversation": conversation,
	}
	rendered, err := prompt.GetRegistry().Render("inspiration", data)
	if err != nil {
		return "", fmt.Errorf("failed to render inspiration prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		JSONMode: true,
	}
	core.GetStrategy(core.TaskWorldBuilding).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Content, nil
}
