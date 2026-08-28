import { Button } from 'antd'

import { createEvent } from './services/eventService'

function App() {
  return (
    <main>
      <Button
        type="primary"
        style={{ position: 'fixed', left: 24, top: 24 }}
        onClick={createEvent}
      >
        Start
      </Button>
    </main>
  )
}

export default App
