package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"ai_novel/internal/model"
	"ai_novel/internal/service/llm/core"
	"ai_novel/internal/service/rag"
	"ai_novel/models"
	"ai_novel/pkg/logger"
	"ai_novel/pkg/prompt"

	"gorm.io/gorm"
)

type DirectorAgent struct {
	llmProvider core.Provider
	db          *gorm.DB
}

func NewDirectorAgent(provider core.Provider, db *gorm.DB) *DirectorAgent {
	return &DirectorAgent{llmProvider: provider, db: db}
}

// InitWorld 根据灵感生成世界观文档
func (a *DirectorAgent) InitWorld(ctx context.Context, description, genre string, chapters int) (*model.WorldConfig, error) {
	// 1. 使用 Registry 渲染动态 Prompt
	data := map[string]interface{}{
		"Description": description,
		"Genre":       genre,
		"Chapters":    chapters,
	}
	rendered, err := prompt.GetRegistry().Render("director", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render director prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:    "",
		JSONMode: false,
	}
	core.GetStrategy(core.TaskWorldBuilding).ApplyToOptions(&options)

	// 2. 调用 LLM
	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("LLM call failed: %w", err)
	}

	// 3. 返回结构化文本内容
	return &model.WorldConfig{
		Content: core.RemoveReasoningContent(resp.Content),
	}, nil
}

// InitWorldStream 根据灵感流式生成世界观文档
func (a *DirectorAgent) InitWorldStream(ctx context.Context, description, genre string, chapters int) (<-chan core.StreamResponse, error) {
	// 1. 使用 Registry 渲染动态 Prompt
	data := map[string]interface{}{
		"Description": description,
		"Genre":       genre,
		"Chapters":    chapters,
	}
	rendered, err := prompt.GetRegistry().Render("director", data)
	if err != nil {
		return nil, fmt.Errorf("failed to render director prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		Model:    "",
		JSONMode: false,
	}
	core.GetStrategy(core.TaskWorldBuilding).ApplyToOptions(&options)

	// 2. 调用 LLM Stream
	streamResp, err := a.llmProvider.StreamChat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("LLM stream call failed: %w", err)
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

// ChatForInspiration 与 LLM 进行灵感头脑风暴
func (a *DirectorAgent) ChatForInspiration(ctx context.Context, history []core.Message) (string, error) {
	systemMsg, err := prompt.GetRegistry().Render("inspiration_chat", nil)
	if err != nil {
		return "", fmt.Errorf("failed to render inspiration_chat prompt: %w", err)
	}

	messages := append([]core.Message{
		{Role: core.RoleSystem, Content: systemMsg},
	}, history...)

	options := core.Options{}
	core.GetStrategy(core.TaskWorldBuilding).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Content, nil
}

// FinalizeInspiration 将对话记录加工成结构化的小说方案
func (a *DirectorAgent) FinalizeInspiration(ctx context.Context, conversation string) (string, error) {
	data := map[string]interface{}{
		"Conversation": conversation,
	}
	rendered, err := prompt.GetRegistry().Render("inspiration", data)
	if err != nil {
		return "", fmt.Errorf("failed to render inspiration prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	options := core.Options{
		JSONMode: true,
	}
	core.GetStrategy(core.TaskWorldBuilding).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return "", fmt.Errorf("LLM call failed: %w", err)
	}

	return resp.Content, nil
}

// =========================================================================
// Phase 3: 双 LLM 管线 — A→B→C 流水线方法
// =========================================================================

// GenerateStoryArc 阶段A：调用远程逻辑大模型生成分卷剧情树 (StoryArc)
func (a *DirectorAgent) GenerateStoryArc(ctx context.Context, bookID uint, startChapter, endChapter int) (*models.StoryArc, error) {
	// 1. 获取书籍信息以构建上下文
	var book models.Book
	if err := a.db.Preload("Characters").First(&book, bookID).Error; err != nil {
		return nil, fmt.Errorf("failed to load book %d: %w", bookID, err)
	}

	// 2. 构建 Prompt
	data := map[string]interface{}{
		"BookTitle":    book.Title,
		"Genre":        book.Genre,
		"WorldSetting": book.WorldSetting.Description,
		"WorldRules":   book.WorldSetting.Rules,
		"StartChapter": startChapter,
		"EndChapter":   endChapter,
	}

	// 拼装角色信息
	var charDescriptions []string
	for _, c := range book.Characters {
		charDescriptions = append(charDescriptions, fmt.Sprintf("- %s (%s): %s", c.Name, c.Role, c.Description))
	}
	data["Characters"] = strings.Join(charDescriptions, "\n")

	promptText := fmt.Sprintf(`你是一个专业的长篇小说剧情策划师。请根据以下信息，为小说 "%s"（%s 题材）生成第 %d 章到第 %d 章的分卷剧情树。

## 世界观
%s

## 世界规则
%s

## 主要角色
%s

请以 JSON 格式返回，包含以下字段：
- main_conflict: 本卷核心冲突 (string)
- turning_points: 关键转折点列表 (JSON array of strings)
- climax: 高潮设计 (string)
- foreshadowing: 需要埋下或回收的伏笔 (string)

仅返回 JSON，不要包含额外说明。`,
		book.Title, book.Genre, startChapter, endChapter,
		book.WorldSetting.Description, book.WorldSetting.Rules,
		strings.Join(charDescriptions, "\n"),
	)

	messages := []core.Message{
		{Role: core.RoleUser, Content: promptText},
	}

	options := core.Options{
		Model:    "", // 由 DynamicProvider 解析为远程逻辑大模型
		JSONMode: true,
	}
	core.GetStrategy(core.TaskPlanning).ApplyToOptions(&options)

	// 3. 调用远程大模型
	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate story arc: %w", err)
	}

	// 4. 解析 JSON 结果
	var arcData struct {
		MainConflict  string   `json:"main_conflict"`
		TurningPoints []string `json:"turning_points"`
		Climax        string   `json:"climax"`
		Foreshadowing string   `json:"foreshadowing"`
	}
	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &arcData); err != nil {
		return nil, fmt.Errorf("failed to parse story arc JSON: %w, raw: %s", err, resp.Content)
	}

	turningPointsJSON, _ := json.Marshal(arcData.TurningPoints)

	// 5. 存入数据库
	arc := &models.StoryArc{
		BookID:        bookID,
		StartChapter:  startChapter,
		EndChapter:    endChapter,
		MainConflict:  arcData.MainConflict,
		TurningPoints: string(turningPointsJSON),
		Climax:        arcData.Climax,
		Foreshadowing: arcData.Foreshadowing,
	}

	if err := a.db.Create(arc).Error; err != nil {
		return nil, fmt.Errorf("failed to save story arc: %w", err)
	}

	logger.Info("Generated StoryArc ID=%d for Book=%d (Chapters %d-%d)", arc.ID, bookID, startChapter, endChapter)
	return arc, nil
}

