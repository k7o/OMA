import { onMount } from 'solid-js'
import { defaultBundle, useData } from '../components/DataContext'
import { useNavigate } from '@solidjs/router'

export const NewEmptyPlayground = () => {
  const { setBundle } = useData()
  const navigate = useNavigate()

  onMount(() => {
    setBundle(JSON.parse(JSON.stringify(defaultBundle)))
    navigate('/play')
  })

  return (
    <div class="flex flex-col h-full w-full items-center justify-center">
      <h1 class="text-4xl text-gray-600">Creating Empty Playground...</h1>
    </div>
  )
}
