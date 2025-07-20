import axios from 'axios'

const BASE_URL = import.meta.env.NODE_ENV === 'production' ? '/api' : 'http://localhost:8081'

export const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('accessToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      const refreshToken = localStorage.getItem('refreshToken')
      if (refreshToken) {
        try {
          const response = await axios.post(`${BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken
          })
          
          const { access_token, refresh_token: newRefreshToken } = response.data
          localStorage.setItem('accessToken', access_token)
          localStorage.setItem('refreshToken', newRefreshToken)
          
          originalRequest.headers.Authorization = `Bearer ${access_token}`
          return api(originalRequest)
        } catch (refreshError) {
          localStorage.removeItem('accessToken')
          localStorage.removeItem('refreshToken')
          window.location.href = '/login'
        }
      }
    }
    
    return Promise.reject(error)
  }
) 