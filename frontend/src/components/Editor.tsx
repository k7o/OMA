import { SplitPane } from 'solid-split-pane'
import { TextEditor } from './TextEditor'
import { useData } from './DataContext'

export const Editor = () => {
  const { data, input, policy, setData, setInput, setPolicy, output } =
    useData()

  return (
    <div class='flex h-screen'>
      <SplitPane gutterClass='gutter gutter-horizontal'>
        <TextEditor title='POLICY' value={policy} onValueChange={setPolicy} />
        <div>
          <SplitPane direction='vertical' gutterClass='gutter gutter-vertical'>
            <TextEditor title='INPUT' value={input} onValueChange={setInput} />
            <TextEditor title='DATA' value={data} onValueChange={setData} />
            <TextEditor title='OUTPUT' value={output} />
          </SplitPane>
        </div>
      </SplitPane>
    </div>
  )
}
