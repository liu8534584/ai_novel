package agent

import (
	"context"

	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/prompt"
)

type SummarizerAgent struct {
	llmProvider core.Provider
}

func NewSummarizerAgent(provider core.Provider) *SummarizerAgent {
	return &SummarizerAgent{llmProvider: provider}
}

// SummarizeWorld 生成世界观摘要
func (a *SummarizerAgent) SummarizeWorld(ctx context.Context, worldView string) (string, error) {
	return a.summarize(ctx, "世界观设定", worldView)
}

// SummarizeOutline 生成大纲摘要
func (a *SummarizerAgent) SummarizeOutline(ctx context.Context, outline string) (string, error) {
	return a.summarize(ctx, "剧情大纲", outline)
}

// SummarizeChapter 生成章节摘要，用于 RAG 存储
func (a *SummarizerAgent) SummarizeChapter(ctx context.Context, title, content string) (string, error) {
	promptText := "请总结以下章节的核心内容，包括发生的关键事件、角色状态变化、重要的伏笔引入或回收。摘要应简洁明了，适合作为后续创作的参考记忆。\n\n章节标题：" + title + "\n章节正文：\n" + content
	
	messages := []core.Message{
		{Role: core.RoleUser, Content: promptText},
	}

	options := core.Options{
		Temperature: 0.3,
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", err
	}

	return core.RemoveReasoningContent(resp.Content), nil
}

func (a *SummarizerAgent) summarize(ctx context.Context, contentType string, content string) (string, error) {
	data := map[string]string{
		"Type":    contentType,
		"Content": content,
	}
	rendered, err := prompt.GetRegistry().Render("summary", data)
	if err != nil {
		return "", err
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.3,
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", err
	}

	return core.RemoveReasoningContent(resp.Content), nil
}
