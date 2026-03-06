CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS books (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  title text NOT NULL,
  author text,
  genre text,
  tags text,
  language text,
  description text,
  total_chapters integer DEFAULT 1,
  status text DEFAULT 'draft',
  world_setting jsonb,
  current_state jsonb,
  llm_config jsonb,
  prompt_bindings jsonb
);

CREATE INDEX IF NOT EXISTS idx_books_title ON books (title);
CREATE INDEX IF NOT EXISTS idx_books_genre ON books (genre);
CREATE INDEX IF NOT EXISTS idx_books_status ON books (status);
CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON books (deleted_at);

CREATE TABLE IF NOT EXISTS chapters (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  title text,
  content text,
  "order" integer,
  objective text,
  summary text,
  user_intent text,
  outline text,
  is_outline_confirmed boolean DEFAULT false,
  current_version integer DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_chapters_book_id ON chapters (book_id);
CREATE INDEX IF NOT EXISTS idx_chapters_deleted_at ON chapters (deleted_at);

CREATE TABLE IF NOT EXISTS chapter_versions (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  chapter_id bigint,
  version integer,
  title text,
  content text,
  word_count integer,
  summary text
);

CREATE INDEX IF NOT EXISTS idx_chapter_versions_chapter_id ON chapter_versions (chapter_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_chapter_versions_chapter_id_version ON chapter_versions (chapter_id, version);
CREATE INDEX IF NOT EXISTS idx_chapter_versions_deleted_at ON chapter_versions (deleted_at);

CREATE TABLE IF NOT EXISTS vector_records (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  chapter_id bigint,
  category text,
  content text,
  embedding text,
  metadata text
);

CREATE INDEX IF NOT EXISTS idx_vector_records_book_id ON vector_records (book_id);
CREATE INDEX IF NOT EXISTS idx_vector_records_chapter_id ON vector_records (chapter_id);
CREATE INDEX IF NOT EXISTS idx_vector_records_category ON vector_records (category);
CREATE INDEX IF NOT EXISTS idx_vector_records_deleted_at ON vector_records (deleted_at);

CREATE TABLE IF NOT EXISTS characters (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  name text,
  role text,
  description text,
  dynamic_state jsonb
);

CREATE INDEX IF NOT EXISTS idx_characters_book_id ON characters (book_id);
CREATE INDEX IF NOT EXISTS idx_characters_deleted_at ON characters (deleted_at);

CREATE TABLE IF NOT EXISTS character_state_records (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  character_id bigint,
  chapter_id bigint,
  state jsonb
);

CREATE INDEX IF NOT EXISTS idx_character_state_records_character_id ON character_state_records (character_id);
CREATE INDEX IF NOT EXISTS idx_character_state_records_chapter_id ON character_state_records (chapter_id);
CREATE INDEX IF NOT EXISTS idx_character_state_records_deleted_at ON character_state_records (deleted_at);

CREATE TABLE IF NOT EXISTS story_events (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  chapter_id bigint,
  chapter_index integer,
  event_type text,
  description text,
  involved_characters text,
  direct_consequence text,
  unresolved_impact text,
  importance integer
);

CREATE INDEX IF NOT EXISTS idx_story_events_book_id ON story_events (book_id);
CREATE INDEX IF NOT EXISTS idx_story_events_chapter_id ON story_events (chapter_id);
CREATE INDEX IF NOT EXISTS idx_story_events_deleted_at ON story_events (deleted_at);

CREATE TABLE IF NOT EXISTS outline_versions (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  version integer,
  world_view text,
  outline text,
  characters text,
  titles text,
  is_selected boolean DEFAULT false,
  is_locked boolean DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_outline_versions_book_id ON outline_versions (book_id);
CREATE INDEX IF NOT EXISTS idx_outline_versions_deleted_at ON outline_versions (deleted_at);

CREATE TABLE IF NOT EXISTS foreshadowings (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  chapter_id bigint,
  chapter_index integer,
  event_type text,
  description text,
  involved_characters text,
  direct_consequence text,
  unresolved_impact text,
  status text,
  importance integer,
  last_referenced_chapter integer,
  resolved_chapter_index integer,
  resolve_reason text
);

CREATE INDEX IF NOT EXISTS idx_foreshadowings_book_id ON foreshadowings (book_id);
CREATE INDEX IF NOT EXISTS idx_foreshadowings_chapter_id ON foreshadowings (chapter_id);
CREATE INDEX IF NOT EXISTS idx_foreshadowings_status ON foreshadowings (status);
CREATE INDEX IF NOT EXISTS idx_foreshadowings_deleted_at ON foreshadowings (deleted_at);

CREATE TABLE IF NOT EXISTS character_anchors (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  character_id bigint,
  personality_labels text,
  core_motivation text,
  behavior_bottom_line text,
  decision_tendency text,
  emotional_triggers text
);

CREATE INDEX IF NOT EXISTS idx_character_anchors_character_id ON character_anchors (character_id);
CREATE INDEX IF NOT EXISTS idx_character_anchors_deleted_at ON character_anchors (deleted_at);

CREATE TABLE IF NOT EXISTS ooc_scores (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  character_id bigint,
  chapter_id bigint,
  personality_consistency double precision,
  motivation_consistency double precision,
  emotional_reasonability double precision,
  cost_missing double precision,
  total_score double precision,
  conclusion text,
  explanation text
);

CREATE INDEX IF NOT EXISTS idx_ooc_scores_character_id ON ooc_scores (character_id);
CREATE INDEX IF NOT EXISTS idx_ooc_scores_chapter_id ON ooc_scores (chapter_id);
CREATE INDEX IF NOT EXISTS idx_ooc_scores_deleted_at ON ooc_scores (deleted_at);

CREATE TABLE IF NOT EXISTS story_contradictions (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  chapter_id bigint,
  type text,
  severity text,
  description text,
  reference text,
  suggestion text
);

CREATE INDEX IF NOT EXISTS idx_story_contradictions_book_id ON story_contradictions (book_id);
CREATE INDEX IF NOT EXISTS idx_story_contradictions_chapter_id ON story_contradictions (chapter_id);
CREATE INDEX IF NOT EXISTS idx_story_contradictions_deleted_at ON story_contradictions (deleted_at);

CREATE TABLE IF NOT EXISTS chapter_health_scores (
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  book_id bigint,
  chapter_id bigint,
  ooc_score double precision,
  event_consistency double precision,
  foreshadowing double precision,
  total_health double precision,
  audit_report text
);

CREATE INDEX IF NOT EXISTS idx_chapter_health_scores_book_id ON chapter_health_scores (book_id);
CREATE INDEX IF NOT EXISTS idx_chapter_health_scores_chapter_id ON chapter_health_scores (chapter_id);
CREATE INDEX IF NOT EXISTS idx_chapter_health_scores_deleted_at ON chapter_health_scores (deleted_at);

CREATE TABLE IF NOT EXISTS prompt_templates (
  id bigserial PRIMARY KEY,
  key text UNIQUE NOT NULL,
  title text NOT NULL,
  category text,
  description text,
  content text,
  updated_at timestamptz,
  source text,
  enabled boolean DEFAULT true
);

INSERT INTO prompt_templates (key, title, category, description, content, updated_at, source, enabled) VALUES
('director', '世界观构建（导演）', 'world', '', $director$
你是一名专业的长篇小说世界观架构师。

请基于以下信息，构建一个【可长期扩展、逻辑自洽、适合连载小说】的世界观。

【用户输入】
- 小说描述：{{.Description}}
- 小说类型：{{.Genre}}
- 预计章节数：{{.Chapters}}

【生成要求】
1. 世界观必须稳定，不允许后续频繁推翻
2. 设定要为剧情冲突服务，而不是空洞背景
3. 避免使用具体剧情细节（留给章节展开）
4. 不要写成小说正文，而是“设定文档”

【输出结构（必须严格遵守）】

一、世界整体设定
- 世界类型：
- 世界规则（物理 / 超自然 / 科技 / 逻辑）：
- 时代背景：
- 核心矛盾：

二、力量 / 系统 / 技术体系（如有）
- 系统来源：
- 运作规则：
- 限制与代价：

三、社会结构
- 主要势力：
- 阶级 / 阵营划分：
- 冲突关系：

四、核心主题
- 小说长期探讨的主题：
- 情绪基调：

五、世界观使用边界
- 允许出现的元素：
- 禁止出现的元素：

【输出格式】
使用清晰的小标题，不要使用任何对话或小说描写语言。
$director$, now(), '系统内置', true),
('planner', '剧情大纲（基础）', 'plan', '', $planner$
你是一名经验丰富的网文策划编辑，擅长设计【长篇连载小说】的整体剧情结构。

【已知世界观】
{{.WorldView}}

【用户输入】
- 小说描述：{{.Description}}
- 小说类型：{{.Genre}}
- 预计章节数：{{.Chapters}}

【生成目标】
生成一份【覆盖全书】的剧情大纲，适合分章节逐步展开。

【要求】
1. 大纲应分为 3~5 个“剧情阶段”
2. 每个阶段都有明确目标、冲突升级与阶段性变化
3. 不写具体章节内容，只写“剧情走向”
4. 保留悬念和伏笔空间

【输出结构（严格遵守）】

一、故事开端（约 {{.ChaptersBegin}} 章）
- 主线目标：
- 初始状态：
- 关键变化：

二、第一轮发展
- 剧情推进重点：
- 主要冲突：
- 世界或角色的变化：

三、中期转折
- 重大事件：
- 主角处境变化：
- 世界观层面的揭示：

四、后期升级
- 冲突升级点：
- 敌我变化：
- 代价与失去：

五、结局阶段
- 最终目标：
- 核心抉择：
- 世界状态收束方向：

【风格要求】
- 用策划语言
- 不要写任何具体对话
- 不要出现章节标题
$planner$, now(), '系统内置', true),
('planner_dark', '剧情大纲（暗黑向）', 'plan', '', $planner_dark$
你是一名擅长黑暗向与高风险叙事的网文策划编辑，强调代价、牺牲与道德灰度。

【已知世界观】
{{.WorldView}}

【用户输入】
- 小说描述：{{.Description}}
- 小说类型：{{.Genre}}
- 预计章节数：{{.Chapters}}

【生成目标】
生成一份偏黑暗走向的【覆盖全书】剧情大纲。

【要求】
1. 主线目标必须包含高代价或极限风险
2. 冲突应更尖锐，包含重大失败或牺牲节点
3. 氛围压迫、阴影感强，但逻辑自洽
4. 不写具体章节内容，只写“剧情走向”

【输出结构（严格遵守）】

一、故事开端（约 {{.ChaptersBegin}} 章）
- 主线目标：
- 初始状态：
- 关键变化：

二、第一轮发展
- 剧情推进重点：
- 主要冲突：
- 世界或角色的变化：

三、中期转折
- 重大事件：
- 主角处境变化：
- 世界观层面的揭示：

四、后期升级
- 冲突升级点：
- 敌我变化：
- 代价与失去：

五、结局阶段
- 最终目标：
- 核心抉择：
- 世界状态收束方向：

【风格要求】
- 用策划语言
- 不要写任何具体对话
- 不要出现章节标题
$planner_dark$, now(), '系统内置', true),
('planner_growth', '剧情大纲（成长向）', 'plan', '', $planner_growth$
你是一名擅长成长向与正向叙事的网文策划编辑，强调成长、希望与团队协作。

【已知世界观】
{{.WorldView}}

【用户输入】
- 小说描述：{{.Description}}
- 小说类型：{{.Genre}}
- 预计章节数：{{.Chapters}}

【生成目标】
生成一份偏成长向的【覆盖全书】剧情大纲。

【要求】
1. 主线目标体现持续成长与阶段性突破
2. 冲突升级以能力提升、关系深化与价值观变化为主
3. 保持积极向上的情绪基调，但仍有挑战
4. 不写具体章节内容，只写“剧情走向”

【输出结构（严格遵守）】

一、故事开端（约 {{.ChaptersBegin}} 章）
- 主线目标：
- 初始状态：
- 关键变化：

二、第一轮发展
- 剧情推进重点：
- 主要冲突：
- 世界或角色的变化：

三、中期转折
- 重大事件：
- 主角处境变化：
- 世界观层面的揭示：

四、后期升级
- 冲突升级点：
- 敌我变化：
- 代价与失去：

五、结局阶段
- 最终目标：
- 核心抉择：
- 世界状态收束方向：

【风格要求】
- 用策划语言
- 不要写任何具体对话
- 不要出现章节标题
$planner_growth$, now(), '系统内置', true),
('planner_twist', '剧情大纲（反转向）', 'plan', '', $planner_twist$
你是一名擅长悬疑与反转叙事的网文策划编辑，强调信息差、误导与关键反转。

【已知世界观】
{{.WorldView}}

【用户输入】
- 小说描述：{{.Description}}
- 小说类型：{{.Genre}}
- 预计章节数：{{.Chapters}}

【生成目标】
生成一份偏逆转向的【覆盖全书】剧情大纲。

【要求】
1. 主线推进依赖信息揭示与多次反转
2. 冲突升级应包含误导与真相揭露
3. 阶段性目标能因新情报被推翻或改写
4. 不写具体章节内容，只写“剧情走向”

【输出结构（严格遵守）】

一、故事开端（约 {{.ChaptersBegin}} 章）
- 主线目标：
- 初始状态：
- 关键变化：

二、第一轮发展
- 剧情推进重点：
- 主要冲突：
- 世界或角色的变化：

三、中期转折
- 重大事件：
- 主角处境变化：
- 世界观层面的揭示：

四、后期升级
- 冲突升级点：
- 敌我变化：
- 代价与失去：

五、结局阶段
- 最终目标：
- 核心抉择：
- 世界状态收束方向：

【风格要求】
- 用策划语言
- 不要写任何具体对话
- 不要出现章节标题
$planner_twist$, now(), '系统内置', true),
('character', '角色设定（基础）', 'character', '', $character$
你是一名小说人物设计师，专门为长篇小说设计“可成长角色”。

【世界观】
{{.WorldView}}

【故事大纲】
{{.Outline}}

【生成目标】
生成小说中的主要角色设定（3–6 名）。

【角色要求】
1. 每个角色都必须有：
   - 明确动机
   - 性格锚点（稳定的性格基石）
   - 行为底线
2. 主角必须具备“长期变化空间”
3. 配角应能推动剧情，而非工具人

【输出格式要求】
必须【仅输出】 JSON 对象格式，严禁包含任何前言、后记、解释性文字、或者是“Thinking Process”（思考过程）。
直接以 "{" 开头并以 "}" 结尾。
JSON 结构如下：
{
  "characters": [
    {
      "name": "角色名",
      "role": "Protagonist / Antagonist / Supporting",
      "description": "身份背景、初始状态与成长方向的综合描述（不可为空）",
      "anchor": {
        "personality_labels": ["标签1", "标签2"],
        "core_motivation": "核心驱动力",
        "behavior_bottom_line": "行为底线",
        "decision_tendency": "决策偏好",
        "emotional_triggers": "情感触发点"
      }
    }
  ]
}
$character$, now(), '系统内置', true),
('chapter_title', '章节标题（全量）', 'plan', '', $chapter_title$
你是一名网络小说章节策划。

【世界观】
{{.WorldView}}

【故事大纲】
{{.Outline}}

【角色设定】
{{.Characters}}

【目标】
为整本小说生成章节标题。

【要求】
1. 共 {{.Chapters}} 章
2. 标题要有推进感，而不是重复描述
3. 避免剧透结局
4. 标题风格统一

【输出格式（必须严格遵守）】

第1章：标题
第2章：标题
……
第{{.Chapters}}章：标题
$chapter_title$, now(), '系统内置', true),
('chapter_title_plan', '章节标题（分阶段指令）', 'plan', '', $chapter_title_plan$
你是一名网络小说章节策划总监。

【故事大纲】
{{.Outline}}

【总章节数】
{{.Chapters}}

【任务】
将大纲整理为“章节标题生成指令”，用于后续分批生成标题。

【要求】
1. 输出必须是面向“生成章节标题”的指令性内容，而不是标题本身
2. 按剧情阶段划分，明确每个阶段对应的大致章节范围（例如：第1-20章）
3. 每个阶段给出关键词、核心冲突、推进目标
4. 衔接顺序清晰，避免跳跃
5. 字数控制在 400-700 字

【输出结构（必须严格遵守）】
一、整体标题风格与基调
二、阶段划分与生成指令
1) 阶段名（章节范围）
   - 关键词：
   - 核心冲突：
   - 推进目标：
2) ...
$chapter_title_plan$, now(), '系统内置', true),
('chapter_title_batch', '章节标题（分批生成）', 'plan', '', $chapter_title_batch$
你是一名网络小说章节策划。

【世界观】
{{.WorldView}}

【故事大纲】
{{.Outline}}

【角色设定】
{{.Characters}}

【当前进度】
- 已有章节数：{{.CurrentCount}}
- 本次生成：从第 {{.StartChapter}} 章开始，生成 {{.BatchSize}} 章

【参考前文章节标题】
{{.PreviousTitles}}

【目标】
为小说生成后续章节标题。

【要求】
1. 从第 {{.StartChapter}} 章开始，连续生成 {{.BatchSize}} 章
2. 标题要紧扣【故事大纲】中对应的剧情阶段
3. 标题要有推进感，承接前文，逻辑连贯
4. 标题风格与前文保持一致

【输出格式（必须严格遵守）】

第{{.StartChapter}}章：标题
...
第{{.EndChapter}}章：标题
$chapter_title_batch$, now(), '系统内置', true),
('chapter_title_batch_plan', '章节标题（指令版分批）', 'plan', '', $chapter_title_batch_plan$
你是一名网络小说章节策划。

【世界观】
{{.WorldView}}

【角色设定】
{{.Characters}}

【章节标题生成指令】
{{.TitlePlan}}

【当前进度】
- 已有章节数：{{.CurrentCount}}
- 本次生成：从第 {{.StartChapter}} 章开始，生成 {{.BatchSize}} 章

【参考前文章节标题】
{{.PreviousTitles}}

【目标】
为小说生成后续章节标题。

【要求】
1. 从第 {{.StartChapter}} 章开始，连续生成 {{.BatchSize}} 章
2. 标题必须紧扣“章节标题生成指令”的阶段划分与推进目标
3. 标题要有推进感，承接前文，逻辑连贯
4. 标题风格统一，避免重复或空泛

【输出格式（必须严格遵守）】

第{{.StartChapter}}章：标题
...
第{{.EndChapter}}章：标题
$chapter_title_batch_plan$, now(), '系统内置', true),
('writer', '章节正文（标准写作）', 'writing', '', $writer$
你是一名正在连载长篇小说的作者。

【世界观】
{{.WorldView}}

【故事大纲（摘要）】
{{.OutlineSummary}}

【角色设定（摘要）】
{{.CharacterSummary}}

【当前章节】
- 章节序号：第 {{.ChapterIndex}} 章
- 章节标题：{{.ChapterTitle}}

【已发生的重要内容（供参考）】
{{.RetrievedContext}}

【写作要求】
1. 本章必须承接前文，不可自相矛盾
2. 推进剧情，而不是重复背景
3. 展现角色行动与选择
4. 避免总结性语言
5. **字数控制**：本章目标字数为 {{.TargetWords}} 字左右，请通过丰富的动作、心理描写和对话来充实内容，避免节奏过快。
6. **分段与节奏**：保持合理的段落长度，通过环境渲染和细节描写控制故事节奏。

【字数要求】
约 {{.TargetWords}} 字

【禁止】
- 不要解释设定
- 不要回顾全文
- 不要使用“这一章我们将会”

【输出】
直接输出小说正文，不要任何额外说明。
$writer$, now(), '系统内置', true),
('outliner', '章节大纲（精细化）', 'plan', '', $outliner$
你是一个顶级的网文大纲架构师，擅长在保持全书逻辑高度一致性的前提下，规划引人入胜的章节剧情。

【核心上下文】
1. **世界观设定**: {{.WorldSetting}}
2. **核心角色设定**: {{.Characters}}
3. **全书主线大纲**: {{.StoryOutline}}
4. **全书目录规划**: {{.AllTitles}}
5. **本章预定标题**: {{.ChapterTitle}} (第 {{.ChapterNum}} 章)
6. **前情提要**: {{.PrevSummary}}
7. **用户特定写作意图**: {{.UserIntent}}

【任务】
请为第 {{.ChapterNum}} 章设计详细大纲。

【大纲生成逻辑指南】
- **紧扣主线**：定位“全书主线大纲”中对应的剧情阶段。本章内容必须是主线大纲中该阶段的具象化展开。
- **标题契合**：大纲内容必须紧紧围绕“本章预定标题”展开，标题即是本章的灵魂和核心事件。
- **前文衔接**：
    - 如果“前情提要”不为空，本章必须逻辑严密地承接前文。
    - **特别注意**：如果“前情提要”为空或为“无”，说明这是开篇或新阶段的开始。此时，你必须完全依据“全书主线大纲”的开篇设定和“本章预定标题”来启动剧情，确保开篇即进入主线。
- **人物驱动**：所有行动必须符合“核心角色设定”，利用人物性格推动剧情走向。
- **世界观约束**：严禁出现违反“世界观设定”的力量体系、社会规则或背景元素。

【约束要求】
1. **起承转合**：本章大纲需包含明确的起因、发展、高潮（冲突）和结果。
2. **场景细化**：将本章拆分为 3-4 个画面感强、节奏紧凑的具体场景。
3. **JSON 格式**：严格输出 JSON。

【输出结构】
{
  "chapter_number": {{.ChapterNum}},
  "title": "{{.ChapterTitle}}",
  "summary": "本章剧情梗概 (需体现如何承接主线并引出本章冲突)",
  "scenes": [
    {
      "order": 1,
      "location": "具体地点",
      "description": "详细动作与情节描述 (包含登场人物及其核心行动)",
      "characters": ["人物1", "人物2"],
      "key_conflict": "本场景的核心矛盾点",
      "outcome": "本场景的产出结果及对后文的铺垫"
    }
  ]
}
$outliner$, now(), '系统内置', true),
('summary', '摘要压缩（300-500字）', 'review', '', $summary$
你是一名专业的文学编辑，擅长从繁杂的设定中提取核心逻辑。

【任务】
请将以下【{{.Type}}】压缩为一份摘要版（300-500字）。

【要求】
1. 只保留核心规则、主要矛盾、关键背景或剧情走向。
2. 剔除细节描述和非必要的修饰词。
3. 确保摘要版能让作者在写作时快速回顾，不产生歧义。
4. 结构清晰，使用简明的小标题。

【待压缩内容】
{{.Content}}
$summary$, now(), '系统内置', true),
('character_dynamic_state', '角色状态抽取', 'review', '', $character_dynamic_state$
你是一名小说编辑，负责维护“角色状态记录表”，用于保证长篇小说中角色行为和设定的一致性。 

【世界观（摘要）】 
{{.WorldViewSummary}} 

【角色初始设定】 
{{.CharacterBaseProfiles}} 

【上一章前的角色状态】 
{{.PreviousCharacterStates}} 

【本章正文】 
{{.ChapterContent}} 

【任务目标】 
从【本章正文】中，提取所有主要角色在本章中发生的【状态变化】。 

【抽取原则】 
1. 只记录对后续剧情有约束力的变化 
2. 如果某个字段没有变化，请写“无变化” 
3. 不要重复角色初始设定 
4. 不要复述剧情细节 
5. 用客观、简洁的编辑语言 

【输出要求（严格遵守）】 
请严格按照以下 JSON 格式输出，以便系统解析。Key 为角色名，Value 为状态对象。

{
  "角色名": {
    "identity_location": "当前身份 / 位置",
    "goal": "当前目标",
    "emotional_state": "当前情绪状态",
    "relationship_changes": "当前关系变化",
    "ability_resource_changes": "能力 / 资源变化",
    "constraints_costs": "新增限制或代价",
    "key_actions": "本章关键行为",
    "conflicts_foreshadowing": "潜在矛盾 / 伏笔"
  }
}

【特别注意】 
- 如果角色在本章未出现，不要输出该角色 
- 禁止使用小说语言 
- 禁止添加推测性的心理描写 
$character_dynamic_state$, now(), '系统内置', true),
('chapter_objective', '章节目标', 'plan', '', $chapter_objective$
你是一名网络小说架构师，擅长规划章节的承载功能。

【任务】
基于大纲和当前进度，为指定章节设计【章节目标】。

【大纲摘要】
{{.OutlineSummary}}

【当前章节】
- 章节序号：第 {{.ChapterIndex}} 章
- 章节标题：{{.ChapterTitle}}

【设计要求】
1. 明确本章的作用：
   - 推进哪个具体的冲突？
   - 服务哪条主线或支线？
   - 角色将面临什么抉择？
2. 目标应具体且具有可操作性，指导 AI 写作。
3. 字数在 150 字左右。

【输出格式】
直接输出目标描述，不要带标题或额外文字。
$chapter_objective$, now(), '系统内置', true),
('writer_layered', '章节正文（分层写作）', 'writing', '', $writer_layered$
[系统设定] 你是网文架构师。

[当前环境]
{{.WorldSummary}}

[核心大纲]
{{.OutlineSummary}}

[角色状态]
{{.CharacterStates}}

[前情回顾]
{{.RetrievedMemories}}

[本次任务]
请基于以上信息，生成第 {{.ChapterIndex}} 章内容。
章节标题：{{.ChapterTitle}}
章节目标：{{.ChapterObjective}}
未回收伏笔：{{.Foreshadowing}}
目标字数：{{.TargetWords}} 字

要求：
1. 严禁 OOC，角色行为必须符合其性格和当前状态。
2. 细节丰富，通过动作、心理和对话描写而非总结性叙述来推进剧情。
3. 严格承接前文，确保逻辑自洽，自然埋下或回收伏笔。
4. 保持文风一致，直接输出小说正文，不要任何额外说明。
$writer_layered$, now(), '系统内置', true),
('state_audit', '状态审计', 'review', '', $state_audit$
你是一个高精度的故事逻辑审计员和 RPG 游戏后台管理员。

【任务】
阅读提供的【小说片段】，对比【当前状态 JSON】，识别出角色属性、物品、技能以及物理状态的任何变化。

【当前状态 JSON】
{{.CurrentState}}

【输入文本（刚写的章节）】
{{.ChapterContent}}

【输出要求】
1. 客观性：只记录文字中明确发生的变动。
2. 格式：必须严格输出 JSON。

{
  "stats_delta": {"属性名": 增量值},
  "new_items": ["物品A"],
  "removed_items": ["物品B"],
  "skill_upgrades": ["技能名"],
  "physical_status": "当前生理状态描述",
  "plot_nodes": ["关键转折点"]
}
$state_audit$, now(), '系统内置', true),
('event_extraction', '关键事件抽取', 'review', '', $event_extraction$
你是一名小说档案管理员，负责记录故事中发生的【关键事件】。

【本章正文】
{{.ChapterContent}}

【任务】
请从正文中提取所有关键事件。

【事件定义】
- 必须是对后续剧情有显著影响的动作、揭示或冲突。
- 排除日常琐碎描写。

【输出格式（严格 JSON）】
[
  {
    "event_type": "主线推进/冲突升级/世界规则揭示/角色转折",
    "description": "简明描述事件经过",
    "involved_characters": "张三, 李四",
    "direct_consequence": "该事件的直接结果",
    "unresolved_impact": "该事件留下的未解决影响（伏笔核心），若无则留空",
    "importance": 1-5
  }
]
$event_extraction$, now(), '系统内置', true),
('foreshadowing_resolution', '伏笔回收判断', 'review', '', $foreshadowing_resolution$
你是一名小说策划编辑，负责核对伏笔是否已被回收。

【未回收伏笔】
- 描述：{{.ForeshadowingDescription}}
- 未解决影响：{{.UnresolvedImpact}}

【本章关键事件】
{{.CurrentEvents}}

【本章角色状态变化】
{{.CharacterStates}}

【任务】
判断该伏笔是否在本章中被明确解决。

【解决标准】
- 伏笔中提到的“未解决影响”已被消除或达成。
- 角色状态变化或新事件直接终结了该因果链。

【输出格式】
必须返回 JSON 格式：
{
  "is_resolved": bool,
  "reason": "简述判定依据，如果已解决请描述解决的具体情节"
}
$foreshadowing_resolution$, now(), '系统内置', true),
('character_anchor_extraction', '性格锚点提取', 'review', '', $character_anchor_extraction$
你是一名小说人物设计师，负责为角色建立“性格锚点”。

【角色初始设定】
{{.CharacterDescription}}

【前10-20章角色状态记录】
{{.CharacterStates}}

【任务】
请基于以上信息，提取该角色的性格锚点。性格锚点应保持相对稳定，作为评估后续行为是否 OOC 的基准。

【输出格式（严格遵守 JSON）】
{
  "personality_labels": "核心性格标签 (3–5 个)",
  "core_motivation": "核心动机 (长期)",
  "behavior_bottom_line": "行为底线 (绝不做的事)",
  "decision_tendency": "决策倾向 (保守 / 激进 / 利己 / 利他)",
  "emotional_triggers": "情绪触发点"
}
$character_anchor_extraction$, now(), '系统内置', true),
('ooc_evaluation', '角色 OOC 评估', 'review', '', $ooc_evaluation$
你是一名小说人物编辑，负责评估角色是否发生“性格崩坏（OOC）”。

【角色性格锚点】
{{.CharacterAnchor}}

【角色历史状态摘要】
{{.CharacterStateHistory}}

【本章中该角色的行为与情绪】
{{.CurrentCharacterBehavior}}

【任务】
评估该角色在本章中的表现，是否偏离其已建立的性格与动机。

【评分维度】
1. 性格一致性
2. 动机一致性
3. 情绪反应合理性
4. 行为是否缺乏应有代价

【输出格式（严格遵守 JSON）】
{
  "personality_consistency": 0,
  "motivation_consistency": 0,
  "emotional_reasonability": 0,
  "cost_missing": 0,
  "total_score": 0,
  "conclusion": "无明显 OOC / 轻度 OOC / 明显 OOC / 严重 OOC",
  "explanation": "简要说明理由，不超过 3 句话"
}

【限制】
- 不要给写作建议
- 不要重写剧情
- 使用编辑判断语言
- 性格一致性评分：偏离越多，分数越高
- 动机一致性评分：偏离越多，分数越高
- 情绪反应合理性评分：越不合理，分数越高
- 行为代价缺失评分：越缺乏代价，分数越高
$ooc_evaluation$, now(), '系统内置', true),
('contradiction_detection', '剧情矛盾检测', 'review', '', $contradiction_detection$
你是一名小说逻辑审计员，负责检测新章节内容是否与已确立的“故事事实”或“世界规则”相冲突。

【世界观规则】
{{.WorldRules}}

【关键历史事实（已发生事件）】
{{.HistoryEvents}}

【主要角色当前状态】
{{.CharacterStates}}

【本章正文】
{{.ChapterContent}}

【检测任务】
请识别本章正文是否包含以下矛盾：
1. 逻辑冲突：违背了之前章节中明确发生的事件结果（如：已死角色复活、丢失物品复现、地理位置瞬间跨越）。
2. 设定冲突：违背了世界观中的物理/超自然/社会规则（如：在禁魔区施法、无视等级压制）。
3. 状态冲突：角色表现出了其当前能力、资源或生理状态无法支持的行为。

【输出格式（严格遵守 JSON 列表）】
请输出一个 JSON 列表，每个对象代表一个矛盾点：
[
  {
    "type": "逻辑冲突 / 设定冲突 / 状态冲突",
    "severity": "low / medium / high",
    "description": "矛盾点的具体描述",
    "reference": "冲突所涉及的历史事实或规则条目",
    "suggestion": "建议的修正方向"
  }
]
若无矛盾，请输出空列表 []。

【限制】
- 若无冲突，请输出空数组 []
- 使用客观、严谨的审计语言
- 不要评价文笔，只关注事实一致性
$contradiction_detection$, now(), '系统内置', true)
ON CONFLICT (key) DO UPDATE SET
  title = EXCLUDED.title,
  category = EXCLUDED.category,
  description = EXCLUDED.description,
  content = EXCLUDED.content,
  updated_at = EXCLUDED.updated_at,
  source = EXCLUDED.source,
  enabled = EXCLUDED.enabled;
