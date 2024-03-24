import { Accessor, For, JSX, Match, Setter, Show, Switch, createSignal } from 'solid-js'
import { MonacoDiffEditor, MonacoEditor } from '../lib/solid-monaco'

import ReplayIcon from '../assets/replay-icon.svg'
import { useData } from './DataContext'
import { DecisionLog } from '../types/DecisionLog'
import { Bundle } from '../types/Bundle'

import ChevronRight from '../assets/chevron-right.svg'
import ChevronDown from '../assets/chevron-down.svg'

type ListItemProps = DecisionLog & {
  bundle?: Bundle
}

const Tabs = ['Input', 'Bundle', 'Result'] as const
type Tabs = (typeof Tabs)[number]

export const ListItem = (props: { item: ListItemProps; previousItem?: ListItemProps }) => {
  const [open, setOpen] = createSignal(false)
  const [tab, setTab] = createSignal<Tabs>('Input')
  const { setInput, setNewBundle } = useData()

  return (
    <li
      class={`py-3 hover:bg-gray-50 flex relative overflow-y-scroll flex-col ${open() && 'h-full'}`}
    >
      <div class="flex items-center" onClick={() => setOpen(!open())}>
        <Show
          when={!open()}
          fallback={<img src={ChevronDown} alt="collapse" class="w-5 h-5 ml-2" />}
        >
          <img src={ChevronRight} alt="expand" class="w-5 h-5 ml-2" />
        </Show>
        <Status item={props.item} />
        <Show when={props.item.bundle}>
          <button
            onClick={(e) => {
              setInput(props.item.input)
              setNewBundle(props.item.bundle!)
              e.stopPropagation()
            }}
            class="px-2 py-1 text-white hover:bg-slate-600 bg-slate-300 rounded mx-4"
          >
            <img src={ReplayIcon} alt="replay" class="w-7 h-7" />
          </button>
        </Show>
        <span class="text-sm">{new Date(props.item.timestamp).toUTCString()}</span>
        <span class="text-sm">{props.item.decision_id}</span>
      </div>

      <Show when={open()}>
        <TabBar tab={tab} setTab={setTab} hasPolicy={props.item.bundle !== undefined} />
        <Switch>
          <Match when={tab() === 'Input'}>
            <SmallEditor
              value={JSON.stringify(JSON.parse(props.item.input), null, 2)}
              previousValue={
                props.previousItem
                  ? JSON.stringify(JSON.parse(props.previousItem.input), null, 2)
                  : undefined
              }
              language="json"
            />
          </Match>
          <Match when={tab() === 'Result'}>
            <SmallEditor
              value={JSON.stringify(JSON.parse(props.item.result), null, 2)}
              previousValue={
                props.previousItem
                  ? JSON.stringify(JSON.parse(props.previousItem.result), null, 2)
                  : undefined
              }
              language="json"
            />
          </Match>
          <Show when={() => tab() === 'Bundle' && props.item.bundle !== undefined}>
            <Match when={tab() === 'Bundle'}>
              <BundleBar bundle={props.item.bundle!} previousBundle={props.previousItem?.bundle}>
                {([filename, content]) => (
                  <SmallEditor
                    value={props.item.bundle![filename]}
                    previousValue={props.previousItem?.bundle?.[filename]}
                    language="rego"
                  />
                )}
              </BundleBar>
            </Match>
          </Show>
        </Switch>
      </Show>
    </li>
  )
}

const BundleBar = (props: {
  bundle: Bundle
  previousBundle?: Bundle
  children: (props: [filename: string, content: string]) => JSX.Element
}) => {
  const [bundleFile, setBundleFile] = createSignal(Object.keys(props.bundle)[0])

  if (
    props.previousBundle !== undefined &&
    Object.keys(props.previousBundle).every((key) => key in props.bundle)
  ) {
    setBundleFile(
      Object.keys(props.bundle).find((key) => props.bundle[key] !== props.previousBundle![key]) ||
        Object.keys(props.bundle)[0],
    )
  }

  return (
    <>
      <div class="flex flex-wrap mt-2 w-full bg-gray-100 px-1">
        <For each={Object.keys(props.bundle)}>
          {(tab) => {
            return (
              <button
                class={`px-4 py-1 my-1 rounded ${
                  bundleFile() === tab ? 'bg-gray-200' : 'bg-gray-100'
                }`}
                onClick={() => setBundleFile(tab)}
              >
                {tab}
              </button>
            )
          }}
        </For>
      </div>
      {props.children([bundleFile(), props.bundle[bundleFile()]])}
    </>
  )
}

const StatusSpan = (props: { text: string; class: string }) => {
  return <span class={`p-2 text-sm mx-2 rounded text-white ${props.class}`}>{props.text}</span>
}

const Status = (props: { item: ListItemProps }) => {
  try {
    const result = JSON.parse(props.item.result)
    const allowed = findAllowedValue(result)
    if (allowed === true) {
      return <StatusSpan text="Allowed" class="bg-green-500" />
    } else if (allowed === false) {
      return <StatusSpan text="Failure" class="bg-red-500" />
    } else if (props.item.result === 'null' || result.errors) {
      return <StatusSpan text="Error" class="bg-amber-500" />
    }
  } catch {}

  return
}

const findAllowedValue = (data: any): any => {
  if (typeof data === 'object' && data !== null) {
    if ('allowed' in data) {
      return data.allowed
    } else if ('allow' in data) {
      return data.allow
    } else {
      for (const key in data) {
        const value = findAllowedValue(data[key])
        if (value !== undefined) {
          return value
        }
      }
    }
  }

  return undefined
}

type TabBarProps = {
  hasPolicy: boolean
  tab: Accessor<Tabs>
  setTab: Setter<Tabs>
}

const TabBar = (props: TabBarProps) => {
  return (
    <div class="flex mt-2 w-full bg-gray-100">
      <For each={Tabs}>
        {(tab, index) => {
          if (tab === 'Bundle' && !props.hasPolicy) {
            return null
          }

          return (
            <button
              class={`${index() !== 0 && 'ml-2'} w-full px-4 rounded-xl ${
                props.tab() === tab ? 'bg-gray-200' : 'bg-gray-100'
              }`}
              onClick={() => props.setTab(tab)}
            >
              {tab}
            </button>
          )
        }}
      </For>
    </div>
  )
}

const SmallEditor = (props: { value: string; language: string; previousValue?: string }) => {
  if (props.previousValue === undefined || props.previousValue === props.value) {
    return (
      <MonacoEditor
        class={`h-full mt-2 flex grow`}
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
          showFoldingControls: 'always',
          readOnly: true,
        }}
      />
    )
  }

  return (
    <MonacoDiffEditor
      class={`mt-2 h-full`}
      originalLanguage={props.language}
      modifiedLanguage={props.language}
      modified={props.value}
      original={props.previousValue}
      onMount={(_, editor) => {
        editor.layout({
          width: 100,
          height: 100,
        })
      }}
      options={{
        scrollBeyondLastLine: false,
        showFoldingControls: 'always',
        readOnly: true,
      }}
    />
  )
}
