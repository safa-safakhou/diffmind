// Transport/access failures are client observations: the server cannot report
// its own disconnection. Preserve provenance, but never retain action authority.
export function connectionReadiness(saved, { loading = false, error = '', denied = false } = {}) {
  if (denied || error || loading) return {
    ...saved,
    runtime: denied ? 'available' : error ? 'disconnected' : 'unknown',
    access: denied ? 'denied' : 'unknown',
    graph: denied ? 'unknown' : saved?.graph || 'unknown',
    actions: { query: false, refresh: false, configure: false, inspect_work: false },
    next_action: denied ? 'request_access' : error ? 'reconnect' : 'wait',
  }
  return saved
}

export function readinessMessage(state) {
  if (!state) return 'Checking workspace readiness…'
  if (state.access === 'denied') return 'Workspace access is unavailable. Request access or return to Projects.'
  if (state.runtime === 'disconnected') return 'Connection unavailable. Showing the last loaded view; reconnect before using context or changing it.'
  if (state.access === 'unknown') return 'Checking workspace access; actions are paused.'
  const graph = state.graph === 'queryable' ? `Saved graph ${state.saved_run_id} remains queryable.` : state.graph === 'unavailable' ? 'Saved graph is unavailable. Inspect the run and restore or rebuild it.' : 'No saved graph yet.'
  const work = state.work?.status || 'unknown'
  const message = {
    empty: 'Context has not been built.', queued: 'Refresh is queued.', running: `Work is running (${state.work?.phase || 'unknown phase'}).`,
    cancelling: 'Work is stopping.', completed: 'Work completed; coverage remains unverified.',
    partial: 'Work finished with partial results. Inspect work before retrying.', failed: 'Work failed. Inspect work before retrying.',
    cancelled: 'Work was cancelled.', interrupted: 'Work was interrupted. Inspect work before retrying.',
  }[work] || 'Current work status is unknown.'
  return `${graph} ${message}`
}
