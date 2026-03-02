package context

import (
	"context"
	"fmt"
	"strings"

	"ai_novel/internal/service/agent"
	"ai_novel/internal/service/rag"
	"ai_novel/models"

	"gorm.io/gorm"
)

type ContextManager struct {
	db         *gorm.DB
	summarizer *agent.SummarizerAgent
	writer     *agent.WriterAgent
	rag        *rag.MemoryRecallService
}

func NewContextManager(db *gorm.DB, summarizer *agent.SummarizerAgent, writer *agent.WriterAgent, rag *rag.MemoryRecallService) *ContextManager {
	return &ContextManager{
		db:         db,
		summarizer: summarizer,
		writer:     writer,
		rag:        rag,
	}
}

// BuildChapterContext 构建分层上下文
func (m *ContextManager) BuildChapterContext(ctx context.Context, bookID uint, chapterID uint) (agent.WriterContext, error) {
	var book models.Book
	if err := m.db.Preload("Chapters").Preload("Characters").First(&book, bookID).Error; err != nil {
		return agent.WriterContext{}, err
	}

	var chapter models.Chapter
	if err := m.db.First(&chapter, chapterID).Error; err != nil {
		return agent.WriterContext{}, err
	}

	// 1. 稳定上下文：世界观摘要 & 大纲摘要
	var selectedPlan models.OutlineVersion
	if err := m.db.Where("book_id = ? AND is_selected = ?", bookID, true).First(&selectedPlan).Error; err == nil {
		if book.WorldSetting.Summary == "" && selectedPlan.WorldView != "" {
			book.WorldSetting.Summary = selectedPlan.WorldView
		}
		if book.CurrentState.OutlineSummary == "" && selectedPlan.Outline != "" {
			book.CurrentState.OutlineSummary = selectedPlan.Outline
		}
	}

	if book.WorldSetting.Summary == "" {
		summary, err := m.summarizer.SummarizeWorld(ctx, book.WorldSetting.Description)
		if err == nil {
			book.WorldSetting.Summary = summary
			m.db.Model(&book).Update("world_setting", book.WorldSetting)
		}
	}

	if book.CurrentState.OutlineSummary == "" {
		summary, err := m.summarizer.SummarizeOutline(ctx, book.Description)
		if err == nil {
			book.CurrentState.OutlineSummary = summary
			m.db.Model(&book).Update("current_state", book.CurrentState)
		}
	}

	// 2. 角色动态状态
	var charStates []string
	for _, char := range book.Characters {
		// 结合静态描述与动态状态
		state := fmt.Sprintf("%s (%s):\n- 设定: %s\n- 当前身份/位置: %s\n- 目标: %s\n- 情绪: %s\n- 关系变化: %s\n- 能力/资源: %s\n- 限制/代价: %s\n- 关键行为: %s\n- 矛盾/伏笔: %s",
			char.Name,
			char.Role,
			char.Description,
			char.DynamicState.IdentityLocation,
			char.DynamicState.Goal,
			char.DynamicState.EmotionalState,
			char.DynamicState.RelationshipChanges,
			char.DynamicState.AbilityResourceChanges,
			char.DynamicState.ConstraintsCosts,
			char.DynamicState.KeyActions,
			char.DynamicState.ConflictsForeshadowing,
		)
		charStates = append(charStates, state)
	}

	// 3. 章节目标
	if chapter.Objective == "" {
		objective, err := m.writer.GenerateChapterObjective(ctx, book.CurrentState.OutlineSummary, chapter.Order, chapter.Title)
		if err == nil {
			chapter.Objective = objective
			m.db.Model(&chapter).Update("objective", objective)
		}
	}

	// 4. 多路向量召回
	query := fmt.Sprintf("%s %s", chapter.Title, chapter.Objective)
	multiMemories, _ := m.rag.MultiRouteRecall(ctx, bookID, query, 5)
	
	// 格式化召回结果
	var recallParts []string
	if items, ok := multiMemories[rag.CollectionCharacters]; ok && len(items) > 0 {
		recallParts = append(recallParts, "相关角色状态:\n"+strings.Join(items, "\n"))
	}
	if items, ok := multiMemories[rag.CollectionHistory]; ok && len(items) > 0 {
		recallParts = append(recallParts, "历史前情回顾:\n"+strings.Join(items, "\n"))
	}
	if items, ok := multiMemories[rag.CollectionOutlines]; ok && len(items) > 0 {
		recallParts = append(recallParts, "大纲主线目标:\n"+strings.Join(items, "\n"))
	}
	if items, ok := multiMemories[rag.CollectionWorldRules]; ok && len(items) > 0 {
		recallParts = append(recallParts, "世界观规则设定:\n"+strings.Join(items, "\n"))
	}
	
	memories := strings.Join(recallParts, "\n\n")
	if memories == "" {
		memories, _ = m.rag.Recall(ctx, bookID, query, 5, "")
	}

	// 5. 获取未回收伏笔
	var openForeshadows []models.Foreshadowing
	var foreshadowingText string
	if err := m.db.Where("book_id = ? AND status = ?", bookID, "open").Order("importance DESC").Limit(10).Find(&openForeshadows).Error; err == nil {
		var lines []string
		for _, f := range openForeshadows {
			lines = append(lines, fmt.Sprintf("- %s (引入: 第%d章, 未解影响: %s)", f.Description, f.ChapterIndex, f.UnresolvedImpact))
		}
		foreshadowingText = strings.Join(lines, "\n")
	}

	// 6. 线性上下文 (上一章结尾)
	var lastChapterTail string
	if chapter.Order > 1 {
		var prevChapter models.Chapter
		if err := m.db.Where("book_id = ? AND `order` = ?", bookID, chapter.Order-1).First(&prevChapter).Error; err == nil {
			content := prevChapter.Content
			if len(content) > 1000 {
				lastChapterTail = content[len(content)-1000:]
			} else {
				lastChapterTail = content
			}
		}
	}

	return agent.WriterContext{
		WorldSummary:      book.WorldSetting.Summary,
		OutlineSummary:    book.CurrentState.OutlineSummary,
		CharacterStates:   strings.Join(charStates, "\n\n"),
		ChapterIndex:      chapter.Order,
		ChapterTitle:      chapter.Title,
		ChapterObjective:  chapter.Objective,
		RetrievedMemories: memories,
		LastChapterTail:   lastChapterTail,
		Foreshadowing:     foreshadowingText,
		TargetWords:       2000,
	}, nil
}
