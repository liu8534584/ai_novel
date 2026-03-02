AI 小说写作软件 · Trae Task Tree

技术栈：Go（后端） + Vue3（前端）
核心能力：多 LLM + 长期记忆（向量库）+ 可控写作流程

Phase 0：项目初始化 & 基础约定（必须先做）
Task 0.1 项目结构初始化

创建 Go 后端项目

创建 Vue3 + Vite 前端项目

前后端分离，通过 HTTP API 通信

Task 0.2 技术选型落地

后端：

Web 框架：Gin

ORM：GORM

配置：Viper

前端：

Vue3 + Composition API

UI：Element Plus / Naive UI

数据库：

主库：SQLite（MVP）→ 可切 MySQL

向量库：Qdrant（Docker）

Task 0.3 统一数据规范

所有 API 返回结构统一：

{
  "code": 0,
  "data": {},
  "message": ""
}

Phase 1：LLM 适配与管理（系统地基）
Task 1.1 LLM Provider 抽象层

定义统一接口：

type LLMProvider interface {
  Chat(ctx context.Context, messages []Message, options Options) (string, error)
}

Task 1.2 OpenAI 协议适配器

实现 OpenAI ChatCompletion

支持：

model

temperature

max_tokens

timeout

Task 1.3 多厂商 Adapter

DeepSeek（OpenAI 兼容）

GLM-4.7

Moonshot

Gemini（HTTP adapter）

Task 1.4 LLM 配置管理 API

CRUD：

新增配置

修改配置

删除配置

测试连接

支持设置默认模型

Task 1.5 前端：模型配置页

列表 + 编辑弹窗

测试按钮（返回是否成功）

Phase 2：书籍管理系统（写作容器）
Task 2.1 书籍数据模型
Book
- id
- title
- genre
- description
- total_chapters
- status
- created_at

Task 2.2 书籍管理 API

创建书籍

更新书籍

删除书籍

查询书籍列表

Task 2.3 前端：书籍列表页

卡片式展示

创建 / 编辑 / 删除

点击进入书籍工作区

Phase 3：小说规划生成（世界观 / 大纲 / 角色）
Task 3.1 Prompt 模板定义

世界观 Prompt

大纲 Prompt

角色 Prompt

Task 3.2 一次生成多版本规划

输入：

描述

类型

章节数

输出：

多个版本（如 3 个）

Task 3.3 规划版本管理
PlanVersion
- id
- book_id
- world_view
- outline
- characters
- is_selected

Task 3.4 规划选择与锁定

用户选中一个版本

设为当前生效版本

其余版本保留

Task 3.5 前端：规划生成页

多版本对比

选中 / 重生成

Phase 4：向量知识库（长期记忆）
Task 4.1 向量库初始化

部署 Qdrant

每本书一个 collection

Task 4.2 Embedding 服务

抽象 Embedding 接口

支持 OpenAI / DeepSeek embedding

Task 4.3 向量写入逻辑

世界观写入

大纲写入

角色写入

章节内容写入

Task 4.4 向量检索 API

Top-K 相似度搜索

支持按章节范围过滤

Phase 5：章节规划（章节标题生成）
Task 5.1 章节数据模型
Chapter
- id
- book_id
- index
- title
- status

Task 5.2 章节标题生成逻辑

输入：

世界观

大纲

角色

输出：

全部章节标题

Task 5.3 前端：章节管理页

列表展示

手动编辑标题

单章重生成

Phase 6：章节内容生成（核心写作）
Task 6.1 Prompt 组装器

拼接：

世界观摘要

大纲摘要

当前章节目标

向量召回内容

Task 6.2 章节内容生成 API

单章生成

批量生成

控制字数

Task 6.3 章节内容版本模型
ChapterVersion
- id
- chapter_id
- content
- created_at

Task 6.4 前端：写作编辑器

Markdown 编辑

版本切换

AI 重写 / 续写

Phase 7：写作流程控制
Task 7.1 章节状态流转

未生成 → 已生成 → 已确认

Task 7.2 上下文一致性校验

检查角色名

检查设定冲突（基础版）

Phase 8：UI & 体验升级（专业感）
Task 8.1 整体布局

左：书籍 & 章节树

中：编辑区

右：AI 控制区

Task 8.2 主题系统

浅色 / 深色

字体大小调节

Phase 9：扩展能力（可选）

导出 epub / txt

多语言写作

AI 校对 / 润色

伏笔追踪 / 角色状态表

给 Trae 的一句执行指令（可直接复制）

按 Phase 顺序逐步实现，每个 Phase 保证可运行、可测试，再进入下一阶段。

如果你愿意，下一步我可以直接帮你：

🔥 写一份 「Trae 专用 Master Prompt」

🧠 把 Prompt 模板（世界观 / 大纲 / 章节）一次性写好

🧱 给你 数据库表结构完整版（可直接建表）