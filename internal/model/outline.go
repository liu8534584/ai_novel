package model

// ChapterOutline 章节大纲
type ChapterOutline struct {
	ChapterNumber int     `json:"chapter_number"`
	Title         string  `json:"title"`
	Summary       string  `json:"summary"`       // 本章核心冲突/目标
	Scenes        []Scene `json:"scenes"`        // 细化到具体的 3-4 个场景
}

// Scene 场景/情节点
type Scene struct {
	Order       int      `json:"order"`
	Location    string   `json:"location"`
	Description string   `json:"description"`  // 场景发生的事
	Characters  []string `json:"characters"`   // 登场人物
	KeyConflict string   `json:"key_conflict"` // 核心冲突（爽点/转折）
	Outcome     string   `json:"outcome"`      // 结果（主角获得了什么或遭遇了什么）
}
