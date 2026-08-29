import { useEffect, useState } from 'react'
import { Button } from 'antd'

import { environment } from './config/environment'
import { createEvent } from './services/eventService'

type SseEvent = {
  event_id: string
}

function App() {
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

  return (
    <main>
      <section className="terminal">
        <div className="terminal-body">
          <section className="panel">
            <div className="command-line">
              <span>$</span>
              <Button
                className="terminal-start"
                onClick={() =>
                  void Promise.all(
                    Array.from({ length: 20 }, () => createEvent()),
                  )
                }
              >
                Start
              </Button>
            </div>
          </section>
          <section className="panel">
            <div className="event-log" aria-label="Events" aria-live="polite">
              {events.length === 0 && <p>waiting for events...</p>}
              {events.map((event) => (
                <div className="event-card" key={event.event_id}>
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
