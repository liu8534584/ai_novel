package agent

import (
	"context"
	"fmt"

	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/prompt"
)

type ChapterTitleAgent struct {
	llmProvider core.Provider
}

func NewChapterTitleAgent(provider core.Provider) *ChapterTitleAgent {
	return &ChapterTitleAgent{llmProvider: provider}
}

func (a *ChapterTitleAgent) GenerateChapterTitlePlan(ctx context.Context, outline string, chapters int) (string, error) {
	data := map[string]interface{}{
		"Outline":  outline,
		"Chapters": chapters,
	}
	rendered, err := prompt.GetRegistry().Render("chapter_title_plan", data)
	if err != nil {
		return "", fmt.Errorf("failed to render chapter title plan prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.5,
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Content, nil
}

// GenerateChapterTitles 生成整本小说的章节标题
func (a *ChapterTitleAgent) GenerateChapterTitles(ctx context.Context, worldView, outline, characters string, chapters int) (string, error) {
	data := map[string]interface{}{
		"WorldView":  worldView,
		"Outline":    outline,
		"Characters": characters,
		"Chapters":   chapters,
	}
	rendered, err := prompt.GetRegistry().Render("chapter_title", data)
	if err != nil {
		return "", fmt.Errorf("failed to render chapter title prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.7,
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Content, nil
}

// GenerateChapterTitlesBatch 分批生成章节标题
func (a *ChapterTitleAgent) GenerateChapterTitlesBatch(ctx context.Context, worldView, titlePlan, characters string, startChapter, batchSize, currentCount int, previousTitles string) (string, error) {
	data := map[string]interface{}{
		"WorldView":      worldView,
		"TitlePlan":      titlePlan,
		"Characters":     characters,
		"StartChapter":   startChapter,
		"BatchSize":      batchSize,
		"EndChapter":     startChapter + batchSize - 1,
		"CurrentCount":   currentCount,
		"PreviousTitles": previousTitles,
	}
	rendered, err := prompt.GetRegistry().Render("chapter_title_batch_plan", data)
	if err != nil {
		return "", fmt.Errorf("failed to render batch chapter title prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.7,
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Content, nil
}
