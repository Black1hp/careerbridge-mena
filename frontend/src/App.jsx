import { BrowserRouter, Routes, Route } from 'react-router-dom'
import SearchPage from './pages/SearchPage'
import DetailPage from './pages/DetailPage'
import './App.css'

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api/v1'

function App() {
  return (
    <BrowserRouter>
      <div className="app">
        <header className="header">
          <a href="/" className="logo">CareerBridge <span className="logo-sub">MENA</span></a>
          <nav>
            <a href="/">Search</a>
          </nav>
        </header>
        <main className="main">
          <Routes>
            <Route path="/" element={<SearchPage apiBase={API_BASE} />} />
            <Route path="/opportunity/:id" element={<DetailPage apiBase={API_BASE} />} />
          </Routes>
        </main>
        <footer className="footer">
          <p>CareerBridge MENA &mdash; Scholarships, internships & competitions across the Arab world</p>
        </footer>
      </div>
    </BrowserRouter>
  )
}

export default App
