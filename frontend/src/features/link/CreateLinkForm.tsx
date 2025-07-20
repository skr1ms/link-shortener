import React, { useState } from 'react'
import { linkApi } from '~/entities/link/api'
import { Link } from '~/shared/types'

interface CreateLinkFormProps {
  onLinkCreated: (link: Link) => void
}

export const CreateLinkForm: React.FC<CreateLinkFormProps> = ({ onLinkCreated }) => {
  const [url, setUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')
    setSuccess('')

    try {
      const newLink = await linkApi.createLink({ url })
      onLinkCreated(newLink)
      setUrl('')
      setSuccess('Ссылка успешно создана!')
      setTimeout(() => setSuccess(''), 3000)
    } catch (error: any) {
      setError(error.response?.data?.message || 'Ошибка создания ссылки')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="card">
      <h3>Создать новую ссылку</h3>
      
      {error && <div className="error mb-20">{error}</div>}
      {success && <div className="success mb-20">{success}</div>}
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="url">URL</label>
          <input
            type="url"
            id="url"
            className="form-control"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://example.com"
            required
          />
        </div>
        
        <button type="submit" className="btn" disabled={loading}>
          {loading ? 'Создание...' : 'Создать ссылку'}
        </button>
      </form>
    </div>
  )
} 