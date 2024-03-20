import { For, Show, createResource, createSignal } from 'solid-js'
import { ListItem } from '../components/ListItem'
import { DecisionLog } from '../types/DecisionLog'
import { DataProvider, useData } from '../components/DataContext'
import { Button } from '../components/Button'

import GearIcon from '../assets/gear-icon.svg'
import RefreshIcon from '../assets/refresh-icon.svg'

async function fetchDecisionLogs() {
  const res = await fetch('http://localhost:8080/api/decision-log/list')
  return (await res.json()) as DecisionLog[]
}

export const DecisionLogs = () => {
  const { applicationSettings, setApplicationSettings } = useData()
  const [decisionLogs, actions] = createResource<DecisionLog[]>(fetchDecisionLogs)
  const [showSettings, setShowSettings] = createSignal(false)

  return (
    <DataProvider>
      <div class="h-full relative">
        <header class="h-14 flex justify-between items-center">
          <h1 class="text-2xl p-2">Decision Logs</h1>
          <div class="flex">
            <Button
              text="Settings"
              icon={GearIcon}
              onClick={() => setShowSettings(!showSettings())}
            />
            <Button text="Refresh" icon={RefreshIcon} onClick={actions.refetch} />
          </div>
        </header>
        <Show when={showSettings()}>
          <div class="grid grid-cols-2 gap-4 p-4 border-y-4">
            <label>OPA Server URL </label>
            <input
              class="mr-2 border-2"
              onChange={(e) => setApplicationSettings('opa_server_url', e.currentTarget.value)}
              value={applicationSettings.opa_server_url}
            />
          </div>
        </Show>
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
