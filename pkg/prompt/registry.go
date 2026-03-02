package prompt

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

// PromptTemplate 定义了 Prompt 的结构
type PromptTemplate struct {
	Name    string
	Version string
	Content string
}

// Registry 负责管理所有的 Prompt 模板
type Registry struct {
	mu        sync.RWMutex
	templates map[string]PromptTemplate
}

var (
	defaultRegistry *Registry
	once            sync.Once
)

// GetRegistry 获取全局单例 Registry
func GetRegistry() *Registry {
	once.Do(func() {
		defaultRegistry = &Registry{
			templates: make(map[string]PromptTemplate),
		}
		// 初始化默认模板
		initDefaultTemplates(defaultRegistry)
	})
	return defaultRegistry
}

// Register 注册一个新的 Prompt 模板
func (r *Registry) Register(tmpl PromptTemplate) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.templates[tmpl.Name] = tmpl
}

// Get 获取指定名称的 Prompt 模板
func (r *Registry) Get(name string) (PromptTemplate, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tmpl, ok := r.templates[name]
	return tmpl, ok
}

// Render 渲染指定的 Prompt 模板
func (r *Registry) Render(name string, data interface{}) (string, error) {
	tmpl, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("prompt template %s not found", name)
	}

	t, err := template.New(name).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}

	return buf.String(), nil
}

// initDefaultTemplates 将硬编码的 Prompt 迁移至 Registry
func initDefaultTemplates(r *Registry) {
	r.Register(PromptTemplate{Name: "director", Version: "v1", Content: DirectorSystemPrompt})
	r.Register(PromptTemplate{Name: "planner", Version: "v1", Content: PlannerSystemPrompt})
	r.Register(PromptTemplate{Name: "planner_dark", Version: "v1", Content: PlannerDarkPrompt})
	r.Register(PromptTemplate{Name: "planner_growth", Version: "v1", Content: PlannerGrowthPrompt})
	r.Register(PromptTemplate{Name: "planner_twist", Version: "v1", Content: PlannerTwistPrompt})
	r.Register(PromptTemplate{Name: "character", Version: "v1", Content: CharacterSystemPrompt})
	r.Register(PromptTemplate{Name: "chapter_title", Version: "v1", Content: ChapterTitleSystemPrompt})
	r.Register(PromptTemplate{Name: "chapter_title_plan", Version: "v1", Content: ChapterTitlePlanSystemPrompt})
	r.Register(PromptTemplate{Name: "chapter_title_batch", Version: "v1", Content: BatchChapterTitleSystemPrompt})
	r.Register(PromptTemplate{Name: "chapter_title_batch_plan", Version: "v1", Content: BatchChapterTitlePlanSystemPrompt})
	r.Register(PromptTemplate{Name: "writer", Version: "v1", Content: WriterSystemPrompt})
	r.Register(PromptTemplate{Name: "outliner", Version: "v1", Content: OutlinerSystemPrompt})
	r.Register(PromptTemplate{Name: "summary", Version: "v1", Content: SummarySystemPrompt})
	r.Register(PromptTemplate{Name: "character_dynamic_state", Version: "v1", Content: CharacterDynamicStatePrompt})
	r.Register(PromptTemplate{Name: "chapter_objective", Version: "v1", Content: ChapterObjectivePrompt})
	r.Register(PromptTemplate{Name: "writer_layered", Version: "v1", Content: WriterLayeredSystemPrompt})
	r.Register(PromptTemplate{Name: "state_audit", Version: "v1", Content: StateAgentSystemPrompt})
	r.Register(PromptTemplate{Name: "event_extraction", Version: "v1", Content: EventExtractionPrompt})
	r.Register(PromptTemplate{Name: "foreshadowing_resolution", Version: "v1", Content: ForeshadowingResolutionPrompt})
	r.Register(PromptTemplate{Name: "character_anchor_extraction", Version: "v1", Content: CharacterAnchorExtractionPrompt})
	r.Register(PromptTemplate{Name: "character_anchor_audit", Version: "v1", Content: CharacterAnchorAuditPrompt})
	r.Register(PromptTemplate{Name: "ooc_evaluation", Version: "v1", Content: OOCEvaluationPrompt})
	r.Register(PromptTemplate{Name: "contradiction_detection", Version: "v1", Content: ContradictionDetectionPrompt})
	r.Register(PromptTemplate{Name: "inspiration_chat", Version: "v1", Content: InspirationChatSystemPrompt})
	r.Register(PromptTemplate{Name: "inspiration", Version: "v1", Content: InspirationSystemPrompt})
}
