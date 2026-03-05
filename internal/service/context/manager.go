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

// BuildChapterContext 构建分层上下文 (已有方法，增加 RecentContext 支持)
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
	charStates := m.buildCharacterStates(book.Characters)

	// 3. 章节目标
	if chapter.Objective == "" {
		objective, err := m.writer.GenerateChapterObjective(ctx, book.CurrentState.OutlineSummary, chapter.Order, chapter.Title)
		if err == nil {
			chapter.Objective = objective
			m.db.Model(&chapter).Update("objective", objective)
		}
	}

	// 4. 多路向量召回
	memories := m.recallMemories(ctx, bookID, chapter.Title, chapter.Objective)

	// 5. 获取未回收伏笔
	foreshadowingText := m.buildForeshadowingText(bookID)

	// 6. 线性上下文 (上一章结尾)
	lastChapterTail := m.getLastChapterTail(bookID, chapter.Order)

	// 7. 第三层：滑动窗口摘要 (新增)
	recentContext := m.buildRecentContext(ctx, bookID, chapter.Order)

	return agent.WriterContext{
		WorldSummary:      book.WorldSetting.Summary,
		OutlineSummary:    book.CurrentState.OutlineSummary,
		CharacterStates:   strings.Join(charStates, "\n\n"),
		ChapterIndex:      chapter.Order,
		ChapterTitle:      chapter.Title,
		ChapterObjective:  chapter.Objective,
		RecentContext:     recentContext,
		RetrievedMemories: memories,
		LastChapterTail:   lastChapterTail,
		Foreshadowing:     foreshadowingText,
		TargetWords:       2000,
	}, nil
}

// AssembleWriterContext 为 Director Pipeline 专用的上下文组装方法 (实现 ContextAssembler 接口)
// 通过 chapterIndex 而非 chapterID 定位章节，并集成 ChapterBlueprint 映射
func (m *ContextManager) AssembleWriterContext(ctx context.Context, bookID uint, chapterIndex int) (agent.WriterContext, error) {
	// 1. 加载书籍
	var book models.Book
	if err := m.db.Preload("Characters").First(&book, bookID).Error; err != nil {
		return agent.WriterContext{}, fmt.Errorf("failed to load book %d: %w", bookID, err)
	}

	// --- 第一层：固定世界观 ---
	worldSummary, outlineSummary := m.buildLayer1(ctx, &book)

	// --- 第二层：动态角色状态 ---
	charStates := m.buildCharacterStates(book.Characters)

	// --- 第三层：最近 3-5 章滑动窗口摘要 ---
	recentContext := m.buildRecentContext(ctx, bookID, chapterIndex)

	// --- 蓝图映射：从 ChapterBlueprint 获取本章目标 ---
	var blueprint models.ChapterBlueprint
	var chapterObjective string
	if err := m.db.Where("book_id = ? AND chapter_index = ?", bookID, chapterIndex).First(&blueprint).Error; err == nil {
		chapterObjective = blueprint.Summary
	}

	// 从 StoryArc 中获取本卷大纲摘要 (如果存在则覆盖)
	var arc models.StoryArc
	if err := m.db.Where("book_id = ? AND start_chapter <= ? AND end_chapter >= ?", bookID, chapterIndex, chapterIndex).
		Order("created_at DESC").First(&arc).Error; err == nil {
		if outlineSummary == "" {
			outlineSummary = fmt.Sprintf("本卷核心冲突：%s\n高潮设计：%s", arc.MainConflict, arc.Climax)
		}
	}

	// --- 查找章节标题 ---
	var chapter models.Chapter
	var chapterTitle string
	if err := m.db.Where("book_id = ? AND `order` = ?", bookID, chapterIndex).First(&chapter).Error; err == nil {
		chapterTitle = chapter.Title
		if chapterObjective == "" {
			chapterObjective = chapter.Objective
		}
	} else if blueprint.Title != "" {
		chapterTitle = blueprint.Title
	}

	// --- RAG 召回 ---
	query := fmt.Sprintf("%s %s", chapterTitle, chapterObjective)
	memories := m.recallMemories(ctx, bookID, chapterTitle, chapterObjective)
	_ = query

	// --- 伏笔 ---
	foreshadowingText := m.buildForeshadowingText(bookID)

	// --- 上一章结尾 ---
	lastChapterTail := m.getLastChapterTail(bookID, chapterIndex)

	return agent.WriterContext{
		WorldSummary:      worldSummary,
		OutlineSummary:    outlineSummary,
		CharacterStates:   strings.Join(charStates, "\n\n"),
		ChapterIndex:      chapterIndex,
		ChapterTitle:      chapterTitle,
		ChapterObjective:  chapterObjective,
		RecentContext:     recentContext,
		RetrievedMemories: memories,
		LastChapterTail:   lastChapterTail,
		Foreshadowing:     foreshadowingText,
		TargetWords:       2000,
	}, nil
}

