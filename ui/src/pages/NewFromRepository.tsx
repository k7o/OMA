import { useParams } from '@solidjs/router'
import { For, createResource } from 'solid-js'

export const NewFromRepository = () => {
  const params = useParams()
  async function fetchRevisions() {
    const res = await fetch(`${window.location.origin}/api/revisions/${params.package_id}`)
    return (await res.json()) as string[]
  }

  const [files, actions] = createResource(fetchRevisions)

  return (
    <div>
      <h1>Package ID: {params.package_id}</h1>
      <h2>Files</h2>
      <For each={files()} fallback={<li>Loading...</li>}>
        {(file) => <li>{file}</li>}
      </For>
    </div>
  )
}
