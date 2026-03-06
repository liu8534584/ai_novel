<template>
  <div class="settings">
    <section class="settings-hero glass-card">
      <div>
        <div class="pill pill-primary">模型配置</div>
        <h1>掌控你的 AI 写作引擎</h1>
        <p>选择模型提供商、配置密钥与基础地址，打造最合适的写作体验。</p>
      </div>
      <div class="hero-actions">
        <el-button type="primary" @click="saveConfig" :loading="saving">保存配置</el-button>
        <el-button @click="testConnection" :loading="testing">测试连接</el-button>
      </div>
    </section>

    <section class="settings-grid">
      <div class="glass-card settings-card">
        <div class="card-header">
          <h2>基础配置</h2>
          <p>支持 OpenAI / DeepSeek / 阶跃星辰等兼容接口。</p>
        </div>
        <el-form :model="form" label-width="120px">
          <el-form-item label="Provider">
            <el-select v-model="form.provider" placeholder="Select provider" @change="handleProviderChange">
              <el-option label="OpenAI" value="openai" />
              <el-option label="DeepSeek" value="deepseek" />
              <el-option label="阶跃星辰 (StepFun)" value="stepfun" />
              <el-option label="LM Studio" value="lmstudio" />
              <el-option label="Ollama" value="ollama" />
              <el-option label="Local (Generic)" value="local" />
              <el-option label="GLM" value="glm" />
            </el-select>
          </el-form-item>

          <el-form-item label="API Key">
            <el-input v-model="form.api_key" type="password" show-password />
          </el-form-item>

          <el-form-item label="Base URL">
            <el-input v-model="form.base_url" placeholder="https://api.openai.com/v1" />
          </el-form-item>

          <el-form-item label="Model">
            <el-input v-model="form.model" />
          </el-form-item>
        </el-form>
      </div>

      <div class="glass-card settings-card">
        <div class="card-header">
          <h2>写作预设</h2>
          <p>推荐的生成节奏与结构设置。</p>
        </div>
        <div class="preset-list">
          <div class="preset-item">
            <div>
              <div class="preset-title">章节大纲强度</div>
              <div class="preset-meta">适合长篇结构化叙事</div>
            </div>
            <el-button plain>标准</el-button>
          </div>
          <div class="preset-item">
            <div>
              <div class="preset-title">文本风格偏好</div>
              <div class="preset-meta">叙述偏电影镜头</div>
            </div>
            <el-button plain>叙事型</el-button>
          </div>
          <div class="preset-item">
            <div>
              <div class="preset-title">输出节奏</div>
              <div class="preset-meta">每 3 分钟生成 1 场景</div>
            </div>
            <el-button plain>快速</el-button>
          </div>
        </div>
        <div class="preset-hint">
          保存配置后，这些设置会自动应用到新的写作流程。
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const saving = ref(false)
const testing = ref(false)
const form = reactive({
  provider: 'openai',
  api_key: '',
  base_url: '',
  model: 'gpt-3.5-turbo'
})

const handleProviderChange = (val: string) => {
  if (val === 'deepseek') {
    form.base_url = 'https://api.deepseek.com'
    form.model = 'deepseek-chat'
  } else if (val === 'stepfun') {
    form.base_url = 'https://api.stepfun.com/v1'
    form.model = 'step-3.5-flash'
  } else if (val === 'lmstudio') {
    form.base_url = 'http://localhost:1234/v1'
    form.model = 'local-model'
  } else if (val === 'ollama') {
    form.base_url = 'http://localhost:11434/v1'
    form.model = 'llama3'
  } else if (val === 'local') {
    form.base_url = 'http://localhost:8000/v1'
    form.model = 'local-model'
  } else if (val === 'glm') {
    form.base_url = 'https://open.bigmodel.cn/api/paas/v4/'
    form.model = 'glm-4.7'
  } else if (val === 'openai') {
    form.base_url = 'https://api.openai.com/v1'
    form.model = 'gpt-4o-mini'
  }
}

const fetchConfig = async () => {
  try {
    const res = await axios.get('/api/config/llm')
    if (res.data.code === 0 && res.data.data) {
      Object.assign(form, res.data.data)
    }
  } catch (error) {
    console.error(error)
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    const res = await axios.put('/api/config/llm', form)
    if (res.data.code === 0) {
      ElMessage.success('Configuration saved')
    } else {
      ElMessage.error('Failed: ' + res.data.message)
    }
  } catch (error) {
    ElMessage.error('Error saving configuration')
  } finally {
    saving.value = false
  }
}

const testConnection = async () => {
  testing.value = true
  try {
    const res = await axios.post('/api/config/llm/test', form)
    if (res.data.code === 0) {
      ElMessage.success('连接成功: ' + res.data.data.reply)
    } else {
      ElMessage.error('连接失败: ' + res.data.message)
    }
  } catch (error) {
    const err = error as any
    const msg = err?.response?.data?.message || '无法连接到服务器'
    ElMessage.error('测试失败: ' + msg)
  } finally {
    testing.value = false
  }
}

onMounted(() => {
  fetchConfig()
})
</script>

<style scoped>
.settings {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.settings-hero {
  padding: 24px 28px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}

.settings-hero h1 {
  margin: 10px 0 6px;
  font-size: 26px;
}

.settings-hero p {
  margin: 0;
  color: #64748b;
}

.hero-actions {
  display: flex;
  gap: 12px;
}

.settings-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 0.8fr);
  gap: 24px;
}

.settings-card {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
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

.preset-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.preset-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 16px;
  border-radius: 14px;
  background: rgba(248, 250, 252, 0.9);
}

.preset-title {
  font-weight: 600;
}

.preset-meta {
  font-size: 12px;
  color: #64748b;
}

.preset-hint {
  font-size: 12px;
  color: #475569;
}

@media (max-width: 1024px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }

  .settings-hero {
    flex-direction: column;
    align-items: flex-start;
  }
}

@media (max-width: 768px) {
  .hero-actions {
    width: 100%;
    flex-direction: column;
  }
}
</style>
