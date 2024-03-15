import { SplitPane } from 'solid-split-pane'
import { useData } from './DataContext'
import { MonacoEditor } from '../lib/solid-monaco'

export const Editor = () => {
  const { data, input, policy, setData, setInput, setPolicy, output } =
    useData()

  return (
    <div class='flex h-screen'>
      <SplitPane gutterClass='gutter gutter-horizontal'>
        <div>
          <h3 class='bg-gray-400 text-white px-2 relative'>POLICY</h3>
          <MonacoEditor
            class='w-full h-full relative'
            language='rego'
            value={policy()}
            onChange={setPolicy}
            options={{
              scrollBeyondLastLine: false,
            }}
          />
        </div>
        <div>
          <SplitPane
            direction='vertical'
            gutterClass='gutter gutter-vertical relative'
          >
            <div>
              <h3 class='bg-gray-400 text-white px-2 relative'>INPUT</h3>
              <MonacoEditor
                class='w-full h-full relative'
                language='rego'
                value={input()}
                onChange={setInput}
                options={{
                  scrollBeyondLastLine: false,
                }}
              />
            </div>
            <div>
              <h3 class='bg-gray-400 text-white px-2 relative'>DATA</h3>
              <MonacoEditor
                class='w-full h-full relative'
                language='rego'
                value={data()}
                onChange={setData}
                options={{
                  scrollBeyondLastLine: false,
                }}
              />
            </div>
            <div>
              <h3 class='bg-gray-400 text-white px-2 relative'>OUTPUT</h3>
              <MonacoEditor
                class='w-full h-full relative'
                language='rego'
                value={output()}
                options={{
                  scrollBeyondLastLine: false,
                  readOnly: true,
                }}
              />
            </div>
          </SplitPane>
        </div>
      </SplitPane>
    </div>
  )
}
