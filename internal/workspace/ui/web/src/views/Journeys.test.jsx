import test from 'node:test'
import assert from 'node:assert/strict'
import { parseHTML } from 'linkedom'
import { render } from 'preact'
import { act } from 'preact/test-utils'
import { Projects } from './Projects.jsx'
import { ImportOrgModal, ReadinessNotice, ProjectWorkspace } from './ProjectWorkspace.jsx'
import { ContractComparison } from './GraphCompare.jsx'
import { PullRequestsView, ImpactDetail } from './PullRequestsView.jsx'
import { EvidenceList, EdgeDetail } from './GraphDetails.jsx'
import { useProjectCapabilities } from '../lib/access.js'
import { AgentConnectionHelp } from '../components/AgentConnectionHelp.jsx'
import { ConfirmDialog } from '../components/Modal.jsx'
import { visualEdges } from './GraphCanvas.jsx'

const settle = () => act(async () => { await new Promise(setImmediate) })

test('visual relationship grouping preserves mixed provenance and original facts', () => {
 const edges=[{from:'a',to:'b',type:'http',details:[{name:'static'}],evidence:[{class:'source_extracted'}]},{from:'a',to:'b',type:'http',details:[{name:'declared'}],evidence:[{class:'pack_declared'}]}]
 const grouped=visualEdges(edges)
 assert.equal(grouped.length,1);assert.deepEqual(grouped[0].evidence.map(e=>e.class),['source_extracted','pack_declared'])
 assert.equal(grouped[0].details.length,2);assert.equal(edges[0].evidence.length,1)
})

test('confirmation failure stays inside the dialog and never retries the mutation', async (t) => {
 const {root,show,click}=await dom(t);let calls=0
 await show(<ConfirmDialog title="Remove repository" message="Review this exact scope" confirmLabel="Remove" onConfirm={async()=>{calls++;throw new Error('Conflict. Reload and review the retained scope.')}} onCancel={()=>{}} />)
 await click('Remove')
 assert.equal(calls,1);assert.match(root.querySelector('[role="dialog"] [role="alert"]').textContent,/Reload and review/)
 await settle();assert.equal(calls,1);assert.ok(root.querySelector('[role="dialog"]'))
})

test('viewer joining handoff separates browser identity, project token and renewal ownership', async (t) => {
 const {root,show}=await dom(t)
 await show(<AgentConnectionHelp pid="approved-company" role="viewer" endpoint="https://context.example/mcp" />)
 assert.match(root.textContent,/approved-company/);assert.match(root.textContent,/https:\/\/context.example\/mcp/)
 assert.match(root.textContent,/browser login does not configure/);assert.match(root.textContent,/cannot import or update/)
 assert.match(root.textContent,/administrator owns token expiry, renewal and support/)
 assert.equal(root.querySelector('input,textarea'),null)
})

test('relationship origins show answering run, unknown scope and correction guidance', async (t) => {
 const {root,show}=await dom(t)
 await show(<EdgeDetail e={{from:'a',to:'external',type:'http',evidence:[{class:'pack_declared',run_id:'saved-old',pack_id:'company',resolution:'approved alias',coverage:'unverified',scope_state:'unknown',revision:{commit:'old',dirty:true}}]}} />)
 assert.match(root.textContent,/pack declared/);assert.match(root.textContent,/saved-old/);assert.match(root.textContent,/approved alias/)
 assert.match(root.textContent,/file scope: unknown/);assert.match(root.textContent,/private improvement gap/);assert.match(root.textContent,/do not establish runtime traffic/)
})

test('PR evidence and eligibility precede a collapsed uncalibrated heuristic', async (t) => {
 const {root,show}=await dom(t)
 await show(<ImpactDetail impact={{pull_request:{repo_name:'app',number:1,title:'Remove route',head:'feature',head_sha:'head',base:'main'},codebase:{changed_files:1,files:[],categories:[],risk_reasons:['file category']},company:{available:false,run_id:'saved-default',freshness:'stale',graph_revision:{commit:'default'},next_action:'Analyze PR head separately; updating main cannot guarantee a match.',limitations:['Deleted surfaces remain unproven.']},risk_score:0,risk_level:'low'}} />)
 const text=root.textContent
 assert.ok(text.indexOf('Codebase impact')<text.indexOf('Attention heuristic:'))
 assert.ok(text.indexOf('Company impact')<text.indexOf('Attention heuristic:'))
 assert.match(text,/saved-default/);assert.match(text,/Deleted surfaces remain unproven/);assert.match(text,/not a probability, merge recommendation or proof of safety/)
 assert.equal(root.querySelector('[aria-label="Attention heuristic"] details').hasAttribute('open'),false)
 assert.match(text,/No automatic code-host comments/)
})
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


