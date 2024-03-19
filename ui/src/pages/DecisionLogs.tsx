import { For, Show, createEffect, createResource } from 'solid-js'
import { ListItem } from '../components/ListItem'
import { DecisionLog } from '../types/DecisionLog'
import { DataProvider } from '../components/DataContext'
import { Button } from '../components/Button'

import RefreshIcon from '../assets/refresh-icon.svg'

async function fetchDecisionLogs() {
  const res = await fetch('http://localhost:8080/api/decision-log/list')
  return await res.json() as DecisionLog[]
}

export const DecisionLogs = () => {
  const [decisionLogs, actions] = createResource<DecisionLog[]>(fetchDecisionLogs)

  return (
    <DataProvider>
      <div class="h-full relative">
        <header class="h-14 flex justify-between items-center">
          <h1 class="text-2xl p-2">Decision Logs</h1>
          <Button text='Refresh' icon={RefreshIcon} onClick={actions.refetch}/>
        </header>
        <ul role="list" class="relative flex flex-col divide-y divide-gray-100">
          <Show when={decisionLogs.loading}>
            <li class="px-2">Loading...</li>
          </Show>
          <For each={decisionLogs()} fallback={<li class="px-2">No decision logs yet</li>}>
            {(log) => <ListItem item={log} />}
          </For>
        </ul>
      </div>
    </DataProvider>
  )
}
