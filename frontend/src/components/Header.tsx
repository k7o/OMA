import OpaIcon from '../assets/opa.svg'
import PlayIcon from '../assets/play.png'
import FormatIcon from '../assets/format-icon.png'
import PublishIcon from '../assets/publish-icon.png'
import { Button } from './Button'
import { useData } from './DataContext'

type ErrorResult = {
  message: string
  code: string
  location: {
    file: string
    row: number
    col: number
  }
}

export const Header = () => {
  const { data, input, policy, setOutput, setCoverage } = useData()

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
        res.json().then((res) => {
          if (res.errors as ErrorResult[]) {
            let output = `${res.errors.length} error${res.errors.length > 1 ? 's' : ''} occurred:\n`
            res.errors.forEach((err: ErrorResult) => {
              output += `policy.rego:${err.location.row}:${err.code} ${err.message}\n`
            })

            setOutput(output)
          } else {
            setOutput(JSON.stringify(res.result, null, 2))
          }

          if (res.coverage) {
            setCoverage(res.coverage)
          }
        })
      }
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <header class='h-14 flex justify-between'>
      <div class='items-center flex mx-2 '>
        <img src={OpaIcon} id='opa-logo' alt='OPA logo' class='h-10' />
        <h3 class='text-2xl hidden md:block'>The Rego Playground</h3>
      </div>
      <div class='items-center flex'>
        <Button text='Evaluate' icon={PlayIcon} onClick={evaluate} />
        <Button text='Format' icon={FormatIcon} />
        <Button text='Publish' icon={PublishIcon} />
      </div>
    </header>
  )
}
