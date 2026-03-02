package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"ai_novel/internal/service/llm/core"
	"ai_novel/models"
	"ai_novel/pkg/prompt"
)

type ConsistencyAgent struct {
	llmProvider core.Provider
}

func NewConsistencyAgent(llmProvider core.Provider) *ConsistencyAgent {
	return &ConsistencyAgent{llmProvider: llmProvider}
}

// ExtractCharacterAnchor 提取角色性格锚点
func (a *ConsistencyAgent) ExtractCharacterAnchor(ctx context.Context, char *models.Character, historyStates string) (*models.CharacterAnchor, error) {
	data := map[string]string{
		"CharacterName":        char.Name,
		"CharacterDescription": char.Description,
		"CharacterStates":      historyStates,
	}
	rendered, err := prompt.GetRegistry().Render("character_anchor_extraction", data)
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
	core.GetStrategy(core.TaskReviewing).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, err
	}

	var anchor models.CharacterAnchor
	if err := json.Unmarshal([]byte(core.ParseJSON(resp.Content)), &anchor); err != nil {
		return nil, fmt.Errorf("failed to parse anchor JSON: %v, resp: %s", err, resp.Content)
	}
	anchor.CharacterID = char.ID

	return &anchor, nil
}

// EvaluateOOC 评估角色 OOC 评分
func (a *ConsistencyAgent) EvaluateOOC(ctx context.Context, anchor *models.CharacterAnchor, historySummary string, currentBehavior string) (*models.OOCScore, error) {
	anchorBytes, _ := json.MarshalIndent(anchor, "", "  ")

	data := map[string]string{
		"CharacterAnchor":           string(anchorBytes),
		"CharacterStateHistory":     historySummary,
		"CurrentCharacterBehavior": currentBehavior,
	}
	rendered, err := prompt.GetRegistry().Render("ooc_evaluation", data)
	if err != nil {
		return nil, err
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model: "",
	}
	core.GetStrategy(core.TaskReviewing).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, err
	}

	var score models.OOCScore
	if err := json.Unmarshal([]byte(core.ParseJSON(resp.Content)), &score); err != nil {
		return nil, fmt.Errorf("failed to parse OOC score JSON: %v, resp: %s", err, resp.Content)
	}
	score.CharacterID = anchor.CharacterID

	return &score, nil
}

// DetectContradictions 检测剧情自相矛盾
func (a *ConsistencyAgent) DetectContradictions(ctx context.Context, worldRules, historyEvents, charStates, currentContent string) ([]models.StoryContradiction, error) {
	data := map[string]string{
		"WorldRules":      worldRules,
		"HistoryEvents":   historyEvents,
		"CharacterStates": charStates,
		"ChapterContent":  currentContent,
	}
	rendered, err := prompt.GetRegistry().Render("contradiction_detection", data)
	if err != nil {
		return nil, err
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:       "",
		Temperature: 0.1,
	}

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, err
	}

	var contradictions []models.StoryContradiction
	if err := json.Unmarshal([]byte(core.ParseJSON(resp.Content)), &contradictions); err != nil {
		return nil, fmt.Errorf("failed to parse contradictions JSON: %v, resp: %s", err, resp.Content)
	}

	return contradictions, nil
}

// EvaluateChapterHealth 综合评估章节健康度
func (a *ConsistencyAgent) EvaluateChapterHealth(ctx context.Context, oocScores []models.OOCScore, contradictions []models.StoryContradiction, openForeshadows []models.Foreshadowing, resolvedForeshadows []models.Foreshadowing) *models.ChapterHealthScore {
	health := &models.ChapterHealthScore{
		EventConsistency: 100.0,
		Foreshadowing:    80.0, // 默认基准分
	}

	// 1. 计算 OOC 平均分 (越低越好，所以用 100 - 平均分)
	if len(oocScores) > 0 {
		var totalOOC float64
		for _, s := range oocScores {
			totalOOC += s.TotalScore
		}
		avgOOC := totalOOC / float64(len(oocScores))
		health.OOCScore = 100.0 - avgOOC
	} else {
		health.OOCScore = 100.0
	}

	// 2. 计算剧情一致性得分 (根据矛盾严重程度扣分)
	for _, c := range contradictions {
		switch c.Severity {
		case "high":
			health.EventConsistency -= 30.0
		case "medium":
			health.EventConsistency -= 15.0
		case "low":
			health.EventConsistency -= 5.0
		}
	}
	if health.EventConsistency < 0 {
		health.EventConsistency = 0
	}

	// 3. 计算伏笔健康度
	// - 回收一个伏笔加分
	// - 存在高重要度未回收伏笔扣分
	health.Foreshadowing += float64(len(resolvedForeshadows) * 10)
	for _, f := range openForeshadows {
		if f.Importance >= 4 {
			health.Foreshadowing -= 5.0
		}
	}
	if health.Foreshadowing > 100 {
		health.Foreshadowing = 100
	}
	if health.Foreshadowing < 0 {
		health.Foreshadowing = 0
	}

	// 4. 综合总分 (加权平均)
	health.TotalHealth = (health.OOCScore * 0.4) + (health.EventConsistency * 0.4) + (health.Foreshadowing * 0.2)

	// 5. 生成简要审计报告
	report := "章节审计报告:\n"
	if health.TotalHealth >= 90 {
		report += "- 整体质量优秀，一致性保持良好。\n"
	} else if health.TotalHealth >= 70 {
		report += "- 整体质量良好，存在少量细节偏差。\n"
	} else {
		report += "- 章节存在明显逻辑或角色偏离，建议修正。\n"
	}

	if len(contradictions) > 0 {
		report += fmt.Sprintf("- 检测到 %d 处剧情矛盾。\n", len(contradictions))
	}
	if len(resolvedForeshadows) > 0 {
		report += fmt.Sprintf("- 成功回收 %d 个伏笔。\n", len(resolvedForeshadows))
	}

	health.AuditReport = report

	return health
}
