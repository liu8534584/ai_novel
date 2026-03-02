<template>
  <div class="writing">
    <section class="writing-header glass-card">
      <div class="header-info">
        <el-button icon="Back" @click="router.back()" circle />
        <div class="book-info" v-if="book">
          <h1>{{ book.title }}</h1>
          <div class="chapter-title-edit" v-if="currentChapter">
            <el-input v-model="currentChapter.title" placeholder="章节标题" size="small" />
          </div>
          <p v-else>未选择章节</p>
        </div>
      </div>
      <div class="header-actions">
        <el-button type="info" plain @click="generateOutline" :loading="outlining">
          {{ currentChapter?.outline ? '重新生成大纲' : '生成本章大纲' }}
        </el-button>
        <el-button 
          v-if="currentChapter?.outline && !currentChapter?.is_outline_confirmed" 
          type="warning" 
          @click="confirmOutline" 
          :loading="confirmingOutline"
        >
          确定使用大纲
        </el-button>
        <el-button 
          type="primary" 
          :loading="writing" 
          @click="writeChapter" 
          :disabled="!currentChapter?.is_outline_confirmed"
          title="必须先确认大纲才能续写"
        >
          自动续写
        </el-button>
        <el-button type="success" plain @click="saveChapter">保存草稿</el-button>
      </div>
    </section>

    <section class="writing-layout">
      <aside class="chapter-sidebar glass-card">
        <div class="sidebar-header">
          <h3>章节目录</h3>
        </div>
        <div class="chapter-list">
          <div
            v-for="chapter in chapters"
            :key="chapter.id"
            :class="['chapter-item', { active: currentChapterId === chapter.id }]"
            @click="selectChapter(chapter)"
          >
            <span class="chapter-order">{{ chapter.order }}</span>
            <span class="chapter-title">{{ chapter.title }}</span>
            <el-icon v-if="chapter.is_outline_confirmed" color="#67C23A" style="margin-left: 5px;"><SuccessFilled /></el-icon>
          </div>
        </div>
      </aside>

      <main class="editor-container glass-card">
        <div v-if="currentChapter" class="editor-wrapper">
          <div v-if="currentChapter.outline && !currentChapter.is_outline_confirmed" class="outline-preview glass-card">
            <div class="outline-header">
              <h4>本章大纲预案 (待确认)</h4>
              <div class="outline-header-actions">
                <el-button size="small" type="primary" link @click="editChapterOutline = !editChapterOutline">
                  {{ editChapterOutline ? '预览大纲' : '编辑大纲' }}
                </el-button>
                <el-tag type="warning">确认后方可开始写作</el-tag>
              </div>
            </div>
            <div class="outline-content">
              <el-input
                v-if="editChapterOutline"
                v-model="currentChapter.outline"
                type="textarea"
                :rows="8"
                placeholder="编辑章节大纲 (JSON 格式)..."
              />
              <pre v-else>{{ formatOutline(currentChapter.outline) }}</pre>
            </div>
          </div>
          <el-input
            v-model="currentChapter.content"
            type="textarea"
            :rows="currentChapter.outline && !currentChapter.is_outline_confirmed ? 15 : 25"
            placeholder="开始写作..."
            class="chapter-editor"
          />
        </div>
        <div v-else class="editor-placeholder">
          <el-empty description="请从左侧选择章节开始写作" />
        </div>
      </main>

      <aside class="context-sidebar glass-card">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="AI 审计" name="audit">
            <div class="context-content audit-panel">
              <div v-if="healthScore" class="audit-section health-section">
                <h4>章节健康度</h4>
                <div class="health-card glass-card">
                  <div class="health-main">
                    <el-progress type="circle" :percentage="Math.round(healthScore.total_health)" :color="getHealthColor(healthScore.total_health)" />
                    <div class="health-stats">
                      <div class="stat-item">
                        <span>OOC 表现:</span>
                        <el-tag size="small" :type="getHealthTag(healthScore.ooc_score)">{{ healthScore.ooc_score.toFixed(1) }}</el-tag>
                      </div>
                      <div class="stat-item">
                        <span>一致性:</span>
                        <el-tag size="small" :type="getHealthTag(healthScore.event_consistency)">{{ healthScore.event_consistency.toFixed(1) }}</el-tag>
                      </div>
                      <div class="stat-item">
                        <span>伏笔健康:</span>
                        <el-tag size="small" :type="getHealthTag(healthScore.foreshadowing)">{{ healthScore.foreshadowing.toFixed(1) }}</el-tag>
                      </div>
                    </div>
                  </div>
                  <div class="audit-report">
                    <h5>审计建议</h5>
                    <p class="report-text">{{ healthScore.audit_report }}</p>
                  </div>
                </div>
              </div>

              <div v-if="foreshadowingAlerts.length > 0" class="audit-section">
                <h4>伏笔预警</h4>
                <div v-for="alert in foreshadowingAlerts" :key="alert.ID" class="audit-card fore-card">
                  <div class="alert-header">
                    <el-tag type="danger">未回收</el-tag>
                    <span class="alert-title">{{ alert.title }}</span>
                  </div>
                  <p class="alert-desc">{{ alert.description }}</p>
                  <div class="alert-meta">
                    <span>重要度: {{ alert.importance }}</span>
                    <span>埋下于第 {{ alert.chapter_index }} 章</span>
                  </div>
                </div>
              </div>

              <div v-if="oocScores.length > 0" class="audit-section">
                <h4>角色 OOC 评分</h4>
                <div v-for="score in oocScores" :key="score.ID" class="audit-card ooc-card">
                  <div class="score-header">
                    <span class="char-name">{{ getCharacterName(score.character_id) }}</span>
                    <el-tag :type="getOOCLevel(score.total_score)">{{ score.conclusion }}</el-tag>
                  </div>
                  <div class="score-details">
                    <div class="score-item">性格一致性: {{ score.personality_consistency }}</div>
                    <div class="score-item">动机一致性: {{ score.motivation_consistency }}</div>
                    <div class="score-item">综合评分: <span class="score-value">{{ score.total_score }}</span></div>
                  </div>
                  <p class="explanation">{{ score.explanation }}</p>
                </div>
              </div>

              <div v-if="contradictions.length > 0" class="audit-section">
                <h4>剧情矛盾检测</h4>
                <div v-for="con in contradictions" :key="con.ID" class="audit-card con-card">
                  <div class="con-header">
                    <el-tag :type="con.severity === '严重' ? 'danger' : 'warning'">{{ con.type }}</el-tag>
                    <span class="severity">{{ con.severity }}</span>
                  </div>
                  <p class="description">{{ con.description }}</p>
                  <div class="suggestion">
                    <strong>修正建议:</strong> {{ con.suggestion }}
                  </div>
                </div>
              </div>

              <div v-if="oocScores.length === 0 && contradictions.length === 0" class="empty-audit">
                <el-empty description="本章暂无审计记录" />
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane label="角色" name="characters">
            <div class="context-content info-panel" v-if="plan">
              <el-input
                v-model="plan.characters"
                type="textarea"
                :rows="20"
                placeholder="编辑角色设定..."
                class="info-editor"
              />
              <div class="info-actions">
                <el-button type="primary" size="small" @click="savePlan" :loading="savingPlan">保存角色设定</el-button>
              </div>
            </div>
            <div v-else class="empty-info">暂无角色设定</div>
          </el-tab-pane>
          <el-tab-pane label="世界观" name="world">
            <div class="context-content info-panel" v-if="plan">
              <el-input
                v-model="plan.world_view"
                type="textarea"
                :rows="20"
                placeholder="编辑世界观设定..."
                class="info-editor"
              />
              <div class="info-actions">
                <el-button type="primary" size="small" @click="savePlan" :loading="savingPlan">保存世界观设定</el-button>
              </div>
            </div>
            <div v-else class="empty-info">暂无世界观设定</div>
          </el-tab-pane>
          <el-tab-pane label="大纲" name="outline">
            <div class="context-content info-panel" v-if="plan">
              <el-input
                v-model="plan.outline"
                type="textarea"
                :rows="20"
                placeholder="编辑剧情大纲..."
                class="info-editor"
              />
              <div class="info-actions">
                <el-button type="primary" size="small" @click="savePlan" :loading="savingPlan">保存大纲设定</el-button>
              </div>
            </div>
            <div v-else class="empty-info">暂无剧情大纲</div>
          </el-tab-pane>
        </el-tabs>
      </aside>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { SuccessFilled } from '@element-plus/icons-vue'
