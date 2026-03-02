package core

type TaskType string

const (
	TaskWorldBuilding   TaskType = "world_building"
	TaskPlanning        TaskType = "planning"
	TaskCharacterDesign TaskType = "character_design"
	TaskWriting         TaskType = "writing"
	TaskReviewing       TaskType = "reviewing"
	TaskStateTracking   TaskType = "state_tracking"
)

type Strategy struct {
	Temperature float32
	MaxTokens   int
}

var strategyTable = map[TaskType]Strategy{
	TaskWorldBuilding: {
		Temperature: 0.8,
		MaxTokens:   4000,
	},
	TaskPlanning: {
		Temperature: 0.5,
		MaxTokens:   4000,
	},
	TaskCharacterDesign: {
		Temperature: 0.7,
		MaxTokens:   2000,
	},
	TaskWriting: {
		Temperature: 0.9,
		MaxTokens:   6000,
	},
	TaskReviewing: {
		Temperature: 0.3,
		MaxTokens:   2000,
	},
	TaskStateTracking: {
		Temperature: 0.1,
		MaxTokens:   500,
	},
}

// GetStrategy returns the strategy for a given task type.
// If the task type is unknown, it returns a default strategy.
func GetStrategy(task TaskType) Strategy {
	if s, ok := strategyTable[task]; ok {
		return s
	}
	// Default strategy
	return Strategy{
		Temperature: 0.7,
		MaxTokens:   1000,
	}
}

// ApplyToOptions applies the strategy to the given options.
func (s Strategy) ApplyToOptions(opt *Options) {
	if opt.Temperature == 0 {
		opt.Temperature = s.Temperature
	}
	if opt.MaxTokens == 0 {
		opt.MaxTokens = s.MaxTokens
	}
}