// GenerateBlueprint 阶段B：读取 StoryArc，调用远程大模型生成单章蓝图 (ChapterBlueprint)
func (a *DirectorAgent) GenerateBlueprint(ctx context.Context, arcID uint, chapterIndex int) (*models.ChapterBlueprint, error) {
	// 1. 读取对应的 StoryArc
	var arc models.StoryArc
	if err := a.db.First(&arc, arcID).Error; err != nil {
		return nil, fmt.Errorf("failed to load story arc %d: %w", arcID, err)
	}

	// 2. 获取书籍角色信息
	var book models.Book
	if err := a.db.Preload("Characters").First(&book, arc.BookID).Error; err != nil {
		return nil, fmt.Errorf("failed to load book %d: %w", arc.BookID, err)
	}

	var charDescriptions []string
	for _, c := range book.Characters {
		charDescriptions = append(charDescriptions, fmt.Sprintf("- %s (%s): %s", c.Name, c.Role, c.Description))
	}

	// 3. 构建 Prompt
	promptText := fmt.Sprintf(`你是一个专业的长篇小说章节策划师。根据以下分卷剧情树和角色信息，为第 %d 章生成详细的章节蓝图。

## 分卷信息 (第 %d 章 ~ 第 %d 章)
核心冲突：%s
关键转折点：%s
高潮设计：%s
伏笔规划：%s

## 角色
%s

请以纯 JSON 格式返回，包含以下字段：
- title: 章节标题 (string)
- summary: 本章预计发生什么 (string, 详细描述)
- character_changes: 预计角色状态变化 (string)
- world_changes: 预计世界线变动 (string)
- new_foreshadowing: 新增伏笔 (string)

仅返回 JSON，不要包含额外说明。`,
		chapterIndex,
		arc.StartChapter, arc.EndChapter,
		arc.MainConflict, arc.TurningPoints, arc.Climax, arc.Foreshadowing,
		strings.Join(charDescriptions, "\n"),
	)

	messages := []core.Message{
		{Role: core.RoleUser, Content: promptText},
	}

	options := core.Options{
		Model:    "", // 远程逻辑大模型
		JSONMode: true,
	}
	core.GetStrategy(core.TaskPlanning).ApplyToOptions(&options)

	// 4. 调用远程大模型
	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate chapter blueprint: %w", err)
	}

	// 5. 解析 JSON 结果
	var bpData struct {
		Title            string `json:"title"`
		Summary          string `json:"summary"`
		CharacterChanges string `json:"character_changes"`
		WorldChanges     string `json:"world_changes"`
		NewForeshadowing string `json:"new_foreshadowing"`
	}
	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &bpData); err != nil {
		return nil, fmt.Errorf("failed to parse blueprint JSON: %w, raw: %s", err, resp.Content)
	}

	// 6. 存入数据库
	blueprint := &models.ChapterBlueprint{
		BookID:           arc.BookID,
		ChapterIndex:     chapterIndex,
		Title:            bpData.Title,
		Summary:          bpData.Summary,
		CharacterChanges: bpData.CharacterChanges,
		WorldChanges:     bpData.WorldChanges,
		NewForeshadowing: bpData.NewForeshadowing,
	}

	if err := a.db.Create(blueprint).Error; err != nil {
		return nil, fmt.Errorf("failed to save chapter blueprint: %w", err)
	}

	logger.Info("Generated ChapterBlueprint ID=%d for Book=%d Chapter=%d", blueprint.ID, arc.BookID, chapterIndex)
	return blueprint, nil
}

