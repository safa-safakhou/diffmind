import { useState } from 'preact/hooks'
import { Modal } from './Modal.jsx'
import './DependencyEvidence.css'

const statusLabel = {
  validated: 'Validated static extraction',
  compatible_unverified: 'Applicable · exact version untested',
  unknown_version: 'Version unknown',
  unsupported_version: 'Version outside support range',
  ambiguous_version: 'Conflicting versions',
  unversioned: 'Source rule without version coverage',
}
const evidenceLabel = {
  lockfile: 'Lockfile', declared_exact: 'Exact declaration', managed_exact: 'Maven management',
  configured: 'Company configuration', constraint: 'Constraint', unresolved: 'Unresolved',
}

export default function DependencyEvidence({ metrics, expanded = false }) {
  const [detailsOpen, setDetailsOpen] = useState(false)
  const [search, setSearch] = useState('')
  const [showManaged, setShowManaged] = useState(false)
  const inventory = metrics?.dependency_inventory
  const dependencies = inventory?.dependencies || []
  const coverage = metrics?.detector_coverage || []
  const declared = dependencies.filter((fact) => fact.scope !== 'management')
  const managedCount = dependencies.length - declared.length
  const filtered = (showManaged ? dependencies : declared).filter((fact) => `${fact.name} ${fact.module} ${fact.version || fact.declared || ''}`.toLowerCase().includes(search.toLowerCase()))
  if (!inventory && !coverage.length) return null
  return <div class="dependency-evidence">
    {!expanded && <button class="btn ghost" onClick={() => setDetailsOpen(true)}>Open version details</button>}
    <details open={expanded}>
      <summary>Framework and library versions ({declared.length})</summary>
      <p>Versions come from saved source inputs. They do not establish what is deployed.</p>
      <label>Find a dependency<input aria-label="Find a dependency" placeholder="Name, module, or version…" value={search} onInput={(e) => setSearch(e.currentTarget.value)} /></label>
      {managedCount > 0 && <label class="managed-versions"><input type="checkbox" checked={showManaged} onChange={(e) => setShowManaged(e.currentTarget.checked)} /> Include {managedCount} managed versions (not proof of installed dependencies)</label>}
      <ul>{filtered.slice(0, 50).map((fact) => <li key={`${fact.ecosystem}:${fact.name}:${fact.module}:${fact.scope}:${fact.sources.join(',')}`}>
        <strong>{fact.name}</strong><br />
        {fact.module} · {fact.version || fact.declared || 'unknown version'} · {evidenceLabel[fact.resolution] || fact.resolution}{fact.scope && <small> · {fact.scope}</small>}
        {fact.conditions && <small> · conditional: {fact.conditions}</small>}
        <details><summary>Version evidence</summary>{fact.declared && <p>Declared: {fact.declared}</p>}<ul>{fact.sources.map((source) => <li key={source}>{source}</li>)}</ul></details>
      </li>)}</ul>
      {!filtered.length && <p>No matching dependency evidence.</p>}
      {filtered.length > 50 && <p>Showing 50 of {filtered.length}. Refine the search to find a specific dependency.</p>}
    </details>
    {!!coverage.length && <details open={expanded}>
      <summary>Detector coverage · {coverage.filter((c) => c.status === 'validated').length} validated / {coverage.length} rules</summary>
      <p>Validation covers static extraction fixtures. Other versions still have source evidence, with the limits below.</p>
      <ul>{coverage.map((item, index) => <li key={`${item.detector_id}:${item.module}:${index}`}>
        <strong>{item.detector_id}</strong><br />{item.module} · {statusLabel[item.status] || item.status}
        <details><summary>Applied rule and versions</summary><p>{item.rule_id} · revision {item.revision}</p><p>{item.reason}</p><ul>{(item.dependencies || []).map((fact, index) => <li key={index}>{fact.name}: {fact.version || fact.declared || 'unknown'}</li>)}</ul></details>
      </li>)}</ul>
    </details>}
    {!!inventory?.limitations?.length && <details><summary>Version resolution limits ({inventory.limitations.length})</summary><ul>{inventory.limitations.map((limit) => <li key={limit}>{limit}</li>)}</ul></details>}
    {detailsOpen && <Modal title="Framework versions and detector coverage" wide onClose={() => setDetailsOpen(false)}><DependencyEvidence metrics={metrics} expanded /></Modal>}
  </div>
}
