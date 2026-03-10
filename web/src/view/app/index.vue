<template>
  <div class="space-y-6">
    <!-- 头部 -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">应用管理</h1>
        <p class="text-slate-500 mt-1">管理所有应用程序</p>
      </div>
      <div class="flex items-center gap-2">
        <!-- 搜索框 -->
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索应用名称或描述"
          clearable
          @input="handleSearch"
          class="w-full md:w-64"
        >
          <template #prefix>
            <n-icon :component="Search" />
          </template>
        </n-input>
        <n-button type="primary" size="medium" class="px-6 rounded-lg shadow-sm" @click="showCreateDialog = true">
          <template #icon>
            <n-icon :component="Plus" />
          </template>
          新建应用
        </n-button>
      </div>
    </div>

    <!-- 加载状态 -->
    <n-skeleton v-if="loading" animated :rows="5" class="rounded-xl" />

    <!-- 错误提示 -->
    <div v-else-if="error" class="rounded-xl">
      <n-alert
        type="error"
        title="加载失败"
        description="无法加载应用数据，请稍后重试"
        show-icon
        class="rounded-xl mb-4"
      />
      <div class="flex justify-end">
        <n-button size="small" @click="fetchApps">重试</n-button>
      </div>
    </div>

    <!-- 应用卡片列表 -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <n-card
        v-for="app in filteredApps"
        :key="app.id"
        :bordered="false"
        class="rounded-xl shadow-sm hover:shadow-md transition-shadow duration-300 overflow-hidden"
      >
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-slate-800">{{ app.name }}</h3>
            <div class="flex items-center gap-2">
              <n-button quaternary circle size="small">
                <n-icon :component="Edit" />
              </n-button>
              <n-button quaternary circle size="small" type="error" @click="handleDelete(app.id, app.name)">
                <n-icon :component="Trash2" />
              </n-button>
            </div>
          </div>
        </template>
        <div class="space-y-4">
          <p class="text-slate-600 text-sm" v-if="app.description">{{ app.description }}</p>
          <p class="text-slate-400 text-sm italic" v-else>无描述</p>
          <div class="pt-4 border-t border-slate-100">
            <div class="flex justify-between text-xs text-slate-500">
              <span>创建时间: {{ formatDate(app.createdAt) }}</span>
              <span>更新时间: {{ formatDate(app.updatedAt) }}</span>
            </div>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && !error && apps.length === 0" class="flex flex-col items-center justify-center py-16 text-center">
      <n-icon :component="FolderOpen" size="48" class="text-slate-300 mb-4" />
      <h3 class="text-lg font-semibold text-slate-700 mb-2">暂无应用</h3>
      <p class="text-slate-500 mb-6">您还没有创建任何应用</p>
      <n-button type="primary" @click="showCreateDialog = true">
        <template #icon>
          <n-icon :component="Plus" />
        </template>
        新建应用
      </n-button>
    </div>

    <!-- 创建应用对话框 -->
    <n-modal
      v-model:show="showCreateDialog"
      preset="card"
      title="创建新应用"
      size="large"
      style="width: 480px;"
    >
      <n-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
      >
        <n-form-item path="name" label="应用名称" required>
          <n-input
            v-model:value="createForm.name"
            placeholder="请输入应用名称"
            :disabled="creating"
          />
        </n-form-item>
        <n-form-item path="description" label="应用描述">
          <n-input
            v-model:value="createForm.description"
            type="textarea"
            placeholder="请输入应用描述"
            :disabled="creating"
            :rows="3"
          />
        </n-form-item>
      </n-form>
      <template #footer>
        <div class="flex justify-end items-center gap-2">
          <n-button @click="showCreateDialog = false" :disabled="creating">
            取消
          </n-button>
          <n-button type="primary" @click="handleCreate" :loading="creating">
            创建
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { NCard, NButton, NInput, NIcon, NSkeleton, NAlert, NModal, NForm, NFormItem, useMessage, useDialog, type FormInst } from 'naive-ui'
import { Search, Plus, Edit, Trash2, FolderOpen } from 'lucide-vue-next'
import { appApi } from '@/api/app'
import type { App } from '@/types/api'

defineOptions({
  name: 'AppView',
})

const message = useMessage()
const dialog = useDialog()

// 状态管理
const apps = ref<App[]>([])
const loading = ref(false)
const error = ref(false)
const searchQuery = ref('')

// 创建应用相关状态
const showCreateDialog = ref(false)
const creating = ref(false)
const createForm = ref({
  name: '',
  description: ''
})
const createFormRef = ref<FormInst | null>(null)
const createRules = ref({
  name: {
    required: true,
    message: '请输入应用名称',
    trigger: ['input', 'blur']
  }
})

// 格式化日期
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 搜索过滤
const filteredApps = computed(() => {
  if (!searchQuery.value) return apps.value

  const query = searchQuery.value.toLowerCase()
  return apps.value.filter(app =>
    app.name.toLowerCase().includes(query) ||
    (app.description && app.description.toLowerCase().includes(query))
  )
})

// 获取应用列表
const fetchApps = async () => {
  loading.value = true
  error.value = false

  try {
    const response = await appApi.getApps()
    apps.value = response.data
  } catch (err) {
    console.error('Failed to fetch apps:', err)
    error.value = true
  } finally {
    loading.value = false
  }
}

// 监听搜索查询变化
const handleSearch = () => {
  // 客户端过滤，不需要重新请求API
}

// 处理创建应用
const handleCreate = async () => {
  if (!createFormRef.value) return

  try {
    await createFormRef.value.validate()
    creating.value = true

    const response = await appApi.createApp(createForm.value)

    if (response.code === 0) {
      message.success('应用创建成功')
      showCreateDialog.value = false
      // 清空表单
      createForm.value = {
        name: '',
        description: ''
      }
      // 重新获取应用列表
      await fetchApps()
    } else {
      message.error(response.message || '创建失败')
    }
  } catch (err) {
    message.error(err instanceof Error ? err.message : '创建失败')
  } finally {
    creating.value = false
  }
}

// 处理删除应用
const handleDelete = (id: string, name: string) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除应用 "${name}" 吗？此操作不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const response = await appApi.deleteApp(id)
        if (response.code === 0) {
          message.success('应用删除成功')
          await fetchApps()
        } else {
          message.error(response.message || '删除失败')
        }
      } catch (err) {
        message.error(err instanceof Error ? err.message : '删除失败')
      }
    }
  })
}

// 组件挂载时获取数据
onMounted(() => {
  fetchApps()
})
</script>

<style scoped>
/* 自定义样式 */
</style>
