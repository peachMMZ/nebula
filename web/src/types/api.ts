// API 类型定义
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface User {
  id: string
  username: string
  email: string
  role: string
  createdAt: string
  updatedAt: string
}

export interface TokenPair {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  user: User
  tokens: TokenPair
}
