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
      <Button
        type="primary"
        style={{ position: 'fixed', left: 24, top: 24 }}
        onClick={createEvent}
      >
        Start
      </Button>
      <div
        aria-label="Events"
        style={{
          display: 'flex',
          flexDirection: 'column',
          gap: 8,
          left: 120,
          position: 'fixed',
          top: 24,
        }}
      >
        {events.map((event) => (
          <div
            key={event.event_id}
            style={{
              border: '1px solid #1677ff',
              borderRadius: 6,
              fontFamily: 'monospace',
              padding: '8px 12px',
            }}
          >
            {event.event_id}
          </div>
        ))}
      </div>
    </main>
  )
}

export default App
