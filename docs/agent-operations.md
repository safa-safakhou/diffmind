# Agent-managed operations

DiffMind supports an agent-operated product workflow, not just agent queries.
The [host-agent playbook](../AGENT_SETUP.md) handles installation and MCP
registration. Users provide intent and access; the agent operates the platform.

## Connection modes and authority

| Connection | Tools | Authority/lifecycle |
| --- | --- | --- |
| Local `diffmind agent` | 14 read + 3 management + 2 host tools | Full local workspace control; starts backend automatically, owns it until disconnect/crash |
| Local `diffmind mcp` | 14 read tools | Original trusted read-only integration; no backend ownership |
| HTTP `/mcp`, viewer | 14 read tools | Read-only, restricted to accessible projects |
| HTTP `/mcp`, editor/admin | Graph and management tools | Same role/membership/host-operation checks as the HTTP API; no local lifecycle/CLI tools |

Local mode is trusted OS-level access, not a sandbox. Its backend binds only
loopback on an automatically reserved port, without shared-server authentication;
do not expose/forward it. Use a secured shared deployment on multi-user or remote
hosts. `--project` is an optional default selector, not a permission boundary.
Admin tokens grant full platform administration; viewer tokens are intentionally
insufficient when the user wants the agent to manage projects.

## Background maintenance and primary actions

A new local agent workspace refreshes its registered repositories on connection
and every 15 minutes while connected. It never expands the imported repository
set. Existing explicit settings are preserved. Runtime status reports the actual
policy; use `agent_runtime` to change it while preserving other settings. Setting
`refresh_interval` to `"0"` and `refresh_on_start` to `false` selects manual updates.
The owning connection controls availability; disconnecting stops background work.

The dashboard's **Update context** refreshes registered repositories and builds
context. In scoped company mode, editors enqueue the saved configuration, while
administrators can also configure/import sources. Viewers explore saved evidence
and ask an editor to update it. **Reload view** only reloads displayed data.
Manual analyzer and graph controls are under **Advanced actions**.

Graph summaries distinguish a saved run, repository analysis freshness and
unverified static coverage. `freshness_basis: live_checkout_status` uses the same read-only check as the dashboard; `freshness_reference: latest_repository_analysis` compares current source against the latest repository analysis, not the selected graph snapshot. Missing or inaccessible local checkouts report unknown.
These states do not assert that the chosen graph incorporates a newer analysis
or that PR-head evidence is eligible.

Relationship `evidence` records preserve each linked object's class, answering
run, recorded revision/source locations, file-scope policy, pack declaration and
identity-resolution reason. `source_extracted`, `source_inferred`, `pack_declared`,
`declared` and `unknown` are distinct; even a runtime-origin label remains
unverified. Missing findings may be excluded or unsupported. Report gaps
privately to the administrator using the existing improvement workflow; use
synthetic positive and negative fixtures before changing detectors/config/packs.

`compare_graphs` returns `inputs_before` and `inputs_after` from the pinned saved
snapshots, including analysis artifact references, revision, schema, analyzer and
recorded file scope. Missing historical inputs stay unknown. `evidence_only`
labels changes limited to recognized provenance/confidence fields; original
before/after evidence is retained. Changed scope, packs or artifacts are possible
inputs, not an established cause. Use `compare_contracts` for request fields.

PR inspection is on demand through `inspect_workspace(operation="pull_request_impact")`.
Read changed files, attention signals, exact caller matches, separate candidates,
eligibility and limitations before the uncalibrated score. A low score is not a
probability or merge recommendation. Refreshing a default branch does not ensure
a PR-head match. With explicit authority, analyze a separate clean PR-head checkout
and select its saved graph via `run_id`; do not switch the user's active branch.
Deleted surfaces require separate matching baseline evidence; internal/transitive
and configuration effects may remain unproven. This query neither captures new
revisions nor posts code-host comments/checks.

Management failures return a conservative `recovery` category and next action
alongside their HTTP status; browser feedback carries the same guidance. Retain
the input, reload/review conflicts, and inspect persisted work after a lost write
response. `retryable:false` prevents interpreting a status as permission to
replay a mutation. A 202 response is still acceptance rather than completion.

## Discover, inspect, mutate

`describe_management` lists the finite operation catalog, method/path,
description, body example and destructive marker. Supply `operation` to inspect
one entry. It documents projects, repositories, imports, ingestion, jobs, packs,
configuration, access, credentials, quotas, graph runs and pull request inspection.
Every existing browser mutation except signed provider webhook delivery is
covered by a catalog operation; an automated route-parity test guards this.

`inspect_workspace` accepts GET operations. `manage_workspace` accepts
mutations. Both take:

