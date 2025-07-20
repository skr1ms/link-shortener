import React from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '~/app/providers/AuthProvider'

interface LayoutProps {
  children: React.ReactNode
}

export const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { logout } = useAuth()

  const handleLogout = () => {
    logout()
  }

  return (
    <div>
      <header className="header">
        <div className="container">
          <div className="header-content">
            <h1>Link Shortener</h1>
            <nav className="nav">
              <Link to="/">Ссылки</Link>
              <Link to="/stats">Статистика</Link>
              <button className="btn btn-secondary" onClick={handleLogout}>
                Выйти
              </button>
            </nav>
          </div>
        </div>
      </header>
      <main className="container">
        {children}
      </main>
    </div>
  )
} 