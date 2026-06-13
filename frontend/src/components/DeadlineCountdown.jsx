import { useState, useEffect } from 'react'

export default function DeadlineCountdown({ deadline }) {
  const [timeLeft, setTimeLeft] = useState('')

  useEffect(() => {
    function calc() {
      const now = new Date()
      const dl = new Date(deadline)
      const diff = dl - now

      if (diff < 0) {
        setTimeLeft('passed')
        return
      }

      const days = Math.floor(diff / (1000 * 60 * 60 * 24))
      const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))

      if (days > 30) {
        setTimeLeft(`${Math.floor(days / 30)}mo left`)
      } else if (days > 0) {
        setTimeLeft(`${days}d ${hours}h left`)
      } else {
        const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
        setTimeLeft(`${hours}h ${mins}m left`)
      }
    }

    calc()
    const interval = setInterval(calc, 60000)
    return () => clearInterval(interval)
  }, [deadline])

  if (!deadline) return null

  const dl = new Date(deadline)
  const diff = dl - new Date()
  const days = diff / (1000 * 60 * 60 * 24)

  let cls = 'deadline-ok'
  if (diff < 0) cls = 'deadline-passed'
  else if (days < 7) cls = 'deadline-urgent'
  else if (days < 30) cls = 'deadline-soon'

  const label = diff < 0 ? 'Passed' : timeLeft

  return (
    <span className={`deadline-badge ${cls}`}>
      {label}
    </span>
  )
}
