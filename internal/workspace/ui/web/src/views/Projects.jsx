import { useEffect, useState } from 'preact/hooks'
import { listProjects, createProject, deleteProject, getSession } from '../lib/api.js'
import { canCreateProject } from '../lib/access.js'
import { navigate } from '../lib/router.js'
import { Modal, ConfirmDialog } from '../components/Modal.jsx'

// Projects is an explicit starting point for creating or joining a workspace.
export function Projects() {
  const [projects, setProjects] = useState(null)
  const [error, setError] = useState('')
  const [showCreate, setShowCreate] = useState(false)
  const [confirmDel, setConfirmDel] = useState(null)
  const [session, setSession] = useState(null)
  const canCreate = canCreateProject(session)

  const refresh = async () => {
    try {
      const [r, identity] = await Promise.all([listProjects(), getSession()])
      setSession(identity)
      setProjects(r.projects || [])
      setError('')
    } catch (e) {
      setError(e.message)

    }
  }
  useEffect(() => { refresh() }, [])

  const onCreated = (p) => {
    setShowCreate(false)
    navigate(`/projects/${p.id}`)
  }

  const doDelete = async (id) => {
    await deleteProject(id)
    setConfirmDel(null)
    refresh()
  }

  return (
    <div class="page">
      <header class="topbar">
        <div>
          <h1>DiffMind</h1>
          <p class="sub">Cross-service dependency graphs</p>
        </div>
        {canCreate && <button class="btn" onClick={() => setShowCreate(true)}>+ New Project</button>}
      </header>

      {error && <div class="banner error" role="alert">{error}</div>}

      <div class="content">
        {projects === null && !error && <p class="muted" role="status">Loading projects…</p>}
        {error && <button class="btn ghost" onClick={refresh}>Retry loading projects</button>}
        {projects && projects.length === 0 && !showCreate && (
          <p class="muted">{canCreate ? 'No projects yet.' : 'No accessible projects. Ask an administrator to grant your user access.'}</p>
        )}
        <div class="card-grid">
          {(projects || []).map((p) => (
            <div class="card" key={p.id}>
              <div class="card-body" onClick={() => navigate(`/projects/${p.id}`)}>
                <h3>{p.name}</h3>
                <code class="muted">{p.id}</code>
                {p.instruction && <p class="muted small">{p.instruction}</p>}
              </div>
              <div class="card-actions">
                <button class="btn ghost tiny" onClick={() => navigate(`/projects/${p.id}`)}>Open</button>
                {session?.role === 'admin' && <button class="btn danger tiny" onClick={() => setConfirmDel(p)}>Delete</button>}
              </div>
            </div>
          ))}
        </div>
      </div>

      {showCreate && (
        <CreateProject
          onClose={() => setShowCreate(false)}
          onCreated={onCreated}
        />
      )}

      {confirmDel && (
        <ConfirmDialog
          title="Delete project?"
          message={`This permanently removes project “${confirmDel.name}” and all its repos, packs, and graph runs from disk. This cannot be undone.`}
          onConfirm={() => doDelete(confirmDel.id)}
          onCancel={() => setConfirmDel(null)}
        />
      )}
    </div>
  )
}

function CreateProject({ onClose, onCreated }) {
  const [name, setName] = useState('DEFAULT')
  const [searchRoots, setSearchRoots] = useState('')
  const [instruction, setInstruction] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  const submit = async () => {
    setError('')
    if (!name.trim()) { setError('Name is required.'); return }
    setBusy(true)
    try {
      const roots = searchRoots.split('\n').map((s) => s.trim()).filter(Boolean)
      const p = await createProject({ name: name.trim(), search_roots: roots, instruction: instruction.trim() })
      onCreated(p)
    } catch (e) {
      setError(e.message)
    } finally {
      setBusy(false)
    }
  }

  return (
    <Modal title="New Project" onClose={onClose}>
      <form onSubmit={(e) => { e.preventDefault(); if (!busy) submit() }}>
      <div class="field">
        <label>Name
        <input value={name} onInput={(e) => setName(e.target.value)} /></label>
      </div>
      <div class="field">
        <label>Repository search roots (one per line, optional)
        <textarea rows="3" value={searchRoots} onInput={(e) => setSearchRoots(e.target.value)} placeholder="/path/to/repos" />
        </label>
      </div>
      <div class="field">
        <label>Default extraction instruction (optional)
        <textarea rows="2" value={instruction} onInput={(e) => setInstruction(e.target.value)} /></label>
      </div>
      {error && <div class="banner error" role="alert">{error}</div>}
      <div class="actions">
        <button class="btn" disabled={busy} type="submit">{busy ? 'Creating…' : 'Create project'}</button>
        <button type="button" class="btn ghost" onClick={onClose}>Cancel</button>
      </div>
      </form>
    </Modal>
  )
}
