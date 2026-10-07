import test from 'node:test'
import assert from 'node:assert/strict'
import { parseHTML } from 'linkedom'
import { render } from 'preact'
import { act } from 'preact/test-utils'
import { FlowReview, FlowReviewDetail, pairedLayout } from './FlowReview.jsx'

const settle = () => act(async () => { await new Promise(setImmediate) })
async function dom(t) {
  const { window, document } = parseHTML('<html><body><div id="root"></div></body></html>')
  const old = { window: globalThis.window, document: globalThis.document, fetch: globalThis.fetch }
  Object.assign(globalThis, { window, document })
  const root = document.getElementById('root')
  t.after(async () => { await act(() => render(null, root)); Object.assign(globalThis, old) })
  return { root, window, show: async (el) => { await act(() => render(el, root)); await settle() } }
}
const node = (id, label, kind = 'db_operation') => ({ id, label, kind, service: 'orders', details: { table: 'audit_log', source_locations: [{ file: 'OrdersService.java', start_line: 40 }] } })
const graph = (nodes) => ({ entry: { service: 'orders', object_id: 'entry' }, nodes, edges: nodes.slice(1).map((n) => ({ from: nodes[0].id, to: n.id, kind: 'write' })), status: 'complete' })
function fixture() {
  const entry = node('entry', 'PUT /orders/{id}', 'http_endpoint'), write = node('write', 'Update orders'), audit = node('audit', 'Write audit_log')
  return { key: 'orders-put', name: entry.label, kind: entry.kind, service: 'orders', change: 'modified', partial: false, before: { entry: { name: entry.label }, graph: graph([entry, write]), callers: [{ from: 'storefront' }] }, after: { entry: { name: entry.label }, graph: graph([entry, write, audit]), callers: [{ from: 'storefront' }] }, changes: [{ kind: 'node', key: 'node:audit', label: audit.label, change: 'added', after: audit }] }
}
const response = (data, ok = true) => ({ ok, status: ok ? 200 : 503, text: async () => JSON.stringify(data) })

test('paired diagrams retain positions, highlight changes and expose keyboard evidence', async (t) => {
  const { root, window, show } = await dom(t)
  const flow = fixture(), layout = pairedLayout(flow.before.graph, flow.after.graph)
  assert.equal(layout.positions.size, 3)
  await show(<FlowReviewDetail flow={flow} from="merge-base-run" to="head-run" />)
  assert.equal(root.querySelectorAll('.fr-diagram-panel').length, 2)
  assert.match(root.textContent, /Added Write audit_log/)
  assert.match(root.textContent, /Nearby callersstorefront/)
  const added = root.querySelector('.fr-node.added')
  assert.ok(added); assert.equal(added.getAttribute('tabindex'), '0')
  const beforeWrite = root.querySelector('[aria-label="Inspect Before Update orders"]'), afterWrite = root.querySelector('[aria-label="Inspect After Update orders"]')
  assert.equal(beforeWrite.getAttribute('transform'), afterWrite.getAttribute('transform'))
  await act(() => { const event = new window.Event('keydown', { bubbles: true }); event.key = 'Enter'; added.dispatchEvent(event) })
  assert.match(root.querySelector('.fr-node-evidence').textContent, /OrdersService.java/)
})

test('removed entrypoint keeps before evidence; self browsing has one diagram', async (t) => {
  const { root, show } = await dom(t)
  const flow = { ...fixture(), after: null, change: 'removed', partial: true }
  await show(<FlowReviewDetail flow={flow} from="before" to="after" />)
  assert.match(root.textContent, /entry point was removed/); assert.match(root.textContent, /Partial evidence/)
  assert.ok(root.querySelector('[aria-label="Before flow"] svg'))
  assert.equal(root.querySelector('[aria-label="After flow"] svg'), null)
  await show(<FlowReviewDetail key="single" flow={fixture()} single from="same" to="same" />)
  assert.equal(root.querySelectorAll('.fr-diagram-panel').length, 1)
})

test('PR missing snapshots and provider failure never imply a completed comparison', async (t) => {
  const { root, window, show } = await dom(t)
  let failure = false
  globalThis.fetch = async () => failure ? response({ error: 'Provider unavailable' }, false) : response({ available: false, base_commit: 'base1234', head_commit: 'head1234', after_run: 'head', next_action: 'Analyze the clean merge-base in a separate checkout.' })
  await show(<FlowReview pid="company" pr={{ repo_id: 'orders', number: 3 }} />)
  assert.match(root.textContent, /Two saved snapshots are needed/); assert.match(root.textContent, /snapshot missing/)
  assert.equal(root.querySelector('svg'), null)
  failure = true
  await act(() => root.querySelector('button').dispatchEvent(new window.Event('click', { bubbles: true }))); await settle()
  assert.match(root.querySelector('[role="alert"]').textContent, /Provider unavailable/)
  assert.doesNotMatch(root.textContent, /PR revisions matched/)
})

