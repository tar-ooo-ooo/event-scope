import { useEffect, useState } from 'react'
import { Button } from 'antd'

import { environment } from './config/environment'
import { createEvent } from './services/eventService'

type SseEvent = {
  event_id: string
  type: string
}

function App() {
  const [events, setEvents] = useState<SseEvent[]>([])

  useEffect(() => {
    const stream = new EventSource(`${environment.apiBaseUrl}/event/stream`)

    stream.addEventListener('event', (message) => {
      setEvents((events) => [
        JSON.parse((message as MessageEvent<string>).data) as SseEvent,
        ...events,
      ])
    })

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
      <ul>
        {events.map((event) => (
          <li key={event.event_id}>
            {event.type}: {event.event_id}
          </li>
        ))}
      </ul>
    </main>
  )
}

export default App
