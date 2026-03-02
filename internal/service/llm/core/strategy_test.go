package core

import (
	"testing"
)

func TestGetStrategy(t *testing.T) {
	tests := []struct {
		name     TaskType
		wantTemp float32
		wantMax  int
	}{
		{TaskWorldBuilding, 0.8, 4000},
		{TaskPlanning, 0.5, 4000},
		{TaskWriting, 0.9, 6000},
		{TaskReviewing, 0.3, 2000},
		{TaskStateTracking, 0.1, 500},
		{"unknown", 0.7, 1000},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			s := GetStrategy(tt.name)
			if s.Temperature != tt.wantTemp {
				t.Errorf("GetStrategy() Temperature = %v, want %v", s.Temperature, tt.wantTemp)
			}
			if s.MaxTokens != tt.wantMax {
				t.Errorf("GetStrategy() MaxTokens = %v, want %v", s.MaxTokens, tt.wantMax)
			}
		})
	}
}

func TestApplyToOptions(t *testing.T) {
	s := Strategy{Temperature: 0.5, MaxTokens: 100}
	
	t.Run("apply to empty options", func(t *testing.T) {
		opt := &Options{}
		s.ApplyToOptions(opt)
		if opt.Temperature != 0.5 {
			t.Errorf("ApplyToOptions() Temperature = %v, want 0.5", opt.Temperature)
		}
		if opt.MaxTokens != 100 {
			t.Errorf("ApplyToOptions() MaxTokens = %v, want 100", opt.MaxTokens)
		}
	})

	t.Run("do not override existing options", func(t *testing.T) {
		opt := &Options{Temperature: 0.8, MaxTokens: 200}
		s.ApplyToOptions(opt)
		if opt.Temperature != 0.8 {
			t.Errorf("ApplyToOptions() Temperature = %v, want 0.8", opt.Temperature)
		}
		if opt.MaxTokens != 200 {
			t.Errorf("ApplyToOptions() MaxTokens = %v, want 200", opt.MaxTokens)
		}
	})
}
