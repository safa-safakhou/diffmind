import { LoadingState } from '../components/LoadingState.jsx'
import { useEffect, useMemo, useState, useRef } from 'preact/hooks'
import dagre from 'dagre'
import { compareFlows, getPullRequestFlows, preparePullRequestFlows } from '../lib/api.js'
import { navigate } from '../lib/router.js'
import './FlowReview.css'

const words = (value) => String(value || '').replaceAll('_', ' ')
const shortCommit = (value) => value ? value.slice(0, 8) : 'unknown'
const changeLabel = { added: 'Added', removed: 'Removed', modified: 'Changed', unchanged: 'Unchanged' }

export function PullRequestFlowPage({ pid, repo, number }) {
  return <div class="trace-page">
    <header class="run-topbar"><button class="btn ghost tiny" onClick={() => navigate(`/projects/${encodeURIComponent(pid)}/pull-requests`)}>← Pull requests</button><h1 class="run-title">PR #{number} · Flow review</h1></header>
    <main class="trace-service-flows"><FlowReview pid={pid} pr={{ repo_id: repo, number }} /></main>
  </div>
}

export function FlowReview({ pid, from, to, service, pr, onPrepared }) {
  const preparationScope = useRef('')
  const [preparing, setPreparing] = useState(false)
  const [prepareError, setPrepareError] = useState('')
  const [offset, setOffset] = useState(0)
  const [retry, setRetry] = useState(0)
  const [state, setState] = useState({ loading: true })
  const [selected, setSelected] = useState('')
  const [search, setSearch] = useState('')
  const scope = JSON.stringify([pid, from, to, service, pr?.repo_id, pr?.number, pr?.head_sha])
  const [pageScope, setPageScope] = useState(scope)
  const currentOffset = pageScope === scope ? offset : 0
  useEffect(() => {
    let alive = true
    preparationScope.current = scope
    setPreparing(false); setPrepareError('')
    setState({ loading: true }); setSelected(''); setSearch('')
    const request = pr
      ? getPullRequestFlows(pid, pr.repo_id, pr.number, { offset: currentOffset })
      : compareFlows(pid, from, to, service, currentOffset)
    request.then((response) => {
      if (!alive) return
      const data = pr ? response.review : response
      setState({ data, context: pr ? response : null })
      setSelected(data?.flows?.find((flow) => flow.change !== 'unchanged')?.key || data?.flows?.[0]?.key || '')
    }).catch((error) => { if (alive) setState({ error: error.message }) })
    return () => { alive = false; preparationScope.current = '' }
  }, [scope, currentOffset, retry])
  const prepare = async () => {
    const startedScope = scope
    setPreparing(true); setPrepareError('')
    try {
      const response = await preparePullRequestFlows(pid, pr.repo_id, pr.number)
      if (preparationScope.current !== startedScope) return
      setState({ data: response.review, context: response })
      setSelected(response.review?.flows?.find((flow) => flow.change !== 'unchanged')?.key || response.review?.flows?.[0]?.key || '')
      go(0); onPrepared?.()
    } catch (error) { if (preparationScope.current === startedScope) setPrepareError(error.message) }
    finally { if (preparationScope.current === startedScope) setPreparing(false) }
  }
  const go = (value) => { setPageScope(scope); setOffset(value) }
  const data = state.data
  const context = state.context
  const single = from === to && !pr
  const flows = data?.flows || []
  const chosen = flows.find((flow) => flow.key === selected)
  const visible = flows.filter((flow) => `${flow.name} ${flow.kind}`.toLowerCase().includes(search.toLowerCase()))
  return <section class="flow-review" aria-label={single ? 'Service flows' : 'Entrypoint flow review'}>
    <header class="fr-heading">
      <div><span class="fr-eyebrow">{single ? 'EXPLORE' : 'REVIEW BY ENTRY POINT'}</span><h2>{single ? 'Service flows' : 'What changed in the flow?'}</h2><p>{single ? 'Follow one entry point at a time, from the request to its dependencies.' : 'Compare each entry point, its data operations, and the services around it.'}</p></div>
      <button class="btn ghost tiny" onClick={() => setRetry((v) => v + 1)} disabled={state.loading || preparing}>Refresh flows</button>
    </header>
    {state.loading && <LoadingState label={pr ? 'Finding merge-base and head snapshots…' : 'Loading service flows…'} detail="Reading saved entry points and their dependencies." />}
    {state.error && <p class="banner error" role="alert">{state.error} <button class="btn ghost tiny" onClick={() => setRetry((v) => v + 1)}>Retry</button></p>}
    {context && <div class="fr-revisions"><span>Merge-base <code title={context.base_commit}>{shortCommit(context.base_commit)}</code> {context.before_run ? '· snapshot found' : '· snapshot missing'}</span><span>PR head <code title={context.head_commit}>{shortCommit(context.head_commit)}</code> {context.after_run ? '· snapshot found' : '· snapshot missing'}</span></div>}
    {context && !context.available && <div class="fr-empty"><h3>Two saved snapshots are needed</h3><p>{context.next_action}</p>{preparing ? <LoadingState label="Preparing PR flows…" detail="Analyzing the merge-base and PR head in isolated checkouts. This can take up to five minutes." /> : <button class="btn primary" onClick={prepare}>Prepare PR flows</button>}{prepareError && <p class="banner error" role="alert">{prepareError}</p>}</div>}
    {data && <>
      <div class="fr-evidence-bar">{!single && <button class="btn ghost tiny" onClick={() => navigate(`/projects/${encodeURIComponent(pid)}/runs/${encodeURIComponent(context?.explore_run || data.to.id)}/trace?${new URLSearchParams({ service: data.service })}`)}>Explore all service flows →</button>}<span>{data.service} · {data.total} {single ? '' : 'changed '}entry point{data.total === 1 ? '' : 's'}</span><span>{single ? 'Saved snapshot' : pr ? 'PR revisions matched' : 'Saved snapshot comparison'} · static evidence</span></div>
      {!single && !!data.environment_changes?.length && <details class="fr-method"><summary>Dependency and framework version changes ({data.environment_changes.length})</summary><p>Review these separately from the flow diagrams. An unchanged flow does not establish compatibility.</p><ul>{data.environment_changes.map((change) => <li key={`${change.ecosystem}:${change.module}:${change.name}:${change.scope}`}><strong>{change.name}</strong> · {change.module}: {(change.before || []).map((fact) => fact.version || fact.declared || 'unknown').join(', ') || 'absent'} → {(change.after || []).map((fact) => fact.version || fact.declared || 'unknown').join(', ') || 'absent'}</li>)}</ul></details>}
      {!flows.length && <p class="fr-notice">{single ? 'No entry points were extracted in this snapshot.' : 'No extracted flow changes.'}</p>}
      {!!flows.length && <div class="fr-workspace">
        <aside class="fr-picker" aria-label="Entry points">
          <label>Find an entry point<input value={search} onInput={(e) => setSearch(e.currentTarget.value)} placeholder="Endpoint or consumer…" /></label>
          {visible.map((flow) => <button class={'fr-entry ' + (flow.key === selected ? 'active' : '')} key={flow.key} onClick={() => setSelected(flow.key)} aria-pressed={flow.key === selected}>
            <span class="fr-entry-kind">{words(flow.kind)}</span><strong>{flow.name}</strong><span class={'fr-status ' + flow.change}>{single ? 'Saved flow' : changeLabel[flow.change]}{flow.partial ? ' · partial evidence' : ''}</span>
          </button>)}
          {!visible.length && <p class="fr-notice">No matches on this page.</p>}
          <nav class="fr-pages" aria-label="Entry point pages"><button class="btn ghost tiny" disabled={!currentOffset} onClick={() => go(Math.max(0, currentOffset - 10))}>Previous</button><span>{currentOffset + 1}–{currentOffset + flows.length} of {data.total}</span><button class="btn ghost tiny" disabled={data.next_offset == null} onClick={() => go(data.next_offset)}>Next</button></nav>
        </aside>
        {chosen && <FlowReviewDetail key={chosen.key} flow={chosen} single={single} from={data.from.id} to={data.to.id} />}
      </div>}
      <details class="fr-method"><summary>Evidence and limits</summary><ul>{(data.notes || []).map((note) => <li key={note}>{note}</li>)}</ul><details><summary>Snapshot inputs</summary><pre>{JSON.stringify({ before: data.inputs_before, after: data.inputs_after }, null, 2)}</pre></details></details>
    </>}
  </section>
}

