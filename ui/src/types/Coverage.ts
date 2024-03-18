type Coverage = {
  covered: Covered[]
  covered_lines: number
  coverage: number
}

type Covered = {
  start: number
  end: number
}
