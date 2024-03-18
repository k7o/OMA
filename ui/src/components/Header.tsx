import OpaIcon from '../assets/opa.svg'
import PlayIcon from '../assets/play.png'
import FormatIcon from '../assets/format-icon.png'
import PublishIcon from '../assets/publish-icon.png'
import { Button } from './Button'
import { useData } from './DataContext'
import { EvalResult } from '../types/EvalResult'
import { createSignal } from 'solid-js'

export const Header = () => {
  const { data, input, policy, setPolicy, setOutput, setCoverage, setLocalHistory } = useData()

  async function evaluate() {
    try {
      const res = await fetch('http://localhost:8080/api/eval', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          data: data(),
          input: input(),
          policy: policy(),
        }),
      })

      if (!res.ok) {
        setOutput(await res.text())
      } else {
        res.json().then((res: EvalResult) => {
          if (res.errors) {
            let output = `${res.errors.length} error${res.errors.length > 1 ? 's' : ''} occurred:\n`
            res.errors.forEach((err) => {
              output += `policy.rego:${err.location.row}:${err.code} ${err.message}\n`
            })

            setOutput(output)
          } else {
            setOutput(JSON.stringify(res.result, null, 2))
          }

          if (res.coverage) {
            setCoverage(res.coverage)
          }

          pushHistory(res)
        })
      }
    } catch (e) {
      console.error(e)
    }
  }

  function pushHistory(evalResult: EvalResult) {
    setLocalHistory((history) => [
      {
        decision_id: evalResult.id,
        policy: policy(),
        input: input(),
        data: data(),
        path: "",
        result: JSON.stringify(evalResult, null, 2),
        timestamp: evalResult.timestamp,
      },
      ...history,
    ])
  }

  async function format() {
    try {
      const res = await fetch('http://localhost:8080/api/format', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          policy: policy(),
        }),
      })

      if (res.ok) {
        let json: FormatResponse = { formatted: policy() }
        try {
          json = (await res.json()) as FormatResponse
        } catch (e) {
          console.error(e)
          setOutput(await res.text())
        }

        setPolicy(json.formatted)
      } else {
        setOutput(await res.text())
      }
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <header class="h-14 flex justify-between">
      <div class="items-center flex mx-2 ">
        <img src={OpaIcon} id="opa-logo" alt="OPA logo" class="h-10" />
        <h3 class="text-2xl hidden md:block">The Rego Playground</h3>
      </div>
      <div class="items-center flex">
        <Button text="Evaluate" icon={PlayIcon} onClick={evaluate} />
        <Button text="Format" icon={FormatIcon} onClick={format} />
        <Button text="Publish" icon={PublishIcon} />
      </div>
    </header>
  )
}
