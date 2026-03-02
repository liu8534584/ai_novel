package model

// StateUpdate 描述了从一段文字中提取的状态变动
type StateUpdate struct {
	StatsDelta     map[string]int `json:"stats_delta"`     // 属性变化，如 {"HP": -10, "Exp": 50}
	NewItems       []string       `json:"new_items"`       // 新获得的物品
	RemovedItems   []string       `json:"removed_items"`   // 消耗或丢失的物品
	SkillUpgrades  []string       `json:"skill_upgrades"`  // 技能习得或升级
	PhysicalStatus string         `json:"physical_status"` // 身体状况描述（如：左臂骨折）
	PlotNodes      []string       `json:"plot_nodes"`      // 关键剧情点记录（用于后续大纲参考）
}
