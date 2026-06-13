import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import DeadlineCountdown from '../components/DeadlineCountdown'

export default function DetailPage({ apiBase }) {
  const { id } = useParams()
  const [opp, setOpp] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch(`${apiBase}/opportunities/${id}`)
      .then(r => {
        if (!r.ok) throw new Error('not found')
        return r.json()
      })
      .then(setOpp)
      .catch(() => setOpp(null))
      .finally(() => setLoading(false))
  }, [apiBase, id])

  if (loading) return <div className="loading">Loading...</div>
  if (!opp) return <div className="empty-state"><h2>Opportunity not found</h2></div>

  const typeLabel = { scholarship: 'Scholarship', internship: 'Internship', competition: 'Competition' }

  return (
    <div className="detail-page">
      <Link to="/" className="back-link">&larr; Back to search</Link>

      <h1>{opp.title}</h1>

      <div className="detail-meta">
        <span className={`badge badge-${opp.type}`}>{typeLabel[opp.type] || opp.type}</span>
        {opp.country && <span style={{ fontSize: 14, color: '#64748b' }}>{opp.country}</span>}
        {opp.source && <span style={{ fontSize: 14, color: '#64748b' }}>Source: {opp.source}</span>}
      </div>

      {opp.deadline && (
        <div className="detail-section">
          <h3>Deadline</h3>
          <DeadlineCountdown deadline={opp.deadline} />
        </div>
      )}

      {opp.description && (
        <div className="detail-section">
          <h3>Description</h3>
          <p>{opp.description}</p>
        </div>
      )}

      {opp.eligibility && (
        <div className="detail-section">
          <h3>Eligibility</h3>
          <p>{opp.eligibility}</p>
        </div>
      )}

      {opp.funding && (
        <div className="detail-section">
          <h3>Funding</h3>
          <p>{opp.funding}</p>
        </div>
      )}

      <div style={{ marginTop: 24 }}>
        <a href={opp.url} target="_blank" rel="noopener noreferrer" className="apply-btn">
          Apply Now &rarr;
        </a>
      </div>
    </div>
  )
}