// ExecuteWriting 阶段C：组装上下文 → 调用本地写作模型生成正文
// 返回值：streamChan (流式输出), finalChan (完整正文, channel关闭时发送), error
func (a *DirectorAgent) ExecuteWriting(
	ctx context.Context,
	ctxMgr ContextAssembler,
	writerAgent *WriterAgent,
	bookID uint,
	chapterID uint,
	chapterIndex int,
	writerModelName string,
) (<-chan string, <-chan string, error) {
	// 1. 组装完整三层记忆上下文
	wCtx, err := ctxMgr.AssembleWriterContext(ctx, bookID, chapterIndex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to assemble writer context: %w", err)
	}

	logger.Info("ExecuteWriting: Book=%d Chapter=%d | Context Algorithm: %s", bookID, chapterIndex, wCtx.SplicingAlgorithm())

	// 2. 渲染 Prompt
	rendered, err := prompt.GetRegistry().Render("writer_layered", wCtx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to render writer prompt: %w", err)
	}

	messages := []core.Message{
		{Role: core.RoleUser, Content: rendered},
	}

	// 3. 显式指定本地写作模型
	options := core.Options{
		Model:       writerModelName, // 如 "qwen3-14b-4bit"，确保正文由本地大模型生成
		Temperature: 0.9,
		MaxTokens:   4000,
	}

	streamResp, err := a.llmProvider.StreamChat(ctx, messages, options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start writing stream: %w", err)
	}

	outputChan := make(chan string)
	finalChan := make(chan string, 1)
	go func() {
		defer close(outputChan)
		defer close(finalChan)
		var sb strings.Builder
		for r := range streamResp {
			if r.Content != "" {
				sb.WriteString(r.Content)
				outputChan <- r.Content
			}
		}
		fullContent := sb.String()
		finalChan <- fullContent
		// 自动落库：防御性持久化，避免流式中断导致数据丢失
		if fullContent != "" {
			a.persistChapterContent(bookID, chapterID, chapterIndex, fullContent)
		}
	}()

	return outputChan, finalChan, nil
}

