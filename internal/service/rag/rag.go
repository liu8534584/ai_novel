package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"ai_novel/internal/service/llm/core"
	"ai_novel/models"
	"ai_novel/pkg/vectorstore"

	"gorm.io/gorm"
)

const (
	CollectionWorldRules = "world_rules"
	CollectionCharacters = "characters"
	CollectionOutlines   = "outlines"
	CollectionHistory    = "history"
)

type MemoryRecallService struct {
	db          *gorm.DB
	llmProvider core.Provider
	vectorStore vectorstore.VectorStore
}

func NewMemoryRecallService(db *gorm.DB, llmProvider core.Provider, vectorStore vectorstore.VectorStore) *MemoryRecallService {
	return &MemoryRecallService{
		db:          db,
		llmProvider: llmProvider,
		vectorStore: vectorStore,
	}
}

// IndexChapter 将章节内容向量化并入库 (历史库)
func (s *MemoryRecallService) IndexChapter(ctx context.Context, bookID uint, chapterID uint, title string, content string) error {
	// 1. 简单的分段处理 (按段落或固定长度)
	paragraphs := strings.Split(content, "\n\n")
	var chunks []string
	var currentChunk string

	for _, p := range paragraphs {
		if len(currentChunk)+len(p) < 1000 {
			currentChunk += p + "\n\n"
		} else {
			if currentChunk != "" {
				chunks = append(chunks, strings.TrimSpace(currentChunk))
			}
			currentChunk = p + "\n\n"
		}
	}
	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSpace(currentChunk))
	}

	var docs []vectorstore.Document
	// 2. 对每个分段进行向量化并存储
	for i, chunk := range chunks {
		embedding, err := s.llmProvider.CreateEmbedding(ctx, chunk, core.Options{})
		if err != nil {
			fmt.Printf("Warning: failed to create embedding for chapter chunk %d: %v\n", i, err)
			continue
		}

		metadata := map[string]interface{}{
			"book_id":    bookID,
			"chapter_id": chapterID,
			"title":      title,
			"chunk_idx":  i,
			"type":       "chapter",
		}

		docs = append(docs, vectorstore.Document{
			Content:  chunk,
			Vector:   embedding,
			Metadata: metadata,
		})

		// 同时保留在本地数据库作为备份/管理
		embJSON, _ := json.Marshal(embedding)
		metaJSON, _ := json.Marshal(metadata)
		record := models.VectorRecord{
			BookID:    bookID,
			ChapterID: chapterID,
			Category:  "chapter",
			Content:   chunk,
			Embedding: string(embJSON),
			Metadata:  string(metaJSON),
		}
		s.db.Create(&record)
	}

	if s.vectorStore != nil && len(docs) > 0 {
		return s.vectorStore.AddDocuments(ctx, CollectionHistory, docs)
	}

	return nil
}

// DeleteChapterIndex 删除指定章节的所有向量索引
func (s *MemoryRecallService) DeleteChapterIndex(ctx context.Context, bookID uint, chapterID uint) error {
	// 1. 从本地数据库删除 (其实已经在 PostWriteProcessor 中清理了，这里为了完整性可以再调一次，或者略过)
	// 但考虑到 VectorRecord 是由 RAG 管理的，这里显式删除更合适。
	if err := s.db.Where("book_id = ? AND chapter_id = ?", bookID, chapterID).Delete(&models.VectorRecord{}).Error; err != nil {
		return err
	}

	// 2. 从远程向量库删除 (如果存在)
	if s.vectorStore != nil {
		filter := map[string]interface{}{
			"book_id":    bookID,
			"chapter_id": chapterID,
		}
		// 需要遍历所有可能的 Collection 进行清理
		collections := []string{CollectionHistory, CollectionCharacters, CollectionOutlines, CollectionWorldRules}
		for _, col := range collections {
			// 只有 history 和 characters 可能包含 specific chapter 的数据
			// outlines 和 world_rules 通常是全局的，或者是按 stage 划分的
			// 但如果有 chapter_id 字段，也可以清理
			s.vectorStore.DeleteDocuments(ctx, col, filter)
		}
	}
	return nil
}

