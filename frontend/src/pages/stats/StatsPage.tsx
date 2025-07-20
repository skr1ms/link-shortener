import React, { useState, useEffect } from 'react'
import { statsApi } from '~/entities/stats/api'
import { StatsResponse } from '~/shared/types'

export const StatsPage: React.FC = () => {
  const [stats, setStats] = useState<StatsResponse | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  
  const [fromDate, setFromDate] = useState(() => {
    const date = new Date()
    date.setMonth(date.getMonth() - 1)
    return date.toISOString().split('T')[0]
  })
  
  const [toDate, setToDate] = useState(() => {
    const date = new Date()
    return date.toISOString().split('T')[0]
  })
  
  const [groupBy, setGroupBy] = useState<'day' | 'month'>('day')

  const loadStats = async () => {
    setLoading(true)
    setError('')
    
    try {
      const response = await statsApi.getStats(fromDate, toDate, groupBy)
      setStats(response)
    } catch (error: any) {
      setError(error.response?.data?.message || 'Ошибка загрузки статистики')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadStats()
  }, [fromDate, toDate, groupBy])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    loadStats()
  }

  return (
    <div>
      <h2 className="mb-20">Статистика переходов</h2>

      <div className="card mb-20">
        <form onSubmit={handleSubmit}>
          <div className="grid">
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '15px' }}>
              <div className="form-group">
                <label htmlFor="fromDate">Дата от</label>
                <input
                  type="date"
                  id="fromDate"
                  className="form-control"
                  value={fromDate}
                  onChange={(e) => setFromDate(e.target.value)}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="toDate">Дата до</label>
                <input
                  type="date"
                  id="toDate"
                  className="form-control"
                  value={toDate}
                  onChange={(e) => setToDate(e.target.value)}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="groupBy">Группировка</label>
                <select
                  id="groupBy"
                  className="form-control"
                  value={groupBy}
                  onChange={(e) => setGroupBy(e.target.value as 'day' | 'month')}
                >
                  <option value="day">По дням</option>
                  <option value="month">По месяцам</option>
                </select>
              </div>
            </div>

            <button type="submit" className="btn" disabled={loading}>
              {loading ? 'Загрузка...' : 'Обновить'}
            </button>
          </div>
        </form>
      </div>

      {error && (
        <div className="card">
          <div className="error">{error}</div>
        </div>
      )}

      {loading ? (
        <div className="loading">Загрузка статистики...</div>
      ) : stats ? (
        <div>
          <div className="grid grid-2 mb-20">
            <div className="card">
              <h3>Общая статистика</h3>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#007bff' }}>
                {stats.total_clicks}
              </div>
              <div>Всего переходов</div>
            </div>

            <div className="card">
              <h3>Период</h3>
              <div>{new Date(fromDate).toLocaleDateString('ru-RU')} - {new Date(toDate).toLocaleDateString('ru-RU')}</div>
            </div>
          </div>

          {stats.stats && stats.stats.length > 0 ? (
            <div className="card">
              <h3>Детальная статистика</h3>
              <table className="table">
                <thead>
                  <tr>
                    <th>Период</th>
                    <th>Переходов</th>
                  </tr>
                </thead>
                <tbody>
                  {stats.stats.map((stat, index) => (
                    <tr key={index}>
                      <td>
                        {stat.period_from === stat.period_to
                          ? new Date(stat.period_from).toLocaleDateString('ru-RU')
                          : `${new Date(stat.period_from).toLocaleDateString('ru-RU')} - ${new Date(stat.period_to).toLocaleDateString('ru-RU')}`
                        }
                      </td>
                      <td>{stat.clicks}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="card text-center">
              <p>Нет данных за выбранный период</p>
            </div>
          )}
        </div>
      ) : null}
    </div>
  )
} 