```json
{
  "operation": "start_ingestion",
  "selectors": {"pid": "PROJECT_ID"},
  "body": {
    "import": {
      "provider": "local",
      "root": "/absolute/path/to/repositories",
      "include": "^(gateway|catalog|billing)$"
    },
    "concurrency": 2
  }
}
```

Selectors correspond to placeholders in the catalog: `pid`, `rid` (repository
or graph run, depending on the operation), `jid`, `pack_id`, `tid`, etc.
Query parameters use a `query` string map. There is no arbitrary URL/header/method
input. Traversal, extra/missing selectors, wrong read/write tool and oversized
bodies are rejected before routing. Maximum body 1 MiB, response 8 MiB; use
focused queries/pagination if necessary.

Responses include `status`, `data`, and `retry_after` when applicable.
HTTP failures become MCP tool errors retaining status and structured error data.
No mutation is automatically retried. Destructive operations require
`confirm` equal to their exact operation name, for example `delete_project`.
Read the current object and preserve unrelated configuration before replacement.
This is not a substitute for the user's authorization.

## End-to-end workflow

1. `list_projects`; `create_project` only if needed.
2. `import_repositories` with `dry_run:true` to preview authorized repositories. Keep its returned `preview_digest` and pass it unchanged inside the approved import request. HTTP 409 requires another preview; never silently expand scope.
3. `start_ingestion` to import/sync/analyze/build, or `body:{}` for incremental
   refresh of registered repositories.
4. Poll `get_ingestion` until terminal; inspect errors and freshness for partial
   results. `get_live_status`, `list_repositories`, `ingestion_history` and
   `list_jobs` explain progress.
5. Use graph tools for source-backed queries. `list_graph_runs` and
   `compare_graphs` inspect saved versions.
6. Queue later work with `enqueue_refresh`; cancel/retry using `cancel_job` and
   `retry_job`. Direct ingestion uses `cancel_ingestion`/`resume_ingestion`.

GitHub credentials come from the server environment or approved GitHub CLI
account, not tool body arguments. Local registrations retain their paths and
are not automatically pulled; the host agent can update authorized local
checkouts using its normal Git tools. Managed clones are synced by DiffMind
and refuse dirty checkout overwrites.

## Local runtime and maintenance tools

`agent_runtime` accepts `status`, `start`, `stop`, `restart`, or `configure`.
Status includes the home, optional dashboard URL and effective settings.
Configure requires the complete settings object; read before editing:

```json
{
  "action": "configure",
  "settings": {
    "refresh_interval": "15m",
    "refresh_on_start": false,
    "refresh_concurrency": 4,
    "project_access": "legacy",
    "repository_workers": 4,
    "job_workers": 2,
    "queue_capacity": 256
  }
}
```

Use `project_access:"scoped"` before administering restricted project tokens;
the local owner retains full OS-trusted authority. Settings persist privately in
`agent-settings.json`; unknown/corrupt/public or
symlink configurations fail closed. Configuration restarts the backend; old
settings are restored if startup/publication fails. `0` or empty interval
disables scheduled refresh. The dashboard port may change after restart; read
status again. Stop does not erase data; subsequent management starts the backend.

One local controller holds a separate lifecycle lock even during maintenance.
Do not run two controllers for one home. Additional clients may use that
backend's HTTP MCP while it lives. Normal disconnect stops the child; an inherited
lifetime pipe also stops it if the controller is killed. Work/history remain
durable; reconnect starts recovery. This is not an always-on OS service: use
shared deployment for refresh independent of developer agent sessions.

`agent_command` invokes only the current installed binary with an argument
array. Allowed command families are `doctor`, `version`, `pack`, `backup`
and `storage`; no shell, arbitrary executable or environment overrides exist.
Commands are serialized with management, limited to five minutes and 1 MiB of
output, and report exit status, truncation and backend-restart outcome.

```json
{"args":["pack","init","/absolute/new/path/my-pack","--id","example.conventions"]}
```

The host agent edits synthetic fixtures and rules, then calls pack lint, test,
explain and install. Project-scoped pack CRUD is also available through
management tools: use the storage key returned by create/list as `pack_id`, and
preserve the manifest's `id` when updating (these can differ for dotted IDs).
See [pack authoring](knowledge-packs.md) for schema and limits.

```json
{"args":["backup","create","--offline","--output","/private/backups/new.tar.gz","--json"]}
```

