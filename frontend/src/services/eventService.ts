import { environment } from '../config/environment'

export function createEvent() {
  return fetch(`${environment.apiBaseUrl}/event`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      event_id: crypto.randomUUID(),
      type: 'manual.start',
    }),
  })
}
