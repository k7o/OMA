import { Button } from './Button'
import { useData } from './DataContext'
import { EvalResult } from '../types/EvalResult'

import OpaIcon from '../assets/opa.svg'
import PlayIcon from '../assets/play-circle.svg'
import FormatIcon from '../assets/format-icon.png'
import PublishIcon from '../assets/publish-icon.png'
import { backend_url } from '../utils/backend_url'

export const Header = () => {
  const {
    data,
    input,
    bundle,
    editingPolicy,
    setBundle,
    setOutput,
    setCoverage,
    setLocalHistory,
    options,
    setOptions,
  } = useData()

  async function evaluate() {
    try {
      const res = await fetch(`${backend_url}/api/eval`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          options: options(),
          data: data(),
          input: input(),
          bundle: bundle,
        }),
      })

      if (!res.ok) {
        setOutput(await res.text())
      } else {
        res.json().then((evalResult: EvalResult) => {
          if (evalResult.errors) {
            let output = `${evalResult.errors.length} error${
              evalResult.errors.length > 1 ? 's' : ''
            } occurred:\n`
            evalResult.errors.forEach((err) => {
              output += `${err.location.file}:${err.location.row} ${err.code}\n ${err.message}\n`
            })

            setOutput(output)
          } else {
            setOutput(JSON.stringify(evalResult.result, null, 2))
          }

          if (evalResult.coverage && options().coverage) {
            setCoverage(evalResult.coverage)
          }

          setLocalHistory((history) => [
            {
              decision_id: evalResult.id,
              bundle: JSON.parse(JSON.stringify(bundle)),
              input: input(),
              data: data(),
              path: '',
              result: evalResult.errors
                ? JSON.stringify(evalResult.errors, null, 2)
                : JSON.stringify(evalResult.result, null, 2),
              is_error: evalResult.errors != null && evalResult.errors.length > 0,
              timestamp: evalResult.timestamp,
            },
            ...history,
          ])
        })
      }
    } catch (e) {
      console.error(e)
    }
  }

  async function format() {
    try {
      const res = await fetch(`${backend_url}/api/format`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          policy: bundle,
        }),
      })

      if (res.ok) {
        let json: FormatResponse = { formatted: bundle[editingPolicy()] }
        try {
          json = (await res.json()) as FormatResponse
        } catch (e) {
          console.error(e)
          setOutput(await res.text())
        }

        setBundle(editingPolicy(), json.formatted)
      } else {
        setOutput(await res.text())
      }
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <header class="h-14 w-full min-h-14 flex justify-between">
      <div class="items-center flex mx-2 ">
        <img src={OpaIcon} id="opa-logo" alt="OPA logo" class="h-10" />
        <h3 class="text-2xl hidden md:block">The Rego Playground</h3>
      </div>
      <div class="items-center flex">
        <input
          type="checkbox"
          class="bg-blue-400 h-5 w-5"
          id="coverage"
          name="Coverage"
          onChange={(e) => {
            setOptions({ coverage: e.target.checked })
          }}
        />
        <label class="text-base px-2">Coverage</label>
        <Button text="Evaluate" icon={PlayIcon} onClick={evaluate} />
        <Button text="Format" icon={FormatIcon} onClick={format} />
        <Button text="Publish" icon={PublishIcon} />
      </div>
    </header>
  )
}