export function FlowReviewDetail({ flow, single = false, from, to }) {
  const [picked, setPicked] = useState(null)
  const changes = (flow.changes || []).filter((change) => !change.evidence_only)
  const layout = useMemo(() => pairedLayout(flow.before?.graph, flow.after?.graph), [flow])
  const meaningful = changes.filter((c) => !c.evidence_only && (c.kind === 'node' || c.kind === 'caller' || c.kind === 'data') && !(c.kind === 'node' && c.label === flow.service))
  const summaries = meaningful.length ? meaningful : changes.filter((c) => c.kind === 'connection')
  const callers = Array.from(new Set([...(flow.before?.callers || []), ...(flow.after?.callers || [])].map((edge) => edge.from)))
  return <article class="fr-detail">
    <header class="fr-detail-heading"><span class="fr-entry-kind">{words(flow.kind)}</span><h3>{flow.name}</h3><span class={'fr-status ' + flow.change}>{single ? 'Saved flow' : changeLabel[flow.change]}</span></header>
    {flow.partial && <p class="fr-partial">Partial evidence. Some steps or dependencies were not extracted; this diagram may be incomplete.</p>}
    {!single && <div class="fr-summary"><strong>{flow.change === 'unchanged' ? 'No extracted flow changes' : flow.change === 'added' ? 'This entry point was added' : flow.change === 'removed' ? 'This entry point was removed' : 'Changes to review'}</strong><ul>{summaries.slice(0, 8).map((change) => <li key={change.key}><span class={'fr-dot ' + change.change} />{changeLabel[change.change]} {change.kind === 'caller' ? 'caller: ' : ''}<b>{change.label}</b>{change.evidence_only ? ' · evidence only' : ''}</li>)}</ul>{summaries.length > 8 && <small>{summaries.length - 8} more changes in the evidence below.</small>}</div>}
    <div class="fr-legend">{!single && <><span><i class="fr-dot added" />Added</span><span><i class="fr-dot removed" />Removed</span><span><i class="fr-dot modified" />Changed</span></>}<span>Dashed arrow · async hop</span><span>Scroll diagrams to explore →</span></div>
    {!single && <DiagramPanel title="Before" run={from} snapshot={flow.before} layout={layout} changes={changes} onSelect={(node) => setPicked({ node, side: 'Before' })} />}
    <DiagramPanel title={single ? 'Flow' : 'After'} run={to} snapshot={flow.after} layout={layout} changes={single ? [] : changes} onSelect={(node) => setPicked({ node, side: single ? 'Flow' : 'After' })} />
    {!!callers.length && <div class="fr-callers"><span>Nearby callers</span>{callers.map((caller) => <span class="fr-caller" key={caller}>{caller}</span>)}<small>Exact matches in either snapshot</small></div>}
    {picked && <div class="fr-node-evidence"><button class="btn ghost tiny" onClick={() => setPicked(null)}>Close evidence</button><strong>{picked.side} · {picked.node.label}</strong><p>{words(picked.node.kind)}{picked.node.service ? ` · ${picked.node.service}` : ''}</p><NodeEvidence details={picked.node.details} /></div>}
    {!single && <details class="fr-method"><summary>Change evidence ({changes.length})</summary>{changes.map((change) => <details key={change.key}><summary>{changeLabel[change.change]} · {change.label}{change.evidence_only ? ' · evidence only' : ''}</summary><div class="fr-raw-pair"><div><strong>Before</strong><pre>{JSON.stringify(change.before ?? null, null, 2)}</pre></div><div><strong>After</strong><pre>{JSON.stringify(change.after ?? null, null, 2)}</pre></div></div></details>)}</details>}
  </article>
}

