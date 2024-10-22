import { For, Show, createResource } from 'solid-js'
import { Revision } from '../types/Revision'
import { A, useNavigate } from '@solidjs/router'
import { backend_url } from '../utils/backend_url'
import { BundleResponse } from '~/types/Bundle'
import { useData } from '~/components/DataContext'

async function fetchRevisions() {
  const res = await fetch(`${backend_url}/api/revisions`)
  return (await res.json()) as Revision[]
}

export const NewFromBundle = () => {
  const [revisions, actions] = createResource(fetchRevisions)
  const { setNewBundle } = useData()
  const navigate = useNavigate()

  return (
    <div class="flex flex-col h-full w-full p-6">
      <Show when={revisions.loading}>
        <h1 class="text-4xl text-gray-600">Loading...</h1>
      </Show>
      <Show when={revisions.error}>
        <h1 class="text-4xl text-red-600">Error loading revisions</h1>
      </Show>
      <Show when={!revisions.loading && !revisions.error}>
        <h1 class="text-4xl text-gray-600">Loaded {revisions()!.length} revisions</h1>
      </Show>
      <table class="min-w-full divide-y divide-gray-300">
        <thead>
          <tr>
            <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
              CreatedAt
            </th>
            <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
              Version
            </th>
            <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
              Name
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <For each={revisions()} fallback={<li class="px-2">No revisions yet</li>}>
            {(revision) => {
              const dateString = new Date(revision.created_at).toISOString()
              return (
                <tr
                  class="hover:bg-slate-200 hover:cursor-pointer"
                  onClick={async () => {
                    if (revision.package_type === 'oci') {
                      const res = await fetch(
                        `${backend_url}/api/revisions/package/${revision.package_id}?package_type=${revision.package_type}`,
                      )
                      const bundle = (await res.json()) as BundleResponse
                      setNewBundle(bundle.files)

                      navigate(
                        `/play?revision_id=${JSON.parse(bundle.files['/.manifest']).revision}`,
                      )
                    } else {
                      navigate(`/new/from-bundle/${revision.package_id}`)
                    }
                  }}
                >
                  <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900">
                    {dateString.slice(2, 10) + ' ' + dateString.slice(11, 16)}
                  </td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                    {revision.version}
                  </td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{revision.name}</td>
                </tr>
              )
            }}
          </For>
        </tbody>
      </table>
    </div>
  )
}
