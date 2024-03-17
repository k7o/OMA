import { SplitPane } from 'solid-split-pane'
import { useData } from './DataContext'
import { MonacoEditor } from '../lib/solid-monaco'
import { For, Show, createEffect, createResource, createSignal } from 'solid-js'
import { editor } from 'monaco-editor'
import type { Monaco } from '@monaco-editor/loader'
import { Lint } from '../types/Lint'
import { c } from 'vite/dist/node/types.d-FdqQ54oU'

export const Editor = () => {
  const [policyInstance, setPolicyInstance] = createSignal<{
    monaco: Monaco
    editor: editor.IStandaloneCodeEditor
  }>()
  const { data, input, policy, setData, setInput, setPolicy, output, coverage, setCoverage } =
    useData()
  const [linting, { refetch: lint }] = createResource<Lint>(async () => {
    try {
      const res = await fetch('http://localhost:8080/api/lint', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ policy: policy() }),
      })

      if (res.ok) {
        try {
          return res.json()
        } catch (error) {
          console.error('Failed to parse lint response', error, res)
        }
      }
    } catch (error) {
      console.error('Failed to fetch lint', error)
    }
  })

  createEffect(() => {
    if (policyInstance()) {
      clearDecorations()

      if (coverage()) {
        let instance = policyInstance()!
        const decorations = coverage()!.covered.map<editor.IModelDeltaDecoration>((covered) => {
          return {
            range: new instance.monaco.Range(covered.start, 1, covered.end, 1),
            options: {
              isWholeLine: true,
              className: 'bg-green-200',
            },
          }
        })

        instance.editor
          .getModel()
          ?.getLinesContent()
          .forEach((line, index) => {
            if (
              line.trim() !== '' &&
              !line.startsWith('package') &&
              !line.startsWith('import') &&
              !line.startsWith('}') &&
              !coverage()!.covered.some((c) => index + 1 >= c.start && index + 1 <= c.end)
            ) {
              decorations.push({
                range: new instance.monaco.Range(index + 1, 1, index + 1, 1),
                options: {
                  isWholeLine: true,
                  className: 'bg-red-200',
                },
              })
            }
          })

        policyInstance()!.editor.createDecorationsCollection(decorations)
      }
    }
  })

  function clearDecorations() {
    if (policyInstance()) {
      const editor = policyInstance()!.editor
      editor.removeDecorations(
        editor
          .getModel()!
          .getAllDecorations()
          .map((d) => d.id),
      )
    }
  }

  return (
    <div class="flex h-screen w-full">
      <SplitPane gutterClass="gutter gutter-horizontal">
        <div>
          <h3 class="bg-gray-400 text-white px-2 relative">POLICY</h3>
          <MonacoEditor
            class="w-full h-full relative"
            language="rego"
            value={policy()}
            onChange={(value) => {
              setCoverage()
              lint()
              setPolicy(value)
            }}
            onMount={(monaco, editor) => {
              setPolicyInstance({ monaco, editor })
            }}
            options={{
              scrollBeyondLastLine: false,
              wordWrap: 'on',
            }}
          />
        </div>
        <div>
          <SplitPane
            direction="vertical"
            gutterClass="gutter gutter-vertical relative"
            sizes={[40, 20, 20, 20]}
          >
            <div>
              <h3 class="bg-gray-400 text-white px-2 relative">INPUT</h3>
              <MonacoEditor
                class="w-full h-full relative"
                language="json"
                value={input()}
                onChange={(value) => {
                  setCoverage()
                  setInput(value)
                }}
                options={{
                  scrollBeyondLastLine: false,
                  wordWrap: 'on',
                }}
              />
            </div>
            <div>
              <h3 class="bg-gray-400 text-white px-2 relative">DATA</h3>
              <MonacoEditor
                class="w-full h-full relative"
                language="json"
                value={data()}
                onChange={setData}
                options={{
                  scrollBeyondLastLine: false,
                  wordWrap: 'on',
                }}
              />
            </div>
            <div>
              <h3 class="bg-gray-400 text-white px-2 relative">OUTPUT</h3>
              <MonacoEditor
                class="w-full h-full relative"
                language="json"
                value={output()}
                options={{
                  scrollBeyondLastLine: false,
                  readOnly: true,
                  minimap: { enabled: false },
                  wordWrap: 'on',
                }}
              />
            </div>
            <div class="block">
              <h3 class="bg-gray-400 text-white px-2 relative">LINT</h3>
              <div class="m-2">
                <Show
                  when={linting()?.errors && linting()!.errors.length > 0}
                  fallback={<span>No linter violations</span>}
                >
                  <p class="px-6 ">{linting()?.message}</p>
                  <ul class="px-6">
                    <For each={linting()?.errors}>
                      {(error) => <li class="list-disc">{error}</li>}
                    </For>
                  </ul>
                </Show>
              </div>
            </div>
          </SplitPane>
        </div>
      </SplitPane>
    </div>
  )
}
