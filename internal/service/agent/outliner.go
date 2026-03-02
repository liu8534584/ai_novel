package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_novel/internal/model"
	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/prompt"
)

type OutlinerAgent struct {
	llmProvider core.Provider
}

func NewOutlinerAgent(provider core.Provider) *OutlinerAgent {
	return &OutlinerAgent{llmProvider: provider}
}

// GenerateOutline 生成章节大纲
func (a *OutlinerAgent) GenerateOutline(
	ctx context.Context,
	worldConfig string,
	characters string,
	storyOutline string,
	chapterTitle string,
	allTitles string,
	currentState string,
	prevSummary string,
	chapterNum int,
	userIntent string,
) (*model.ChapterOutline, error) {

	// 1. 使用 Registry 渲染动态 Prompt
	data := map[string]interface{}{
		"WorldSetting": worldConfig,
		"Characters":   characters,
		"StoryOutline": storyOutline,
		"ChapterTitle": chapterTitle,
		"AllTitles":    allTitles,
		"CurrentState": currentState,
		"PrevSummary":  prevSummary,
		"ChapterNum":   chapterNum,
		"UserIntent":   userIntent,
	}
	rendered, err := prompt.GetRegistry().Render("outliner", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render outliner prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.7,
	}

	// 2. 调用 LLM
	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// 3. 解析结果
	var outline model.ChapterOutline
	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &outline); err != nil {
		return nil, fmt.Errorf("failed to parse outline JSON: %w, raw: %s", err, resp.Content)
	}

	return &outline, nil
}
