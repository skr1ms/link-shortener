import React, { useState, useEffect } from 'react'
import { linkApi } from '~/entities/link/api'
import { Link } from '~/shared/types'
import { CreateLinkForm } from '~/features/link/CreateLinkForm'
import { LinksList } from '~/features/link/LinksList'

export const DashboardPage: React.FC = () => {
  const [links, setLinks] = useState<Link[]>([])
  const [totalCount, setTotalCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [currentPage, setCurrentPage] = useState(1)
  const limit = 10

  const loadLinks = async () => {
    setLoading(true)
    try {
      const offset = (currentPage - 1) * limit
      const response = await linkApi.getLinks(limit, offset)
      setLinks(response.links || [])
      setTotalCount(response.count)
    } catch (error) {
      console.error('Error loading links:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadLinks()
  }, [currentPage])

  const handleLinkCreated = (newLink: Link) => {
    setLinks([newLink, ...links])
    setTotalCount(totalCount + 1)
  }

  const handleLinkDeleted = (id: number) => {
    setLinks(links.filter(link => link.ID !== id))
    setTotalCount(totalCount - 1)
  }

  const handleLinkUpdated = (updatedLink: Link) => {
    setLinks(links.map(link => 
      link.ID === updatedLink.ID ? updatedLink : link
    ))
  }

  const totalPages = Math.ceil(totalCount / limit)

  return (
    <div>
      <h2 className="mb-20">Управление ссылками</h2>
      
      <div className="grid grid-2">
        <div>
          <CreateLinkForm onLinkCreated={handleLinkCreated} />
        </div>
        
        <div>
          <div className="card">
            <h3>Статистика</h3>
            <p>Всего ссылок: {totalCount}</p>
          </div>
        </div>
      </div>

      {loading ? (
        <div className="loading">Загрузка...</div>
      ) : (
        <>
          <LinksList
            links={links}
            onLinkDeleted={handleLinkDeleted}
            onLinkUpdated={handleLinkUpdated}
          />
          
          {totalPages > 1 && (
            <div className="text-center mt-20">
              <div className="actions">
                <button
                  className="btn btn-secondary"
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                >
                  Назад
                </button>
                <span style={{ margin: '0 20px' }}>
                  Страница {currentPage} из {totalPages}
                </span>
                <button
                  className="btn btn-secondary"
                  onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                  disabled={currentPage === totalPages}
                >
                  Вперёд
                </button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  )
} 