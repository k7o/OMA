import { onMount } from 'solid-js'
import { defaultBundle, defaultInput, useData } from '../components/DataContext'
import { useNavigate } from '@solidjs/router'

export const NewEmptyPlayground = () => {
  const { setNewBundle } = useData()
  const navigate = useNavigate()

  onMount(() => {
    setNewBundle(
      JSON.parse(JSON.stringify(defaultBundle)),
      defaultInput,
    )
    navigate('/play')
  })

  return (
    <div class="flex flex-col h-full w-full items-center justify-center">
      <h1 class="text-4xl text-gray-600">Creating Empty Playground...</h1>
    </div>
  )
}
