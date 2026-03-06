<template>
  <div class="plans">
    <section class="plans-hero glass-card">
      <div>
        <div class="pill pill-primary">规划生成</div>
        <h1>多版本世界观与剧情规划</h1>
        <p>一次生成多套方案，选择最适合的故事走向。</p>
      </div>
      <div class="hero-actions">
        <el-button type="primary" @click="generatePlans" :loading="loading">生成规划</el-button>
        <el-button @click="refreshPlans">刷新列表</el-button>
      </div>
    </section>

    <section class="plans-grid">
      <div class="glass-card plans-card">
        <div class="card-header">
          <h2>基础信息</h2>
          <p>选择书籍并设置生成参数。</p>
        </div>
        <el-form :model="form" label-width="90px">
          <el-form-item label="书籍">
            <el-select v-model="selectedBookId" placeholder="请选择书籍" @change="handleBookChange">
              <el-option v-for="book in books" :key="book.id" :label="book.title" :value="book.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="题材">
            <el-input v-model="form.genre" placeholder="例如：科幻 / 玄幻" />
          </el-form-item>
          <el-form-item label="简介">
            <el-input v-model="form.description" type="textarea" :rows="4" placeholder="故事简述与核心创意" />
          </el-form-item>
          <el-form-item label="章节数">
            <el-input-number v-model="form.chapters" :min="1" />
          </el-form-item>
          <el-form-item label="版本数">
            <el-input-number v-model="form.count" :min="1" :max="5" />
          </el-form-item>
        </el-form>
      </div>

      <div class="glass-card characters-card" v-if="selectedBookId && characters.length > 0">
        <div class="card-header">
          <h2>角色库</h2>
          <p>已生成的角色列表。</p>
        </div>
        <div class="characters-list">
          <div v-for="char in characters" :key="char.id" class="character-item">
            <div class="char-header">
              <span class="char-name">{{ char.name }}</span>
              <el-tag size="small" type="info">{{ char.role }}</el-tag>
            </div>
            <div class="char-desc">{{ char.description }}</div>
          </div>
        </div>
      </div>

      <div class="plans-list">
        <div class="glass-card plan-card" v-for="plan in plans" :key="plan.id" :class="{ selected: plan.is_selected }">
          <div class="plan-header">
            <div class="plan-badge">方案 {{ plan.id }}</div>
            <el-tag v-if="plan.is_selected" type="success" effect="dark">已选择</el-tag>
          </div>
          
          <el-collapse v-model="activeNames">
            <el-collapse-item title="世界观设定" name="world">
              <div class="plan-text">{{ plan.world_view }}</div>
            </el-collapse-item>
            <el-collapse-item title="剧情大纲" name="outline">
              <div class="plan-text">{{ plan.outline }}</div>
            </el-collapse-item>
            <el-collapse-item title="角色设定" name="characters">
              <div class="plan-text" v-if="plan.characters">
                <template v-if="Array.isArray(plan.characters)">
                  <div v-for="(char, idx) in plan.characters" :key="idx" class="plan-char-item">
                    <strong>{{ char.name }}</strong> ({{ char.role }}): {{ char.description }}
                  </div>
                </template>
                <template v-else>{{ plan.characters }}</template>
                <div class="regen-actions" v-if="plan.is_selected && !plan.is_locked">
                  <el-divider />
                  <el-button
                    type="primary"
                    size="small"
                    :loading="generatingCharacters"
                    @click="generateCharacters(plan)"
                  >
                    重新生成角色设定
                  </el-button>
                </div>
              </div>
              <div class="plan-empty" v-else>
                <el-button
                  v-if="plan.is_selected"
                  type="primary"
                  size="small"
                  :loading="generatingCharacters"
                  @click="generateCharacters(plan)"
                >
                  生成角色设定
                </el-button>
                <span v-else>请先选择此方案以生成角色</span>
              </div>
            </el-collapse-item>
            <el-collapse-item title="章节标题" name="titles">
              <div class="plan-text" v-if="plan.titles">
                {{ plan.titles }}
                <div class="regen-actions" v-if="plan.is_selected && plan.characters && !plan.is_locked">
                  <el-divider />
                  <el-button
                    type="primary"
                    size="small"
                    :loading="generatingChapters"
                    @click="generateChapters(plan)"
                  >
                    重新生成章节标题
                  </el-button>
                </div>
              </div>
              <div class="plan-empty" v-else>
                <el-button
                  v-if="plan.is_selected && plan.characters"
                  type="primary"
                  size="small"
                  :loading="generatingChapters"
                  @click="generateChapters(plan)"
                >
                  生成章节标题
                </el-button>
                <span v-else-if="!plan.is_selected">请先选择此方案</span>
                <span v-else>请先生成角色设定</span>
              </div>
            </el-collapse-item>
          </el-collapse>

          <div class="plan-actions">
            <template v-if="!plan.is_selected">
              <el-button type="primary" class="select-btn" @click="selectPlan(plan)">使用此方案</el-button>
            </template>
            <template v-else>
              <el-button v-if="!plan.is_locked" type="success" @click="confirmPlan(plan)">确定使用（锁定后不可更改）</el-button>
              <el-button v-else type="info" plain @click="unlockPlan(plan)">方案已锁定 (点击解锁)</el-button>
            </template>
          </div>
        </div>
        <div v-if="!loading && plans.length === 0" class="empty-state glass-card">
          <el-empty description="暂无规划方案" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