test('access checks retain prior data on transient failure and clear it on denial', async (t) => {
  const { root, show } = await dom(t)
  const originalTimer = globalThis.setTimeout
  const originalClear = globalThis.clearTimeout
  let poll, status = 200
  globalThis.setTimeout = (fn, ms, ...args) => ms === 3000 ? (poll = fn, -1) : originalTimer(fn, ms, ...args)
  globalThis.clearTimeout = (id) => { if (id !== -1) originalClear(id) }
  t.after(() => { globalThis.setTimeout = originalTimer; globalThis.clearTimeout = originalClear })
  globalThis.fetch = async () => ({ ok: status === 200, status, text: async () => JSON.stringify(status === 200 ? { can_refresh: true } : { error: 'Access check failed' }) })
  function Probe() { const state = useProjectCapabilities('company'); return <pre>{JSON.stringify(state)}</pre> }
  await show(<Probe />)
  assert.equal(JSON.parse(root.textContent).data.can_refresh, true)
  status = 503; await act(async () => poll()); await settle()
  let state = JSON.parse(root.textContent)
  assert.equal(state.data.can_refresh, true); assert.equal(state.unavailable, false); assert.ok(state.error)
  status = 403; await act(async () => poll()); await settle()
  state = JSON.parse(root.textContent)
  assert.equal(state.data, null); assert.equal(state.unavailable, true)
  status = 200; await act(async () => poll()); await settle()
  state = JSON.parse(root.textContent)
  assert.equal(state.data.can_refresh, true); assert.equal(state.error, ''); assert.equal(state.unavailable, false)
})


test('PR providers distinguish local-only, query failure and successfully empty results', async (t) => {
  const { root, show } = await dom(t)
  let repository = { repo_id: 'local', repo_name: 'Local checkout', status: 'local_only', message: 'No Git remote configured.', open_count: 0, pull_requests: [] }
  globalThis.fetch = async () => ({ ok: true, status: 200, text: async () => JSON.stringify({ repositories: [repository], checked_count: repository.status === 'ok' ? 1 : 0, repo_count: 1 }) })
  await show(<PullRequestsView pid="local-company" />)
  assert.match(root.textContent, /Local checkout/)
  assert.match(root.textContent, /No Git remote configured/)
  assert.match(root.textContent, /Some repository providers are unavailable/)
  assert.doesNotMatch(root.textContent, /No open pull requests match this view/)
  repository = { ...repository, status: 'ok', message: '' }
  await show(<PullRequestsView pid="ready-company" />)
  assert.match(root.textContent, /No open pull requests match this view/)
  globalThis.fetch = async () => ({ ok: false, status: 503, text: async () => JSON.stringify({ error: 'Provider service unavailable' }) })
  await show(<PullRequestsView pid="failed-company" />)
  assert.match(root.querySelector('[role="alert"]').textContent, /Provider service unavailable/)
  assert.match(root.textContent, /Pull requests could not be loaded/)
  assert.doesNotMatch(root.textContent, /No open pull requests match this view/)
})


test('late provider responses cannot replace a different workspace', async (t) => {
  const { root, show } = await dom(t)
  let older
  const response = (name) => ({ ok: true, status: 200, text: async () => JSON.stringify({ checked_count: 0, repositories: [{ repo_id: name, repo_name: name, status: 'local_only', open_count: 0, pull_requests: [] }] }) })
  globalThis.fetch = async (url) => url.includes('/older/') ? new Promise((resolve) => { older = resolve }) : response('Current workspace repository')
  await show(<PullRequestsView pid="older" />)
  await show(<PullRequestsView pid="current" />)
  assert.match(root.textContent, /Current workspace repository/)
  await act(async () => older(response('Previous workspace repository'))); await settle()
  assert.match(root.textContent, /Current workspace repository/)
  assert.doesNotMatch(root.textContent, /Previous workspace repository/)
})


test('readiness notice separates saved graph from partial, failed and completed work', async (t) => {
 const { root, show } = await dom(t)
 for (const status of ['partial', 'failed', 'completed']) {
  await show(<ReadinessNotice readiness={{ graph: 'queryable', saved_run_id: 'older', saved_at: '2026-10-01T12:00:00Z', work: { status } }} />)
  const notice = root.querySelector('[aria-label="Workspace readiness"]')
  assert.match(notice.textContent, /older remains queryable/)
  assert.match(notice.textContent, /Saved at 2026-10-01/)
  if (status === 'completed') assert.match(notice.textContent, /coverage remains unverified/)
  else assert.match(notice.textContent, /Inspect work before retrying/)
 }
})


test('a late workspace read cannot replace the newly selected project', async (t) => {
 const { root, show } = await dom(t)
 const response = (data) => ({ ok: true, status: 200, text: async () => JSON.stringify(data) })
 let resolveOld
 globalThis.fetch = async (url) => {
  if (url.includes('/capabilities')) return response({ role: 'admin', mode: 'legacy', can_refresh: true, can_configure: true })
  if (url.includes('/ingestion')) return response({ status: 'not_started' })
  if (url.includes('/old/workspace')) return new Promise((resolve) => { resolveOld = resolve })
  return response({ project: { name: 'New workspace' }, repos: [], teams: [], readiness: { graph: 'missing', work: { status: 'empty' }, actions: { configure: true } } })
 }
 await show(<ProjectWorkspace pid="old" />)
 assert.ok(resolveOld)
 await show(<ProjectWorkspace pid="new" />)
 assert.match(root.querySelector('h1').textContent, /New workspace/)
 resolveOld(response({ project: { name: 'Old workspace' }, repos: [], teams: [] }))
 await settle()
 assert.match(root.querySelector('h1').textContent, /New workspace/)
 assert.doesNotMatch(root.textContent, /Old workspace/)
})
