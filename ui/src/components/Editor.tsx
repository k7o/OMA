import { SplitPane } from 'solid-split-pane'
import { useData } from './DataContext'
import { MonacoEditor } from '../lib/solid-monaco'
import { For, Show, createEffect, createResource, createSignal } from 'solid-js'
import { editor } from 'monaco-editor'
import type { Monaco } from '@monaco-editor/loader'
import { Lint } from '../types/Lint'
import { ListItem } from './ListItem'
import { useSearchParams } from '@solidjs/router'

export const Editor = () => {
  const [policyInstance, setPolicyInstance] = createSignal<{
    monaco: Monaco
    editor: editor.IStandaloneCodeEditor
  }>()
  const {
    data,
    input,
    bundle,
    editingPolicy,
    setData,
    setInput,
    setBundle,
    setEditingPolicy,
    output,
    coverage,
    setCoverage,
    localHistory,
  } = useData()
  const [searchParams] = useSearchParams()

  const [linting, { refetch: lint }] = createResource<Lint>(async () => {
    try {
      const res = await fetch('http://localhost:8080/api/lint', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ policy: bundle[editingPolicy()] }),
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

  console.log(searchParams)

  return (
    <div class="flex h-full w-full">
      <SplitPane gutterClass="gutter gutter-horizontal" sizes={[70, 30]}>
        <div>
          <SplitPane gutterClass="gutter gutter-vertical" direction="vertical" sizes={[60, 40]}>
            <div class="flex">
              <SplitPane
                direction="horizontal"
                gutterClass="gutter gutter-horizontal"
                sizes={[35, 65]}
              >
                <div class="flex flex-col h-full w-full">
                  <h3 class="bg-gray-400 text-white px-2 flex">FILES</h3>
                  <ul class="h-full mt-2">
                    <For each={Object.keys(bundle)} fallback={<li class="px-2 pt-4">No files</li>}>
                      {(file) => (
                        <li
                          class={`px-4 py-1 break-words m-2 rounded hover:bg-slate-300 bg-gray-100 ${
                            editingPolicy() === file && 'bg-gray-300'
                          }`}
                          onClick={() => setEditingPolicy(file)}
                        >
                          {file}
                        </li>
                      )}
                    </For>
                  </ul>
                </div>
                <div class="flex flex-col h-full w-full flex-grow">
                  <h3 class="bg-gray-400 text-white px-2 relative overflow-hidden">POLICY</h3>
                  <MonacoEditor
                    class="w-full h-full relative"
                    language="rego"
                    value={bundle[editingPolicy()]}
                    onChange={(value) => {
                      setCoverage()
                      lint()
                      setBundle(editingPolicy(), value)
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
              </SplitPane>
            </div>
            <div class="flex flex-col relative">
              <h3 class="bg-gray-400 text-white px-2 relative">HISTORY</h3>
              <ul class="h-full">
                <For fallback={<li class="px-2 pt-4">No history</li>} each={localHistory()}>
                  {(item, index) => {
                    if (localHistory().length === 1 || localHistory().length === index() + 1) {
                      return <ListItem item={item} />
                    }

                    return <ListItem item={item} previousItem={localHistory()[index() + 1]} />
                  }}
                </For>
              </ul>
            </div>
          </SplitPane>
        </div>
        {/* <div class="flex flex-col flex-wrap"> */}
        <div>
          <SplitPane
            direction="vertical"
            gutterClass="gutter gutter-vertical relative"
            sizes={[40, 20, 20, 20]}
          >
            <div class="flex flex-col h-full w-full">
              <h3 class="bg-gray-400 text-white px-2">INPUT</h3>
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
                  minimap: { enabled: false },
                  autoIndent: 'full',
                  autoClosingBrackets: 'always',
                  folding: true,
                }}
              />
            </div>
            <div>
              <h3 class="bg-gray-400 text-white px-2 relative">DATA</h3>
              <MonacoEditor
                class="w-full h-full relative overflow-hidden"
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
                class="w-full h-full relative overflow-hidden"
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
            <div class="block overflow-hidden">
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