// PostProcessChapter 阶段D：调用远程逻辑大模型提取状态，更新记忆系统
func (a *DirectorAgent) PostProcessChapter(
	ctx context.Context,
	ragService *rag.MemoryRecallService,
	stateAgent *StateAgent,
	bookID uint,
	chapterID uint,
	chapterContent string,
) error {
	// 1. 获取章节与角色信息
	var chapter models.Chapter
	if err := a.db.First(&chapter, chapterID).Error; err != nil {
		return fmt.Errorf("failed to load chapter %d: %w", chapterID, err)
	}

	var book models.Book
	if err := a.db.Preload("Characters").First(&book, bookID).Error; err != nil {
		return fmt.Errorf("failed to load book %d: %w", bookID, err)
	}

	// 2. 构建状态提取 Prompt，请求远程逻辑大模型
	var charNames []string
	for _, c := range book.Characters {
		charNames = append(charNames, c.Name)
	}

	promptText := fmt.Sprintf(`阅读以下最新章节正文，提取主角与重要配角的最新状态变化和新事件。

## 角色列表
%s

## 章节正文
%s

请以纯 JSON 格式返回，包含以下三个顶层字段：
1. "character_states": 对象，key 为角色名，value 包含 identity_location, goal, emotional_state, relationship_changes, ability_resource_changes, constraints_costs, key_actions, conflicts_foreshadowing
2. "new_events": 数组，每项包含 event_type (主线推进/冲突升级/世界规则揭示/角色转折), description, involved_characters, direct_consequence, unresolved_impact, importance (1-5)
3. "litrpg_state_changes": 对象，key 为角色名，value 包含：
   - inventory_changes: { "new_items": ["物品A"], "removed_items": ["物品B"] } (物品增减)
   - stats_delta: { "HP": -10, "Exp": 50 } (数值属性增量变化)
   如果角色没有物品或数值变化，该角色可省略。

仅返回 JSON。`,
		strings.Join(charNames, ", "),
		chapterContent,
	)

	messages := []core.Message{
		{Role: core.RoleUser, Content: promptText},
	}

	options := core.Options{
		Model:    "", // 远程逻辑大模型 (由 DynamicProvider 解析)
		JSONMode: true,
	}
	core.GetStrategy(core.TaskStateTracking).ApplyToOptions(&options)

	resp, err := a.llmProvider.Chat(ctx, messages, options)
	if err != nil {
		return fmt.Errorf("post-process LLM call failed: %w", err)
	}

	// 3. 解析返回的 JSON
	var postResult struct {
		CharacterStates map[string]models.CharacterDynamicState `json:"character_states"`
		NewEvents       []struct {
			EventType          string `json:"event_type"`
			Description        string `json:"description"`
			InvolvedCharacters string `json:"involved_characters"`
			DirectConsequence  string `json:"direct_consequence"`
			UnresolvedImpact   string `json:"unresolved_impact"`
			Importance         int    `json:"importance"`
		} `json:"new_events"`
		LitRPGStateChanges map[string]struct {
			InventoryChanges struct {
				NewItems     []string `json:"new_items"`
				RemovedItems []string `json:"removed_items"`
			} `json:"inventory_changes"`
			StatsDelta map[string]int `json:"stats_delta"`
		} `json:"litrpg_state_changes"`
	}

	cleanJSON := core.ParseJSON(resp.Content)
	if err := json.Unmarshal([]byte(cleanJSON), &postResult); err != nil {
		return fmt.Errorf("failed to parse post-process JSON: %w, raw: %s", err, resp.Content)
	}

	// 4. 更新角色状态到 RAG 向量库
	for charName, state := range postResult.CharacterStates {
		if err := ragService.IndexCharacterState(ctx, bookID, chapterID, charName, state); err != nil {
			logger.Info("Warning: failed to index character state for %s: %v", charName, err)
		}

		// 同时更新数据库中的 Character.DynamicState
		a.db.Model(&models.Character{}).
			Where("book_id = ? AND name = ?", bookID, charName).
			Update("dynamic_state", state)
	}

	// 5. 存储新事件到 RAG
	for _, evt := range postResult.NewEvents {
		event := models.StoryEvent{
			BookID:             bookID,
			ChapterID:          chapterID,
			ChapterIndex:       chapter.Order,
			EventType:          evt.EventType,
			Description:        evt.Description,
			InvolvedCharacters: evt.InvolvedCharacters,
			DirectConsequence:  evt.DirectConsequence,
			UnresolvedImpact:   evt.UnresolvedImpact,
			Importance:         evt.Importance,
		}

		if err := a.db.Create(&event).Error; err != nil {
			logger.Info("Warning: failed to save event: %v", err)
			continue
		}

		if err := ragService.IndexEvent(ctx, bookID, chapterID, event); err != nil {
			logger.Info("Warning: failed to index event: %v", err)
		}
	}

	// 6. 将章节正文切块入库，作为未来的历史记忆
	if err := ragService.IndexChapter(ctx, bookID, chapterID, chapter.Title, chapterContent); err != nil {
		logger.Info("Warning: failed to index chapter content: %v", err)
	}

	// 7. 持久化 LitRPG 状态变更（物品 & 数值）
	for charName, changes := range postResult.LitRPGStateChanges {
		var char models.Character
		if err := a.db.Where("book_id = ? AND name = ?", bookID, charName).First(&char).Error; err != nil {
			logger.Info("Warning: character %s not found for LitRPG state update: %v", charName, err)
			continue
		}
		// 初始化
		if char.Inventory == nil {
			char.Inventory = []string{}
		}
		if char.Stats == nil {
			char.Stats = map[string]int{}
		}
		// 追加新物品
		char.Inventory = append(char.Inventory, changes.InventoryChanges.NewItems...)
		// 移除消耗物品
		for _, rm := range changes.InventoryChanges.RemovedItems {
			for i, item := range char.Inventory {
				if item == rm {
					char.Inventory = append(char.Inventory[:i], char.Inventory[i+1:]...)
					break
				}
			}
		}
		// 累加数值变化
		for k, v := range changes.StatsDelta {
			char.Stats[k] += v
		}
		// 持久化
		a.db.Model(&char).Updates(map[string]interface{}{
			"inventory": char.Inventory,
			"stats":     char.Stats,
		})
		logger.Info("LitRPG state updated for %s: Inventory=%v Stats=%v", charName, char.Inventory, char.Stats)
	}

	logger.Info("PostProcessChapter completed: Book=%d Chapter=%d | States=%d Events=%d LitRPG=%d",
		bookID, chapterID, len(postResult.CharacterStates), len(postResult.NewEvents), len(postResult.LitRPGStateChanges))

	return nil
}

