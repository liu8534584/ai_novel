package model

type Character struct {
	Name        string          `json:"name"`
	Role        string          `json:"role"` // Protagonist, Antagonist, Supporting
	Description string          `json:"description"`
	Anchor      CharacterAnchor `json:"anchor"`
	Stats       map[string]int  `json:"stats"`
	Inventory   []string        `json:"inventory"`
	Skills      []string        `json:"skills"`
}

type CharacterAnchor struct {
	PersonalityLabels  string `json:"personality_labels"`   // 核心性格标签
	CoreMotivation     string `json:"core_motivation"`      // 核心动机
	BehaviorBottomLine string `json:"behavior_bottom_line"` // 行为底线
	DecisionTendency   string `json:"decision_tendency"`    // 决策倾向
	EmotionalTriggers  string `json:"emotional_triggers"`   // 情绪触发点
}
