package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel defines the common fields for all models with custom JSON tags.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// WorldSetting represents the world setting, stored as JSON.
type WorldSetting struct {
	Genre       string `json:"genre"`
	Tone        string `json:"tone"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
	Summary     string `json:"summary"` // 摘要版世界观 (300-500字)
}

// BookLLMConfig represents the LLM configuration for a specific book.
type BookLLMConfig struct {
	Provider string `json:"provider"`
	APIKey   string `json:"api_key"`
	BaseURL  string `json:"base_url"`
	Model    string `json:"model"`
}

type InspirationChat struct {
	BaseModel
	Messages string `gorm:"type:text"` // JSON string of chat messages
}

type PromptBindings struct {
	WorldView      []string `json:"world_view"`
	Plan           []string `json:"plan"`
	Character      []string `json:"character"`
	ChapterTitle   []string `json:"chapter_title"`
	ChapterOutline []string `json:"chapter_outline"`
	Writing        []string `json:"writing"`
	Review         []string `json:"review"`
}

// CharacterDynamicState 角色动态状态表
type CharacterDynamicState struct {
	IdentityLocation       string `json:"identity_location"`       // 当前身份 / 位置
	Goal                   string `json:"goal"`                    // 当前目标
	EmotionalState         string `json:"emotional_state"`         // 当前情绪状态
	RelationshipChanges    string `json:"relationship_changes"`    // 当前关系变化
	AbilityResourceChanges string `json:"ability_resource_changes"` // 能力 / 资源变化
	ConstraintsCosts       string `json:"constraints_costs"`       // 新增限制或代价
	KeyActions             string `json:"key_actions"`             // 本章关键行为
	ConflictsForeshadowing string `json:"conflicts_foreshadowing"` // 潜在矛盾 / 伏笔
}

// CharacterState represents the dynamic state of the character.
type CharacterState struct {
	HP        int      `json:"hp"`
	MaxHP     int      `json:"max_hp"`
	MP        int      `json:"mp"`
	MaxMP     int      `json:"max_mp"`
	Exp       int      `json:"exp"`
	Level     int      `json:"level"`
	Inventory []string `json:"inventory"`
	Skills    []string `json:"skills"`
	Status    []string `json:"status"` // e.g., "Injured", "Poisoned"
	Location  string   `json:"location"`
}

// CurrentState represents the current state of the novel generation, stored as JSON.
type CurrentState struct {
	ChapterIndex   int            `json:"chapter_index"`
	Summary        string         `json:"summary"`
	OutlineSummary string         `json:"outline_summary"` // 全局大纲摘要
	Notes          string         `json:"notes"`
	Character      CharacterState `json:"character"`
}

// Book represents a novel.
type Book struct {
	BaseModel
	Title         string `json:"title" gorm:"index;not null"`
	Author        string         `json:"author"`
	Genre         string         `json:"genre" gorm:"index"`
	Tags          string         `json:"tags"`     // 标签，逗号分隔
	Language      string         `json:"language"` // 创作语言
	Description   string         `json:"description" gorm:"type:text"`
	TotalChapters int            `json:"total_chapters" gorm:"default:1"`
	Status        string         `json:"status" gorm:"index;default:'draft'"` // draft, planning, writing, completed

	// WorldSetting and CurrentState are stored as JSON strings in the database
	// GORM's serializer:json tag automatically handles marshalling/unmarshalling
	WorldSetting WorldSetting `gorm:"serializer:json" json:"world_setting"`
	CurrentState CurrentState `gorm:"serializer:json" json:"current_state"`
	LLMConfig    BookLLMConfig `gorm:"serializer:json" json:"llm_config"`
	PromptBindings PromptBindings `gorm:"serializer:json" json:"prompt_bindings"`

	Chapters   []Chapter   `json:"chapters" gorm:"foreignKey:BookID"`
	Characters []Character `json:"characters" gorm:"foreignKey:BookID"`
}

// ChapterVersion represents a version of a chapter content.
type ChapterVersion struct {
	BaseModel
	ChapterID uint   `json:"chapter_id" gorm:"index"`
	Version   int    `json:"version"`
	Title     string `json:"title"`
	Content   string `json:"content" gorm:"type:text"`
	WordCount int    `json:"word_count"`
	Summary   string `json:"summary" gorm:"type:text"` // 自动生成的版本摘要
}

// VectorRecord represents a vectorized chunk of content for RAG.
type VectorRecord struct {
	BaseModel
	BookID    uint    `json:"book_id" gorm:"index"`
	ChapterID uint    `json:"chapter_id" gorm:"index"`
	Category  string  `json:"category" gorm:"index"`         // chapter, event, character
	Content   string  `json:"content" gorm:"type:text"`
	Embedding string  `json:"embedding" gorm:"type:text"` // 存储为 JSON 字符串
	Metadata  string  `json:"metadata" gorm:"type:text"`     // 存储为 JSON 字符串
}

// Chapter represents a chapter in the book.
type Chapter struct {
	BaseModel
	BookID         uint             `json:"book_id"`
	Title          string           `json:"title"`
	Content        string           `json:"content" gorm:"type:text"`
	Order          int              `json:"order"`
	Objective      string           `json:"objective" gorm:"type:text"`   // 章节目标
	Summary        string           `json:"summary" gorm:"type:text"`     // 章节摘要
	UserIntent         string           `json:"user_intent" gorm:"type:text"` // 用户写作意图
	Outline            string           `json:"outline" gorm:"type:text"`     // 章节大纲 (JSON string)
	IsOutlineConfirmed bool             `json:"is_outline_confirmed" gorm:"default:false"`
	CurrentVersion     int              `json:"current_version" gorm:"default:1"`
	Versions       []ChapterVersion `json:"versions" gorm:"foreignKey:ChapterID"`
}

// Character represents a character in the book (Static info).
type Character struct {
	BaseModel
	BookID       uint                  `json:"book_id"`
	Name         string                `json:"name"`
	Role         string                `json:"role"` // e.g., Protagonist, Antagonist
	Description  string                `json:"description"`
	DynamicState CharacterDynamicState `gorm:"serializer:json" json:"dynamic_state"` // 角色动态状态表
}

// CharacterStateRecord 角色状态变更记录 (历史轨迹)
type CharacterStateRecord struct {
	BaseModel
	CharacterID uint                  `json:"character_id" gorm:"index"`
	ChapterID   uint                  `json:"chapter_id" gorm:"index"`
	State       CharacterDynamicState `json:"state" gorm:"serializer:json"`
}

// StoryEvent 关键剧情事件 (因果链)
type StoryEvent struct {
	BaseModel
	BookID             uint   `json:"book_id" gorm:"index"`
	ChapterID          uint   `json:"chapter_id" gorm:"index"`
	ChapterIndex       int    `json:"chapter_index"`
	EventType          string `json:"event_type"`                   // 主线推进, 冲突升级, 世界规则揭示, 角色转折
	Description        string `json:"description" gorm:"type:text"` // 事件描述
	InvolvedCharacters string `json:"involved_characters"`          // 涉及角色
	DirectConsequence  string `json:"direct_consequence" gorm:"type:text"`  // 直接后果
	UnresolvedImpact   string `json:"unresolved_impact" gorm:"type:text"`    // 潜在影响 (伏笔核心)
	Importance         int    `json:"importance"`                   // 重要程度: 1 ~ 5
}

type OutlineVersion struct {
	BaseModel
	BookID     uint   `json:"book_id" gorm:"index"`
	Version    int    `json:"version"` // 版本号
	WorldView  string `json:"world_view" gorm:"type:text"`
	Outline    string `json:"outline" gorm:"type:text"`
	Characters string `json:"characters" gorm:"type:text"`
	Titles     string `json:"titles" gorm:"type:text"`
	IsSelected bool   `json:"is_selected" gorm:"default:false"`
	IsLocked   bool   `json:"is_locked" gorm:"default:false"` // 锁定机制：锁定后不可删除/修改
}

// StateUpdate represents a partial update to the state.
type StateUpdate struct {
	StatsDelta map[string]int `json:"stats_delta"` // e.g., {"hp": -10, "exp": 50}
	NewItems   []string       `json:"new_items"`
	LostItems  []string       `json:"lost_items"`
	NewStatus  []string       `json:"new_status"`
	Summary    string         `json:"summary"` // Update plot summary
}

// ApplyStateUpdate applies changes to the book's current state.
func (b *Book) ApplyStateUpdate(update StateUpdate) error {
	// 1. Update Stats
	if val, ok := update.StatsDelta["hp"]; ok {
		b.CurrentState.Character.HP += val
		if b.CurrentState.Character.HP < 0 {
			b.CurrentState.Character.HP = 0
		}
		if b.CurrentState.Character.HP > b.CurrentState.Character.MaxHP {
			b.CurrentState.Character.HP = b.CurrentState.Character.MaxHP
		}
	}
	if val, ok := update.StatsDelta["mp"]; ok {
		b.CurrentState.Character.MP += val
		if b.CurrentState.Character.MP < 0 {
			b.CurrentState.Character.MP = 0
		}
		if b.CurrentState.Character.MP > b.CurrentState.Character.MaxMP {
			b.CurrentState.Character.MP = b.CurrentState.Character.MaxMP
		}
	}
	if val, ok := update.StatsDelta["exp"]; ok {
		b.CurrentState.Character.Exp += val
	}

	// 2. Update Inventory
	for _, item := range update.NewItems {
		b.CurrentState.Character.Inventory = append(b.CurrentState.Character.Inventory, item)
	}

	if len(update.LostItems) > 0 {
		// Simple remove logic
		newInv := []string{}
		for _, item := range b.CurrentState.Character.Inventory {
			keep := true
			for _, lost := range update.LostItems {
				if item == lost {
					keep = false
					break
				}
			}
			if keep {
				newInv = append(newInv, item)
			}
		}
		b.CurrentState.Character.Inventory = newInv
	}

	// 3. Update Status
	for _, s := range update.NewStatus {
		exists := false
		for _, existing := range b.CurrentState.Character.Status {
			if existing == s {
				exists = true
				break
			}
		}
		if !exists {
			b.CurrentState.Character.Status = append(b.CurrentState.Character.Status, s)
		}
	}

	// 4. Update Summary
	if update.Summary != "" {
		b.CurrentState.Summary = update.Summary
	}

	return nil
}

// Foreshadowing represents a plot foreshadowing (伏笔).
type Foreshadowing struct {
	BaseModel
	BookID                uint   `json:"book_id" gorm:"index"`
	ChapterID             uint   `json:"chapter_id" gorm:"index"`      // 引入伏笔的章节ID
	ChapterIndex          int    `json:"chapter_index"`                // 引入伏笔的章节序号
	EventType             string `json:"event_type"`                   // 事件类型：主线推进, 冲突升级, 世界规则揭示...
	Description           string `json:"description" gorm:"type:text"` // 伏笔描述
	InvolvedCharacters    string `json:"involved_characters"`          // 涉及角色
	DirectConsequence     string `json:"direct_consequence"`           // 直接后果
	UnresolvedImpact      string `json:"unresolved_impact" gorm:"type:text"` // 未解决的影响 (伏笔核心)
	Status                string `json:"status" gorm:"index"`          // 状态: open, resolved, deprecated
	Importance            int    `json:"importance"`                   // 重要程度: 1 ~ 5
	LastReferencedChapter int    `json:"last_referenced_chapter"`      // 最近一次被提及或相关的章节序号
	ResolvedChapterIndex  int    `json:"resolved_chapter_index"`       // 回收伏笔的章节序号
	ResolveReason         string `json:"resolve_reason" gorm:"type:text"` // 回收说明
}

// CharacterAnchor 角色性格锚点 (稳定，慢变)
type CharacterAnchor struct {
	BaseModel
	CharacterID        uint   `json:"character_id"`
	PersonalityLabels  string `json:"personality_labels"`   // 核心性格标签 (3–5 个)
	CoreMotivation     string `json:"core_motivation"`      // 核心动机 (长期)
	BehaviorBottomLine string `json:"behavior_bottom_line"` // 行为底线 (绝不做的事)
	DecisionTendency   string `json:"decision_tendency"`    // 决策倾向 (保守 / 激进 / 利己 / 利他)
	EmotionalTriggers  string `json:"emotional_triggers"`   // 情绪触发点
}

// OOCScore 角色 OOC 评分
type OOCScore struct {
	BaseModel
	CharacterID            uint    `json:"character_id"`
	ChapterID              uint    `json:"chapter_id"`
	PersonalityConsistency float64 `json:"personality_consistency"` // 性格一致性偏离评分 (0–100)
	MotivationConsistency  float64 `json:"motivation_consistency"`  // 动机一致性偏离评分 (0–100)
	EmotionalReasonability float64 `json:"emotional_reasonability"` // 情绪反应合理性评分 (0–100)
	CostMissing            float64 `json:"cost_missing"`            // 行为代价缺失评分 (0–100)
	TotalScore             float64 `json:"total_score"`             // 综合 OOC 评分 (0–100)
	Conclusion             string  `json:"conclusion"`              // 结论: 无明显 OOC / 轻度 OOC / 明显 OOC / 严重 OOC
	Explanation            string  `json:"explanation" gorm:"type:text"` // 说明
}

// StoryContradiction 剧情矛盾记录
type StoryContradiction struct {
	BaseModel
	BookID      uint   `json:"book_id"`
	ChapterID   uint   `json:"chapter_id"`
	Type        string `json:"type"`        // 逻辑冲突 / 设定冲突 / 状态冲突
	Severity    string `json:"severity"`    // 轻微 / 显著 / 严重
	Description string `json:"description" gorm:"type:text"`
	Reference   string `json:"reference" gorm:"type:text"`
	Suggestion  string `json:"suggestion" gorm:"type:text"`
}

// ChapterHealthScore 章节健康度总评分
type ChapterHealthScore struct {
	BaseModel
	BookID           uint    `json:"book_id"`
	ChapterID        uint    `json:"chapter_id"`
	OOCScore         float64 `json:"ooc_score"`          // OOC 平均分
	EventConsistency float64 `json:"event_consistency"` // 剧情一致性得分 (100 - 扣分)
	Foreshadowing    float64 `json:"foreshadowing"`    // 伏笔活跃度/回收质量得分
	TotalHealth      float64 `json:"total_health"`      // 综合健康度评分
	AuditReport      string  `json:"audit_report" gorm:"type:text"` // 审计建议
}