import { watch } from 'vue'

const route = useRoute()
const router = useRouter()
const bookId = route.params.bookId

const book = ref<any>(null)
const chapters = ref<any[]>([])
const characters = ref<any[]>([])
const plan = ref<any>(null)
const currentChapterId = ref<number | null>(null)
const writing = ref(false)
const outlining = ref(false)
const confirmingOutline = ref(false)
const savingPlan = ref(false)
const editChapterOutline = ref(false)
const activeTab = ref('audit')
const oocScores = ref<any[]>([])
const contradictions = ref<any[]>([])
const healthScore = ref<any>(null)
const foreshadowingAlerts = ref<any[]>([])

const currentChapter = computed(() => chapters.value.find(c => c.id === currentChapterId.value))

const getOOCLevel = (score: number) => {
  if (score < 30) return 'success'
  if (score < 50) return 'warning'
  return 'danger'
}

const getHealthColor = (score: number) => {
  if (score >= 90) return '#67C23A'
  if (score >= 70) return '#E6A23C'
  return '#F56C6C'
}

const getHealthTag = (score: number) => {
  if (score >= 90) return 'success'
  if (score >= 70) return 'warning'
  return 'danger'
}

const fetchAuditData = async (chapterId: number) => {
  try {
    const [oocRes, conRes, healthRes, foreRes] = await Promise.all([
      axios.get(`/api/chapters/${chapterId}/ooc-scores`),
      axios.get(`/api/chapters/${chapterId}/contradictions`),
      axios.get(`/api/chapters/${chapterId}/health`),
      axios.get(`/api/books/${bookId}/foreshadowing-alerts`)
    ])
    oocScores.value = oocRes.data.data || []
    contradictions.value = conRes.data.data || []
    healthScore.value = healthRes.data.data || null
    foreshadowingAlerts.value = foreRes.data.data || []
  } catch (error) {
    console.error('获取审计数据失败', error)
  }
}