The owned backend is stopped before maintenance and restored afterward, including
command failure/cancellation. Other external writers/old read-only MCP clients
may be using the same home. Stopping cancels active work and persists recovery
state; inspect running jobs before intentionally interrupting them. Those writers
may still block an exclusive lease: do not bypass their locks.
`backup restore`, `backup rotate` and `storage migrate` require `confirm`
equal to that two-word command. Existing offline flags, non-overwrite/path
guards and retention semantics still apply. For restore, use an absent target;
automatic workspace relocation is not implemented.
[Backup/recovery](backup-recovery.md) explains preservation and confidentiality.

## Security and tests

Remote mutation callbacks use the **current HTTP tool request's credentials**,
not cached initialization credentials. They re-enter the same authentication,
role/project checks, mutation guards, queue quotas and audit middleware as the
UI. Identity changes or token revocation cannot inherit a previous admin's
management rights. Tool annotations describe effects; they do not authorize them.

Viewer tokens remain query-only. Scoped editors can refresh assigned projects,
but cannot import, configure host paths, change packs, issue tokens or administer
access. Administrators can perform these operations. Local host tools are not
available on remote MCP; a deployment agent uses authorized host tooling for
installation, service lifecycle and offline maintenance.

`TestAgentAcceptance` launches the actual installed binary over stdio and performs
empty-workspace onboarding, three-repository graph extraction/query, incremental
reuse, pack testing/installation, failed-command recovery, backup, SQLite
migration, persistent scheduling, reconnect and controller-crash cleanup through
MCP. Permission, identity-switch, mutation-route parity, request bounds and setup
tests complement the existing company acceptance/race suites. Native release
gates run the same agent acceptance test against their installed archive.


## Repository scope and PR provider context

Import previews return candidate `source_type`, effective `default_branch`, local
`analysis_paths`, the project ID and requested scope alongside `preview_digest`.
Inspect these before approval. Local repositories are analyzed in place without
Git pull. Managed repository file configuration remains unknown until checkout;
the preview does not promise complete extraction or freeze future source edits.

An imported GitHub `api_base` persists as repository `git_api_base` and is reused
for PR listing, PR files/impact and live status. Explicit add/update operations can
also configure it. API URLs require HTTPS, except loopback HTTP for local
integrations. Credentials come from server environment or GitHub CLI using the
approved API hostname. Redirects to another origin are rejected. Changing
`git_url` clears a prior custom API endpoint unless the same update explicitly
approves its replacement.

PR lists distinguish local-only, missing remote, unsupported provider, unavailable
configuration and provider request failure from a successfully queried empty
list. `checked_count` counts successful provider queries; `repo_count` counts all
registered sources. Partial provider availability does not prove an absence of
PRs. Authentication, access and rate-limit failures give safe next steps without
reflecting provider response bodies. Public GitHub remains the automatic endpoint
for public GitHub repository URLs; existing custom sources without an approved
API endpoint require configuration rather than a guessed endpoint.

## Shared readiness before work

Call native MCP `get_readiness` (or management `inspect_workspace` with operation `get_readiness`) before setup, refresh or investigation. It works before a graph exists and separates the saved run and its input provenance from current work, checkout freshness, coverage limits and currently permitted actions. HTTP clients use `GET /api/v1/projects/{pid}/readiness`; browser workspace metadata includes the same contract. Query-only connections never advertise mutation authority. Follow `next_action`, recheck after acceptance or connection failures, and pin `saved_run_id` for evidence queries. See [the complete state table](shared-readiness.md).


## Installed-client and PR-head validation boundaries

The scoped company route was exercised through an installed Codex CLI app-server:
14 viewer read tools, project isolation, source-backed dependencies, replacement
credential verification and explicit old-token revocation. The client configuration
references a bearer-token environment variable; it does not contain the secret.
Browser login remains separate. See the exact candidate record in
[completion verification](research/2026-10-03/completion/README.md).

These were automated client operations without a model turn. They do not establish
unprompted agent adoption or independent human comprehension. Full-management local
stdio has a single lifecycle owner per home. Hosts that initialize several local
MCP processes for one home can encounter competing-controller startup failure;
use the existing service's authenticated HTTP connection for shared clients rather
than deleting locks or starting competing writers.

Exact PR callers require extracted origin and their own recorded clean analysis
revision matching the linked fact revision. This is a saved-source claim, not
runtime traffic. Dirty, older, declared and missing caller provenance stays outside
exact scoring even when the changed service matches the PR head. Explicitly local
sources stay in place when remote/provider metadata is supplied for PR retrieval;
that metadata does not authorize managed cloning, pulling or branch changes.

Use [current validation commands](validation.md) rather than historical extractor
evaluator instructions. The independently source-labeled public PR projections are
controlled provider fixtures over public source; they do not measure actual
upstream review outcomes, calibrated scores or universal framework coverage.
