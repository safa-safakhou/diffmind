# Workspace guide

The workspace is part of the single DiffMind repository and binary. It owns
projects, repository imports, analysis scheduling, graph construction, the
versioned query API, MCP access, and the web interface.

Start it from an installed release:

```bash
diffmind doctor
diffmind
```

Or from a source checkout:

```bash
make build
./bin/diffmind
```

Open `http://127.0.0.1:8090`, create a project, and choose **Import & build**.
The guided operation discovers a GitHub organization or local repository tree,
imports and syncs it, runs deterministic analysis, and builds the project graph.
Progress is persisted under the project, so reloads show the current phase and
the last result. The individual import, analysis, and graph actions remain
available for advanced workflows. All state defaults to `~/.diffmind`; set
`DIFFMIND_HOME` to choose another workspace, not to relocate existing data.

For the complete first-use and agent integration flow, see
[personal setup](../personal-setup.md). For implementation details,
see [the architecture guide](../ARCHITECTURE.md).

## Ingestion API

The web workflow uses a project-scoped asynchronous endpoint:

```text
POST /api/projects/{project}/ingestion
GET  /api/projects/{project}/ingestion
```

The POST body may contain an `import` object using the same GitHub/local fields
as repository import, plus `concurrency` and deterministic analysis `options`.
Omit `import` to rebuild all repositories already registered in the project.
POST returns `202 Accepted`; GET reports `running`, `completed`, `partial`, or
`failed`, the current phase, counters, graph run ID, and errors. A project
ingestion is exclusive with fleet refresh and conflicting project mutations.

## Entry-point flows and PR review

Select a service in the graph and choose **Browse service flows** to inspect one
HTTP endpoint, queue consumer, job, webhook, RPC endpoint, or command at a time.
The flow shows extracted local steps, data operations and matched downstream
services. Nodes open the saved evidence. Large diagrams scroll at a readable
size; conditional and async paths retain their recorded meaning.

**PR impact** includes a before/after flow review. It queries the provider for
immutable base/head SHAs and the actual merge-base, then finds completed graph
snapshots containing clean analyses of that repository at those exact commits.
The latest 500 runs are searched. Missing snapshots are reported with the
required commits; this read does not switch branches or start analysis. Capture
those revisions in separate checkouts using the existing analysis/build workflow.
Other repositories and analyzer/pack/file-scope inputs can differ between the
snapshots, so the comparison retains those inputs and does not claim causation.
Removed entry points keep their baseline evidence and callers.

The same per-entrypoint result is available to agents:

```text
GET /api/v1/projects/{pid}/graph/flows?from=BASE_RUN&to=HEAD_RUN&service=SERVICE
GET /api/projects/{pid}/pull-requests/{repo_id}/{number}/flows
MCP compare_flows(project, from, to, service, offset?, limit?)
```

Use identical `from`/`to` run IDs to browse a service. Use **Compare graphs** and
its service selector to explore flow changes between arbitrary saved snapshots.
Both HTTP and MCP return `flows` with before/after diagrams, source evidence,
changes, partial status, and snapshot inputs. Follow `next_offset` to retrieve
all results; the default page size is 10 and the maximum is 50. Comparisons
return only added, removed, or semantically modified flows, with pagination
after filtering. Metadata-only changes are omitted. Service browsing includes
all entry points. The PR screen links to this separate service browser. Each flow
is bounded to four service hops and 200 diagram nodes. Entrypoints pair by exact
service, kind and name; renames are removal plus addition, and ambiguous
identities require better extraction rather than a guessed pairing.

PR-specific access is exposed as the read-only `pull_request_flows` operation in
the agent management catalog. Optional `from`/`to` pin older graph run IDs and
are validated against the provider's current merge-base and head revisions.
Static evidence, including an unchanged or partial flow, never establishes
runtime coverage or merge safety.

### Preparing PR flow evidence

The PR flow screen offers **Prepare PR flows** when merge-base/head evidence is
missing. Agents can call `prepare_pull_request_flows` through the management
catalog, or POST `/api/projects/{pid}/pull-requests/{repo_id}/{number}/flows/prepare`.
The request runs deterministic analysis in a private clone at the provider's
exact merge-base and head commits, with a five-minute limit. It preserves local
branches, dirty files, registered repository analysis state, and company graph
history. Both sides use the same recorded company context. Only tracked source
is analyzed; untracked hints and project graph supplements are excluded.

Prepared evidence is cached separately under `pr-reviews` by project, repository,
head revision and analyzer identity. A changed merge-base requires preparation
again. The ordinary GET remains read-only and reports missing snapshots.
The impact endpoint uses prepared head evidence when available; explicit saved
`run_id` selections retain their original semantics. Evidence-only revision or
location changes stay visible in the evidence details without marking a flow as
behaviorally changed. An entrypoint without a local extracted path remains
partial and does not inherit the repository-wide dependency graph.
