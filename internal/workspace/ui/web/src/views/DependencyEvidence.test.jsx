import test from 'node:test'
import assert from 'node:assert/strict'
import { parseHTML } from 'linkedom'
import { render } from 'preact'
import { act } from 'preact/test-utils'
import DependencyEvidence from '../components/DependencyEvidence.jsx'

test('dependency evidence retains module versions, uncertainty, and searchable source evidence', async (t) => {
  const { window, document } = parseHTML('<html><body><div id="root"></div></body></html>')
  const old = { window: globalThis.window, document: globalThis.document }
  Object.assign(globalThis, { window, document })
  const root = document.getElementById('root')
  t.after(async () => { await act(() => render(null, root)); Object.assign(globalThis, old) })
  const metrics = {
    dependency_inventory: { dependencies: [
      { ecosystem: 'npm', name: 'express', module: 'apps/new', version: '5.1.0', resolution: 'lockfile', sources: ['pnpm-lock.yaml'] },
      { ecosystem: 'npm', name: 'express', module: 'apps/legacy', declared: '^4', resolution: 'constraint', sources: ['apps/legacy/package.json'] },
      { ecosystem: 'maven', name: 'example:managed', module: '.', version: '2.0.0', resolution: 'managed_exact', scope: 'management', sources: ['pom.xml'] },
    ], limitations: ['Inherited version unavailable'] },
    detector_coverage: [{ detector_id: 'javascript.http.express', module: 'apps/legacy', status: 'unknown_version', rule_id: 'source-fallback', revision: 'test', reason: 'Version not established' }],
  }
  await act(() => render(<DependencyEvidence metrics={metrics} />, root))
  assert.match(root.textContent, /apps\/new · 5.1.0 · Lockfile/)
  assert.match(root.textContent, /apps\/legacy · \^4 · Constraint/)
  assert.match(root.textContent, /0 validated \/ 1 rules/)
  assert.match(root.textContent, /Version unknown/)
  assert.match(root.textContent, /pnpm-lock.yaml/)
  assert.doesNotMatch(root.textContent, /example:managed/)
  await act(() => { const checkbox = root.querySelector('input[type=checkbox]'); checkbox.checked = true; checkbox.dispatchEvent(new window.Event('change', { bubbles: true })) })
  assert.match(root.textContent, /example:managed/)
  await act(() => { const input = root.querySelector('input'); input.value = 'legacy'; input.dispatchEvent(new window.Event('input', { bubbles: true })) })
  assert.doesNotMatch(root.textContent, /apps\/new/)
  assert.match(root.textContent, /apps\/legacy/)
})