type Book = {
  id: number
  title: string
  genre: string
  description: string
  total_chapters: number
}

type PlanVersion = {
  id: number
  world_view: string
  outline: string
  characters: string | any[]
  titles: string
  is_selected: boolean
  is_locked: boolean
}

type Character = {
  id: number
  name: string
  role: string
  description: string
}

const route = useRoute()
const router = useRouter()

const books = ref<Book[]>([])
const plans = ref<PlanVersion[]>([])
const characters = ref<Character[]>([])
const loading = ref(false)
const fetchingCharacters = ref(false)
const generatingCharacters = ref(false)
const generatingChapters = ref(false)
const selectedBookId = ref<number | null>(null)
const activeNames = ref(['world', 'outline'])

const form = reactive({
  genre: '',
  description: '',
  chapters: 1,
  count: 3
})

const currentBook = computed(() => books.value.find((item) => item.id === selectedBookId.value) || null)

const streamPostSSE = async (
  url: string,
  body: Record<string, any> | null,
  onMessage: (text: string) => void,
  onStateUpdate?: (payload: any) => void
) => {
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: body ? JSON.stringify(body) : undefined
  })

  if (!response.ok) {
    let errorMsg = `请求失败 (${response.status})`
    try {
      const payload = await response.json()
      if (payload?.error) errorMsg = payload.error
    } catch {
      const text = await response.text()
      if (text) errorMsg = text
    }
    throw new Error(errorMsg)
  }

  if (!response.body) {
    throw new Error('无流式响应')
  }

  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let currentEvent = ''

  while (true) {
    const { value, done } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })
    const lines = buffer.split('\n')
    buffer = lines.pop() || ''

    for (const line of lines) {
      const trimmedLine = line.trim()
      if (!trimmedLine) continue

      if (trimmedLine.startsWith('event:')) {
        currentEvent = trimmedLine.slice(6).trim()
      } else if (trimmedLine.startsWith('data:')) {
        const dataStr = trimmedLine.slice(5).trim()
        let data: any = dataStr
        try {
          data = JSON.parse(dataStr)
        } catch {
          data = dataStr
        }

        if (currentEvent === 'message') {
          onMessage(String(data))
        } else if (currentEvent === 'state_update') {
          onStateUpdate?.(data)
        } else if (currentEvent === 'error') {
          throw new Error(String(data))
        } else if (currentEvent === 'end') {
          return
        }
      }
    }
  }
}

const fetchBooks = async () => {
  try {
    const res = await axios.get('/api/books')
    books.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取书籍失败')
  }
}