// persistChapterContent 防御性落库：将流式输出的完整正文保存到 chapters 表并创建版本记录
func (a *DirectorAgent) persistChapterContent(bookID, chapterID uint, chapterIndex int, fullContent string) {
	wordCount := len([]rune(fullContent))

	var err error
	const maxVersionRetries = 3
	for attempt := 0; attempt < maxVersionRetries; attempt++ {
		err = a.db.Transaction(func(tx *gorm.DB) error {
			// 获取最新版本号
			var lastVersion models.ChapterVersion
			nextVersionNum := 1
			if err := tx.Where("chapter_id = ?", chapterID).Order("version desc").First(&lastVersion).Error; err == nil {
				nextVersionNum = lastVersion.Version + 1
			}

			// 创建版本记录（并发下依赖唯一索引兜底）
			version := models.ChapterVersion{
				ChapterID: chapterID,
				Version:   nextVersionNum,
				Content:   fullContent,
				WordCount: wordCount,
			}
			if err := tx.Create(&version).Error; err != nil {
				return err
			}

			// 更新章节主表
			if err := tx.Model(&models.Chapter{}).Where("id = ?", chapterID).Updates(map[string]interface{}{
				"content":         fullContent,
				"current_version": nextVersionNum,
			}).Error; err != nil {
				return err
			}

			// 更新书籍当前状态
			if err := tx.Model(&models.Book{}).Where("id = ?", bookID).Update("current_state", models.CurrentState{
				ChapterIndex: chapterIndex,
				Summary:      fmt.Sprintf("Pipeline completed Chapter %d (Version: %d, Words: %d)", chapterIndex, nextVersionNum, wordCount),
			}).Error; err != nil {
				return err
			}

			return nil
		})
		if err == nil {
			break
		}
		if !isDuplicateVersionPersistErr(err) {
			break
		}
	}

	if err != nil {
		logger.Info("Warning: failed to persist chapter content: %v", err)
	} else {
		logger.Info("Auto-persisted chapter content: Book=%d Chapter=%d Words=%d", bookID, chapterID, wordCount)
	}
}

func isDuplicateVersionPersistErr(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint") || strings.Contains(msg, "duplicate key")
}

// ContextAssembler 接口，用于解耦 DirectorAgent 与 ContextManager
type ContextAssembler interface {
	AssembleWriterContext(ctx context.Context, bookID uint, chapterIndex int) (WriterContext, error)
}
