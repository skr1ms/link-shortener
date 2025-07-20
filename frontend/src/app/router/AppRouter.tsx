import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from '../providers/AuthProvider'
import { LoginPage } from '~/pages/auth/LoginPage'
import { RegisterPage } from '~/pages/auth/RegisterPage'
import { DashboardPage } from '~/pages/dashboard/DashboardPage'
import { StatsPage } from '~/pages/stats/StatsPage'
import { Layout } from '~/widgets/layout/Layout'

export const AppRouter: React.FC = () => {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return (
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    )
  }

  return (
    <Layout>
      <Routes>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/stats" element={<StatsPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Layout>
  )
} 