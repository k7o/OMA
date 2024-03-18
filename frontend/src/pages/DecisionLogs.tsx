import { For, createResource } from 'solid-js'
import { ListItem } from '../components/ListItem'
import { DecisionLog } from '../types/DecisionLog'

export const DecisionLogs = () => {
  const [decisionLogs] = createResource<DecisionLog[]>(async () => {
    try {
      const res = await fetch('http://localhost:8080/api/decision_logs')

      if (res.ok) {
        return res.json()
      }
    } catch (error) {
      console.error('Failed to fetch decision logs', error)
    }
  })

  return (
    <div class="h-full relative">
      <h1 class="text-2xl p-2">Decision Logs</h1>
      <ul role="list" class="relative flex flex-col divide-y divide-gray-100">
        <For each={decisionLogs()} fallback={<li>Loading...</li>}>
          {(log) => <ListItem item={log} />}
        </For>
      </ul>
    </div>
  )
}