// IndexWorldRule 将世界观规则入库
func (s *MemoryRecallService) IndexWorldRule(ctx context.Context, bookID uint, content string) error {
	embedding, err := s.llmProvider.CreateEmbedding(ctx, content, core.Options{})
	if err != nil {
		return err
	}

	doc := vectorstore.Document{
		Content: content,
		Vector:  embedding,
		Metadata: map[string]interface{}{
			"book_id": bookID,
			"type":    "world_rule",
		},
	}

	if s.vectorStore != nil {
		return s.vectorStore.AddDocuments(ctx, CollectionWorldRules, []vectorstore.Document{doc})
	}
	return nil
}

// IndexCharacter 将角色信息入库
func (s *MemoryRecallService) IndexCharacter(ctx context.Context, bookID uint, charName string, content string) error {
	embedding, err := s.llmProvider.CreateEmbedding(ctx, content, core.Options{})
	if err != nil {
		return err
	}

	doc := vectorstore.Document{
		Content: content,
		Vector:  embedding,
		Metadata: map[string]interface{}{
			"book_id":   bookID,
			"char_name": charName,
			"type":      "character",
		},
	}

	if s.vectorStore != nil {
		return s.vectorStore.AddDocuments(ctx, CollectionCharacters, []vectorstore.Document{doc})
	}
	return nil
}

// IndexOutline 将大纲阶段入库
func (s *MemoryRecallService) IndexOutline(ctx context.Context, bookID uint, stage string, content string) error {
	embedding, err := s.llmProvider.CreateEmbedding(ctx, content, core.Options{})
	if err != nil {
		return err
	}

	doc := vectorstore.Document{
		Content: content,
		Vector:  embedding,
		Metadata: map[string]interface{}{
			"book_id": bookID,
			"stage":   stage,
			"type":    "outline",
		},
	}

	if s.vectorStore != nil {
		return s.vectorStore.AddDocuments(ctx, CollectionOutlines, []vectorstore.Document{doc})
	}
	return nil
}

// IndexFullPlan 将整个计划入库（世界观、角色、大纲）
func (s *MemoryRecallService) IndexFullPlan(ctx context.Context, bookID uint, worldView, characters, outline string) error {
	// 1. 世界观切片入库
	rules := strings.Split(worldView, "\n\n")
	for _, rule := range rules {
		if strings.TrimSpace(rule) == "" {
			continue
		}
		s.IndexWorldRule(ctx, bookID, rule)
	}

	// 2. 角色库入库
	// 简单按行分割角色信息，理想情况应按角色对象
	charList := strings.Split(characters, "\n\n")
	for _, char := range charList {
		if strings.TrimSpace(char) == "" {
			continue
		}
		s.IndexCharacter(ctx, bookID, "unknown", char)
	}

	// 3. 大纲库入库
	stages := strings.Split(outline, "\n\n")
	for i, stage := range stages {
		if strings.TrimSpace(stage) == "" {
			continue
		}
		s.IndexOutline(ctx, bookID, fmt.Sprintf("stage_%d", i), stage)
	}

	return nil
}

// IndexEvent 将剧情事件向量化并入库
func (s *MemoryRecallService) IndexEvent(ctx context.Context, bookID uint, chapterID uint, event models.StoryEvent) error {
	content := fmt.Sprintf("事件描述: %s\n涉及角色: %s\n直接后果: %s\n潜在影响: %s", 
		event.Description, event.InvolvedCharacters, event.DirectConsequence, event.UnresolvedImpact)

	embedding, err := s.llmProvider.CreateEmbedding(ctx, content, core.Options{})
	if err != nil {
		return fmt.Errorf("failed to create embedding for event: %w", err)
	}

	embJSON, _ := json.Marshal(embedding)
	metadata := map[string]interface{}{
		"book_id":      bookID,
		"chapter_id":   chapterID,
		"event_id":     event.ID,
		"event_type":   event.EventType,
		"importance":   event.Importance,
		"type":         "event",
	}
	metaJSON, _ := json.Marshal(metadata)

	record := models.VectorRecord{
		BookID:    bookID,
		ChapterID: chapterID,
		Category:  "event",
		Content:   content,
		Embedding: string(embJSON),
		Metadata:  string(metaJSON),
	}

	return s.db.Create(&record).Error
}

