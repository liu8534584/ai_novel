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
	messages, options, err := a.buildChapterTitlePlanRequest(outline, chapters)
	if err != nil {
		return "", err
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return core.RemoveReasoningContent(resp.Content), nil
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

	return core.RemoveReasoningContent(resp.Content), nil
}

// GenerateChapterTitlesBatch 分批生成章节标题
func (a *ChapterTitleAgent) GenerateChapterTitlesBatch(ctx context.Context, worldView, titlePlan, characters string, startChapter, batchSize, currentCount int, previousTitles string) (string, error) {
	messages, options, err := a.buildChapterTitlesBatchRequest(worldView, titlePlan, characters, startChapter, batchSize, currentCount, previousTitles)
	if err != nil {
		return "", err
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return core.RemoveReasoningContent(resp.Content), nil
}

func (a *ChapterTitleAgent) GenerateChapterTitlePlanStream(ctx context.Context, outline string, chapters int) (<-chan core.StreamResponse, error) {
	messages, options, err := a.buildChapterTitlePlanRequest(outline, chapters)
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
				outputChan <- core.StreamResponse{Error: r.Error}
				return
			}
			
			// Process content through filter
			filteredContent := thinkFilter.Process(r.Content)
			
			if filteredContent != "" || r.FinishReason != "" {
				newResp := r
				newResp.Content = filteredContent
				outputChan <- newResp
			}
		}
		if rest := thinkFilter.Flush(); rest != "" {
			outputChan <- core.StreamResponse{Content: rest}
		}
	}()

	return outputChan, nil
}

func (a *ChapterTitleAgent) GenerateChapterTitlesBatchStream(ctx context.Context, worldView, titlePlan, characters string, startChapter, batchSize, currentCount int, previousTitles string) (<-chan core.StreamResponse, error) {
	messages, options, err := a.buildChapterTitlesBatchRequest(worldView, titlePlan, characters, startChapter, batchSize, currentCount, previousTitles)
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
				outputChan <- core.StreamResponse{Error: r.Error}
				return
			}
			
			// Process content through filter
			filteredContent := thinkFilter.Process(r.Content)
			
			if filteredContent != "" || r.FinishReason != "" {
				newResp := r
				newResp.Content = filteredContent
				outputChan <- newResp
			}
		}
		if rest := thinkFilter.Flush(); rest != "" {
			outputChan <- core.StreamResponse{Content: rest}
		}
	}()

	return outputChan, nil
}

func (a *ChapterTitleAgent) buildChapterTitlePlanRequest(outline string, chapters int) ([]core.Message, core.Options, error) {
	data := map[string]interface{}{
		"Outline":  outline,
		"Chapters": chapters,
	}
	rendered, err := prompt.GetRegistry().Render("chapter_title_plan", data)
	if err != nil {
		return nil, core.Options{}, fmt.Errorf("failed to render chapter title plan prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.5,
	}
	return messages, options, nil
}

func (a *ChapterTitleAgent) buildChapterTitlesBatchRequest(worldView, titlePlan, characters string, startChapter, batchSize, currentCount int, previousTitles string) ([]core.Message, core.Options, error) {
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
		return nil, core.Options{}, fmt.Errorf("failed to render batch chapter title prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.7,
	}
	return messages, options, nil
}