function diagramNodes(graph) {
  const groups = new Map()
  for (const node of graph?.nodes || []) {
    if (node.kind === 'service' && node.service === graph.entry?.service) continue
    const previous = groups.get(node.id)
    groups.set(node.id, previous ? { ...node, details: { ...node.details, occurrences: [...(previous.details?.occurrences || [previous.details]), node.details] } } : node)
  }
  return Array.from(groups.values())
}

export function pairedLayout(before, after) {
  const nodes = new Map([...diagramNodes(before), ...diagramNodes(after)].map((node) => [node.id, node]))
  const g = new dagre.graphlib.Graph({ multigraph: true })
  g.setGraph({ rankdir: 'LR', nodesep: 34, ranksep: 64, marginx: 24, marginy: 30 })
  g.setDefaultEdgeLabel(() => ({}))
  for (const id of Array.from(nodes.keys()).sort()) g.setNode(id, { width: 184, height: 82 })
  const edges = new Map([...(before?.edges || []), ...(after?.edges || [])].map((edge) => [JSON.stringify([edge.from, edge.to, edge.kind]), edge]))
  for (const key of Array.from(edges.keys()).sort()) { const edge = edges.get(key); if (nodes.has(edge.from) && nodes.has(edge.to)) g.setEdge(edge.from, edge.to, { weight: edge.cycle ? 0 : 1 }, key) }
  dagre.layout(g)
  return { positions: new Map(g.nodes().map((id) => [id, g.node(id)])), width: g.graph().width || 280, height: g.graph().height || 150 }
}