const loadPlans = async () => {
  if (!selectedBookId.value) return
  loading.value = true
  try {
    const res = await axios.get(`/api/books/${selectedBookId.value}/plans`)
    plans.value = res.data.data || []
    // 尝试解析 characters 如果是字符串的话
    plans.value.forEach(plan => {
      if (typeof plan.characters === 'string' && plan.characters.startsWith('[')) {
        try {
          plan.characters = JSON.parse(plan.characters)
        } catch (e) {
          console.error('Failed to parse characters JSON', e)
        }
      }
    })
  } catch (error) {
    ElMessage.error('获取规划失败')
  } finally {
    loading.value = false
  }
}

const fetchCharacters = async () => {
  if (!selectedBookId.value) return
  fetchingCharacters.value = true
  try {
    const res = await axios.get(`/api/books/${selectedBookId.value}/characters`)
    characters.value = res.data.data || []
  } catch (error) {
    console.error('获取角色失败', error)
  } finally {
    fetchingCharacters.value = false
  }
}

const handleBookChange = () => {
  if (selectedBookId.value) {
    router.push(`/plans/${selectedBookId.value}`)
    if (currentBook.value) {
      form.genre = currentBook.value.genre || ''
      form.description = currentBook.value.description || ''
      form.chapters = currentBook.value.total_chapters || 1
    }
    loadPlans()
    fetchCharacters()
  }
}

const generatePlans = async () => {
  if (!selectedBookId.value) {
    ElMessage.warning('请先选择书籍')
    return
  }
  loading.value = true
  const versionCount = form.count > 0 ? form.count : 3
  plans.value = Array.from({ length: versionCount }).map((_, idx) => ({
    id: -(idx + 1),
    world_view: '',
    outline: '',
    characters: '',
    titles: '',
    is_selected: false,
    is_locked: false
  }))
  let currentVersionIndex = 0
  let currentPhase = ''

  try {
    await streamPostSSE(
      `/api/books/${selectedBookId.value}/plans/generate`,
      {
        description: form.description,
        genre: form.genre,
        chapters: form.chapters,
        count: form.count
      },
      (text) => {
        if (currentPhase === 'world_start') {
          plans.value.forEach(p => {
            if (p) p.world_view += text
          })
          return
        }
        
        const currentPlan = plans.value[currentVersionIndex]
        if (!currentPlan) return
        currentPlan.outline += text
      },
      (payload) => {
        if (payload?.phase) {
          currentPhase = payload.phase
        }
        if (payload?.phase === 'plan_version_start' && typeof payload.index === 'number') {
          currentVersionIndex = Math.max(0, payload.index - 1)
          const currentPlan = plans.value[currentVersionIndex]
          if (!currentPlan) {
            plans.value[currentVersionIndex] = {
              id: -(currentVersionIndex + 1),
              world_view: plans.value[0]?.world_view || '',
              outline: '',
              characters: '',
              titles: '',
              is_selected: false,
              is_locked: false
            }
          } else {
            const worldView = plans.value[0]?.world_view || ''
            if (worldView && !currentPlan.world_view) {
              currentPlan.world_view = worldView
            }
          }
        }
      }
    )
    await loadPlans()
    ElMessage.success('规划生成完成')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '规划生成失败')
  } finally {
    loading.value = false
  }
}

const selectPlan = async (plan: PlanVersion) => {
  if (!selectedBookId.value) return
  try {
    await axios.put(`/api/books/${selectedBookId.value}/plans/${plan.id}/select`)
    plans.value = plans.value.map((item) => ({ ...item, is_selected: item.id === plan.id }))
    ElMessage.success('已选择方案')
  } catch (error) {
    ElMessage.error('选择方案失败')
  }
}

const confirmPlan = async (plan: PlanVersion) => {
  if (!selectedBookId.value) return
  try {
    await axios.put(`/api/books/${selectedBookId.value}/plans/${plan.id}/lock`, { locked: true })
    plan.is_locked = true
    ElMessage.success('方案已锁定，现在可以开始创作了')
  } catch (error) {
    ElMessage.error('锁定方案失败')
  }
}

