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
		name := "planner"
		if i < len(promptNames) {
			name = promptNames[i]
		}
		rendered, err := prompt.GetRegistry().Render(name, data)
		if err != nil {
			return nil, fmt.Errorf("failed to render planner prompt: %w", err)
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
			return nil, fmt.Errorf("LLM call failed: %w", err)
		}

		versions = append(versions, model.OutlineVersion{
			WorldView: worldView,
			Outline:   resp.Content,
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
