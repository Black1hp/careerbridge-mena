import { useState, useEffect } from 'react'
import OpportunityCard from '../components/OpportunityCard'

export default function SearchPage({ apiBase }) {
  const [query, setQuery] = useState('')
  const [type, setType] = useState('')
  const [country, setCountry] = useState('')
  const [countries, setCountries] = useState([])
  const [results, setResults] = useState([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetch(`${apiBase}/countries`)
      .then(r => r.json())
      .then(setCountries)
      .catch(() => {})
  }, [apiBase])

  useEffect(() => {
    setLoading(true)
    const params = new URLSearchParams()
    if (query) params.set('q', query)
    if (type) params.set('type', type)
    if (country) params.set('country', country)
    params.set('page', page)
    params.set('limit', '20')

    const timeout = setTimeout(() => {
      fetch(`${apiBase}/search?${params}`)
        .then(r => r.json())
        .then(data => {
          setResults(data.results || [])
          setTotal(data.total || 0)
          setLoading(false)
        })
        .catch(() => {
          setResults([])
          setTotal(0)
          setLoading(false)
        })
    }, query ? 300 : 0)

    return () => clearTimeout(timeout)
  }, [apiBase, query, type, country, page])

  return (
    <div>
      <div className="hero">
        <h1>Find Your Opportunity</h1>
        <p>Scholarships, internships, and competitions across the Arab world</p>
      </div>

      <div className="filters">
        <input
          type="text"
          placeholder="Search by title or keyword..."
          value={query}
          onChange={e => { setQuery(e.target.value); setPage(1) }}
        />
        <select value={type} onChange={e => { setType(e.target.value); setPage(1) }}>
          <option value="">All Types</option>
          <option value="scholarship">Scholarships</option>
          <option value="internship">Internships</option>
          <option value="competition">Competitions</option>
        </select>
        <select value={country} onChange={e => { setCountry(e.target.value); setPage(1) }}>
          <option value="">All Countries</option>
          {countries.map(c => <option key={c} value={c}>{c}</option>)}
        </select>
      </div>

      {!loading && total > 0 && (
        <p className="results-info">{total} result{total !== 1 ? 's' : ''} found</p>
      )}

      {loading ? (
        <div className="loading">Loading...</div>
      ) : results.length === 0 ? (
        <div className="empty-state">
          <h2>No opportunities found</h2>
          <p>Try adjusting your filters or search terms.</p>
        </div>
      ) : (
        results.map(opp => (
          <OpportunityCard key={opp.id} opportunity={opp} />
        ))
      )}

      {total > 20 && (
        <div style={{ display: 'flex', justifyContent: 'center', gap: 8, marginTop: 20 }}>
          <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>Previous</button>
          <span style={{ padding: '8px 12px', fontSize: 14, color: '#64748b' }}>Page {page}</span>
          <button disabled={results.length < 20} onClick={() => setPage(p => p + 1)}>Next</button>
        </div>
      )}
    </div>
  )
}
