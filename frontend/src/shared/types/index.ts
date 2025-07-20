export interface User {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  Email: string
  Name: string
}

export interface Link {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  OriginalURL: string
  Hash: string
}

export interface Stats {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  LinkId: number
  ClickCount: number
  Date: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface RefreshTokenResponse {
  access_token: string
  refresh_token: string
}

export interface CreateLinkRequest {
  url: string
}

export interface UpdateLinkRequest {
  url: string
  hash: string
}

export interface GetLinksResponse {
  links: Link[]
  count: number
}

export interface StatsPayload {
  period_from: string
  period_to: string
  clicks: number
}

export interface StatsResponse {
  stats: StatsPayload[]
  total_clicks: number
} 