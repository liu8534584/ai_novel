<template>
  <div class="master-outline" v-if="selectedBookId">
    <!-- 顶部提示栏 -->
    <div class="outline-banner glass-card" v-if="blueprints.length > 0">
      <span class="banner-icon">📋</span>
      <span>编辑最看重开头了！放心，我帮你把前{{ blueprints.length }}章的情绪点和钩子都安排好了。</span>
    </div>

    <div class="outline-container">
      <!-- 左侧章节列表 -->
      <aside class="chapter-sidebar glass-card">
        <div class="sidebar-header" @click="showBookInfo = !showBookInfo">
          <span class="sidebar-title">{{ book?.title || '作品简介' }}</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <div class="book-info" v-if="showBookInfo && book">
          <p class="book-desc">{{ book.description }}</p>
          <el-tag size="small" type="info">{{ book.genre }}</el-tag>
        </div>

        <div class="chapter-list">
          <div
            v-for="bp in blueprints"
            :key="bp.chapter_index"
            class="chapter-item"
            :class="{ active: selectedChapter === bp.chapter_index }"
            @click="selectChapter(bp.chapter_index)"
          >
            <span class="chapter-label">第{{ bp.chapter_index }}章</span>
            <span class="chapter-title">{{ bp.title || '未命名' }}</span>
          </div>
        </div>

        <div class="sidebar-actions">
          <el-button type="primary" size="small" @click="generateBatch" :loading="generating" :disabled="!selectedBookId">
            {{ blueprints.length > 0 ? '继续生成' : '开始生成总纲' }}
          </el-button>
        </div>
      </aside>

      <!-- 右侧详情面板 -->
      <main class="detail-panel glass-card" v-if="currentBlueprint">
        <div class="detail-header">
          <h2 class="detail-title">{{ currentBlueprint.title || '未命名章节' }}</h2>
          <div class="auto-save-toggle">
            <span class="save-label">自动保存{{ autoSave ? '已开启' : '已关闭' }}</span>
            <el-switch v-model="autoSave" size="small" />
          </div>
        </div>

        <div class="detail-form">
          <div class="form-row">
            <label class="form-label">主角动机：</label>
            <el-input
              v-model="editData.protagonist_motivation"
              placeholder="本章主角的核心驱动力"
              @change="onFieldChange('protagonist_motivation')"
            />
          </div>
          <div class="form-row">
            <label class="form-label">关键伏笔：</label>
            <el-input
              v-model="editData.key_foreshadowing"
              placeholder="本章埋下或需回收的伏笔"
              @change="onFieldChange('key_foreshadowing')"
            />
          </div>
          <div class="form-row">
            <label class="form-label">出场人物：</label>
            <el-input
              v-model="editData.appearing_characters"
              placeholder="楚天阳、苏清月（背景提及）、墨无痕（暗线）"
              @change="onFieldChange('appearing_characters')"
            />
          </div>
          <div class="form-row">
            <label class="form-label">审视亮点：</label>
            <el-input
              v-model="editData.highlight"
              placeholder="逆袭打脸爽感+势力暗流初现，节奏张弛有度"
              @change="onFieldChange('highlight')"
            />
          </div>
          <div class="form-row">
            <label class="form-label">核心事件：</label>
            <el-input
              v-model="editData.core_events"
              placeholder="本章最重要的情节转折"
              @change="onFieldChange('core_events')"
            />
          </div>
          <div class="form-row">
            <label class="form-label">面临挑战：</label>
            <el-input
              v-model="editData.challenges"
              placeholder="主角本章面临的困难或冲突"
              @change="onFieldChange('challenges')"
            />
          </div>
          <div class="form-row form-row-textarea">
            <label class="form-label">章节大纲：</label>
            <el-input
              v-model="editData.summary"
              type="textarea"
              :rows="4"
              placeholder="详细描述本章剧情走向"
              @change="onFieldChange('summary')"
            />
          </div>
        </div>

        <div class="detail-actions">
          <el-button type="warning" @click="regenerateChapter" :loading="regenerating">
            重新生成本章
          </el-button>
        </div>
      </main>

      <!-- 空状态 -->
      <main class="detail-panel glass-card empty-panel" v-else>
        <el-empty :description="blueprints.length === 0 ? '还没有生成总纲，请点击左侧按钮开始' : '请从左侧选择一个章节'">
          <el-button v-if="blueprints.length === 0" type="primary" @click="generateBatch" :loading="generating">
            生成总纲
          </el-button>
        </el-empty>
      </main>
    </div>
  </div>

  <!-- 未选择书籍 -->
  <div class="master-outline" v-else>
    <section class="outline-hero glass-card">
      <div>
        <div class="pill pill-primary">总纲管理</div>
        <h1>章节蓝图编排</h1>
        <p>一键生成每章的标题、大纲、人物、伏笔与核心事件。</p>
      </div>
    </section>
    <div class="glass-card" style="padding: 24px">
      <el-form label-width="60px">
        <el-form-item label="书籍">
          <el-select v-model="selectedBookId" placeholder="请选择书籍" @change="handleBookChange">
            <el-option v-for="b in books" :key="b.id" :label="b.title" :value="b.id" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

