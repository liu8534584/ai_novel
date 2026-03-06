package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_novel/internal/service/llm/core"
	"ai_novel/models"
	"ai_novel/pkg/prompt"
)

type StateAgent struct {
	llmProvider core.Provider
}

func NewStateAgent(provider core.Provider) *StateAgent {
	return &StateAgent{llmProvider: provider}
}

// ExtractDynamicStateChanges 提取角色动态状态变化
func (a *StateAgent) ExtractDynamicStateChanges(ctx context.Context, worldSummary, baseProfiles, previousStates, chapterContent string) (map[string]models.CharacterDynamicState, error) {
	data := map[string]string{
		"WorldViewSummary":         worldSummary,
		"CharacterBaseProfiles":    baseProfiles,
		"PreviousCharacterStates":  previousStates,
		"ChapterContent":           chapterContent,
	}
	rendered, err := prompt.GetRegistry().Render("character_dynamic_state", data)
	if err != nil {
		return nil, err
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:    "",
		JSONMode: true,
	}
	core.GetStrategy(core.TaskStateTracking).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, err
	}

	var result map[string]models.CharacterDynamicState
	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return nil, fmt.Errorf("failed to parse character dynamic state: %w, raw: %s", err, resp.Content)
	}

	return result, nil
}

/*
// AnalyzeAndSyncState 分析文本并同步数据库状态 (Deprecated: Use ExtractDynamicStateChanges instead)
func (a *StateAgent) AnalyzeAndSyncState(ctx context.Context, currentJSON string, chapterContent string) (string, error) {
	// 1. 构造 Prompt
	data := map[string]string{
		"CurrentState":   currentJSON,
		"ChapterContent": chapterContent,
	}
	rendered, err := prompt.GetRegistry().Render("state_audit", data)
	if err != nil {
		return "", fmt.Errorf("failed to render state prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:    "",
		JSONMode: true,
	}
	core.GetStrategy(core.TaskStateTracking).ApplyToOptions(&options)

	// 2. 调用 LLM 提取变动
	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	// 3. 解析变动 JSON
	var update model.StateUpdate
	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &update); err != nil {
		return "", fmt.Errorf("state parse error: %v, raw: %s", err, resp.Content)
	}

	// 4. 合并状态 (Merge Logic)
	return a.mergeState(currentJSON, update)
}

// mergeState 将变动应用到原始 JSON 字符串中，返回新的完整 JSON
func (a *StateAgent) mergeState(oldJSON string, update model.StateUpdate) (string, error) {
	var character model.Character
	// This logic is likely broken or incomplete, so commenting out.
	return oldJSON, nil
}

	if err := json.Unmarshal([]byte(oldJSON), &character); err != nil {
		return "", fmt.Errorf("failed to unmarshal current state: %w", err)
	}

	// 2. Apply Stats Delta
	if character.Stats == nil {
		character.Stats = make(map[string]int)
	}
	for k, v := range update.StatsDelta {
		character.Stats[k] += v
	}

	// 3. Apply New Items
	if len(update.NewItems) > 0 {
		character.Inventory = append(character.Inventory, update.NewItems...)
	}

	// 4. Apply Removed Items
	if len(update.RemovedItems) > 0 {
		for _, itemToRemove := range update.RemovedItems {
			for i, item := range character.Inventory {
				if item == itemToRemove {
					character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
					break
				}
			}
		}
	}

	// 5. Apply Skill Upgrades
	existingSkills := make(map[string]bool)
	for _, s := range character.Skills {
		existingSkills[s] = true
	}
	for _, s := range update.SkillUpgrades {
		if !existingSkills[s] {
			character.Skills = append(character.Skills, s)
			existingSkills[s] = true
		}
	}

	// 6. Serialize back to JSON
	newJSONBytes, err := json.Marshal(character)
	if err != nil {
		return "", fmt.Errorf("failed to marshal new state: %w", err)
	}

	return string(newJSONBytes), nil
}
*/
