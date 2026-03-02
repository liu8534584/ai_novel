# InfinityWriter 后端开发指南 (Go + Gin)

既然前端已经有了展示“属性面板”和“流式内容”的容器，后端 Go (Gin) 接口的任务就是把 **智能体 (Agents)** 的逻辑串联起来，并以 **SSE (Server-Sent Events)** 的方式把数据喂给前端。

## 1. 路由设计 (Router)

在 `internal/handler/router.go` 中，我们需要定义以下核心端点：

```go
func RegisterRoutes(r *gin.Engine, h *NovelHandler) {
    api := r.Group("/api")
    {
        api.POST("/books", h.CreateBook)          // 触发 Director Agent：初始化世界和主角
        api.GET("/books/:id/state", h.GetState)   // 获取当前主角 HP/MP/物品 JSON
        
        api.POST("/chapters/:id/outline", h.GenerateOutline) // 触发 Outliner Agent：生成本章场景
        api.POST("/chapters/:id/write", h.WriteChapter)     // 触发 Writer + State Agent：流式写稿并更新状态
    }
}
```

## 2. 核心：流式写稿接口 (WriteChapter)

这是最复杂的一个接口。它需要：
1.  调用 `WriterAgent` 产生文字流。
2.  实时通过 **SSE 推送** 给前端。
3.  等文字结束时，调用 `StateAgent` 分析结果并存入数据库。

**Cursor 开发指令：**
> “请在 internal/handler/novel.go 中实现 WriteChapter 方法。
>
> 1. 设置响应头为 `text/event-stream`，支持 SSE。
> 2. 从数据库读取当前章节的大纲（Scenes）和书籍的 CurrentState。
> 3. 循环遍历每个 Scene：
>    * 调用 `WriterAgent.WriteSceneStream` 获得一个 channel。
>    * 通过 `c.Stream` 将文字片段实时发送给前端。
> 4. 闭环逻辑：在一个场景完成后，立即将生成的文本传给 `StateAgent.AnalyzeAndSyncState`。
> 5. 获取 StateAgent 返回的新状态 JSON，更新数据库中的 `books.current_state`。
> 6. 在流结束前，发送一个特殊的 `state_update` 消息给前端，让左侧面板更新。”

## 3. 总导演接口 (CreateBook)

这个接口负责从“灵感”到“结构化数据”的飞跃。

**Cursor 开发指令：**
> “请实现 CreateBook 处理器：
>
> 1. 接收 JSON body: `{"idea": "一句话灵感"}`。
> 2. 调用 `DirectorAgent.InitWorld(idea)`。
> 3. 将返回的 `WorldConfig` 存入数据库 books 表的 `world_setting` 和 `current_state` 字段。
> 4. 返回创建成功的书籍 ID 和初始化的 JSON 数据。”

## 4. 详细的系统整合逻辑 (The "Glue" Code)

在 Go 中，处理 JSON 与数据库的映射推荐使用 GORM。由于你要处理大量的动态 JSON（HP、物品等），这里有个技巧：

```go
// 建议在 model 中加入这个方法，方便合并 HP 变动
func (b *Book) ApplyStateUpdate(update model.StateUpdate) error {
    var stats map[string]int
    json.Unmarshal([]byte(b.CurrentState), &stats)

    // 1. 更新数值 (HP/MP/Exp)
    for key, delta := range update.StatsDelta {
        stats[key] += delta
    }

    // 2. 处理物品 (此处逻辑可由 Cursor 补全：append 或 remove)
    // ...

    // 3. 序列化回字符串
    newJSON, _ := json.Marshal(stats)
    b.CurrentState = string(newJSON)
    return nil
}
```

## 5. 如何在 Cursor 中一步步推进？

你可以直接按顺序复制以下三条指令给 Cursor：

### 第一步：基础框架与中间件
> “我正在用 Go 和 Gin 开发 AI 写作后端。请帮我写一个基础的 main.go。要求：
> 1. 引入 Gin 框架。
> 2. 添加 CORS 中间件，允许前端 Vue3 项目（通常是 localhost:5173）跨域访问。
> 3. 初始化 GORM 连接到 MySQL。”

### 第二步：实现 SSE 封装
> “由于 AI 写作需要流式显示，请在 internal/pkg/sse 中封装一个简单的 SSE 消息发送器。要求能方便地发送 text 类型（正文）和 json 类型（状态更新）的消息。”

### 第三步：串联 Agent 业务逻辑
> “现在编写 NovelService。这个 Service 需要组合之前写好的 DirectorAgent, OutlinerAgent, WriterAgent 和 StateAgent。 请实现一个核心方法 `ProcessChapter(chapterID uint)`： 它是这个系统的总控。逻辑是：获取大纲 -> 循环场景写稿 -> 实时分析状态 -> 更新数据库。 请注意错误处理，如果 LLM 接口超时，需要能捕获并返回给前端。”

## 🚀 为什么这套接口设计是“完全自动化”的？

因为你把决策权交给了数据：
1.  **前端不需要知道逻辑**，它只管订阅 `/api/chapters/:id/write`。
2.  **后端通过 StateAgent 的分析结果**，动态修改数据库里的数值。
3.  **下一次请求时**，Go 会自动从数据库取出最新的“残血”或“满级”状态喂给 AI。
