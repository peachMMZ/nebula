<template>
  <div class="space-y-6">
    <!-- 头部欢迎语 -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-slate-800">仪表盘 / Dashboard</h1>
        <p class="text-slate-500 mt-1">欢迎回来，Admin</p>
      </div>
      <div class="flex items-center gap-2">
        <n-button type="primary" size="medium" class="px-6 rounded-lg shadow-sm">
          <template #icon>
            <n-icon :component="Plus" />
          </template>
          发布新版本
        </n-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
      <!-- 应用总数 -->
      <n-card
        :bordered="false"
        class="rounded-xl shadow-sm hover:shadow-md transition-shadow duration-300"
      >
        <template #header>
          <span class="text-sm font-medium text-slate-500">应用总数</span>
        </template>
        <template #header-extra>
          <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600">
            <n-icon size="20" :component="Smartphone" />
          </div>
        </template>
        <div class="mt-2">
          <span class="text-3xl font-bold text-slate-800">12</span>
          <span
            class="ml-2 text-xs font-medium text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full"
            >+2 本周</span
          >
        </div>
      </n-card>

      <!-- 活跃版本 -->
      <n-card
        :bordered="false"
        class="rounded-xl shadow-sm hover:shadow-md transition-shadow duration-300"
      >
        <template #header>
          <span class="text-sm font-medium text-slate-500">活跃版本</span>
        </template>
        <template #header-extra>
          <div class="p-2 bg-emerald-50 rounded-lg text-emerald-600">
            <n-icon size="20" :component="Rocket" />
          </div>
        </template>
        <div class="mt-2">
          <span class="text-3xl font-bold text-slate-800">48</span>
        </div>
      </n-card>

      <!-- 下载次数 -->
      <n-card
        :bordered="false"
        class="rounded-xl shadow-sm hover:shadow-md transition-shadow duration-300"
      >
        <template #header>
          <span class="text-sm font-medium text-slate-500">总下载量</span>
        </template>
        <template #header-extra>
          <div class="p-2 bg-blue-50 rounded-lg text-blue-600">
            <n-icon size="20" :component="Download" />
          </div>
        </template>
        <div class="mt-2">
          <span class="text-3xl font-bold text-slate-800">1.2k</span>
          <span class="ml-2 text-xs font-medium text-rose-600 bg-rose-50 px-2 py-0.5 rounded-full"
            >-5% 同比</span
          >
        </div>
      </n-card>

      <!-- 存储用量 -->
      <n-card
        :bordered="false"
        class="rounded-xl shadow-sm hover:shadow-md transition-shadow duration-300"
      >
        <template #header>
          <span class="text-sm font-medium text-slate-500">存储用量</span>
        </template>
        <template #header-extra>
          <div class="p-2 bg-amber-50 rounded-lg text-amber-600">
            <n-icon size="20" :component="HardDrive" />
          </div>
        </template>
        <div class="mt-2">
          <span class="text-3xl font-bold text-slate-800">4.2GB</span>
          <span class="ml-2 text-xs font-medium text-slate-400">/ 10GB</span>
        </div>
      </n-card>
    </div>

    <!-- 最近更新表格 -->
    <n-card title="最近发布动态" :bordered="false" class="rounded-xl shadow-sm mt-6">
      <n-data-table :columns="columns" :data="data" :bordered="false" class="mt-4" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NCard, NButton, NDataTable, NIcon } from 'naive-ui'
import { Smartphone, Rocket, Download, HardDrive, Plus } from 'lucide-vue-next'

defineOptions({
  name: 'HomeView',
})

const data = ref([
  { app: 'Nebula App', version: 'v1.0.2', time: '2026-03-09 14:00', status: '已发布' },
  { app: 'Nebula Admin', version: 'v2.1.0', time: '2026-03-08 10:00', status: '审核中' },
])

const columns = ref([
  { title: '应用名称', key: 'app' },
  { title: '版本', key: 'version' },
  { title: '发布时间', key: 'time' },
  { title: '状态', key: 'status' },
])
</script>
