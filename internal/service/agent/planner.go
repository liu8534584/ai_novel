package agent

import (
	"context"
	"fmt"

	"ai_novel/internal/model"
	"ai_novel/internal/service/llm/core"
	"ai_novel/pkg/prompt"
)

type PlanAgent struct {
	llmProvider core.Provider
}

func NewPlanAgent(provider core.Provider) *PlanAgent {
	return &PlanAgent{llmProvider: provider}
}

// GeneratePlanVersions 生成多个大纲版本
func (a *PlanAgent) GeneratePlanVersions(ctx context.Context, description string, genre string, worldView string, chapters int, count int) ([]model.OutlineVersion, error) {
	if count <= 0 {
		count = 3
	}
	data := map[string]interface{}{
		"WorldView":     worldView,
		"Description":   description,
		"Genre":         genre,
		"Chapters":      chapters,
		"ChaptersBegin": int(float64(chapters) * 0.2),
	}
	promptNames := []string{"planner_dark", "planner_growth", "planner_twist"}
	versions := make([]model.OutlineVersion, 0, count)

	for i := 0; i < count; i++ {
		messages, options, err := a.buildPlanRequest(data, i, promptNames)
		if err != nil {
			return nil, err
		}

		resp, err := a.llmProvider.Chat(ctx, messages, options)
		if err != nil {
			return nil, fmt.Errorf("LLM call failed: %w", err)
		}

		versions = append(versions, model.OutlineVersion{
			WorldView: worldView,
			Outline:   core.RemoveReasoningContent(resp.Content),
		})
	}

	return versions, nil
}

// GenerateSinglePlan 生成单个版本的大纲
func (a *PlanAgent) GenerateSinglePlan(ctx context.Context, description, genre, worldView string, chapters int) (string, error) {
	versions, err := a.GeneratePlanVersions(ctx, description, genre, worldView, chapters, 1)
	if err != nil {
		return "", err
	}
	return versions[0].Outline, nil
}

func (a *PlanAgent) GeneratePlanVersionStream(ctx context.Context, description string, genre string, worldView string, chapters int, idx int) (<-chan core.StreamResponse, error) {
	data := map[string]interface{}{
		"WorldView":     worldView,
		"Description":   description,
		"Genre":         genre,
		"Chapters":      chapters,
		"ChaptersBegin": int(float64(chapters) * 0.2),
	}
	promptNames := []string{"planner_dark", "planner_growth", "planner_twist"}
	messages, options, err := a.buildPlanRequest(data, idx, promptNames)
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

func (a *PlanAgent) buildPlanRequest(data map[string]interface{}, idx int, promptNames []string) ([]core.Message, core.Options, error) {
	name := "planner"
	if idx < len(promptNames) {
		name = promptNames[idx]
	}
	rendered, err := prompt.GetRegistry().Render(name, data)
	if err != nil {
		return nil, core.Options{}, fmt.Errorf("failed to render planner prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}
	options := core.Options{
		Model: "",
	}
	core.GetStrategy(core.TaskPlanning).ApplyToOptions(&options)
	return messages, options, nil
}
