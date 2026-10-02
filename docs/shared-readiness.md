# Shared workspace readiness

`GET /api/v1/projects/{pid}/readiness`, `workspace.readiness`, and the MCP `get_readiness` tool use the same read-only projection. Readiness works before any graph exists. The management catalog also exposes `get_readiness` through `inspect_workspace`.

The response separates these dimensions:

- `runtime`: the query process is available. A client can observe `disconnected`; the server cannot report its own absence.
- `maintenance` and `connection_mode`: managed connections observe a running backend; query-only stdio connections report maintenance as unknown and cannot refresh or configure. Availability does not promise an enabled background schedule.
- `access`: permitted at observation time. HTTP continues returning ordinary 401/403/404 errors for denied or unavailable scopes; it does not reveal project metadata. Browser and MCP error projections pause actions and provide recovery guidance.
- `graph`: queryable, missing, unavailable or unknown. Queryable requires a completed run and a readable, non-null, matching graph artifact. A completed manifest alone is insufficient. Missing artifacts can fall back to an older readable completed snapshot. Malformed or inaccessible artifacts report unavailable rather than an empty workspace.
- `saved_run_id`, `saved_at` and `inputs`: immutable saved provenance. An absent completion timestamp is unknown, not the current time. Inputs bind repository IDs to analysis IDs used for that graph.
- `sources`: last stored analysis/head observations and current analysis freshness. Stored revisions explicitly carry `revision_basis=stored_repository_status`; reads do not fetch remote heads. They are distinct from graph inputs.
- `evidence`: the shared checkout freshness calculation compares the latest repository analysis. It does not prove that the saved graph incorporates that analysis. Coverage remains unverified and evidence remains static source evidence.
- `pr_head_eligibility`: unknown. Use the existing PR investigation surface to assess a specific revision; neither work completion nor checkout freshness establishes PR-head eligibility.
- `work`: observed kind, ID, phase, status and timestamp. Active execution takes priority over queued work and terminal history; timestamps break ties. Job success normalizes to completed but preserves a correlated partial ingestion.
- `actions` and `next_action`: currently permitted query/refresh/configure/inspect operations and the recommended next step. The server rechecks authorization, quotas and work admission on every mutation. This response is not an authorization token or a guarantee of future acceptance.

## Cross-interface state table

| Observation | Saved graph | Work | Browser | MCP | Next step |
| --- | --- | --- | --- | --- | --- |
| New project | Missing | Empty | No saved graph; setup controls only for a permitted configurator | Same state without requiring graph summary success | Review import, or request an editor |
| Accepted queued refresh | Retained if present | Queued | Disable duplicate updates/configuration; retained graph remains usable | Same state; never infer readiness from 202 | Inspect work |
| Running import, analysis or build | Retained if present | Running and phase | Activity and saved snapshot coexist; disable conflicting actions | Same phase and saved provenance | Inspect work |
| Completed work | Queryable only if artifact validates | Completed | Show saved ID and timestamp, explicit unverified coverage | Same provenance and limits | Query saved graph, or refresh if no graph exists |
| Partial refresh | Retained or newly built from available artifacts | Partial | Keep graph visible and explain partial results | Same distinction; success of an envelope cannot erase partial results | Inspect work before retry |
| Failed refresh | Older graph stays queryable if intact | Failed | Keep graph visible; explain failure independently | Same saved graph and failed work | Inspect work before retry |
| Cancelled/interrupted work | Retained if intact | Cancelled/interrupted | Explain terminal outcome without claiming an empty project | Same persisted outcome | Inspect work before retry |
| Damaged saved artifact | Unavailable | Independent | Do not present its completed manifest as a usable snapshot | Same graph limitation | Inspect run; restore or rebuild |
| Temporary connection/access-check failure | Last loaded view retained; current query authority paused | Last observation only | Explain connection failure and pause actions until successful readback | Managed tool returns `IsError` with disconnected/unknown readiness and no authority | Reconnect/recheck |
| Confirmed access loss | Unknown to caller; hide workspace | Unknown to caller | Remove graph view and give generic access guidance | `IsError` with denied readiness; omit identifiers, provenance and actions | Request access or select a visible project |
| Query-only stdio | Queryable if readable locally | Persisted observation, not a worker liveness guarantee | Not applicable | Mutation actions always false; maintenance unknown | Query saved evidence or request an editor |

Readiness reads do not register/import repositories, mutate snapshots, sync source trees, fetch PRs, retry jobs or extend grants. Browser requests from a previous project or an older poll cannot replace a newer workspace. HTTP artifact validation reuses the existing file-size/mtime cache so repeated polls do not repeatedly parse large graph files.

## Verification

See [the 3 October validation record](research/2026-10-03/shared-readiness/README.md). This contract completes the shared-readiness foundation (MNI-175); it does not complete the background contention, enterprise deployment, independent corpus, installed-agent adoption or release gates tracked separately.
