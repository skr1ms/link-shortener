import React from 'react'
import { BrowserRouter } from 'react-router-dom'
import { AuthProvider } from './providers/AuthProvider'
import { AppRouter } from './router/AppRouter'

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRouter />
      </AuthProvider>
    </BrowserRouter>
  )
} 