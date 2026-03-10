import type { ApiResponse, App } from '@/types/api'
import { apiClient } from '@/api/client'

export interface GetAppsParams {
  name?: string
  description?: string
}

export interface CreateAppRequest {
  name: string
  description?: string
}

export interface UpdateAppRequest {
  name?: string
  description?: string
}

export const appApi = {
  async getApps(params?: GetAppsParams): Promise<ApiResponse<App[]>> {
    const queryString = params ? '?' + new URLSearchParams(params as Record<string, string>).toString() : ''
    return apiClient.request<App[]>(`/apps${queryString}`)
  },

  async getApp(id: string): Promise<ApiResponse<App>> {
    return apiClient.request<App>(`/apps/${id}`)
  },

  async createApp(data: CreateAppRequest): Promise<ApiResponse<App>> {
    return apiClient.request<App>('/apps', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async updateApp(id: string, data: UpdateAppRequest): Promise<ApiResponse<App>> {
    return apiClient.request<App>(`/apps/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  async deleteApp(id: string): Promise<ApiResponse<{ message: string }>> {
    return apiClient.request<{ message: string }>(`/apps/${id}`, {
      method: 'DELETE',
    })
  },
}