function DiagramPanel({ title, run, snapshot, layout, changes, onSelect }) {
  const graph = snapshot?.graph
  const nodes = diagramNodes(graph)
  const present = new Set(nodes.map((node) => node.id))
  const states = new Map(changes.filter((c) => c.kind === 'node' && !c.evidence_only).map((change) => [change.key.slice(5), change.change]))
  const edgeStates = new Map(changes.filter((c) => c.kind === 'edge' && !c.evidence_only).map((change) => [change.key.slice(5), change.change]))
  const marker = `fr-arrow-${title.toLowerCase()}`
  return <section class="fr-diagram-panel" aria-label={`${title} flow`}>
    <header><strong>{title}</strong><code>{run}</code>{graph && <small>{nodes.length} steps · {graph.status === 'complete' ? 'extracted flow' : `${graph.status} evidence`}</small>}</header>
    {!snapshot ? <p class="fr-notice">Entry point {title === 'Before' ? 'not present in the baseline' : 'removed in this snapshot'}.</p> : <div class="fr-diagram-scroll"><svg width={layout.width} height={Math.max(154, layout.height)} role="img" aria-label={`${title}: ${snapshot.entry.name}`}>
      <defs><marker id={marker} viewBox="0 0 10 10" refX="9" refY="5" markerWidth="6" markerHeight="6" orient="auto"><path d="M0 0L10 5L0 10" fill="#7385a4" /></marker></defs>
      {(graph.edges || []).filter((edge) => present.has(edge.from) && present.has(edge.to)).map((edge, i) => {
        const a = layout.positions.get(edge.from), b = layout.positions.get(edge.to)
        const x1 = a.x + a.width / 2, x2 = b.x - b.width / 2, mid = (x1 + x2) / 2
        const d = edge.cycle ? `M${x1},${a.y} C${x1 + 40},${a.y - 70} ${x2 - 40},${b.y - 70} ${x2},${b.y}` : `M${x1},${a.y} C${mid},${a.y} ${mid},${b.y} ${x2},${b.y}`
        const state = edgeStates.get(JSON.stringify([edge.from, edge.to, edge.kind])) || 'unchanged'
        const condition = edge.condition?.summary || (['conditional', 'may'].includes(edge.reachability) ? edge.reachability : '')
        return <g key={i} class={`fr-edge ${state}` + (edge.async ? ' async' : '')}><path d={d} marker-end={`url(#${marker})`} /><title>{words(edge.kind)}{edge.reachability ? ` · ${edge.reachability}` : ''}</title>{condition && <text x={mid} y={(a.y + b.y) / 2 - 9} text-anchor="middle">{String(condition).slice(0, 24)}</text>}</g>
      })}
      {nodes.map((node) => {
        const p = layout.positions.get(node.id), state = states.get(node.id) || 'unchanged'
        const label = node.label || node.id, lines = wrapLabel(label)
        return <g key={node.id} class={`fr-node ${state}`} transform={`translate(${p.x - p.width / 2},${p.y - p.height / 2})`} role="button" tabindex="0" aria-label={`Inspect ${title} ${label}`} onClick={() => onSelect(node)} onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); onSelect(node) } }}>
          <title>{label} · {words(node.kind)}</title><rect width={p.width} height={p.height} rx="10" /><text class="fr-node-kind" x="12" y="20">{words(node.kind).slice(0, 25)}</text>{lines.map((line, i) => <text key={i} class="fr-node-label" x="12" y={43 + i * 17}>{line}</text>)}
        </g>
      })}
    </svg></div>}
  </section>
}

function wrapLabel(label) {
  if (label.length <= 24) return [label]
  let split = label.lastIndexOf(' ', 24)
  if (split < 8) split = 24
  const rest = label.slice(split).trim()
  return [label.slice(0, split), rest.length > 24 ? rest.slice(0, 23) + '…' : rest]
}

function NodeEvidence({ details }) {
  const values = details || {}
  const fields = ['summary', 'operation', 'table', 'database', 'method', 'path', 'url', 'topic', 'queue', 'platform'].filter((key) => values[key] && typeof values[key] !== 'object')
  const locations = Array.isArray(values.source_locations) ? values.source_locations : []
  return <div>{fields.length > 0 && <dl class="fr-evidence-fields">{fields.map((key) => <div key={key}><dt>{words(key)}</dt><dd>{String(values[key])}</dd></div>)}</dl>}
    {locations.map((location, i) => <p class="fr-source" key={i}>{location.file || location.path}{location.start_line || location.line ? `:${location.start_line || location.line}` : ''}</p>)}
    {!fields.length && !locations.length && <p>No additional source evidence was saved for this step.</p>}
    {!!Object.keys(values).length && <details><summary>Full saved evidence</summary><pre>{JSON.stringify(values, null, 2)}</pre></details>}
  </div>
}
