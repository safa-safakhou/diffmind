import {afterEach, describe, expect, it, vi} from 'vitest'
import {loadRunForm} from './RunForm.jsx'

afterEach(() => vi.unstubAllGlobals())

describe('selected repository run defaults', () => {
  it('uses the selected path for a first run', () => {
    vi.stubGlobal('localStorage', {getItem: () => null, removeItem: () => {}})
    expect(loadRunForm({repo_path: '/selected/repository'}).repo_path).toBe('/selected/repository')
  })
  it('keeps settings while replacing a remembered target with the selected repository', () => {
    vi.stubGlobal('localStorage', {getItem: () => JSON.stringify({repo_path: '/previous/repository', runtime: {workers: 2}, quality: {min_confidence: 0.8}}), removeItem: () => {}})
    const result = loadRunForm({repo_path: '/selected/repository'})
    expect(result.repo_path).toBe('/selected/repository')
    expect(result.runtime.workers).toBe(2)
    expect(result.quality.min_confidence).toBe(0.8)
  })
  it('preserves a remembered target only for an unbound form', () => {
    vi.stubGlobal('localStorage', {getItem: () => JSON.stringify({repo_path: '/previous/repository'}), removeItem: () => {}})
    expect(loadRunForm({}).repo_path).toBe('/previous/repository')
  })
})
