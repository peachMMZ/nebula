<template>
  <n-layout class="h-screen w-full" has-sider>
    <!-- 侧边栏 -->
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="240"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
      class="h-full"
    >
      <div class="flex h-[64px] items-center justify-center overflow-hidden whitespace-nowrap px-4 py-2">
        <div class="flex items-center gap-2">
            <div class="h-8 w-8 flex-shrink-0 rounded-lg bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center text-white font-bold">N</div>
            <span v-show="!collapsed" class="text-lg font-bold text-slate-700 transition-opacity duration-300">Nebula</span>
        </div>
      </div>
      
      <n-menu
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="22"
        :options="menuOptions"
        :value="activeKey"
        @update:value="handleMenuUpdate"
      />
    </n-layout-sider>

    <n-layout class="h-full flex flex-col">
      <!-- 顶部 Header -->
      <n-layout-header bordered class="h-16 flex items-center justify-between px-6 bg-white z-10">
        <!-- 左侧：面包屑等 -->
        <div class="flex items-center">
            <!-- 可以在这里放面包屑 -->
        </div>

        <!-- 右侧：用户菜单 -->
        <div class="flex items-center gap-4">
             <n-dropdown :options="userOptions" @select="handleUserSelect">
                <div class="flex items-center gap-2 cursor-pointer hover:bg-gray-100 px-2 py-1 rounded transition-colors">
                    <n-avatar round size="small" class="bg-emerald-500 text-white">
                        A
                    </n-avatar>
                    <span class="text-sm font-medium text-slate-700">Admin</span>
                </div>
            </n-dropdown>
        </div>
      </n-layout-header>

      <!-- 主要内容区域 -->
      <n-layout-content content-style="padding: 24px; min-height: calc(100vh - 64px);" class="bg-gray-50/50">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
             <component :is="Component" />
          </transition>
        </router-view>
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { ref, h, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  NLayout, 
  NLayoutSider, 
  NLayoutHeader, 
  NLayoutContent, 
  NMenu, 
  NIcon, 
  NButton,
  NAvatar,
  NDropdown,
  type MenuOption 
} from 'naive-ui'
import { LayoutDashboard, AppWindow, Settings, LogOut } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

function renderIcon(icon: any) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)

// 菜单配置
const menuOptions: MenuOption[] = [
  {
    label: '概览',
    key: 'home',
    icon: renderIcon(LayoutDashboard)
  },
  {
    label: '应用管理', // Uncommented and updated
    key: 'apps',
    icon: renderIcon(AppWindow)
  },
  {
    label: '系统设置', // Uncommented and updated
    key: 'settings',
    icon: renderIcon(Settings)
  }
]

// 当前选中的菜单项，根据路由自动匹配
const activeKey = computed(() => {
    return (route.name as string) || 'home'
})

const handleMenuUpdate = (key: string) => {
    if (key === 'home') {
        router.push('/')
    } else {
        router.push({ name: key })
    }
}

// 用户下拉菜单
const userOptions = [
    {
        label: '退出登录',
        key: 'logout',
        icon: renderIcon(LogOut)
    }
]

const handleUserSelect = (key: string) => {
    if (key === 'logout') {
        authStore.logout()
        router.push('/login')
    }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
