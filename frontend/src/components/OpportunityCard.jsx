import { Link } from 'react-router-dom'
import DeadlineCountdown from './DeadlineCountdown'

export default function OpportunityCard({ opportunity: opp }) {
  const typeLabel = { scholarship: 'Scholarship', internship: 'Internship', competition: 'Competition' }

  return (
    <Link to={`/opportunity/${opp.id}`} className="card">
      <div className="card-header">
        <span className="card-title">{opp.title}</span>
        {opp.deadline && <DeadlineCountdown deadline={opp.deadline} />}
      </div>
      <div className="card-meta">
        <span className={`badge badge-${opp.type}`}>{typeLabel[opp.type] || opp.type}</span>
        {opp.country && <span>{opp.country}</span>}
        <span>{opp.source}</span>
      </div>
      {opp.description && <p className="card-desc">{opp.description}</p>}
    </Link>
  )
}
