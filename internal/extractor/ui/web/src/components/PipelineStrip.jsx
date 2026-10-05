import { stages, runMeta } from '../lib/store.js'

// Deterministic pipeline stages in event-emission order.
const ORDER = ['ast_index', 'deterministic_discovery', 'connections', 'reconcile']
const PRETTY = {
  ast_index: 'AST Index',
  deterministic_discovery: 'Deterministic',
  connections: 'Connections',
  reconcile: 'Reconcile',
}

export function PipelineStrip() {
  const map = stages.value
  const meta = runMeta.value
  const terminal = ['completed', 'failed', 'cancelled'].includes(meta?.status)
  const started = Date.parse(meta?.startedAt)
  const finished = Date.parse(meta?.finishedAt)
  const duration = started > 0 ? (terminal ? finished : Date.now()) - started : NaN
  const elapsed = Number.isFinite(duration) && duration >= 0 ? humanDuration(duration) : '\u2013'

  return (
    <div class="pipeline-strip" tabIndex={0} aria-label="Extraction stages">
      {ORDER.map((id) => {
        const s = map.get(id) || { name: id, status: terminal ? 'unrecorded' : 'pending', summary: {} }
        const summary = formatSummary(s.summary)
        const status = terminal && s.status === 'pending' ? 'unrecorded' : s.status
        return (
          <div class={'pipeline-stage ' + status} key={id}>
            <div class="name">{PRETTY[id]}</div>
            <div class="count">{statusLabel(status)}</div>
            <div class="stage-tip">{s.tip}</div>
            <div class="stage-tokens">{summary}</div>
          </div>
        )
      })}
      <div class="pipeline-stage" style="flex: 0 0 160px; border-right: none;">
        <div class="name">Elapsed</div>
        <div class="count">{elapsed}</div>
        <div class="stage-tip">{meta?.status || 'waiting'}</div>
        <div class="stage-tokens" />
      </div>
    </div>
  )
}

function statusLabel(status) {
  if (status === 'unrecorded') return 'not recorded'
  if (status === 'success') return 'done'
  if (status === 'running') return 'running'
  if (status === 'failed') return 'failed'
  if (status === 'cancelled') return 'cancelled'
  if (status === 'skipped') return 'skipped'
  return 'pending'
}

function formatSummary(summary) {
  if (!summary || typeof summary !== 'object') return ''
  const pairs = Object.entries(summary)
    .filter(([k, v]) => typeof v === 'number' && Number.isFinite(v) && !k.endsWith('_ms'))
    .slice(0, 3)
  return pairs.map(([k, v]) => `${k}: ${v}`).join(' · ')
}

function humanDuration(ms) {
  if (!Number.isFinite(ms) || ms < 0) return '–'
  const s = Math.floor(ms / 1000)
  if (s < 60) return s + 's'
  const m = Math.floor(s / 60)
  const r = s % 60
  if (m < 60) return m + 'm ' + r + 's'
  const h = Math.floor(m / 60)
  return h + 'h ' + (m % 60) + 'm'
}
