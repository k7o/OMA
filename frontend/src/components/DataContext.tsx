import { createSignal, createContext, useContext, JSX } from 'solid-js'

function createInitialState() {
  const [policy, setPolicy] = createSignal(defaultPolicy)
  const [input, setInput] = createSignal(defaultInput)
  const [data, setData] = createSignal('')
  const [output, setOutput] = createSignal('')
  const [coverage, setCoverage] = createSignal<Coverage | undefined>()

  return {
    policy,
    setPolicy,
    input,
    setInput,
    data,
    setData,
    output,
    setOutput,
    coverage,
    setCoverage,
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
