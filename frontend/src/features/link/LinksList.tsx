import React, { useState } from 'react'
import { linkApi } from '~/entities/link/api'
import { Link } from '~/shared/types'

interface LinksListProps {
  links: Link[]
  onLinkDeleted: (id: number) => void
  onLinkUpdated: (link: Link) => void
}

export const LinksList: React.FC<LinksListProps> = ({ links, onLinkDeleted, onLinkUpdated }) => {
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editUrl, setEditUrl] = useState('')
  const [editHash, setEditHash] = useState('')
  const [loading, setLoading] = useState<number | null>(null)

  const baseUrl = import.meta.env.NODE_ENV === 'production' ? '/api/link/' : 'http://localhost:8081/link/'

  const handleEdit = (link: Link) => {
    setEditingId(link.ID)
    setEditUrl(link.OriginalURL)
    setEditHash(link.Hash)
  }

  const handleCancelEdit = () => {
    setEditingId(null)
    setEditUrl('')
    setEditHash('')
  }

  const handleUpdate = async (id: number) => {
    setLoading(id)
    try {
      const updatedLink = await linkApi.updateLink(id, { url: editUrl, hash: editHash })
      onLinkUpdated(updatedLink)
      setEditingId(null)
      setEditUrl('')
      setEditHash('')
    } catch (error) {
      console.error('Error updating link:', error)
    } finally {
      setLoading(null)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Вы уверены, что хотите удалить эту ссылку?')) {
      return
    }

    setLoading(id)
    try {
      await linkApi.deleteLink(id)
      onLinkDeleted(id)
    } catch (error) {
      console.error('Error deleting link:', error)
    } finally {
      setLoading(null)
    }
  }

  const copyToClipboard = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text)
      alert('Ссылка скопирована!')
    } catch (error) {
      console.error('Error copying to clipboard:', error)
    }
  }

  if (links.length === 0) {
    return (
      <div className="card text-center">
        <p>У вас пока нет созданных ссылок</p>
      </div>
    )
  }

  return (
    <div className="card">
      <h3>Ваши ссылки</h3>
      <table className="table">
        <thead>
          <tr>
            <th>Короткая ссылка</th>
            <th>Оригинальная ссылка</th>
            <th>Создана</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {links.map((link) => (
            <tr key={link.ID}>
              <td>
                <div>
                  {editingId === link.ID ? (
                    <input
                      type="text"
                      className="form-control"
                      value={editHash}
                      onChange={(e) => setEditHash(e.target.value)}
                    />
                  ) : (
                    <div>
                      <a
                        href={baseUrl + link.Hash}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {baseUrl + link.Hash}
                      </a>
                      <button
                        className="btn btn-secondary"
                        onClick={() => copyToClipboard(baseUrl + link.Hash)}
                        style={{ marginLeft: '10px', padding: '4px 8px', fontSize: '12px' }}
                      >
                        Копировать
                      </button>
                    </div>
                  )}
                </div>
              </td>
              <td>
                {editingId === link.ID ? (
                  <input
                    type="url"
                    className="form-control"
                    value={editUrl}
                    onChange={(e) => setEditUrl(e.target.value)}
                  />
                ) : (
                  <a href={link.OriginalURL} target="_blank" rel="noopener noreferrer">
                    {link.OriginalURL}
                  </a>
                )}
              </td>
              <td>{new Date(link.CreatedAt).toLocaleDateString('ru-RU')}</td>
              <td>
                <div className="actions">
                  {editingId === link.ID ? (
                    <>
                      <button
                        className="btn"
                        onClick={() => handleUpdate(link.ID)}
                        disabled={loading === link.ID}
                      >
                        {loading === link.ID ? 'Сохранение...' : 'Сохранить'}
                      </button>
                      <button
                        className="btn btn-secondary"
                        onClick={handleCancelEdit}
                        disabled={loading === link.ID}
                      >
                        Отмена
                      </button>
                    </>
                  ) : (
                    <>
                      <button
                        className="btn btn-secondary"
                        onClick={() => handleEdit(link)}
                        disabled={loading === link.ID}
                      >
                        Изменить
                      </button>
                      <button
                        className="btn btn-danger"
                        onClick={() => handleDelete(link.ID)}
                        disabled={loading === link.ID}
                      >
                        {loading === link.ID ? 'Удаление...' : 'Удалить'}
                      </button>
                    </>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
} 