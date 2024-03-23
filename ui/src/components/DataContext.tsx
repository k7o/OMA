import { createSignal, createContext, useContext, JSX, onMount } from 'solid-js'
import { DecisionLog } from '../types/DecisionLog'
import { createStore } from 'solid-js/store'
import { makePersisted } from '@solid-primitives/storage'

function createInitialState() {
  const [bundle, setBundle] = createStore<Bundle>({
    'policy.rego': defaultPolicy,
  })
  const [editingPolicy, setEditingPolicy] = createSignal<string>(Object.keys(bundle)[0])
  const [input, setInput] = createSignal(defaultInput)
  const [data, setData] = createSignal('')
  const [output, setOutput] = createSignal('')
  const [coverage, setCoverage] = createSignal<Coverage | undefined>()
  const [localHistory, setLocalHistory] = createSignal<DecisionLog[]>([])
  const [options, setOptions] = createSignal<EvalOptions>({
    coverage: false,
  })

  const [applicationSettings, setApplicationSettings] = makePersisted(
    createStore<ApplicationSettings>(
      {
        opa_server_url: 'http://localhost:8181',
        bundle_server_url: 'https://gitlab.com/api/v4/projects/55642500/packages/generic/bundle',
      },
      { name: 'createInitialState' },
    ),
  )

  onMount(() => {
    fetch('http://localhost:8080/api/download', {
      headers: {
        'Content-Type': 'application/json',
      },
      method: 'POST',
      body: JSON.stringify({ application_settings: applicationSettings }),
    })
      .then((res) => res.json())
      .then((data) => setBundle(() => data.files))
  })

  return {
    bundle,
    setBundle,
    editingPolicy,
    setEditingPolicy,
    input,
    setInput,
    data,
    setData,
    output,
    setOutput,
    coverage,
    setCoverage,
    localHistory,
    setLocalHistory,
    options,
    setOptions,
    applicationSettings,
    setApplicationSettings,
  } as const
}

const DataContext = createContext<ReturnType<typeof createInitialState>>()

export const DataProvider = (props: { children: JSX.Element }) => {
  return <DataContext.Provider value={createInitialState()}>{props.children}</DataContext.Provider>
}

export const useData = () => useContext(DataContext)!

export const defaultPolicy = `package example.authz

import rego.v1

default allow := false

allow if {
    input.method == "GET"
    input.path == ["salary", input.subject.user]
}

allow if is_admin

is_admin if "admin" in input.subject.groups
`

export const defaultInput = `{
    "method": "GET",
    "path": ["salary", "bob"],
    "subject": {
        "user": "bob",
        "groups": ["sales", "marketing"]
    }
}`
