import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '~/entities/auth/api'

export const RegisterPage: React.FC = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setSuccess('')

    try {
      await authApi.register({ email, password, name })
      setSuccess('Регистрация успешна! Теперь вы можете войти')
      setTimeout(() => navigate('/login'), 2000)
    } catch (error: any) {
      setError(error.response?.data?.message || 'Ошибка регистрации')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container">
      <div style={{ marginTop: '100px' }}>
        <form className="form" onSubmit={handleSubmit}>
          <h2 className="text-center mb-20">Регистрация</h2>
          
          {error && <div className="error mb-20">{error}</div>}
          {success && <div className="success mb-20">{success}</div>}
          
          <div className="form-group">
            <label htmlFor="name">Имя</label>
            <input
              type="text"
              id="name"
              className="form-control"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              minLength={3}
            />
          </div>
          
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
              minLength={8}
            />
          </div>
          
          <button type="submit" className="btn" disabled={loading} style={{ width: '100%' }}>
            {loading ? 'Загрузка...' : 'Зарегистрироваться'}
          </button>
          
          <div className="text-center mt-20">
            <Link to="/login">Уже есть аккаунт? Войти</Link>
          </div>
        </form>
      </div>
    </div>
  )
} 