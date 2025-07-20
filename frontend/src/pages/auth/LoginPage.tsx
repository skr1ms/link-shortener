import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '~/app/providers/AuthProvider'
import { authApi } from '~/entities/auth/api'

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      const response = await authApi.login({ email, password })
      login(response.access_token, response.refresh_token)
      navigate('/')
    } catch (error: any) {
      setError(error.response?.data?.message || 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container">
      <div style={{ marginTop: '100px' }}>
        <form className="form" onSubmit={handleSubmit}>
          <h2 className="text-center mb-20">Вход</h2>
          
          {error && <div className="error mb-20">{error}</div>}
          
          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              className="form-control"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">Пароль</label>
            <input
              type="password"
              id="password"
              className="form-control"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          
          <button type="submit" className="btn" disabled={loading} style={{ width: '100%' }}>
            {loading ? 'Загрузка...' : 'Войти'}
          </button>
          
          <div className="text-center mt-20">
            <Link to="/register">Нет аккаунта? Регистрация</Link>
          </div>
        </form>
      </div>
    </div>
  )
} 