import { api } from '~/shared/api'
import { Link, CreateLinkRequest, UpdateLinkRequest, GetLinksResponse } from '~/shared/types'

export const linkApi = {
  getLinks: async (limit: number, offset: number): Promise<GetLinksResponse> => {
    const response = await api.get(`/link?limit=${limit}&offset=${offset}`)
    return response.data
  },

  createLink: async (data: CreateLinkRequest): Promise<Link> => {
    const response = await api.post('/link', data)
    return response.data
  },

  updateLink: async (id: number, data: UpdateLinkRequest): Promise<Link> => {
    const response = await api.patch(`/link/${id}`, data)
    return response.data
  },

  deleteLink: async (id: number): Promise<void> => {
    await api.delete(`/link/${id}`)
  }
} 