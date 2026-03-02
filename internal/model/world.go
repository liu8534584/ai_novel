package model

// WorldConfig represents the new structured world setting document.
type WorldConfig struct {
	Content string `json:"content"` // The raw markdown-like setting document
}

// Old structures (commented or kept for migration if needed)
/*
type PowerSystem struct {
	Tiers       []string `json:"tiers"`
	Currency    string   `json:"currency"`
	UpgradePath string   `json:"upgrade_path"`
}

type Character struct {
	Name         string         `json:"name"`
	Personality  string         `json:"personality"`
	Stats        map[string]int `json:"stats"`
	Inventory    []string       `json:"inventory"`
	Skills       []string       `json:"skills"`
	CheatAbility string         `json:"cheat_ability"`
}
*/
