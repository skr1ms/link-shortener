import { api } from '~/shared/api'
import { LoginRequest, RegisterRequest, LoginResponse, RefreshTokenRequest, RefreshTokenResponse, User } from '~/shared/types'

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', data)
    return response.data
  },

  register: async (data: RegisterRequest): Promise<User> => {
    const response = await api.post('/auth/register', data)
    return response.data
  },

  refresh: async (data: RefreshTokenRequest): Promise<RefreshTokenResponse> => {
    const response = await api.post('/auth/refresh', data)
    return response.data
  }
} 