test('new scopes ignore late responses and pagination uses exact saved IDs', async (t) => {
  const { root, window, show } = await dom(t)
  let finishOld
  const paths = []
  globalThis.fetch = (path) => {
    paths.push(path)
    if (path.includes('old-company')) return new Promise((resolve) => { finishOld = resolve })
    return Promise.resolve(response({ service: 'orders', from: { id: 'base' }, to: { id: 'head' }, total: 11, next_offset: path.includes('offset=10') ? null : 10, flows: [fixture()], notes: [] }))
  }
  await show(<FlowReview pid="old-company" from="base" to="head" service="orders" />)
  await show(<FlowReview pid="new-company" from="base" to="head" service="orders" />)
  finishOld(response({ service: 'private-old-service', flows: [], total: 0 })); await settle()
  assert.doesNotMatch(root.textContent, /private-old-service/)
  const next = [...root.querySelectorAll('button')].find((button) => button.textContent === 'Next')
  await act(() => next.dispatchEvent(new window.Event('click', { bubbles: true }))); await settle()
  assert.match(paths.at(-1), /from=base&to=head&service=orders&offset=10/)
})


test('repeated semantic steps render once and retain their occurrence evidence', async (t) => {
  const { root, window, show } = await dom(t)
  const flow = fixture()
  flow.after.graph.nodes.push({ ...flow.after.graph.nodes[2], details: { table: 'audit_log', operation: 'INSERT' } })
  await show(<FlowReviewDetail flow={flow} single />)
  assert.equal(root.querySelectorAll('[aria-label="Inspect Flow Write audit_log"]').length, 1)
  await act(() => { const event = new window.Event('keydown', { bubbles: true }); event.key = 'Enter'; root.querySelector('[aria-label="Inspect Flow Write audit_log"]').dispatchEvent(event) })
  assert.match(root.querySelector('.fr-node-evidence').textContent, /occurrences/)
  assert.match(root.querySelector('.fr-node-evidence').textContent, /OrdersService.java/)
})

test('prepare PR flows shows progress and refreshes impact only after exact snapshots complete', async (t) => {
 const { root, window, show } = await dom(t)
 let complete, prepared = 0
 globalThis.fetch = async (path, opts) => opts?.method === 'POST' ? new Promise((resolve) => { complete = resolve }) : response({ available: false, base_commit: 'base', head_commit: 'head', next_action: 'Prepare isolated snapshots.' })
 await show(<FlowReview pid="company" pr={{ repo_id: 'orders', number: 7 }} onPrepared={() => prepared++} />)
 const prepare = [...root.querySelectorAll('button')].find((b) => b.textContent === 'Prepare PR flows')
 await act(() => prepare.dispatchEvent(new window.Event('click', { bubbles: true }))); await settle()
 assert.match(root.querySelector('[role="status"]').textContent, /Preparing PR flows/)
 assert.equal(prepared, 0)
 complete(response({ available: true, before_run: 'base', after_run: 'head', review: { service: 'orders', from: { id: 'base' }, to: { id: 'head' }, total: 1, flows: [fixture()] } })); await settle()
 assert.equal(prepared, 1); assert.equal(root.querySelectorAll('.fr-diagram-panel').length, 2)
 assert.doesNotMatch(root.textContent, /snapshots are needed/)
})

test('empty PR diff shows no diagrams and offers separate service exploration', async (t) => {
 const { root, show } = await dom(t)
 globalThis.fetch = async () => response({ available: true, explore_run: 'company-run', review: { service: 'orders', from: { id: 'base' }, to: { id: 'head' }, total: 0, flows: [] } })
 await show(<FlowReview pid="company" pr={{ repo_id: 'orders', number: 7 }} />)
 assert.match(root.textContent, /0 changed entry points/)
 assert.match(root.textContent, /No extracted flow changes/)
 assert.match(root.textContent, /Explore all service flows/)
 assert.equal(root.querySelectorAll('.fr-diagram-panel').length, 0)
 assert.equal(root.querySelectorAll('.fr-entry').length, 0)
})
