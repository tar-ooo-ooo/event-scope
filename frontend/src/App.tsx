import { type FormEvent, useEffect, useState } from 'react'

import { environment } from './config/environment'
import { createEvent } from './services/eventService'

type SseEvent = {
  event_id: string
  success: boolean
}

function App() {
  const [command, setCommand] = useState('')
  const [commandOutput, setCommandOutput] = useState<string[]>([])
  const [events, setEvents] = useState<SseEvent[]>([])

  useEffect(() => {
    const stream = new EventSource(`${environment.apiBaseUrl}/event/stream`)

    stream.onmessage = (message) => {
      setEvents((events) => [
        ...events,
        JSON.parse(message.data) as SseEvent,
      ])
    }

    return () => stream.close()
  }, [])

  function runCommand(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const value = command.trim()
    if (!value) return

    const [name, countText] = value.split(/\s+/)
    const count = Number(countText)

    if (name === 'start' && !countText) {
      setCommandOutput((output) => [...output, 'start: missing event count'])
    } else if (name === 'start' && Number.isInteger(count) && count > 0) {
      void Promise.all(Array.from({ length: count }, () => createEvent()))
      setCommandOutput((output) => [...output, `started ${count} events`])
    } else {
      setCommandOutput((output) => [...output, `${value}: command not found`])
    }

    setCommand('')
  }

  return (
    <main>
      <section className="terminal">
        <div className="terminal-body">
          <section className="panel">
            <div className="command-output">
              {commandOutput.map((output, index) => (
                <p key={`${output}-${index}`}>{output}</p>
              ))}
            </div>
            <form className="command-line" onSubmit={runCommand}>
              <span>$</span>
              <input
                autoFocus
                className="terminal-input"
                onChange={(event) => setCommand(event.target.value)}
                placeholder="type start <count>"
                value={command}
              />
            </form>
          </section>
          <section className="panel">
            <div className="event-log" aria-label="Events" aria-live="polite">
              {events.length === 0 && <p>waiting for events...</p>}
              {events.map((event) => (
                <div
                  className={
                    event.success === false
                      ? 'event-card event-card--failed'
                      : 'event-card'
                  }
                  key={event.event_id}
                >
                  <code>{event.event_id}</code>
                </div>
              ))}
            </div>
          </section>
          <section className="panel" />
          <section className="panel" />
        </div>
      </section>
    </main>
  )
}

export default App