// IndexCharacterState 将角色状态变更向量化并入库
func (s *MemoryRecallService) IndexCharacterState(ctx context.Context, bookID uint, chapterID uint, charName string, state models.CharacterDynamicState) error {
	content := fmt.Sprintf("角色: %s\n当前目标: %s\n关键行为: %s\n情绪状态: %s\n位置/身份: %s", 
		charName, state.Goal, state.KeyActions, state.EmotionalState, state.IdentityLocation)

	embedding, err := s.llmProvider.CreateEmbedding(ctx, content, core.Options{})
	if err != nil {
		return fmt.Errorf("failed to create embedding for character state: %w", err)
	}

	embJSON, _ := json.Marshal(embedding)
	metadata := map[string]interface{}{
		"book_id":    bookID,
		"chapter_id": chapterID,
		"char_name":  charName,
		"type":       "character",
	}
	metaJSON, _ := json.Marshal(metadata)

	record := models.VectorRecord{
		BookID:    bookID,
		ChapterID: chapterID,
		Category:  "character",
		Content:   content,
		Embedding: string(embJSON),
		Metadata:  string(metaJSON),
	}

	return s.db.Create(&record).Error
}

// Recall 相关记忆召回 (兼容旧接口，返回 string)
func (s *MemoryRecallService) Recall(ctx context.Context, bookID uint, query string, topK int, category string) (string, error) {
	records, err := s.recallRecords(ctx, bookID, query, topK, category)
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "无相关历史记忆", nil
	}
	// 拼接时增加 Category 标识
	var formatted []string
	for _, r := range records {
		formatted = append(formatted, fmt.Sprintf("[%s] %s", r.Category, r.Content))
	}
	return strings.Join(formatted, "\n---\n"), nil
}

// recallRecords 内部实现：返回 []models.VectorRecord
func (s *MemoryRecallService) recallRecords(ctx context.Context, bookID uint, query string, topK int, category string) ([]models.VectorRecord, error) {
	// 获取查询向量
	queryEmb, err := s.llmProvider.CreateEmbedding(ctx, query, core.Options{})
	if err != nil {
		return nil, err
	}

	// 从数据库中获取所有符合条件的记录
	var records []models.VectorRecord
	queryBuilder := s.db.Where("book_id = ?", bookID)
	if category != "" {
		queryBuilder = queryBuilder.Where("category = ?", category)
	}
	if err := queryBuilder.Find(&records).Error; err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	// 计算余弦相似度并排序
	type scoredRecord struct {
		record models.VectorRecord
		score  float32
	}
	var scored []scoredRecord

	for _, r := range records {
		var emb []float32
		if err := json.Unmarshal([]byte(r.Embedding), &emb); err != nil {
			continue
		}
		score := s.cosineSimilarity(queryEmb, emb)
		scored = append(scored, scoredRecord{record: r, score: score})
	}

	// 按分数降序排列
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// 取前 TopK
	limit := topK
	if len(scored) < limit {
		limit = len(scored)
	}

	var results []models.VectorRecord
	for i := 0; i < limit; i++ {
		results = append(results, scored[i].record)
	}

	return results, nil
}

// MultiRouteRecall 多路召回：同时搜索角色、历史、大纲库
func (s *MemoryRecallService) MultiRouteRecall(ctx context.Context, bookID uint, query string, topK int) (map[string][]string, error) {
	if s.vectorStore == nil {
		// 回退到原有的 SQLite 搜索，但模拟分类返回
		results := make(map[string][]string)
		collections := []string{CollectionCharacters, CollectionHistory, CollectionOutlines, CollectionWorldRules}
		
		for _, col := range collections {
			// 调用内部逻辑 recallRecords
			records, err := s.recallRecords(ctx, bookID, query, topK, col)
			if err == nil && len(records) > 0 {
				var items []string
				for _, r := range records {
					items = append(items, r.Content)
				}
				results[col] = items
			}
		}
		return results, nil
	}

	queryEmb, err := s.llmProvider.CreateEmbedding(ctx, query, core.Options{})
	if err != nil {
		return nil, err
	}

	results := make(map[string][]string)
	collections := []string{CollectionCharacters, CollectionHistory, CollectionOutlines, CollectionWorldRules}

	for _, col := range collections {
		docs, err := s.vectorStore.Search(ctx, col, queryEmb, topK)
		if err != nil {
			fmt.Printf("Warning: search failed for collection %s: %v\n", col, err)
			continue
		}

		var items []string
		for _, doc := range docs {
			// 过滤 book_id
			if bid, ok := doc.Metadata["book_id"].(float64); ok && uint(bid) != bookID {
				continue
			}
			items = append(items, doc.Content)
		}
		results[col] = items
	}

	return results, nil
}

func (s *MemoryRecallService) cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float32
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
