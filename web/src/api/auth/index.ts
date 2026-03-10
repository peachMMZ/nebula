import type { ApiResponse, LoginRequest, LoginResponse, User } from '@/types/api'
import { apiClient } from '@/api/client'

export const authApi = {
  async login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return apiClient.request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async getProfile(): Promise<ApiResponse<User>> {
    return apiClient.request<User>('/auth/profile')
  },

  async refreshToken(refreshToken: string): Promise<ApiResponse<{ accessToken: string; refreshToken: string; expiresIn: number }>> {
    return apiClient.request<{ accessToken: string; refreshToken: string; expiresIn: number }>('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refreshToken }),
    })
  },
}