import { Accessor, For, Match, Setter, Show, Switch, createSignal } from 'solid-js'
import ChevronRight from '../assets/chevron-right.svg'
import ChevronDown from '../assets/chevron-down.svg'
import { MonacoEditor } from '../lib/solid-monaco'

import ReplayIcon from '../assets/replay-icon.svg'
import { useData } from './DataContext'

export const ListItem = (props: { item: PolicyRun }) => {
  const { setPolicy, setInput } = useData()
  const [open, setOpen] = createSignal(false)
  const [tab, setTab] = createSignal<Tabs>('Input')

  return (
    <li class={`py-3 hover:bg-gray-50 flex relative flex-col h-full`}>
      <div class="flex items-center" onClick={() => setOpen(!open())}>
        <Show
          when={!open()}
          fallback={<img src={ChevronDown} alt="collapse" class="w-5 h-5 ml-2" />}
        >
          <img src={ChevronRight} alt="expand" class="w-5 h-5 ml-2" />
        </Show>
        <button
          onClick={(e) => {
            e.stopPropagation()
          }}
          class="px-2 py-1 text-white hover:bg-slate-600 bg-slate-300 rounded mx-4"
        >
          <img src={ReplayIcon} alt="replay" class="w-10 h-10" />
        </button>
        <span class="text-sm">{props.item.timestamp.toUTCString()}</span>
        <span class="text-sm">{props.item.id}</span>
      </div>

      <Show when={open()}>
        <TabBar tab={tab} setTab={setTab} />
        <Switch
          fallback={
            <Match when={tab() === 'Input'}>
              <SmallEditor value={props.item.input} language="json" />
            </Match>
          }
        >
          <Match when={tab() === 'Input'}>
            <SmallEditor value={props.item.input} language="json" />
          </Match>
          <Match when={tab() === 'Output'}>
            <SmallEditor value={props.item.output} language="json" />
          </Match>
          <Match when={tab() === 'Policy'}>
            <SmallEditor value={props.item.policy} language="rego" />
          </Match>
        </Switch>
      </Show>
    </li>
  )
}

const Tabs = ['Input', 'Policy', 'Output'] as const
type Tabs = (typeof Tabs)[number]

type TabBarProps = {
  tab: Accessor<Tabs>
  setTab: Setter<Tabs>
}

const TabBar = (props: TabBarProps) => {
  return (
    <div class="flex mt-2 w-full bg-gray-100">
      <For each={Tabs}>
        {(tab, index) => (
          <button
            class={`${index() !== 0 && 'ml-2'} px-4 py- rounded-xl ${
              props.tab() === tab ? 'bg-gray-200' : 'bg-gray-100'
            }`}
            onClick={() => props.setTab(tab)}
          >
            {tab}
          </button>
        )}
      </For>
    </div>
  )
}

const SmallEditor = (props: { value: string; language: string }) => {
  return (
    <MonacoEditor
      class={`h-full mt-2 flex`}
      language={props.language}
      value={props.value}
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
  )
}
