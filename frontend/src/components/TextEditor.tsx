import { createCodeMirror, createEditorControlledValue } from 'solid-codemirror'
import { lineNumbers } from '@codemirror/view'
import { Accessor, Setter } from 'solid-js'

type TextEditorProps = {
  title: string
  value?: Accessor<string>
  onValueChange?: Setter<string>
}

export const TextEditor = ({
  title,
  value,
  onValueChange,
}: TextEditorProps) => {
  const { ref, createExtension, editorView } = createCodeMirror({
    onValueChange: onValueChange,
  })
  if (value !== undefined) {
    createEditorControlledValue(editorView, value)
  }
  createExtension(lineNumbers())

  return (
    <div class='w-full h-full'>
      <h3 class='bg-gray-400 text-white px-2'>{title}</h3>
      <div ref={ref} />
    </div>
  )
}
