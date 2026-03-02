# AI Novel (InfinityWriter)

AI 驱动的长篇小说辅助创作系统。

## 📁 项目结构

```text
.
├── frontend/           # Vue 3 前端项目 (Vite + Element Plus)
├── internal/           # 后端核心业务逻辑 (Go)
│   ├── config/         # 配置加载
│   ├── handler/        # API 处理器
│   └── service/        # AI Agent 与核心服务
├── models/             # 数据库模型
├── pkg/                # 公共工具包
├── main.go             # 后端入口
├── config.yaml         # 配置文件 (需要根据 config.example.yaml 创建)
└── Makefile            # 基础启动脚本
```

## 🛠️ 环境要求

- Go 1.20+
- Node.js 18+
- SQLite 3

## 🚀 快速开始

### 1. 后端启动

```bash
# 安装依赖
go mod tidy

# 运行后端
go run main.go
```

### 2. 前端启动

```bash
cd frontend
# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 3. 使用 Makefile (推荐)

```bash
# 启动后端服务
make dev-back

# 启动前端服务
make dev-front
```

## 📖 主要功能

- **分层上下文控制**: 确保长篇创作逻辑自洽。
- **角色 OOC 审计**: 实时监控角色性格崩坏。
- **伏笔追踪系统**: 自动记录并提醒未回收伏笔。
- **剧情矛盾检测**: 审计世界观与历史事实冲突。
