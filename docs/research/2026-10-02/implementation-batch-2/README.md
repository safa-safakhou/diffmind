# Implementation batch 2 — reviewed scope and dependable continuity

This batch connects import approval, feedback, source freshness, background work and workspace access. It advances T01, T03, T04, T06, T07 and T12 together, without marking their broader coordinated acceptance criteria complete.

## Committed changes

| Commit | Outcome |
| --- | --- |
| `c4bb22f` | Checkpoint of the preceding research and implementation batch 1. |
| `0816d70` | Server-checked import review, captured ingestion candidates, field-associated import errors and matching agent workflow. |
| `a5a3882` | Shared read-only freshness, overdue startup, persisted failure backoff, capability failure continuity and MCP guidance. |
| `998aeb3` | Eligibility follows completion order when a manually retried older job finishes after newer jobs. |

The compiled candidate and retained evidence are identified by [manifest.json](manifest.json). Documentation and Linear records are committed separately after verification. No changes were pushed or deployed.

## Connected behavior

**Review → approval → ingestion.** Both browser and agent preview requests return `preview_digest`. The browser carries it into approval, and the MCP catalog/playbook tells agents to do the same. The server compares the project, candidate registrations, import settings, effective GitHub branch and local analysis include/exclude boundaries. Changed scope returns HTTP 409 before registration or ingestion acceptance. The browser retains the draft and requires a new preview. After acceptance, ingestion uses the captured candidate list instead of discovering additional repositories in its worker.

The digest checks consistency and does not grant authority. Existing direct clients that omit it remain compatible; mandatory review enforcement is still open. Preview does not freeze repository contents. Local file-analysis boundaries are checked during approval but are not shown as an expanded review table or pinned against later configuration edits. Managed provider behavior has automated fixtures; this batch did not exercise a live company GitHub import.

**Feedback → recovery.** Invalid directory and regex errors remain inside the dialog and associate with the affected field through `aria-invalid` and `aria-describedby`. A review conflict clears approval rather than silently retrying. Accepted ingestion still reports progress independently of a completed graph. Other modal mutations and a complete provider/auth recovery taxonomy remain open.

**Current source → saved evidence.** UI and MCP use the same read-only checkout comparison against the latest repository analysis. An inaccessible checkout or unavailable analysis artifacts produces unknown freshness. Summary metadata explicitly says `freshness_reference: latest_repository_analysis`; a fresh repository analysis does not prove that the selected graph incorporates it. Saved snapshots remain queryable and unchanged by these reads. Managed repository status uses the last observed remote/head revision; this read does not fetch the remote or establish live deployment state.

**Connection → bounded maintenance.** Successful jobs and completed manual ingestions persist the last maintenance completion. Startup/scheduled work skips a project until it becomes overdue, and status supplies the reason and next eligible time. Jobs retain their existing three-attempt budget. After exhausted failures, automatic jobs wait twice the greater of the interval and one minute, doubling consecutive failures up to six hours. Explicit manual refresh bypasses the delay. A successful retry of an older job resets eligibility according to completion order. Local agent disconnect still stops its backend. The 15-minute default is a policy choice, not a measured optimum.

**Temporary failure → access loss.** A transient capability check failure retains the last loaded workspace and pauses its actions. A confirmed 401/403/404 removes the workspace view and gives generic return/request-access guidance. Successful capability checks restore normal behavior. Existing server authorization remains authoritative.

## Verification and live observations

- Full `go test ./...` passed; retained package output identifies the tested scope.
- Frontend suite passed: 16 library and 17 component tests. Production Vite output is embedded in the candidate.
- Go race checks passed for UI, query, agent host and MCP server. Focused cases cover approval conflicts, GitHub branch changes, changed file boundaries, persisted cadence, manual override, failure backoff, retry ordering and inaccessible analysis artifacts.
- Real Chromium invalid-directory submission associated the error with its field. Adding a repository after preview caused a 409, retained the directory draft, disabled approval and registered zero repositories. A new preview approved exactly two repositories.
- Real MCP stdio initialized with readiness/review instructions, carried the digest into ingestion, observed 202 acceptance and waited for a completed graph. A local edit produced identical UI/MCP evidence with one dirty repository while preserving the saved run.
- Browser capability responses were controlled at the network boundary: 503 retained the workspace; 403 removed it. These responses were injected for the observer, not changes to actual company membership.
- Disconnect/reconnect retained one completed refresh job and reported `not_due`. A disposable server using a failing analyzer exhausted three attempts, retained one job during subsequent shortened scheduled checks and reported `failure_backoff`; manual refresh created a second job.

Expected 400/409/503/403 browser resource errors are deliberate negative cases, not JavaScript exceptions. Each harness stops its owned browser/backend processes. Fixtures and DiffMind homes are disposable; personal agent settings, company workspaces and source repositories are not modified.

## Task disposition and remaining evidence

| Task | Progress in this batch | Still open |
| --- | --- | --- |
| T01 / MNI-175 | Equivalent live freshness, explicit reference and unknown/fallback semantics. | Combined runtime/job/action state contract and full partial/disconnected matrix. |
| T03 / MNI-177 | Server review check across UI/MCP; captured candidate acceptance. | Expanded branch/file-scope review, legacy-client enforcement and live managed-provider negatives. |
| T04 / MNI-178 | Field association, retained drafts and repreview recovery. | Other mutations, provider/auth errors and common recovery taxonomy. |
| T06 / MNI-180 | Empty/sole/multiple/query-only guidance, approved scope and acceptance versus readiness. | Shared-state dependencies and actual unprompted installed-host adoption under T27. |
| T07 / MNI-181 | Overdue reconnect, failure backoff, manual override and retry completion ordering. | Managed dirty-source, busy/two-agent/resource-limit matrix on the final release candidate. |
| T12 / MNI-185 | Transient access continuity and confirmed denial. | Provider-not-ready versus zero PR and broader partial/loading journeys. |

[Linear readback](linear-status.json) confirms all six remain In Progress. Earlier historical descriptions refer to then-uncommitted work; the new updates identify the commits above.

This is an automated observer study with real runtime actions on WSL/Linux amd64 and headless Chromium. It does not establish recruited-human usability, unprompted host adoption, native platform distribution readiness, large-company query latency, comprehensive dependency coverage or release readiness. The project still requires T26/T27/T28 corpus, adoption and candidate gates.

## Reproduction

Retained [harness](harness/) files are exact observers with `.txt` suffixes to avoid accidental execution/discovery. Copy them to a fresh disposable directory and remove the suffix, place a freshly built `diffmind` binary beside them, and supply the original disposable Demo Shop fixture identified in the manifest. `helpers.mjs` uses the installed Playwright runtime path from this session; adapt that path for another machine. Run `review-live.mjs` followed by `continuity-live.mjs` with Node. Both fail on a mismatched assertion and retain timestamped events/screenshots under their own `evidence` directory.

Logs expose only synthetic fixtures, project names and disposable paths. They contain no company credentials. Timing uses a 200ms schedule only in the isolated failure trial; production defaults remain unchanged.
