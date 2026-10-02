import test from 'node:test'
import assert from 'node:assert/strict'
import { connectionReadiness, readinessMessage } from './readiness.js'

test('disconnection preserves provenance and pauses every action; denial loses graph claims', () => {
  const saved = { graph: 'queryable', saved_run_id: 'old', actions: { query: true, refresh: true, configure: true } }
  const disconnected = connectionReadiness(saved, { error: '503' })
  assert.equal(disconnected.saved_run_id, 'old')
  assert.equal(disconnected.runtime, 'disconnected')
  assert.ok(Object.values(disconnected.actions).every((value) => !value))
  assert.equal(disconnected.next_action, 'reconnect')
  const denied = connectionReadiness(saved, { denied: true })
  assert.equal(denied.access, 'denied')
  assert.equal(denied.graph, 'unknown')
  assert.equal(denied.next_action, 'request_access')
  assert.equal(connectionReadiness(saved), saved)
})

test('partial and failed maintenance preserve the saved graph explanation', () => {
  for (const status of ['partial', 'failed']) {
    const message = readinessMessage({ graph: 'queryable', saved_run_id: 'old', work: { status } })
    assert.match(message, /old remains queryable/)
    assert.match(message, /Inspect work before retrying/)
  }
  assert.match(readinessMessage({ graph: 'missing', work: { status: 'completed' } }), /No saved graph yet.*coverage remains unverified/)
})