const fetchBookData = async () => {
  try {
    const [bookRes, chaptersRes, plansRes, charactersRes] = await Promise.all([
      axios.get(`/api/books/${bookId}`),
      axios.get(`/api/books/${bookId}/chapters`),
      axios.get(`/api/books/${bookId}/plans`),
      axios.get(`/api/books/${bookId}/characters`)
    ])
    book.value = bookRes.data.data
    chapters.value = chaptersRes.data.data || []
    characters.value = charactersRes.data.data || []
    const allPlans = plansRes.data.data || []
    plan.value = allPlans.find((p: any) => p.is_selected)

    // 自动选择上次编辑的章节（如果有本地草稿）
    if (!currentChapterId.value && chapters.value.length > 0) {
      const draftKey = `novel_draft_${bookId}`
      const savedDraft = localStorage.getItem(draftKey)
      if (savedDraft) {
        try {
          const { chapterId } = JSON.parse(savedDraft)
          const draftChapter = chapters.value.find(c => c.id === chapterId)
          if (draftChapter) {
            selectChapter(draftChapter)
          }
        } catch (e) {
          console.error('解析本地草稿失败', e)
        }
      }
    }
  } catch (error) {
    ElMessage.error('获取书籍数据失败')
  }
}

const getCharacterName = (id: number) => {
  const char = characters.value.find(c => c.id === id)
  return char ? char.name : `角色 ${id}`
}

// 自动保存草稿到 localStorage
watch(() => currentChapter.value?.content, (newContent) => {
  if (currentChapter.value && newContent !== undefined) {
    const draftKey = `novel_draft_${bookId}`
    const draftData = {
      chapterId: currentChapter.value.id,
      content: newContent
    }
    localStorage.setItem(draftKey, JSON.stringify(draftData))
  }
})

const clearLocalDraft = () => {
  localStorage.removeItem(`novel_draft_${bookId}`)
}

const loadLocalDraft = (chapterId: number) => {
  const draftKey = `novel_draft_${bookId}`
  const savedDraft = localStorage.getItem(draftKey)
  if (savedDraft) {
    try {
      const { chapterId: savedId, content } = JSON.parse(savedDraft)
      if (savedId === chapterId && content) {
        return content
      }
    } catch (e) {
      console.error('解析本地草稿失败', e)
    }
  }
  return null
}

const selectChapter = (chapter: any) => {
  currentChapterId.value = chapter.id
  fetchAuditData(chapter.id)

  const localContent = loadLocalDraft(chapter.id)
  if (localContent && localContent !== chapter.content) {
    chapter.content = localContent
    ElMessage.info('已恢复未保存的本地草稿')
  }
}

const generateOutline = async () => {
  if (!currentChapter.value) return
  outlining.value = true
  try {
    await axios.post(`/api/chapters/${currentChapter.value.id}/outline`)
    ElMessage.success('本章大纲生成成功')
    await fetchBookData() // 刷新以获取更新的大纲和确认状态
  } catch (error) {
    ElMessage.error('生成大纲失败')
  } finally {
    outlining.value = false
  }
}

