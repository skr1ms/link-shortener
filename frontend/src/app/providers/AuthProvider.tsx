import React, { createContext, useContext, useState, useEffect } from 'react'

interface AuthContextType {
  isAuthenticated: boolean
  accessToken: string | null
  refreshToken: string | null
  login: (accessToken: string, refreshToken: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [accessToken, setAccessToken] = useState<string | null>(null)
  const [refreshToken, setRefreshToken] = useState<string | null>(null)

  useEffect(() => {
    const storedAccessToken = localStorage.getItem('accessToken')
    const storedRefreshToken = localStorage.getItem('refreshToken')
    
    if (storedAccessToken && storedRefreshToken) {
      setAccessToken(storedAccessToken)
      setRefreshToken(storedRefreshToken)
    }
  }, [])

  const login = (newAccessToken: string, newRefreshToken: string) => {
    localStorage.setItem('accessToken', newAccessToken)
    localStorage.setItem('refreshToken', newRefreshToken)
    setAccessToken(newAccessToken)
    setRefreshToken(newRefreshToken)
  }

  const logout = () => {
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    setAccessToken(null)
    setRefreshToken(null)
  }

  const value: AuthContextType = {
    isAuthenticated: !!accessToken,
    accessToken,
    refreshToken,
    login,
    logout
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
} 