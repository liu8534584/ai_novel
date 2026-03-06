<template>
  <div class="app-shell">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">IW</div>
        <div class="brand-text">
          <div class="brand-title">InfinityWriter</div>
          <div class="brand-subtitle">AI Novel Studio</div>
        </div>
      </div>
      <nav class="nav">
        <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
          <el-icon><House /></el-icon>
          <span>工作台</span>
        </router-link>
        <router-link
          to="/plans"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/plans') }"
        >
          <el-icon><Notebook /></el-icon>
          <span>规划生成</span>
        </router-link>
        <router-link
          to="/outline"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/outline') }"
        >
          <el-icon><List /></el-icon>
          <span>总纲管理</span>
        </router-link>
        <router-link
          to="/prompts"
          class="nav-item"
          :class="{ active: $route.path.startsWith('/prompts') }"
        >
          <el-icon><Document /></el-icon>
          <span>提示词管理</span>
        </router-link>
        <router-link to="/settings" class="nav-item" :class="{ active: $route.path === '/settings' }">
          <el-icon><Setting /></el-icon>
          <span>模型设置</span>
        </router-link>
      </nav>
      <div class="sidebar-footer glass-card">
        <div class="footer-title">创作节奏</div>
        <div class="footer-metrics">
          <div class="metric">
            <div class="metric-value">3</div>
            <div class="metric-label">进行中</div>
          </div>
          <div class="metric">
            <div class="metric-value">14</div>
            <div class="metric-label">章节草稿</div>
          </div>
        </div>
        <el-button type="primary" size="small" class="footer-button">开始新的故事</el-button>
      </div>
    </aside>
    <main class="main">
      <header class="topbar glass-card">
        <div class="topbar-left">
          <div class="topbar-title">AI 写小说工作台</div>
          <div class="topbar-subtitle">沉浸式世界观搭建与章节生成</div>
        </div>
        <div class="topbar-actions">
          <el-input class="topbar-search" placeholder="搜索书籍 / 章节 / 灵感" />
          <el-button type="warning">生成新章节</el-button>
        </div>
      </header>
      <div class="content">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { Document, House, List, Notebook, Setting } from '@element-plus/icons-vue'
</script>

<style scoped>
.app-shell {
  display: grid;
  grid-template-columns: 260px 1fr;
  min-height: 100vh;
  gap: 24px;
  padding: 24px;
}

.sidebar {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
}

.brand-mark {
  width: 46px;
  height: 46px;
  border-radius: 16px;
  background: linear-gradient(135deg, #3b82f6, #60a5fa);
  color: #fff;
  display: grid;
  place-items: center;
  font-weight: 700;
  font-size: 18px;
}

.brand-title {
  font-weight: 700;
  font-size: 18px;
}

.brand-subtitle {
  font-size: 12px;
  color: #64748b;
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 14px;
  color: #334155;
  font-weight: 500;
  transition: background 0.2s ease, color 0.2s ease;
}

.nav-item.active {
  background: rgba(59, 130, 246, 0.12);
  color: #1d4ed8;
}

.nav-item:hover {
  background: rgba(148, 163, 184, 0.16);
}

.sidebar-footer {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.footer-title {
  font-size: 14px;
  color: #475569;
  font-weight: 600;
}

.footer-metrics {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.metric {
  background: rgba(248, 250, 252, 0.7);
  border-radius: 12px;
  padding: 12px;
}

.metric-value {
  font-size: 20px;
  font-weight: 700;
  color: #0f172a;
}

.metric-label {
  font-size: 12px;
  color: #64748b;
}

.footer-button {
  align-self: flex-start;
}

.main {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  gap: 16px;
}

.topbar-title {
  font-size: 20px;
  font-weight: 700;
}

.topbar-subtitle {
  font-size: 12px;
  color: #64748b;
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.topbar-search {
  width: 280px;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

@media (max-width: 1024px) {
  .app-shell {
    grid-template-columns: 1fr;
  }

  .sidebar {
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
  }

  .nav {
    flex-direction: row;
  }

  .sidebar-footer {
    display: none;
  }
}

@media (max-width: 768px) {
  .topbar {
    flex-direction: column;
    align-items: flex-start;
  }

  .topbar-actions {
    width: 100%;
    flex-direction: column;
    align-items: stretch;
  }

  .topbar-search {
    width: 100%;
  }
}
</style>