const confirmOutline = async () => {
  if (!currentChapter.value) return
  confirmingOutline.value = true
  try {
    await axios.post(`/api/chapters/${currentChapter.value.id}/outline/confirm`)
    ElMessage.success('大纲已确认，现在可以开始续写了')
    await fetchBookData()
  } catch (error) {
    ElMessage.error('确认大纲失败')
  } finally {
    confirmingOutline.value = false
  }
}

const formatOutline = (outlineStr: string) => {
  try {
    const outline = JSON.parse(outlineStr)
    let formatted = `【${outline.title}】\n\n`
    formatted += `摘要：${outline.summary}\n\n`
    if (outline.scenes) {
      outline.scenes.forEach((scene: any) => {
        formatted += `场景 ${scene.order}: ${scene.location}\n`
        formatted += `描述: ${scene.description}\n`
        formatted += `冲突: ${scene.key_conflict}\n`
        formatted += `结果: ${scene.outcome}\n\n`
      })
    }
    return formatted
  } catch (e) {
    return outlineStr
  }
}

const saveChapter = async () => {
  if (!currentChapter.value) return
  try {
    await axios.put(`/api/chapters/${currentChapter.value.id}`, {
      content: currentChapter.value.content,
      title: currentChapter.value.title,
      outline: currentChapter.value.outline
    })
    ElMessage.success('保存成功')
    clearLocalDraft()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const savePlan = async () => {
  if (!plan.value) return
  savingPlan.value = true
  try {
    await axios.put(`/api/books/${bookId}/plans/${plan.value.id}`, {
      world_view: plan.value.world_view,
      outline: plan.value.outline,
      characters: plan.value.characters
    })
    ElMessage.success('设定保存成功')
  } catch (error) {
    ElMessage.error('保存设定失败')
  } finally {
    savingPlan.value = false
  }
}

const writeChapter = async () => {
  if (!currentChapter.value) return
  writing.value = true
  currentChapter.value.content = '' // 清空内容准备续写

  try {
    const response = await fetch(`/api/chapters/${currentChapter.value.id}/write`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
    })

    if (!response.body) return

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { value, done } = await reader.read()
      if (done) break
      
      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || '' // 保留最后一个不完整的行
      
      let currentEvent = ''
      for (const line of lines) {
        const trimmedLine = line.trim()
        if (!trimmedLine) continue

        if (trimmedLine.startsWith('event:')) {
          currentEvent = trimmedLine.slice(6).trim()
        } else if (trimmedLine.startsWith('data:')) {
          const dataStr = trimmedLine.slice(5).trim()
          try {
            const data = JSON.parse(dataStr)
            if (currentEvent === 'message') {
              currentChapter.value.content += data
            } else if (currentEvent === 'end') {
              ElMessage.success('生成完成')
              fetchAuditData(currentChapter.value.id)
            } else if (currentEvent === 'error') {
              ElMessage.error('生成出错: ' + data)
            }
          } catch (e) {
            // 如果不是 JSON，尝试直接使用内容
            if (currentEvent === 'message') {
              currentChapter.value.content += dataStr
            }
          }
        }
      }
    }
  } catch (error) {
    ElMessage.error('生成失败')
  } finally {
    writing.value = false
  }
}

onMounted(() => {
  fetchBookData()
})
</script>

<style scoped>
.writing {
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: calc(100vh - 100px);
}

.writing-header {
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.chapter-title-edit {
  margin-top: 4px;
  width: 200px;
}

.book-info h1 {
  margin: 0;
  font-size: 1.5rem;
  color: #0f172a;
}

.book-info p {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 0.875rem;
}

.writing-layout {
  display: grid;
  grid-template-columns: 260px 1fr 320px;
  gap: 20px;
  flex: 1;
  min-height: 0;
}

.chapter-sidebar, .context-sidebar {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid rgba(148, 163, 184, 0.1);
}

.sidebar-header h3 {
  margin: 0;
  font-size: 1rem;
  color: #1e293b;
}

.chapter-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.chapter-item {
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.2s ease;
  margin-bottom: 4px;
}

.chapter-item:hover {
  background: rgba(59, 130, 246, 0.05);
}

.chapter-item.active {
  background: rgba(59, 130, 246, 0.1);
  color: #2563eb;
}

.outline-preview {
  margin-bottom: 20px;
  padding: 20px;
  background: rgba(245, 158, 11, 0.05);
  border: 1px solid rgba(245, 158, 11, 0.2);
}

.outline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.outline-header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.outline-header h4 {
  margin: 0;
  color: #b45309;
}

.outline-content {
  background: white;
  padding: 15px;
  border-radius: 8px;
  max-height: 200px;
  overflow-y: auto;
}

.outline-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: inherit;
  font-size: 0.9rem;
  line-height: 1.6;
  color: #475569;
}

