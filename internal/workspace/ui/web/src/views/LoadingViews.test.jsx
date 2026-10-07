import test from 'node:test'
import assert from 'node:assert/strict'
import { parseHTML } from 'linkedom'
import { render } from 'preact'
import { act } from 'preact/test-utils'
import { ProjectWorkspace } from './ProjectWorkspace.jsx'
import { RunView } from './RunView.jsx'

const settle = () => act(async () => { await new Promise(setImmediate) })
const response = (data, status = 200) => ({ ok: status === 200, status, text: async () => JSON.stringify(data) })
const deferred = () => { let resolve; const promise = new Promise((r) => { resolve = r }); return { promise, resolve } }
async function dom(t) {
  const { window, document } = parseHTML('<html><body><div id="root"></div></body></html>')
  const previous = { window: globalThis.window, document: globalThis.document, fetch: globalThis.fetch, EventSource: globalThis.EventSource }
  Object.assign(globalThis, { window, document, EventSource: class { addEventListener() {} close() {} } })
  const root = document.getElementById('root')
  t.after(async () => { await act(() => render(null, root)); Object.assign(globalThis, previous) })
  return { root, show: async (element) => { await act(() => render(element, root)); await settle() } }
}
const workspace = { project: { name: 'Orders' }, repos: [], teams: [], readiness: { graph: 'missing', actions: {} } }
function workspaceFetch(metadata, graph) {
  return async (url) => url.includes('/capabilities') ? response({ role: 'viewer' }) : url.includes('/ingestion') ? response({ status: 'not_started' }) : url.includes('/archgraph') ? graph() : metadata()
}

test('project metadata shows loading instead of premature empty graph; failure settles and retry recovers', async (t) => {
  const { root, show } = await dom(t), request = deferred()
  globalThis.fetch = workspaceFetch(() => request.promise)
  await show(<ProjectWorkspace pid="orders" />)
  const board = root.querySelector('.workspace-board')
  assert.match(board.querySelector('[role="status"]').textContent, /Loading workspace/)
  assert.equal(board.getAttribute('aria-busy'), 'true')
  assert.doesNotMatch(root.textContent, /No saved graph yet|no graph yet/)
  await act(() => request.resolve(response({ error: 'Service unavailable' }, 503))); await settle()
  assert.equal(board.querySelector('.loading-state'), null)
  assert.doesNotMatch(board.textContent, /No saved graph yet/)
  globalThis.fetch = workspaceFetch(() => response(workspace))
  await act(() => [...root.querySelectorAll('button')].find((b) => b.textContent === 'Retry loading workspace').click()); await settle()
  assert.equal(board.getAttribute('aria-busy'), 'false')
  assert.match(board.textContent, /No saved graph yet/)
})

test('saved graph loading ends on failure; retry displays progress again', async (t) => {
  const { root, show } = await dom(t), graph = deferred()
  globalThis.fetch = workspaceFetch(() => response({ ...workspace, latest_run: { id: 'saved', status: 'completed' } }), () => graph.promise)
  await show(<ProjectWorkspace pid="orders" />)
  const board = root.querySelector('.workspace-board')
  assert.match(board.querySelector('[role="status"]').textContent, /Loading graph/)
  assert.doesNotMatch(board.textContent, /No saved graph yet/)
  await act(() => graph.resolve(response({ error: 'Graph unavailable' }, 503))); await settle()
  assert.equal(board.querySelector('[role="status"]'), null)
  assert.equal(board.getAttribute('aria-busy'), 'false')
  assert.match(board.textContent, /Could not load the saved graph/)
  const retry = deferred()
  globalThis.fetch = workspaceFetch(() => response(workspace), () => retry.promise)
  await act(() => [...root.querySelectorAll('button')].find((b) => b.textContent === 'Retry loading graph').click()); await settle()
  assert.match(board.querySelector('[role="status"]').textContent, /Loading graph/)
  await act(() => retry.resolve(response({ error: 'Still unavailable' }, 503))); await settle()
  assert.equal(board.querySelector('[role="status"]'), null)
})

test('run view distinguishes a pending run and graph request from an unavailable graph', async (t) => {
  const { root, show } = await dom(t), run = deferred(), graph = deferred()
  globalThis.fetch = async (url) => url.includes('/capabilities') ? response({ role: 'viewer' }) : url.includes('/archgraph') ? graph.promise : run.promise
  await show(<RunView pid="orders" rid="saved" />)
  assert.match(root.querySelector('.run-stage [role="status"]').textContent, /Loading run/)
  assert.doesNotMatch(root.textContent, /No graph available/)
  await act(() => run.resolve(response({ run: { status: 'completed' } }))); await settle()
  assert.match(root.querySelector('.run-stage [role="status"]').textContent, /Loading graph/)
  await act(() => graph.resolve(response({ error: 'Read failed' }, 503))); await settle()
  assert.equal(root.querySelector('.run-stage [role="status"]'), null)
  assert.match(root.querySelector('[role="alert"]').textContent, /Read failed/)
  assert.doesNotMatch(root.textContent, /No graph available/)
})
