export const DecisionLogs = () => {
  return (
    <div class='h-full relative'>
      <h1 class='text-2xl p-2'>Decision Logs</h1>
      <ul role='list' class='relative flex flex-col divide-y divide-gray-100'>
        <ListItem
          item={{
            id: '1',
            policy: 'policy',
            input: 'input',
            data: 'data',
            output: 'output',
            timestamp: new Date(),
          }}
        />
        <ListItem
          item={{
            id: '1',
            policy: 'policy',
            input: 'input',
            data: 'data',
            output: 'output',
            timestamp: new Date(),
          }}
        />
      </ul>
    </div>
  )
}

type PolicyRun = {
  id: string
  policy: string
  input: string
  data: string
  output: string
  timestamp: Date
}

import { Show, createSignal } from 'solid-js'
import ChevronRight from '../assets/chevron-right.svg'
import ChevronDown from '../assets/chevron-down.svg'
import { MonacoEditor } from '../lib/solid-monaco'

const ListItem = (props: { item: PolicyRun }) => {
  const [open, setOpen] = createSignal(false)

  return (
    <li class={`py-3 hover:bg-gray-50 flex relative flex-col h-full`}>
      <div class='flex items-center' onClick={() => setOpen(!open())}>
        <Show
          when={!open()}
          fallback={
            <img src={ChevronDown} alt='collapse' class='w-5 h-5 ml-2' />
          }
        >
          <img src={ChevronRight} alt='expand' class='w-5 h-5 ml-2' />
        </Show>
        <button class=' p-2 text-white hover:bg-slate-600 bg-slate-300 rounded mx-4'>
          replay
        </button>
        <span>{props.item.timestamp.toUTCString()}</span>
        <span>{props.item.id}</span>
      </div>

      <Show when={open()}>
        <MonacoEditor
          class={`h-full mt-2 flex`}
          language='json'
          value={JSON.stringify(props.item, null, 2)}
          onMount={(_, editor) => {
            editor.layout({
              width: editor.getScrollWidth(),
              height: Math.min(editor.getContentHeight(), 500),
            })
          }}
          options={{
            scrollBeyondLastLine: false,
            readOnly: true,
          }}
        />
      </Show>
    </li>
  )
}