.chapter-order {
  font-family: 'Fira Code', monospace;
  font-size: 0.75rem;
  opacity: 0.5;
  width: 24px;
}

.editor-container {
  display: flex;
  flex-direction: column;
  padding: 24px;
}

.chapter-editor :deep(.el-textarea__inner) {
  border: none;
  background: transparent;
  box-shadow: none;
  font-size: 1.125rem;
  line-height: 1.8;
  padding: 0;
  color: #1e293b;
  resize: none;
}

.audit-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 16px;
}

.audit-section h4 {
  margin: 0 0 12px;
  font-size: 1rem;
  color: #334155;
}

.health-card {
  padding: 16px;
  margin-bottom: 20px;
}

.health-main {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 16px;
}

.health-stats {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
  color: #64748b;
}

.audit-report h5 {
  margin: 0 0 8px;
  font-size: 0.9rem;
  color: #1e293b;
}

.report-text {
  font-size: 0.875rem;
  color: #475569;
  line-height: 1.6;
  white-space: pre-wrap;
  margin: 0;
}

.fore-card {
  border-left: 4px solid #f87171;
  padding: 12px;
  margin-bottom: 12px;
}

.alert-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.alert-title {
  font-weight: 600;
  color: #1e293b;
}

.alert-desc {
  font-size: 0.875rem;
  color: #475569;
  margin: 0 0 8px;
}

.alert-meta {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: #94a3b8;
}

.info-panel {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-editor :deep(.el-textarea__inner) {
  font-size: 0.875rem;
  line-height: 1.6;
  color: #334155;
  background: rgba(255, 255, 255, 0.5);
}

.info-actions {
  display: flex;
  justify-content: flex-end;
}

.info-text {
  font-size: 0.875rem;
  color: #334155;
  line-height: 1.8;
  white-space: pre-wrap;
}

.empty-info {
  padding: 40px;
  text-align: center;
  color: #94a3b8;
  font-size: 0.875rem;
}

.score-header, .con-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.char-name {
  font-weight: 600;
  font-size: 0.875rem;
}

.score-details {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  font-size: 0.75rem;
  color: #64748b;
  margin-bottom: 8px;
}

.score-value {
  font-weight: 700;
  color: #1e293b;
}

.explanation, .description {
  font-size: 0.8125rem;
  line-height: 1.5;
  margin: 0;
  color: #475569;
}

.suggestion {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px dashed rgba(148, 163, 184, 0.2);
  font-size: 0.75rem;
  color: #059669;
}

.context-content {
  height: 100%;
  overflow-y: auto;
}

.context-content pre {
  white-space: pre-wrap;
  font-family: inherit;
  font-size: 0.875rem;
  line-height: 1.6;
  margin: 0;
  padding: 16px;
  color: #475569;
}

.book-info p {
  margin: 4px 0 0;
  font-size: 14px;
  color: #64748b;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.writing-layout {
  display: grid;
  grid-template-columns: 260px 1fr 300px;
  gap: 20px;
  flex: 1;
  min-height: 0;
}

.chapter-sidebar {
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.chapter-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.batch-generate-action {
  padding: 16px 8px;
}

.batch-controls {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.batch-controls :deep(.el-input-number) {
  width: 100%;
}

.chapter-item {
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  gap: 12px;
  margin-bottom: 4px;
  transition: all 0.2s;
}

.chapter-item:hover {
  background: rgba(255, 255, 255, 0.5);
}

.chapter-item.active {
  background: #3b82f6;
  color: white;
}

.chapter-order {
  opacity: 0.6;
  font-size: 12px;
  width: 20px;
}

.editor-container {
  display: flex;
  flex-direction: column;
  padding: 24px;
}

.editor-wrapper {
  height: 100%;
}

.chapter-editor :deep(.el-textarea__inner) {
  border: none;
  background: transparent;
  box-shadow: none;
  font-size: 18px;
  line-height: 1.8;
  padding: 0;
  resize: none;
}

.editor-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.context-sidebar {
  padding: 16px;
}

.context-content {
  font-size: 13px;
  line-height: 1.6;
  color: #475569;
  white-space: pre-wrap;
  height: calc(100vh - 300px);
  overflow-y: auto;
}

pre {
  font-family: inherit;
  margin: 0;
}
</style>
