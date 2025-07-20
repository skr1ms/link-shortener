import { api } from '~/shared/api'
import { StatsResponse } from '~/shared/types'

export const statsApi = {
  getStats: async (from: string, to: string, by: 'day' | 'month'): Promise<StatsResponse> => {
    const response = await api.get(`/stats?from=${from}&to=${to}&by=${by}`)
    return response.data
  }
} 