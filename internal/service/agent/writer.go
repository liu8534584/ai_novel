package agent

import (
	"context"
	"fmt"

	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/prompt"
)

type WriterAgent struct {
	llmProvider core.Provider
}

func NewWriterAgent(provider core.Provider) *WriterAgent {
	return &WriterAgent{llmProvider: provider}
}

// WriterContext 写作上下文结构
type WriterContext struct {
	WorldSummary      string
	OutlineSummary    string
	CharacterStates   string
	ChapterIndex      int
	ChapterTitle      string
	ChapterObjective  string
	RetrievedMemories string // 长期记忆 (RAG)
	LastChapterTail   string // 短期记忆 (线性接戏)
	Foreshadowing     string // 未回收伏笔 (Open Events)
	TargetWords       int
}

// SplicingAlgorithm 返回拼接后的上下文描述（用于调试或日志）
func (c WriterContext) SplicingAlgorithm() string {
	return fmt.Sprintf("World + Outline + Characters + Recent(Memories+Tail) + OpenEvents(Foreshadowing)")
}

// GenerateChapterObjective 生成本章节的写作目标和策略
func (a *WriterAgent) GenerateChapterObjective(ctx context.Context, outlineSummary string, chapterIndex int, chapterTitle string) (string, error) {
	data := map[string]interface{}{
		"OutlineSummary": outlineSummary,
		"ChapterIndex":   chapterIndex,
		"ChapterTitle":   chapterTitle,
	}
	rendered, err := prompt.GetRegistry().Render("chapter_objective", data)
	if err != nil {
		return "", err
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model: "",
	}
	core.GetStrategy(core.TaskPlanning).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

// WriteChapterStream 使用分层上下文流式生成正文
func (a *WriterAgent) WriteChapterStream(ctx context.Context, wCtx WriterContext) (<-chan string, error) {
	rendered, err := prompt.GetRegistry().Render("writer_layered", wCtx)
	if err != nil {
		return nil, err
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.9,   // 写作需要较高的创意度
		MaxTokens:   4000,  // 增加 MaxTokens 以支持长章节生成
	}
	core.GetStrategy(core.TaskWriting).ApplyToOptions(&options)

	streamResp, err := a.llmProvider.StreamChat(ctx, messages, options)
	if err != nil {
		return nil, err
	}

	outputChan := make(chan string)
	go func() {
		defer close(outputChan)
		for r := range streamResp {
			if r.Content != "" {
				outputChan <- r.Content
			}
		}
	}()

	return outputChan, nil
}

// WriteSceneStream 保留兼容性
func (a *WriterAgent) WriteSceneStream(
	ctx context.Context,
	worldSetting string,
	currentState string,
	sceneDesc string,
	lastPara string,
) (<-chan string, error) {
	// ... (implementation remains same but uses old prompt if needed, or I can just refactor callers)
	return nil, fmt.Errorf("deprecated: use WriteChapterStream instead")
}
