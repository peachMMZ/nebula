<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-100 p-4">
    <!-- 背景装饰圆 -->
    <div class="fixed top-0 left-0 w-full h-full overflow-hidden pointer-events-none z-0">
      <div
        class="absolute -top-[10%] -left-[5%] w-[600px] h-[600px] bg-indigo-200 rounded-full blur-[120px] opacity-40"
      ></div>
      <div
        class="absolute top-[40%] -right-[5%] w-[500px] h-[500px] bg-emerald-200 rounded-full blur-[100px] opacity-40"
      ></div>
    </div>

    <!-- 主卡片容器 -->
    <div
      class="relative z-10 flex w-full max-w-5xl overflow-hidden rounded-3xl bg-white shadow-2xl shadow-slate-200/50"
    >
      <!-- 左侧装饰区 -->
      <div
        class="hidden w-1/2 flex-col justify-between bg-slate-900 p-12 text-white lg:flex relative overflow-hidden"
      >
        <!-- 装饰背景 -->
        <div
          class="absolute inset-0 bg-linear-to-br from-slate-800 via-slate-900 to-black z-0"
        ></div>
        <div
          class="absolute top-0 right-0 w-[400px] h-[400px] bg-linear-to-br from-emerald-500/20 to-teal-500/20 rounded-full blur-[80px] pointer-events-none z-0 transform translate-x-1/3 -translate-y-1/3"
        ></div>

        <!-- 内容 -->
        <div class="relative z-10">
          <div class="flex items-center gap-3 mb-8">
            <div
              class="w-8 h-8 rounded-lg bg-linear-to-br from-emerald-400 to-teal-600 flex items-center justify-center"
            >
              <span class="font-bold text-white">N</span>
            </div>
            <span class="text-xl font-bold tracking-wider">NEBULA</span>
          </div>

          <h2 class="text-4xl font-bold leading-tight mb-6">
            Manage Your Updates <br />
            <span class="text-transparent bg-clip-text bg-linear-to-r from-emerald-400 to-teal-300"
              >Efficiently & Securely</span
            >
          </h2>

          <p class="text-slate-400 leading-relaxed max-w-sm">
            Nebula 是一个现代化的应用更新管理平台，为您提供稳定、高效的版本分发服务。
          </p>
        </div>

        <div class="relative z-10 mt-12 grid grid-cols-2 gap-6">
          <div class="space-y-1">
            <h3 class="text-2xl font-bold text-white">0.0k+</h3>
            <p class="text-xs text-slate-500 uppercase tracking-wide">Active Users</p>
          </div>
          <div class="space-y-1">
            <h3 class="text-2xl font-bold text-emerald-400">99.9%</h3>
            <p class="text-xs text-slate-500 uppercase tracking-wide">Uptime</p>
          </div>
        </div>
      </div>

      <!-- 右侧表单区 -->
      <div class="flex w-full flex-col justify-center bg-white p-8 lg:w-1/2 lg:p-16">
        <div class="mb-8">
          <h1 class="text-2xl font-bold text-slate-800 lg:text-3xl">欢迎回来</h1>
          <p class="mt-2 text-sm text-slate-500">请输入您的账号密码登录系统</p>
        </div>

        <n-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          size="large"
          @submit.prevent="handleSubmit"
        >
          <n-form-item path="username" label="用户名">
            <n-input
              v-model:value="formData.username"
              :disabled="loading"
              @keydown.enter.prevent
            >
              <template #prefix>
                <n-icon><User /></n-icon>
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="password" label="密码">
            <n-input
              v-model:value="formData.password"
              type="password"
              show-password-on="click"
              :disabled="loading"
              @keydown.enter.prevent
              @keypress.enter="handleSubmit"
            >
              <template #prefix>
                <n-icon><Lock /></n-icon>
              </template>
            </n-input>
          </n-form-item>

          <div class="flex items-center justify-between mb-6">
            <div class="text-sm text-slate-500">
              <span class="cursor-pointer hover:text-emerald-600 transition-colors">忘记密码?</span>
            </div>
          </div>

          <div v-if="error" class="mb-6">
            <n-alert type="error" closable :show-icon="true" class="text-sm">
              {{ error }}
            </n-alert>
          </div>

          <n-button
            type="primary"
            block
            size="large"
            :loading="loading"
            @click="handleSubmit"
            class="h-12 text-base font-semibold shadow-emerald-500/20 shadow-lg"
          >
            登 录
          </n-button>
        </n-form>

        <div class="mt-8 text-center text-xs text-slate-400">默认测试账号：admin / 123456</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'
import {
  NForm,
  NFormItem,
  NInput,
  NButton,
  NAlert,
  NIcon,
  type FormInst,
  type FormRules,
} from 'naive-ui'
import { User, Lock } from 'lucide-vue-next'

defineOptions({
  name: 'LoginView',
})

const router = useRouter()
const authStore = useAuthStore()
const formRef = ref<FormInst | null>(null)

const formData = ref({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const rules: FormRules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: ['input', 'blur'],
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: ['input', 'blur'],
  },
}

const handleSubmit = async () => {
  error.value = ''
  await formRef.value?.validate()
  loading.value = true

  try {
    const response = await authApi.login(formData.value)

    if (response.code === 0) {
      // 保存登录状态
      authStore.login(response.data.user, response.data.tokens)

      // 跳转到首页
      router.push('/')
    } else {
      error.value = response.message || '登录失败'
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : '网络错误，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>