// =========================================================================
// 内部辅助方法
// =========================================================================

// buildLayer1 构建第一层记忆：世界观 + 大纲
func (m *ContextManager) buildLayer1(ctx context.Context, book *models.Book) (worldSummary, outlineSummary string) {
	// 尝试从选中的 Plan 加载
	var selectedPlan models.OutlineVersion
	if err := m.db.Where("book_id = ? AND is_selected = ?", book.ID, true).First(&selectedPlan).Error; err == nil {
		if book.WorldSetting.Summary == "" && selectedPlan.WorldView != "" {
			book.WorldSetting.Summary = selectedPlan.WorldView
		}
		if book.CurrentState.OutlineSummary == "" && selectedPlan.Outline != "" {
			book.CurrentState.OutlineSummary = selectedPlan.Outline
		}
	}

	// 如果仍然为空，通过 Summarizer 生成
	if book.WorldSetting.Summary == "" && book.WorldSetting.Description != "" {
		summary, err := m.summarizer.SummarizeWorld(ctx, book.WorldSetting.Description)
		if err == nil {
			book.WorldSetting.Summary = summary
			m.db.Model(book).Update("world_setting", book.WorldSetting)
		}
	}

	if book.CurrentState.OutlineSummary == "" && book.Description != "" {
		summary, err := m.summarizer.SummarizeOutline(ctx, book.Description)
		if err == nil {
			book.CurrentState.OutlineSummary = summary
			m.db.Model(book).Update("current_state", book.CurrentState)
		}
	}

	return book.WorldSetting.Summary, book.CurrentState.OutlineSummary
}

// buildCharacterStates 构建第二层记忆：角色动态状态
func (m *ContextManager) buildCharacterStates(characters []models.Character) []string {
	var charStates []string
	for _, char := range characters {
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
	return charStates
}

// buildRecentContext 构建第三层记忆：最近 3-5 章滑动窗口摘要
func (m *ContextManager) buildRecentContext(ctx context.Context, bookID uint, currentChapterIndex int) string {
	windowSize := 3 // 默认回看 3 章
	startIndex := currentChapterIndex - windowSize
	if startIndex < 1 {
		startIndex = 1
	}

	var recentChapters []models.Chapter
	if err := m.db.Where("book_id = ? AND `order` >= ? AND `order` < ?", bookID, startIndex, currentChapterIndex).
		Order("`order` ASC").Find(&recentChapters).Error; err != nil {
		return ""
	}

	if len(recentChapters) == 0 {
		return ""
	}

	var summaries []string
	for _, ch := range recentChapters {
		// 优先使用已有摘要
		if ch.Summary != "" {
			summaries = append(summaries, fmt.Sprintf("【第%d章 - %s】%s", ch.Order, ch.Title, ch.Summary))
			continue
		}

		// 如果没有摘要但有正文，则通过 Summarizer 生成
		if ch.Content != "" {
			summary, err := m.summarizer.SummarizeChapter(ctx, ch.Title, ch.Content)
			if err == nil && summary != "" {
				// 回写摘要到数据库
				m.db.Model(&ch).Update("summary", summary)
				summaries = append(summaries, fmt.Sprintf("【第%d章 - %s】%s", ch.Order, ch.Title, summary))
			}
		}
	}

	return strings.Join(summaries, "\n\n")
}

// recallMemories 多路 RAG 召回
func (m *ContextManager) recallMemories(ctx context.Context, bookID uint, chapterTitle, chapterObjective string) string {
	query := fmt.Sprintf("%s %s", chapterTitle, chapterObjective)
	multiMemories, _ := m.rag.MultiRouteRecall(ctx, bookID, query, 5)

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

	return memories
}

// buildForeshadowingText 获取未回收伏笔
func (m *ContextManager) buildForeshadowingText(bookID uint) string {
	var openForeshadows []models.Foreshadowing
	if err := m.db.Where("book_id = ? AND status = ?", bookID, "open").Order("importance DESC").Limit(10).Find(&openForeshadows).Error; err == nil {
		var lines []string
		for _, f := range openForeshadows {
			lines = append(lines, fmt.Sprintf("- %s (引入: 第%d章, 未解影响: %s)", f.Description, f.ChapterIndex, f.UnresolvedImpact))
		}
		return strings.Join(lines, "\n")
	}
	return ""
}

// getLastChapterTail 获取上一章结尾
func (m *ContextManager) getLastChapterTail(bookID uint, currentOrder int) string {
	if currentOrder <= 1 {
		return ""
	}

	var prevChapter models.Chapter
	if err := m.db.Where("book_id = ? AND `order` = ?", bookID, currentOrder-1).First(&prevChapter).Error; err == nil {
		content := prevChapter.Content
		if len(content) > 1000 {
			return content[len(content)-1000:]
		}
		return content
	}
	return ""
}
