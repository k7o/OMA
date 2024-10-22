import { describe, it, expect } from 'vitest'
import { buildTree } from './FileTree'

describe('buildTree', () => {
  it('should build a tree', () => {
    expect(buildTree({ 'a/b/c.ts': 'package a.b.c' })).toEqual({
      files: [],
      directories: [
        {
          name: 'a',
          files: [],
          directories: [
            {
              name: 'b',
              files: ['c.ts'],
              directories: [],
            },
          ],
        },
      ],
    })
  })

  it('should build a tree with nested directories', () => {
    expect(buildTree({ 'a/b/c.ts': 'package a.b.c', 'a/b/d.ts': 'package a.b.d' })).toEqual({
      files: [],
      directories: [
        {
          name: 'a',
          files: [],
          directories: [{ name: 'b', files: ['c.ts', 'd.ts'], directories: [] }],
        },
      ],
    })
  })
})
