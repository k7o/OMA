import { ListItem } from "../components/ListItem"

export const DecisionLogs = () => {
  return (
    <div class='h-full relative'>
      <h1 class='text-2xl p-2'>Decision Logs</h1>
      <ul role='list' class='relative flex flex-col divide-y divide-gray-100'>
        <ListItem
          item={{
            id: '1',
            policy: 'policy',
            input: 'input',
            data: 'data',
            output: 'output',
            timestamp: new Date(),
          }}
        />
        <ListItem
          item={{
            id: '1',
            policy: 'policy',
            input: 'input',
            data: 'data',
            output: 'output',
            timestamp: new Date(),
          }}
        />
      </ul>
    </div>
  )
}

