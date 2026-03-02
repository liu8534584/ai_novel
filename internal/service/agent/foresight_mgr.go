package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ai_novel/internal/service/llm/core"
	"ai_novel/models"
	"ai_novel/pkg/prompt"

	"gorm.io/gorm"
)

type ForesightAgent struct {
	db          *gorm.DB
	llmProvider core.Provider
}

func NewForesightAgent(db *gorm.DB, provider core.Provider) *ForesightAgent {
	return &ForesightAgent{
		db:          db,
		llmProvider: provider,
	}
}

// ExtractEvents 从章节内容中提取关键事件
func (a *ForesightAgent) ExtractEvents(ctx context.Context, chapterContent string) ([]models.StoryEvent, error) {
	data := map[string]string{
		"ChapterContent": chapterContent,
	}
	rendered, err := prompt.GetRegistry().Render("event_extraction", data)
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

	cleanJSON := core.ParseJSON(resp.Content)
	var events []models.StoryEvent
	if err := json.Unmarshal([]byte(cleanJSON), &events); err != nil {
		return nil, fmt.Errorf("failed to parse events JSON: %w, raw: %s", err, resp.Content)
	}

	return events, nil
}

// UpdateForeshadowing 伏笔自动追踪流程
func (a *ForesightAgent) UpdateForeshadowing(ctx context.Context, bookID uint, chapterID uint, chapterIndex int, chapterContent string, newEvents []models.StoryEvent, newStates string) error {
	// 1. 获取所有 open 的伏笔
	var openForeshadows []models.Foreshadowing
	if err := a.db.Where("book_id = ? AND status = ?", bookID, "open").Find(&openForeshadows).Error; err != nil {
		return err
	}

	// 2. 标记新伏笔 (从关键事件中筛选具有潜在影响的)
	for i := range newEvents {
		event := &newEvents[i]
		if event.UnresolvedImpact != "" || event.Importance >= 3 {
			foreshadow := models.Foreshadowing{
				BookID:                 bookID,
				ChapterID:              chapterID,
				ChapterIndex:           chapterIndex,
				Description:            event.Description,
				InvolvedCharacters:     event.InvolvedCharacters,
				DirectConsequence:      event.DirectConsequence,
				UnresolvedImpact:       event.UnresolvedImpact,
				Importance:             event.Importance,
				Status:                 "open",
				LastReferencedChapter: chapterIndex,
			}
			if err := a.db.Create(&foreshadow).Error; err != nil {
				return err
			}
		}
	}

	// 3. 检测旧伏笔是否被提及或回收
	for _, f := range openForeshadows {
		// A. 提及检测 (简单规则检测：关键词匹配)
		mentioned := false
		if strings.Contains(chapterContent, f.InvolvedCharacters) || strings.Contains(chapterContent, f.Description) {
			mentioned = true
		}

		if mentioned {
			a.db.Model(&f).Update("last_referenced_chapter", chapterIndex)
		}

		// B. 回收检测 (重要程度高的调用 LLM)
		isResolved := false
		reason := ""

		if f.Importance >= 4 {
			resolved, r, err := a.LLMResolveCheck(ctx, f, newEvents, newStates)
			if err == nil {
				isResolved = resolved
				reason = r
			}
		} else {
			// 低重要程度伏笔使用简单逻辑判定 (暂时留空或根据 EventType 匹配)
			for _, e := range newEvents {
				if strings.Contains(e.Description, f.UnresolvedImpact) || strings.Contains(e.DirectConsequence, f.UnresolvedImpact) {
					isResolved = true
					reason = "被新事件直接覆盖"
					break
				}
			}
		}

		if isResolved {
			a.db.Model(&f).Updates(map[string]interface{}{
				"status":                 "resolved",
				"resolved_chapter_index": chapterIndex,
				"resolve_reason":         reason,
			})
		}
	}

	return nil
}

// LLMResolveCheck 伏笔回收检测
func (a *ForesightAgent) LLMResolveCheck(ctx context.Context, f models.Foreshadowing, currentEvents []models.StoryEvent, characterStates string) (bool, string, error) {
	eventsJSON, _ := json.Marshal(currentEvents)
	data := map[string]string{
		"ForeshadowingDescription": f.Description,
		"UnresolvedImpact":         f.UnresolvedImpact,
		"CurrentEvents":            string(eventsJSON),
		"CharacterStates":          characterStates,
	}
	rendered, err := prompt.GetRegistry().Render("foreshadowing_resolution", data)
	if err != nil {
		return false, "", err
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
		return false, "", err
	}

	var result struct {
		IsResolved bool   `json:"is_resolved"`
		Reason     string `json:"reason"`
	}
	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &result); err != nil {
		return false, "", fmt.Errorf("failed to parse resolution JSON: %w, raw: %s", err, resp.Content)
	}

	return result.IsResolved, result.Reason, nil
}

// ForeshadowingAlert 伏笔风险预警
type ForeshadowingAlert struct {
	Foreshadowing models.Foreshadowing `json:"foreshadowing"`
	OverdueCount  int                  `json:"overdue_count"` // 超出章节数
	RiskLevel     string               `json:"risk_level"`    // low, medium, high
}

// GetOpenForeshadowingWithAlerts 获取未回收伏笔并附带风险预警
func (a *ForesightAgent) GetOpenForeshadowingWithAlerts(bookID uint, currentChapterIndex int) ([]ForeshadowingAlert, error) {
	var openForeshadows []models.Foreshadowing
	if err := a.db.Where("book_id = ? AND status = ?", bookID, "open").Order("chapter_index ASC").Find(&openForeshadows).Error; err != nil {
		return nil, err
	}

	var alerts []ForeshadowingAlert
	for _, f := range openForeshadows {
		// 计算超时章节数
		// 阈值定义：
		// 1. 重要程度 5：超过 5 章未提及/回收 -> High
		// 2. 重要程度 4：超过 10 章未提及/回收 -> Medium
		// 3. 其他：超过 20 章未提及/回收 -> Low
		
		overdue := currentChapterIndex - f.ChapterIndex
		riskLevel := "none"
		
		if f.Importance >= 5 {
			if overdue > 5 {
				riskLevel = "high"
			} else if overdue > 3 {
				riskLevel = "medium"
			} else {
				riskLevel = "low"
			}
		} else if f.Importance >= 4 {
			if overdue > 10 {
				riskLevel = "high"
			} else if overdue > 5 {
				riskLevel = "medium"
			} else {
				riskLevel = "low"
			}
		} else {
			if overdue > 20 {
				riskLevel = "high"
			} else if overdue > 10 {
				riskLevel = "medium"
			} else {
				riskLevel = "low"
			}
		}

		alerts = append(alerts, ForeshadowingAlert{
			Foreshadowing: f,
			OverdueCount:  overdue,
			RiskLevel:     riskLevel,
		})
	}

	return alerts, nil
}