type Book = {
  id: number
  title: string
  genre: string
  description: string
  total_chapters: number
}

type Blueprint = {
  id: number
  book_id: number
  chapter_index: number
  title: string
  summary: string
  protagonist_motivation: string
  key_foreshadowing: string
  appearing_characters: string
  highlight: string
  core_events: string
  challenges: string
}

const route = useRoute()
const router = useRouter()

const books = ref<Book[]>([])
const book = ref<Book | null>(null)
const blueprints = ref<Blueprint[]>([])
const selectedBookId = ref<number | null>(null)
const selectedChapter = ref<number | null>(null)
const showBookInfo = ref(false)
const autoSave = ref(true)
const generating = ref(false)
const regenerating = ref(false)
let saveTimer: ReturnType<typeof setTimeout> | null = null

const editData = reactive<Partial<Blueprint>>({})

const currentBlueprint = computed(() =>
  blueprints.value.find(bp => bp.chapter_index === selectedChapter.value) || null
)

const fetchBooks = async () => {
  try {
    const res = await axios.get('/api/books')
    books.value = res.data.data || []
  } catch { ElMessage.error('获取书籍失败') }
}

const loadOutline = async () => {
  if (!selectedBookId.value) return
  try {
    const res = await axios.get(`/api/books/${selectedBookId.value}/master-outline`)
    const data = res.data.data
    book.value = data.book
    blueprints.value = data.blueprints || []
    if (blueprints.value.length > 0 && !selectedChapter.value) {
      const firstBlueprint = blueprints.value[0]
      if (firstBlueprint) {
        selectChapter(firstBlueprint.chapter_index)
      }
    }
  } catch { ElMessage.error('获取总纲失败') }
}

const selectChapter = (idx: number) => {
  selectedChapter.value = idx
  const bp = blueprints.value.find(b => b.chapter_index === idx)
  if (bp) {
    Object.assign(editData, {
      protagonist_motivation: bp.protagonist_motivation,
      key_foreshadowing: bp.key_foreshadowing,
      appearing_characters: bp.appearing_characters,
      highlight: bp.highlight,
      core_events: bp.core_events,
      challenges: bp.challenges,
      summary: bp.summary,
    })
  }
}

const onFieldChange = (field: string) => {
  // 同步到本地 blueprints 数组
  const bp = blueprints.value.find(b => b.chapter_index === selectedChapter.value)
  if (bp) {
    (bp as any)[field] = (editData as any)[field]
  }

  if (autoSave.value) {
    if (saveTimer) clearTimeout(saveTimer)
    saveTimer = setTimeout(() => saveField(field), 800)
  }
}

const saveField = async (field: string) => {
  if (!selectedBookId.value || !selectedChapter.value) return
  try {
    await axios.put(
      `/api/books/${selectedBookId.value}/master-outline/${selectedChapter.value}`,
      { [field]: (editData as any)[field] }
    )
  } catch {
    ElMessage.error('保存失败')
  }
}

const generateBatch = async () => {
  if (!selectedBookId.value) return
  generating.value = true
  const lastBlueprint = blueprints.value.length > 0 ? blueprints.value[blueprints.value.length - 1] : undefined
  const startChapter = lastBlueprint ? lastBlueprint.chapter_index + 1 : 1
  try {
    const res = await axios.post(`/api/books/${selectedBookId.value}/master-outline/generate`, {
      start_chapter: startChapter,
      batch_size: 10,
    })
    const newBps = res.data.data.blueprints || []
    // 追加或替换
    for (const nb of newBps) {
      const idx = blueprints.value.findIndex(b => b.chapter_index === nb.chapter_index)
      if (idx >= 0) {
        blueprints.value[idx] = nb
      } else {
        blueprints.value.push(nb)
      }
    }
    blueprints.value.sort((a, b) => a.chapter_index - b.chapter_index)
    if (newBps.length > 0) {
      selectChapter(newBps[0].chapter_index)
    }
    ElMessage.success(`成功生成 ${newBps.length} 章蓝图`)
  } catch {
    ElMessage.error('总纲生成失败')
  } finally {
    generating.value = false
  }
}

