export type EvalResult = {
  id: string
  result: any
  errors: Error[]
  coverage: Coverage
  timestamp: string
}

type Error = {
  message: string
  code: string
  location: {
    file: string
    row: number
    col: number
  }
}

export type Coverage = {
  covered: Covered[]
  covered_lines: number
  coverage: number
}

export type Covered = {
  start: number
  end: number
}
