package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ai_novel/internal/service/agent"
	"ai_novel/internal/service/rag"
	"ai_novel/models"
	"ai_novel/pkg/logger"

	"gorm.io/gorm"
)

// PostWriteProcessor 统一后处理器：合并所有写作后处理逻辑
// 确保 Handler 路径 (WriteChapter) 和 Pipeline 路径 (ExecuteWriting) 使用完全一致的后处理链
type PostWriteProcessor struct {
	DB          *gorm.DB
	State       *agent.StateAgent
	Foresight   *agent.ForesightAgent
	Consistency *agent.ConsistencyAgent
	Summarizer  *agent.SummarizerAgent
	RAG         *rag.MemoryRecallService
}

// Process 执行完整的后处理链（应当异步调用）
func (p *PostWriteProcessor) Process(ctx context.Context, bookID, chapterID uint, chapterIndex int, fullContent string) {
	logger.Info("PostWriteProcessor: starting for Book=%d Chapter=%d", bookID, chapterID)

	// =========================================================================
	// 0. 数据清理与状态回滚 (处理重写场景)
	// =========================================================================
	// 清理本章产生的历史记录，防止重复堆积
	p.DB.Where("chapter_id = ?", chapterID).Delete(&models.StoryEvent{})
	p.DB.Where("chapter_id = ?", chapterID).Delete(&models.CharacterStateRecord{})
	p.DB.Where("chapter_id = ?", chapterID).Delete(&models.StoryContradiction{})
	p.DB.Where("chapter_id = ?", chapterID).Delete(&models.ChapterHealthScore{})
	p.DB.Where("chapter_id = ?", chapterID).Delete(&models.OOCScore{})
	
	// 清理向量索引 (本地 + 远程)
	if err := p.RAG.DeleteChapterIndex(ctx, bookID, chapterID); err != nil {
		logger.Info("PostWriteProcessor: failed to delete vector index: %v", err)
	}

	// 伏笔特殊处理：
	// 1. 删除本章引入的新伏笔
	p.DB.Where("book_id = ? AND chapter_index = ?", bookID, chapterIndex).Delete(&models.Foreshadowing{})
	// 2. 回滚本章解决的伏笔状态 (resolved -> open)
	p.DB.Model(&models.Foreshadowing{}).
		Where("book_id = ? AND resolved_chapter_index = ?", bookID, chapterIndex).
		Updates(map[string]interface{}{"status": "open", "resolved_chapter_index": 0})

	// =========================================================================
	// 1. 加载数据
	// =========================================================================
	var book models.Book
	if err := p.DB.Preload("Characters").First(&book, bookID).Error; err != nil {
		logger.Info("PostWriteProcessor: failed to load book %d: %v", bookID, err)
		return
	}

	// 1. 角色动态状态提取
	var baseProfiles, previousStates []string
	for _, char := range book.Characters {
		baseProfiles = append(baseProfiles, fmt.Sprintf("%s: %s (%s)", char.Name, char.Role, char.Description))
		stateJSON, _ := json.Marshal(char.DynamicState)
		previousStates = append(previousStates, fmt.Sprintf("%s: %s", char.Name, string(stateJSON)))
	}

	updates, err := p.State.ExtractDynamicStateChanges(
		ctx,
		book.WorldSetting.Summary,
		strings.Join(baseProfiles, "\n"),
		strings.Join(previousStates, "\n"),
		fullContent,
	)

	var events []models.StoryEvent

	if err != nil {
		logger.Info("PostWriteProcessor: state extraction failed: %v", err)
	} else {
		// 2. 更新角色状态 + LitRPG 合并
		for _, char := range book.Characters {
			if update, ok := updates[char.Name]; ok {
				char.DynamicState = update
				p.DB.Model(&char).Update("dynamic_state", char.DynamicState)

				// 记录历史轨迹
				stateRecord := models.CharacterStateRecord{
					CharacterID: char.ID,
					ChapterID:   chapterID,
					State:       update,
				}
				p.DB.Create(&stateRecord)
			}
		}

		// 3. 事件抽取（只做一次）
		events, err = p.Foresight.ExtractEvents(ctx, fullContent)
		if err != nil {
			logger.Info("PostWriteProcessor: event extraction failed: %v", err)
		} else {
			for i := range events {
				events[i].BookID = bookID
				events[i].ChapterID = chapterID
				events[i].ChapterIndex = chapterIndex
				p.DB.Create(&events[i])
			}

			// 4. 伏笔追踪
			updatesJSON, _ := json.Marshal(updates)
			p.Foresight.UpdateForeshadowing(ctx, bookID, chapterID, chapterIndex, fullContent, events, string(updatesJSON))
		}

		// 5. OOC 评分
		for _, char := range book.Characters {
			if update, ok := updates[char.Name]; ok {
				var anchor models.CharacterAnchor
				if err := p.DB.Where("character_id = ?", char.ID).First(&anchor).Error; err != nil {
					newAnchor, err := p.Consistency.ExtractCharacterAnchor(ctx, &char, "")
					if err == nil {
						anchor = *newAnchor
						p.DB.Create(&anchor)
					}
				}
				if anchor.ID != 0 {
					behavior := fmt.Sprintf("目标: %s\n行为: %s\n情绪: %s", update.Goal, update.KeyActions, update.EmotionalState)
					score, err := p.Consistency.EvaluateOOC(ctx, &anchor, "", behavior)
					if err == nil {
						score.ChapterID = chapterID
						p.DB.Create(score)
					}
				}
			}
		}
	}

	// 6. 矛盾检测
	recallQuery := fmt.Sprintf("分析当前章节内容是否存在与历史事实冲突: 第%d章", chapterIndex)
	historyMemory, err := p.RAG.Recall(ctx, bookID, recallQuery, 10, "")
	if err != nil {
		historyMemory = "无法召回历史记忆"
	}

	contradictions, err := p.Consistency.DetectContradictions(
		ctx,
		book.WorldSetting.Rules,
		historyMemory,
		strings.Join(previousStates, "\n"),
		fullContent,
	)
	if err == nil {
		for _, con := range contradictions {
			con.BookID = bookID
			con.ChapterID = chapterID
			p.DB.Create(&con)
		}
	}

	// 7. 健康度评估
	var oocScores []models.OOCScore
	p.DB.Where("chapter_id = ?", chapterID).Find(&oocScores)

	var openForeshadows, resolvedForeshadows []models.Foreshadowing
	p.DB.Where("book_id = ? AND status = ?", bookID, "open").Find(&openForeshadows)
	p.DB.Where("book_id = ? AND status = ? AND resolved_chapter_index = ?", bookID, "resolved", chapterIndex).Find(&resolvedForeshadows)

	healthScore := p.Consistency.EvaluateChapterHealth(ctx, oocScores, contradictions, openForeshadows, resolvedForeshadows)
	healthScore.BookID = bookID
	healthScore.ChapterID = chapterID
	p.DB.Create(healthScore)

	// 8. 章节摘要生成 + 回写
	var chapter models.Chapter
	if err := p.DB.First(&chapter, chapterID).Error; err == nil {
		summary, err := p.Summarizer.SummarizeChapter(ctx, chapter.Title, fullContent)
		if err == nil && summary != "" {
			p.DB.Model(&chapter).Update("summary", summary)
		}
	}

	// 9. RAG 向量化（同步，因为外层已经是异步的了）
	if chapter.Title != "" {
		// 章节摘要入库
		p.RAG.IndexChapter(ctx, bookID, chapterID, chapter.Title, fullContent)

		// 事件向量化
		for _, ev := range events {
			p.RAG.IndexEvent(ctx, bookID, chapterID, ev)
		}

		// 角色状态向量化
		if updates != nil {
			for name, state := range updates {
				charContent := fmt.Sprintf("角色: %s\n最新状态: %s\n当前目标: %s\n性格/情绪: %s",
					name, state.IdentityLocation, state.Goal, state.EmotionalState)
				p.RAG.IndexCharacter(ctx, bookID, name, charContent)
				p.RAG.IndexCharacterState(ctx, bookID, chapterID, name, state)
			}
		}
	}

	logger.Info("PostWriteProcessor: completed for Book=%d Chapter=%d", bookID, chapterID)
}
