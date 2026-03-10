import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User, TokenPair } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)

  // 初始化时从 localStorage 恢复状态
  const init = () => {
    const storedUser = localStorage.getItem('user')
    const storedAccessToken = localStorage.getItem('accessToken')
    const storedRefreshToken = localStorage.getItem('refreshToken')

    if (storedUser && storedAccessToken) {
      user.value = JSON.parse(storedUser)
      accessToken.value = storedAccessToken
      refreshToken.value = storedRefreshToken
    }
  }

  const login = (userData: User, tokens: TokenPair) => {
    user.value = userData
    accessToken.value = tokens.accessToken
    refreshToken.value = tokens.refreshToken

    // 保存到 localStorage
    localStorage.setItem('user', JSON.stringify(userData))
    localStorage.setItem('accessToken', tokens.accessToken)
    localStorage.setItem('refreshToken', tokens.refreshToken)
  }

  const logout = () => {
    user.value = null
    accessToken.value = null
    refreshToken.value = null

    localStorage.removeItem('user')
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
  }

  const isAuthenticated = () => {
    return !!accessToken.value
  }

  init()

  return {
    user,
    accessToken,
    refreshToken,
    login,
    logout,
    isAuthenticated,
  }
})
