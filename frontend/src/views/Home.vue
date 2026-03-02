<template>
  <div class="home">
    <section class="hero glass-card">
      <div class="hero-content">
        <div class="pill pill-primary">创作中枢</div>
        <h1>把世界观写进未来</h1>
        <p>统一管理灵感、人物、章节与剧情推进，实时同步世界状态与角色成长。</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" @click="openCreate">创建书籍</el-button>
          <el-button size="large" @click="openInspiration">灵感模式</el-button>
          <el-button size="large">导入灵感</el-button>
        </div>
        <div class="hero-tags">
          <span class="pill pill-primary" v-if="currentBook && currentBook.llm_config">
            模型：{{ currentBook.llm_config.model || '未配置' }}
          </span>
          <span class="pill pill-primary">今日目标：完成 2 章</span>
        </div>
      </div>
      <div class="hero-panel" v-if="currentBook">
        <div class="panel-card">
          <div class="panel-title">当前项目</div>
          <div class="panel-name">「{{ currentBook.title }}」</div>
          <div class="panel-meta">
            {{ currentBook.genre || '未分类' }} · {{ currentBook.total_chapters }} 章 · {{ statusLabel(currentBook.status) }}
          </div>
          <div class="panel-stats">
            <div>
              <div class="stat-value">{{ currentBook.chapters_count || 0 }}</div>
              <div class="stat-label">已创章节</div>
            </div>
            <div>
              <div class="stat-value">{{ currentBook.total_chapters }}</div>
              <div class="stat-label">目标章节</div>
            </div>
          </div>
          <el-progress :percentage="Math.round(((currentBook.chapters_count || 0) / currentBook.total_chapters) * 100)" color="#3b82f6" />
        </div>
      </div>
    </section>

    <section class="content-grid">
      <div class="column">
        <div class="glass-card section-card">
          <div class="section-header">
            <div>
              <h2>新建故事</h2>
              <p>输入灵感，让导演智能体生成世界观。</p>
            </div>
            <el-button type="primary" @click="openCreate">开始生成</el-button>
          </div>
          <el-input
            type="textarea"
            :rows="4"
            placeholder="写下你的故事灵感，例如：在冰封的巨型城市里，一位少年发现了会说话的陨石……"
          />
          <div class="prompt-suggestions">
            <span class="pill pill-primary">奇幻群像</span>
            <span class="pill pill-primary">克苏鲁都市</span>
            <span class="pill pill-primary">赛博修真</span>
          </div>
        </div>

        <div class="glass-card section-card">
          <div class="section-header">
            <div>
              <h2>最近项目</h2>
              <p>继续你最熟悉的世界。</p>
            </div>
            <el-button @click="openCreate">新建书籍</el-button>
          </div>
          <div class="project-list" v-loading="loading">
            <div v-if="books.length === 0" class="empty-state">
              <el-empty description="还没有书籍，立即创建" />
            </div>
            <div v-for="book in books" :key="book.id" class="project-item">
              <div>
                <div class="project-title">{{ book.title }}</div>
                <div class="project-meta">
                  {{ book.genre || '未分类' }} · {{ book.total_chapters }} 章 · {{ statusLabel(book.status) }}
                </div>
              </div>
              <div class="project-actions">
                <el-tag :type="statusType(book.status)">{{ statusLabel(book.status) }}</el-tag>
                <el-button size="small" type="primary" @click="goWrite(book)">写作</el-button>
                <el-button size="small" @click="goPlan(book)">规划</el-button>
                <el-button size="small" @click="openEdit(book)">编辑</el-button>
                <el-button size="small" type="danger" plain @click="deleteBook(book)">删除</el-button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="column">
        <div class="glass-card section-card" v-if="currentBook">
          <div class="section-header">
            <div>
              <h2>最近章节</h2>
              <p>该书籍的最新章节动态。</p>
            </div>
            <el-button @click="goWrite(currentBook)">进入写作</el-button>
          </div>
          <div class="state-grid">
            <div class="state-item">
              <div class="state-label">已完成章节</div>
              <div class="state-value">{{ currentBook.chapters_count || 0 }}</div>
            </div>
            <div class="state-item">
              <div class="state-label">总计划章节</div>
              <div class="state-value">{{ currentBook.total_chapters }}</div>
            </div>
          </div>
        </div>

        <div class="glass-card section-card">
          <div class="section-header">
            <div>
              <h2>创作灵感</h2>
              <p>记录并管理你的故事灵感。</p>
            </div>
            <el-button type="warning" plain @click="openCreate">新建书籍</el-button>
          </div>
          <div class="timeline">
            <div class="timeline-item" v-for="book in books.slice(0, 3)" :key="'time-'+book.id">
              <div class="timeline-dot" :class="{ active: book.status === 'writing' }"></div>
              <div>
                <div class="timeline-title">《{{ book.title }}》</div>
                <div class="timeline-meta">{{ statusLabel(book.status) }} · {{ book.genre }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <el-dialog v-model="dialogVisible" :title="isEditing ? '编辑书籍' : '创建书籍'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="书名">
          <el-input v-model="form.title" placeholder="请输入书名" />
        </el-form-item>
        <el-form-item label="题材">
          <el-input v-model="form.genre" placeholder="例如：科幻 / 玄幻 / 悬疑" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="故事简介或灵感摘要" />
        </el-form-item>
        <el-form-item label="章节数">
          <el-input-number v-model="form.total_chapters" :min="1" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status">
            <el-option label="草稿" value="draft" />
            <el-option label="写作中" value="writing" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">模型配置 (LLM)</el-divider>
        <el-form-item label="Provider">
          <el-select v-model="form.llm_config.provider" placeholder="选择提供商" @change="handleProviderChange">
            <el-option label="OpenAI" value="openai" />
            <el-option label="DeepSeek" value="deepseek" />
            <el-option label="GLM" value="glm" />
          </el-select>
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.llm_config.api_key" type="password" show-password placeholder="输入 API Key" />
        </el-form-item>
        <el-form-item label="Base URL">
          <el-input v-model="form.llm_config.base_url" placeholder="例如: https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="Model">
          <el-input v-model="form.llm_config.model" placeholder="例如: gpt-4o / deepseek-chat" />
        </el-form-item>

        <el-divider content-position="left">提示词模板绑定</el-divider>
        <el-form-item v-for="step in promptSteps" :key="step.key" :label="step.label">
          <el-select
            v-model="form.prompt_bindings[step.key]"
            multiple
            filterable
            collapse-tags
            collapse-tags-tooltip
            placeholder="选择一个或多个提示词模板"
          >
            <el-option v-for="item in promptOptions" :key="item.key" :label="item.title" :value="item.key" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveBook">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="inspirationDialogVisible" title="灵感模式" width="600px" custom-class="inspiration-dialog" @close="saveInspirationChat">
      <div class="chat-container">
        <div class="chat-messages" ref="chatScrollRef">
          <div v-for="(msg, index) in chatMessages" :key="index" :class="['chat-bubble', msg.role]">
            <div class="bubble-content">
              {{ msg.content }}
            </div>
          </div>
          <div v-if="isGenerating" class="chat-bubble assistant">
            <div class="bubble-content generating">
              <span>正在构思...</span>
              <el-icon class="is-loading"><Loading /></el-icon>
            </div>
          </div>
          <div v-if="inspirationResult" class="chat-bubble assistant">
            <div class="bubble-content">
              <div class="inspiration-result">
                <div class="result-header">✨ 已生成小说方案</div>
                <div class="result-item"><strong>书名：</strong>{{ inspirationResult.title }}</div>
                <div class="result-item"><strong>题材：</strong>{{ inspirationResult.genre.join(' / ') }}</div>
                <div class="result-item"><strong>核心主题：</strong>{{ inspirationResult.theme }}</div>
                <div class="result-item"><strong>简介：</strong>{{ inspirationResult.description }}</div>
                <el-button type="primary" size="small" @click="applyInspiration" style="margin-top: 10px;">
                  以此创建小说
                </el-button>
              </div>
            </div>
          </div>
        </div>
        <div class="chat-input-area">
          <el-input
            v-model="userInput"
            type="textarea"
            :rows="3"
            placeholder="与 AI 助手对话，完善你的灵感..."
            @keyup.enter.ctrl="sendInspiration"
          />
          <div class="input-actions">
            <div class="left">
              <span class="hint">Ctrl + Enter 发送</span>
            </div>
            <div class="right">
              <el-button type="warning" :loading="isFinalizing" :disabled="chatMessages.length < 2" @click="finalizeInspiration">
                生成小说方案
              </el-button>
              <el-button type="primary" :loading="isGenerating" @click="sendInspiration">
                发送
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'

type Book = {
  id: number
  title: string
  genre: string
  description: string
  total_chapters: number
  status: string
  created_at: string
  chapters_count?: number
  llm_config?: {
    provider: string
    api_key: string
    base_url: string
    model: string
  }
  prompt_bindings?: {
    world_view: string[]
    plan: string[]
    character: string[]
    chapter_title: string[]
    chapter_outline: string[]
    writing: string[]
    review: string[]
  }
}

const books = ref<Book[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const inspirationDialogVisible = ref(false)
const isEditing = ref(false)
const isGenerating = ref(false)
const isFinalizing = ref(false)
const router = useRouter()

const chatMessages = ref<Array<{ role: 'user' | 'assistant', content: string }>>([
  { role: 'assistant', content: '你好！我是你的小说创作顾问。告诉我一些你对新故事的想法，哪怕只是一个片段或一个概念，我会帮你把它完善成一个完整的方案。' }
])
const userInput = ref('')
const inspirationResult = ref<any>(null)
const chatScrollRef = ref<HTMLElement | null>(null)

const promptOptions = [
  { key: 'director', title: '世界观构建（导演）' },
  { key: 'planner', title: '剧情大纲（基础）' },
  { key: 'planner_dark', title: '剧情大纲（暗黑向）' },
  { key: 'planner_growth', title: '剧情大纲（成长向）' },
  { key: 'planner_twist', title: '剧情大纲（反转向）' },
  { key: 'character', title: '角色设定（基础）' },
  { key: 'chapter_title', title: '章节标题（全量）' },
  { key: 'chapter_title_plan', title: '章节标题（分阶段指令）' },
  { key: 'chapter_title_batch', title: '章节标题（分批生成）' },
  { key: 'chapter_title_batch_plan', title: '章节标题（指令版分批）' },
  { key: 'writer', title: '章节正文（标准写作）' },
  { key: 'outliner', title: '章节大纲（精细化）' },
  { key: 'summary', title: '摘要压缩（300-500字）' },
  { key: 'character_dynamic_state', title: '角色状态抽取' },
  { key: 'chapter_objective', title: '章节目标' },
  { key: 'writer_layered', title: '章节正文（分层写作）' },
  { key: 'state_audit', title: '状态审计' },
  { key: 'event_extraction', title: '关键事件抽取' },
  { key: 'foreshadowing_resolution', title: '伏笔回收判断' },
  { key: 'character_anchor_extraction', title: '性格锚点提取' },
  { key: 'ooc_evaluation', title: '角色 OOC 评估' },
  { key: 'contradiction_detection', title: '剧情矛盾检测' }
]

const createDefaultBindings = () => ({
  world_view: ['director'] as string[],
  plan: ['planner'] as string[],
  character: ['character'] as string[],
  chapter_title: ['chapter_title'] as string[],
  chapter_outline: ['outliner'] as string[],
  writing: ['writer_layered'] as string[],
  review: ['state_audit', 'event_extraction', 'foreshadowing_resolution', 'ooc_evaluation', 'contradiction_detection'] as string[]
})

type PromptBindingKey = keyof ReturnType<typeof createDefaultBindings>

const promptSteps: Array<{ key: PromptBindingKey; label: string }> = [
  { key: 'world_view', label: '世界观' },
  { key: 'plan', label: '大纲' },
  { key: 'character', label: '角色' },
  { key: 'chapter_title', label: '章节标题' },
  { key: 'chapter_outline', label: '章节大纲' },
  { key: 'writing', label: '正文写作' },
  { key: 'review', label: '复盘审计' }
]

const form = reactive({
  id: 0,
  title: '',
  genre: '',
  description: '',
  total_chapters: 1,
  status: 'draft',
  llm_config: {
    provider: 'openai',
    api_key: '',
    base_url: '',
    model: ''
  },
  prompt_bindings: createDefaultBindings()
})

const currentBook = computed(() => books.value[0] || null)

const scrollToBottom = async () => {
  await nextTick()
  if (chatScrollRef.value) {
    chatScrollRef.value.scrollTop = chatScrollRef.value.scrollHeight
  }
}
const loadInspirationChat = async () => {
  try {
    const res = await axios.get('/api/books/inspiration/chat')
    if (res.data.code === 0 && res.data.data) {
      const messages = JSON.parse(res.data.data)
      if (Array.isArray(messages) && messages.length > 0) {
        chatMessages.value = messages
        return true
      }
    }
  } catch (error) {
    console.error('Failed to load inspiration chat', error)
  }
  return false
}

const saveInspirationChat = async () => {
  try {
    await axios.post('/api/books/inspiration/chat/save', {
      messages: JSON.stringify(chatMessages.value)
    })
  } catch (error) {
    console.error('Failed to save inspiration chat', error)
  }
}

const openInspiration = async () => {
  inspirationDialogVisible.value = true
  inspirationResult.value = null
  const hasHistory = await loadInspirationChat()
  if (!hasHistory) {
    chatMessages.value = [
      { role: 'assistant', content: '你好！我是你的小说创作顾问。告诉我一些你对新故事的想法，哪怕只是一个片段或一个概念，我会帮你把它完善成一个完整的方案。' }
    ]
  }
  nextTick(() => {
    scrollToBottom()
  })
}

const sendInspiration = async () => {
  if (!userInput.value.trim() || isGenerating.value) return

  const userContent = userInput.value
  chatMessages.value.push({ role: 'user', content: userContent })
  userInput.value = ''
  isGenerating.value = true
  scrollToBottom()

  try {
    const res = await axios.post('/api/books/inspiration/chat', {
      messages: chatMessages.value
    })
    
    chatMessages.value.push({
      role: 'assistant',
      content: res.data.data
    })
    await saveInspirationChat()
  } catch (error) {
    ElMessage.error('获取灵感失败')
    chatMessages.value.push({ role: 'assistant', content: '抱歉，构思灵感时出了点问题。' })
  } finally {
    isGenerating.value = false
    scrollToBottom()
  }
}

const finalizeInspiration = async () => {
  if (isFinalizing.value) return
  isFinalizing.value = true
  
  // 将整个对话拼接成文本
  const conversation = chatMessages.value
    .map(m => `${m.role === 'user' ? '用户' : '顾问'}: ${m.content}`)
    .join('\n\n')

  try {
    const res = await axios.post('/api/books/inspiration/finalize', {
      conversation: conversation
    })
    
    let content = res.data.data
    if (typeof content === 'string') {
      const jsonMatch = content.match(/\{[\s\S]*\}/)
      if (jsonMatch) {
        inspirationResult.value = JSON.parse(jsonMatch[0])
      } else {
        inspirationResult.value = JSON.parse(content)
      }
    } else {
      inspirationResult.value = content
    }
    
    await saveInspirationChat()
    scrollToBottom()
  } catch (error) {
    ElMessage.error('加工灵感失败')
  } finally {
    isFinalizing.value = false
  }
}

const applyInspiration = async () => {
  const parsed = inspirationResult.value
  if (!parsed) return

  resetForm()
  await fetchGlobalLLMConfig()

  form.title = parsed.title
  form.genre = parsed.genre.join(' / ')
  form.description = parsed.description
  if (parsed.theme) {
    form.description += `\n\n核心主题：${parsed.theme}`
  }
  
  inspirationDialogVisible.value = false
  dialogVisible.value = true
}

const handleProviderChange = (val: string) => {
  if (val === 'deepseek') {
    form.llm_config.base_url = 'https://api.deepseek.com'
    form.llm_config.model = 'deepseek-chat'
  } else if (val === 'glm') {
    form.llm_config.base_url = 'https://open.bigmodel.cn/api/paas/v4/'
    form.llm_config.model = 'glm-4.7'
  } else if (val === 'openai') {
    form.llm_config.base_url = 'https://api.openai.com/v1'
    form.llm_config.model = 'gpt-4o-mini'
  }
}

const statusLabel = (status: string) => {
  if (status === 'completed') return '已完成'
  if (status === 'writing') return '写作中'
  return '草稿'
}

const statusType = (status: string) => {
  if (status === 'completed') return 'success'
  if (status === 'writing') return 'warning'
  return 'info'
}

const resetForm = () => {
  form.id = 0
  form.title = ''
  form.genre = ''
  form.description = ''
  form.total_chapters = 1
  form.status = 'draft'
  form.llm_config = {
    provider: 'openai',
    api_key: '',
    base_url: '',
    model: ''
  }
  form.prompt_bindings = createDefaultBindings()
}

const fetchGlobalLLMConfig = async () => {
  try {
    const res = await axios.get('/api/config/llm')
    if (res.data.code === 0 && res.data.data) {
      const config = res.data.data
      // 只有当全局配置有效时才覆盖
      if (config.provider) {
        form.llm_config.provider = config.provider
        form.llm_config.api_key = config.api_key || ''
        form.llm_config.base_url = config.base_url || ''
        form.llm_config.model = config.model || ''
      }
    }
  } catch (error) {
    console.error('Failed to fetch global LLM config', error)
  }
}

const fetchBooks = async () => {
  loading.value = true
  try {
    const res = await axios.get('/api/books')
    books.value = res.data.data || []
    // 自动刷新当前书籍的最新状态
    if (books.value.length > 0) {
      const firstBook = books.value[0]
      if (firstBook) {
        const stateRes = await axios.get(`/api/books/${firstBook.id}/state`)
        if (stateRes.data.code === 0 && stateRes.data.data) {
          Object.assign(firstBook, stateRes.data.data)
        }
      }
    }
  } catch (error) {
    ElMessage.error('获取书籍失败')
  } finally {
    loading.value = false
  }
}

const openCreate = async () => {
  isEditing.value = false
  resetForm()
  await fetchGlobalLLMConfig()
  dialogVisible.value = true
}

const openEdit = (book: Book) => {
  isEditing.value = true
  form.id = book.id
  form.title = book.title
  form.genre = book.genre
  form.description = book.description
  form.total_chapters = book.total_chapters
  form.status = book.status
  if (book.llm_config) {
    form.llm_config = { ...book.llm_config }
  } else {
    form.llm_config = {
      provider: 'openai',
      api_key: '',
      base_url: '',
      model: ''
    }
  }
  form.prompt_bindings = book.prompt_bindings ? { ...createDefaultBindings(), ...book.prompt_bindings } : createDefaultBindings()
  dialogVisible.value = true
}

const saveBook = async () => {
  const payload = {
    title: form.title,
    genre: form.genre,
    description: form.description,
    total_chapters: form.total_chapters,
    status: form.status,
    llm_config: form.llm_config,
    prompt_bindings: form.prompt_bindings
  }
  try {
    if (isEditing.value) {
      await axios.put(`/api/books/${form.id}`, payload)
      ElMessage.success('已更新书籍')
    } else {
      await axios.post('/api/books', payload)
      ElMessage.success('已创建书籍')
    }
    dialogVisible.value = false
    await fetchBooks()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const deleteBook = async (book: Book) => {
  try {
    await ElMessageBox.confirm(`确认删除《${book.title}》？`, '删除书籍', {
      type: 'warning'
    })
    await axios.delete(`/api/books/${book.id}`)
    ElMessage.success('已删除书籍')
    await fetchBooks()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const goPlan = (book: any) => {
  router.push(`/plans/${book.id}`)
}

const goWrite = (book: any) => {
  router.push(`/writing/${book.id}`)
}

onMounted(() => {
  fetchBooks()
})
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.hero {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(0, 1fr);
  gap: 24px;
  padding: 28px;
}

.hero-content h1 {
  font-size: 32px;
  margin: 12px 0;
}

.hero-content p {
  color: #475569;
  margin-bottom: 24px;
}

.hero-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.hero-tags {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.hero-panel {
  display: flex;
  align-items: stretch;
}

.panel-card {
  background: rgba(248, 250, 252, 0.85);
  border-radius: 18px;
  padding: 20px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.panel-title {
  color: #64748b;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.panel-name {
  font-size: 22px;
  font-weight: 700;
}

.panel-meta {
  color: #64748b;
  font-size: 13px;
}

.panel-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
}

.stat-label {
  font-size: 12px;
  color: #64748b;
}

.content-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 1fr);
  gap: 24px;
}

.column {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.section-card {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.section-header h2 {
  margin: 0 0 6px 0;
  font-size: 18px;
}

.section-header p {
  margin: 0;
  color: #64748b;
  font-size: 13px;
}

.prompt-suggestions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.project-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.project-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 14px;
  background: rgba(248, 250, 252, 0.9);
}

.project-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.empty-state {
  padding: 12px 0;
}

.project-title {
  font-weight: 600;
}

.project-meta {
  font-size: 12px;
  color: #64748b;
}

.state-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.state-item {
  background: rgba(248, 250, 252, 0.9);
  padding: 14px 16px;
  border-radius: 14px;
}

.state-label {
  font-size: 12px;
  color: #64748b;
}

.state-value {
  font-size: 16px;
  font-weight: 600;
  margin-top: 6px;
}

.timeline {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.timeline-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}

.timeline-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #cbd5f5;
  margin-top: 6px;
}

.timeline-dot.active {
  background: #f97316;
}

.timeline-title {
  font-weight: 600;
}

.timeline-meta {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}

@media (max-width: 1024px) {
  .hero {
    grid-template-columns: 1fr;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .hero-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .project-item {
    flex-direction: column;
    align-items: flex-start;
  }

  .project-actions {
    width: 100%;
    justify-content: flex-start;
    flex-wrap: wrap;
  }
}
.inspiration-dialog :deep(.el-dialog__body) {
  padding: 0;
}

.chat-container {
  display: flex;
  flex-direction: column;
  height: 500px;
  background: #f8fafc;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.chat-bubble {
  max-width: 85%;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 14px;
  line-height: 1.6;
}

.chat-bubble.user {
  align-self: flex-end;
  background: #2563eb;
  color: white;
  border-bottom-right-radius: 2px;
}

.chat-bubble.assistant {
  align-self: flex-start;
  background: white;
  color: #1e293b;
  border-bottom-left-radius: 2px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.bubble-content.generating {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #64748b;
}

.inspiration-result {
  background: #fff;
  border: 1px solid #e1e4e8;
  border-radius: 8px;
  padding: 15px;
  margin-top: 5px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.inspiration-result .result-header {
  font-weight: bold;
  color: #409eff;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px dashed #eee;
}

.result-item {
  font-size: 13px;
  line-height: 1.6;
  margin-bottom: 8px;
  color: #444;
}

.result-item strong {
  color: #111;
  width: 70px;
  display: inline-block;
}

.chat-input-area {
  padding: 15px;
  border-top: 1px solid #eee;
  background: #f9f9f9;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.input-actions .right {
  display: flex;
  gap: 10px;
}

.hint {
  font-size: 12px;
  color: #999;
}
</style>