const unlockPlan = async (plan: PlanVersion) => {
  if (!selectedBookId.value) return
  try {
    await axios.put(`/api/books/${selectedBookId.value}/plans/${plan.id}/lock`, { locked: false })
    plan.is_locked = false
    ElMessage.success('方案已解锁')
  } catch (error) {
    ElMessage.error('解锁方案失败')
  }
}

const generateCharacters = async (plan: PlanVersion) => {
  if (!selectedBookId.value) return
  generatingCharacters.value = true
  plan.characters = ''
  try {
    await streamPostSSE(
      `/api/books/${selectedBookId.value}/plans/characters`,
      null,
      (text) => {
        if (typeof plan.characters !== 'string') {
          plan.characters = ''
        }
        plan.characters += text
      }
    )
    await loadPlans()
    await fetchCharacters() // 刷新角色列表
    ElMessage.success('角色设定生成完成')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '角色设定生成失败')
  } finally {
    generatingCharacters.value = false
  }
}

const generateChapters = async (plan: PlanVersion) => {
  if (!selectedBookId.value) return
  generatingChapters.value = true
  plan.titles = ''
  try {
    await streamPostSSE(
      `/api/books/${selectedBookId.value}/plans/chapters`,
      null,
      (text) => {
        plan.titles += text
      }
    )
    await loadPlans()
    ElMessage.success('章节标题生成完成')
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '章节标题生成失败')
  } finally {
    generatingChapters.value = false
  }
}

const refreshPlans = () => {
  loadPlans()
}

watch(
  () => route.params.bookId,
  (value) => {
    const id = value ? Number(value) : null
    if (id && id !== selectedBookId.value) {
      selectedBookId.value = id
      handleBookChange()
    }
  }
)

onMounted(async () => {
  await fetchBooks()
  const initialId = route.params.bookId ? Number(route.params.bookId) : null
  if (initialId) {
    selectedBookId.value = initialId
    handleBookChange()
  }
})
</script>

<style scoped>
.plans {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.plans-hero {
  padding: 24px 28px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}

.plans-hero h1 {
  margin: 10px 0 6px;
  font-size: 26px;
}

.plans-hero p {
  margin: 0;
  color: #64748b;
}

.hero-actions {
  display: flex;
  gap: 12px;
}

.plans-grid {
  display: grid;
  grid-template-columns: minmax(0, 0.9fr) minmax(0, 1.1fr);
  gap: 24px;
}

.plans-card {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: fit-content;
}

.characters-card {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  height: fit-content;
}

.characters-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.character-item {
  padding-bottom: 12px;
  border-bottom: 1px solid #f1f5f9;
}

.character-item:last-child {
  border-bottom: none;
}

.char-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.char-name {
  font-weight: 600;
  color: #1e293b;
}

.char-desc {
  font-size: 13px;
  color: #64748b;
  line-height: 1.5;
}

.plan-char-item {
  margin-bottom: 8px;
}

.plan-char-item:last-child {
  margin-bottom: 0;
}

.card-header h2 {
  margin: 0 0 6px 0;
  font-size: 18px;
}

.card-header p {
  margin: 0;
  color: #64748b;
  font-size: 13px;
}

.plans-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.plan-card {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  transition: all 0.3s ease;
  border: 1px solid rgba(148, 163, 184, 0.1);
}

.plan-card.selected {
  border-color: #2563eb;
  background: rgba(37, 99, 235, 0.05);
}

.plan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.plan-badge {
  background: #2563eb;
  color: white;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 600;
}

.plan-text {
  font-size: 0.875rem;
  color: #475569;
  line-height: 1.6;
  white-space: pre-wrap;
}

.plan-empty {
  padding: 12px;
  background: #f1f5f9;
  border-radius: 8px;
  text-align: center;
  font-size: 0.875rem;
  color: #64748b;
}

.plan-actions {
  margin-top: 8px;
  display: flex;
  justify-content: flex-end;
}

.select-btn {
  width: 100%;
}

.empty-state {
  padding: 20px;
}

@media (max-width: 1024px) {
  .plans-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .plans-hero {
    flex-direction: column;
    align-items: flex-start;
  }

  .hero-actions {
    width: 100%;
    flex-direction: column;
  }
}
</style>