const regenerateChapter = async () => {
  if (!selectedBookId.value || !selectedChapter.value) return
  regenerating.value = true
  try {
    const res = await axios.post(
      `/api/books/${selectedBookId.value}/master-outline/${selectedChapter.value}/regenerate`
    )
    const bp = res.data.data.blueprint
    const idx = blueprints.value.findIndex(b => b.chapter_index === selectedChapter.value)
    if (idx >= 0) {
      blueprints.value[idx] = { ...blueprints.value[idx], ...bp }
    }
    selectChapter(selectedChapter.value)
    ElMessage.success('重新生成完成')
  } catch {
    ElMessage.error('重新生成失败')
  } finally {
    regenerating.value = false
  }
}

const handleBookChange = () => {
  if (selectedBookId.value) {
    router.push(`/outline/${selectedBookId.value}`)
    selectedChapter.value = null
    loadOutline()
  }
}

watch(() => route.params.bookId, (value) => {
  const id = value ? Number(value) : null
  if (id && id !== selectedBookId.value) {
    selectedBookId.value = id
    loadOutline()
  }
})

onMounted(async () => {
  await fetchBooks()
  const id = route.params.bookId ? Number(route.params.bookId) : null
  if (id) {
    selectedBookId.value = id
    loadOutline()
  }
})
</script>

<style scoped>
.master-outline {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.outline-hero {
  padding: 24px 28px;
}
.outline-hero h1 {
  margin: 10px 0 6px;
  font-size: 26px;
}
.outline-hero p {
  margin: 0;
  color: #64748b;
}

/* 顶部提示栏 */
.outline-banner {
  padding: 14px 20px;
  background: linear-gradient(135deg, #fff8e8 0%, #fff3d6 100%);
  border-left: 4px solid #f59e0b;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: #92400e;
  border-radius: 12px;
}
.banner-icon {
  font-size: 18px;
}

/* 主体容器 */
.outline-container {
  display: grid;
  grid-template-columns: 240px 1fr;
  gap: 16px;
  min-height: 600px;
}

/* 左侧栏 */
.chapter-sidebar {
  padding: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  padding: 16px 16px 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  border-bottom: 1px solid #f1f5f9;
}
.sidebar-title {
  font-weight: 700;
  font-size: 15px;
  color: #c2410c;
}

.book-info {
  padding: 8px 16px 12px;
  border-bottom: 1px solid #f1f5f9;
}
.book-desc {
  font-size: 12px;
  color: #64748b;
  margin: 0 0 6px;
  line-height: 1.5;
}

.chapter-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.chapter-item {
  padding: 10px 16px;
  cursor: pointer;
  display: flex;
  gap: 6px;
  align-items: baseline;
  font-size: 13px;
  color: #475569;
  transition: all 0.15s ease;
  border-left: 3px solid transparent;
}
.chapter-item:hover {
  background: #f8fafc;
}
.chapter-item.active {
  background: linear-gradient(90deg, #fff7ed 0%, #fff 100%);
  border-left-color: #ea580c;
  color: #c2410c;
  font-weight: 600;
}
.chapter-label {
  white-space: nowrap;
  font-weight: 600;
}
.chapter-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sidebar-actions {
  padding: 12px 16px;
  border-top: 1px solid #f1f5f9;
}
.sidebar-actions .el-button {
  width: 100%;
}

/* 右侧面板 */
.detail-panel {
  padding: 24px 28px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.empty-panel {
  display: flex;
  align-items: center;
  justify-content: center;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.detail-title {
  margin: 0;
  font-size: 20px;
  font-weight: 700;
  color: #1e293b;
}
.auto-save-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
}
.save-label {
  font-size: 12px;
  color: #94a3b8;
}

.detail-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.form-row {
  display: flex;
  align-items: center;
  gap: 12px;
}
.form-row-textarea {
  align-items: flex-start;
}
.form-label {
  min-width: 72px;
  font-size: 14px;
  font-weight: 600;
  color: #334155;
  text-align: right;
  flex-shrink: 0;
}

.detail-actions {
  display: flex;
  justify-content: center;
  padding-top: 8px;
}

@media (max-width: 1024px) {
  .outline-container {
    grid-template-columns: 1fr;
  }
  .chapter-sidebar {
    max-height: 200px;
  }
}
</style>
