import test from 'node:test'
import assert from 'node:assert/strict'
import { parseHTML } from 'linkedom'
import { render } from 'preact'
import { act } from 'preact/test-utils'
import { Projects } from './Projects.jsx'
import { ImportOrgModal } from './ProjectWorkspace.jsx'
import { ContractComparison } from './GraphCompare.jsx'
import { EvidenceList } from './GraphDetails.jsx'

const settle = () => act(async () => { await new Promise(setImmediate) })
async function dom(t) {
  const { window, document } = parseHTML('<html><body><div id="root"></div><button id="outside">Outside</button></body></html>')
  const old = { window: globalThis.window, document: globalThis.document, fetch: globalThis.fetch }
  Object.assign(globalThis, { window, document })
  const root = document.getElementById('root')
  t.after(async () => { await act(() => render(null, root)); Object.assign(globalThis, old) })
  const show = async (element) => { await act(() => render(element, root)); await settle() }
  const button = (name) => [...root.querySelectorAll('button')].find((node) => node.textContent === name)
  const click = async (name) => { const node = button(name); assert.ok(node, name); assert.equal(node.disabled, false, name); await act(async () => node.dispatchEvent(new window.Event('click', { bubbles: true }))); await settle() }
  const input = async (label, value) => {
    const node = [...root.querySelectorAll('label')].find((node) => node.textContent.trim() === label)?.querySelector('input')
    assert.ok(node, label)
    await act(() => { node.value = value; node.dispatchEvent(new window.Event('input', { bubbles: true })) })
  }
  return { root, document, show, button, click, input }
}

test('first use stays dismissible; project load failure is not an empty workspace', async (t) => {
  const { root, show, click } = await dom(t)
  let failed = true
  globalThis.fetch = async (url) => ({ ok: !failed, status: failed ? 503 : 200, text: async () => JSON.stringify(failed ? { error: 'Workspace unavailable' } : url.endsWith('/session') ? { role: 'admin', mode: 'legacy' } : { projects: [] }) })
  await show(<Projects />)
  assert.match(root.querySelector('[role="alert"]').textContent, /unavailable/)
  assert.doesNotMatch(root.textContent, /No projects yet/)
  assert.equal(root.querySelector('[role="dialog"]'), null)
  failed = false
  await click('Retry loading projects')
  assert.match(root.textContent, /No projects yet/)
  await click('+ New Project')
  assert.ok(root.querySelector('[role="dialog"][aria-modal="true"]'))
  await click('Cancel')
  assert.equal(root.querySelector('[role="dialog"]'), null)
})

test('import retains reviewed candidates and draft, invalidates changed scope and shows local errors', async (t) => {
  const { root, show, button, click, input } = await dom(t)
  const calls = []
  let fail = false
  const onImport = async (body) => { calls.push(body); if (fail) throw new Error('Invalid include regex'); return { preview_digest: 'reviewed-scope', count: 1, results: [{ name: 'checkout', path: '/company/checkout', status: 'preview' }] } }
  const view = (open = true) => <ImportOrgModal open={open} onImport={onImport} onClose={() => {}} />
  await show(view())
  await click('Local directory'); await input('Root directory', '/company')
  assert.equal(button('Import and build graph').disabled, true)
  await click('Preview repositories')
  assert.equal(calls[0].dry_run, true)
  assert.match(root.querySelector('[aria-label="Import preview"]').textContent, /checkout/)
  assert.equal(button('Import and build graph').disabled, false)
  await show(view(false)); assert.equal(root.querySelector('[role="dialog"]'), null)
  await show(view()); assert.equal(button('Import and build graph').disabled, false)
  await input('Include regex', 'checkout')
  assert.equal(button('Import and build graph').disabled, true)
  await click('Preview repositories'); await click('Import and build graph')
  assert.equal(calls.at(-1).dry_run, false)
  assert.equal(calls.at(-1).root, '/company')
  assert.equal(calls.at(-1).preview_digest, 'reviewed-scope')
  fail = true
  await click('Preview repositories')
  assert.match(root.querySelector('[role="dialog"] [role="alert"]').textContent, /Invalid include regex/)
  assert.equal(button('Import and build graph').disabled, true)
})

test('all extracted evidence remains reachable past the initial display limit', async (t) => {
  const { root, show, click } = await dom(t)
  await show(<EvidenceList title="Dependencies" items={Array.from({ length: 163 }, (_, i) => i)} renderItem={(item) => <p data-fact={item}>{item}</p>} />)
  assert.equal(root.querySelectorAll('[data-fact]').length, 80)
  await click('Show more dependencies (80 of 163)')
  assert.equal(root.querySelectorAll('[data-fact]').length, 160)
  await click('Show more dependencies (160 of 163)')
  assert.equal(root.querySelectorAll('[data-fact]').length, 163)
  assert.equal(root.querySelector('button'), null)
})

test('contract comparison pins direction, exposes evidence, and never claims complete compatibility', async (t) => {
  const { root, show, click } = await dom(t)
  let fail = true
  const paths = []
  globalThis.fetch = async (url) => { paths.push(url); return { ok: !fail, status: fail ? 503 : 200, text: async () => JSON.stringify(fail ? { error: 'Evidence unavailable' } : { changes: [{ key: 'k', compatibility: 'potentially_breaking', change: 'modified', before: { service: 'checkout', name: 'id', type: 'string' }, after: { service: 'checkout', name: 'id', type: 'integer', source_file: 'openapi.yaml', source_line: 12 } }] }) } }
  await show(<ContractComparison pid="company a" from="before" to="after" />)
  assert.match(root.querySelector('[role="alert"]').textContent, /unavailable/)
  fail = false; await click('Retry contracts')
  assert.equal(paths[0], '/api/v1/projects/company%20a/contracts/compare?from=before&to=after')
  assert.match(root.textContent, /potentially breaking/)
  assert.match(root.textContent, /source_line/)
  assert.match(root.textContent, /not proof of compatibility/)
})
