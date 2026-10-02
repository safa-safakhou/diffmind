# Integrated implementation batch 1

This batch implements connected improvements across 15 planned tasks. Six tasks
meet their recorded local acceptance criteria; nine have partial implementation
and remain in progress. This is an uncommitted implementation, not a published
release or completion of the whole Linear project.

## What the connected path now does

A first-time browser user sees an explicit project starting point rather than a
mandatory modal. Failed project loading stays a failure with retry; it does not
become an empty workspace. Creating a project works with labeled fields and
Enter. Shared dialogs trap keyboard focus, support Escape, restore the trigger's
focus, and focus Cancel before a destructive confirmation.

Import starts with a non-mutating preview showing the backend's candidate names,
paths and count. Closing/reopening retains the draft and preview within the
workspace session. Changing inputs invalidates continuation. Failed previews
remove the previous approval and display an error inside the active dialog.
Continuation imports the configured scope and exposes ingestion progress and its
terminal result. The dialog explicitly says that the backend scans again at
commit time: this is not a frozen candidate-list approval.

The dashboard has one primary **Update context** action. **Reload view** fetches
data without starting analysis. Configuration/run/build controls are grouped
under **Advanced actions** and hidden from roles that cannot use them. Company
viewers receive a read-only explanation, while editors enqueue the saved
configuration without import authority. Access-loss guidance is generic.

Workers and analyzer cache identity use the executable provided by the running
CLI, with an explicit `DIFFMIND_BINARY` override still taking precedence. Launching
the UI by absolute path works even when the binary directory is absent from PATH.

New agent homes refresh registered repositories at connection/startup and every
15 minutes while connected. Explicit existing settings, including manual mode,
are preserved. Local working trees are analyzed in place. This does not discover
new repositories or promise service availability after the owner disconnects.
MCP initialization now explains project selection, read-only queries, pagination,
static evidence and limits on interpreting missing findings.

Configured repository-relative file globs now govern AST source/config input,
OpenAPI enrichment, pack extraction/detectors and repository metrics. Test and
fixture exclusions use the shared source policy. Invalid globs fail validation;
excluded alternate OpenAPI files cannot enrich included routes. Effective
inherited service configuration contributes to cache identity. Saved protocol
metadata records the explicit file scope and filter version.

UI and MCP summary responses expose a shared evidence description: saved graph,
repository-analysis freshness, static basis and unverified coverage. They identify
the freshness source explicitly: the dashboard checks current checkouts, while
MCP uses stored repository status. These are not yet one live readiness service;
fresh analysis does not prove that the selected graph includes it.

The overview starts at a minimum 0.75 scale, producing 12 CSS-pixel service labels
from 16-pixel text in the six-service desktop trial. Services, resources and
relationships support keyboard focus and Enter/Space inspection; focus on an
offscreen item pans it into view. Evidence lists now expose continuation beyond
80 objects, 120 traces and 40 resource facts. Snapshot comparison also shows the
existing request-contract compatibility results, with before/after evidence and
explicit limitations when extracted fields are absent.

## Task disposition

| Task | Linear | State after batch | Result / remaining work |
| --- | --- | --- | --- |
| T01 | MNI-175 | In Progress | Shared evidence dimensions implemented; combined job/runtime/action state table and equivalent live freshness remain. |
| T02 | MNI-176 | Done | Primary update, separate reload, advanced controls and backend-verified role routing. |
| T03 | MNI-177 | In Progress | Preview review and invalidation implemented; backend candidate freezing and cross-surface scope contract remain. |
| T04 | MNI-178 | In Progress | Import errors propagate into the dialog; other mutations and field-associated recovery taxonomy remain. |
| T05 | MNI-179 | Done | CLI executable used for workers and identity; PATH-absent live ingestion and existing doctor/bootstrap tests pass; platform limits documented. |
| T06 | MNI-180 | In Progress | MCP onboarding/evidence instructions added; fuller terminal-state guidance and actual installed-host adoption evaluation remain. |
| T07 | MNI-181 | In Progress | Default maintenance and preserved settings validated; overdue-only startup and repeated-failure policy still need work. |
| T09 | MNI-183 | Done | Shared dialog focus, labels, Escape, cancellation and destructive-confirmation focus verified in Chromium. |
| T10 | MNI-187 | Done | Six-candidate preview, retained draft, invalid root, non-mutation and completed graph verified live. |
| T11 | MNI-184 | In Progress | Readable overview, keyboard inspection and evidence continuation implemented; larger graphs and broader viewport validation remain. |
| T12 | MNI-185 | In Progress | Failed loading, explicit first use, role controls and generic access loss improved; provider/PR-state work remains. |
| T13 | MNI-186 | Done | Scope applied across extraction/enrichment/packs/metrics; inherited-config invalidation, saved scope and public-repo regression verified. |
| T14 | MNI-188 | In Progress | Summary evidence basis and coverage limits exposed; relationship-level evidence classes/correction paths remain. |
| T17 | MNI-191 | Done | Same four field changes in MCP/dashboard alongside five graph fact changes; missing/unsupported fields never imply full compatibility. |
| T25 | MNI-200 | In Progress | Entry/playbook/operation/file-scope docs updated; all historical evaluation/release guidance and platform checks remain. |

Other planned tasks retain their previous state. Dependency edges and epic scope
are retained. A completed child task does not close a prerequisite's broader scope
or the parent epic.

## Verification of the exact candidate

The candidate is identified by [manifest.json](manifest.json), which records the
source HEAD, all 45 changed/deleted/new implementation input files, their content
hashes, the executable hash and evidence hashes. The source is based on
`5d2548514bfc5d920dfd40def5fcacd80162e9ae` with local modifications. The binary's
SHA-256 is `598187fe5bca29403a32cfe178694c8d8f16b7595c253cec5f17cae4884bed9c`.

- `go test ./...`: 65 packages passed, including extractor, workspace, MCP,
  command/agent acceptance and setup scripts. Other packages have no tests.
- `go test -race` for UI, agent host, query and knowledge packages: passed.
- Frontend: 16 library tests and 16 component tests passed; production Vite build
  passed and its assets are included in the embedded Go application.
- New regressions cover glob precedence/safety, source/config scope, OpenAPI scope,
  pack extraction/detection scope, metrics, inherited-policy cache invalidation,
  executable selection, preserved explicit manual settings, unknown coverage,
  onboarding/import errors, evidence continuation and contract comparison.
- `git diff --check`: passed.

The live trial executed the built binary and real Chromium/MCP interactions:

| Live check | Observed result | Evidence |
| --- | --- | --- |
| First use and dialogs | No forced modal; labeled Name; Enter creates; Tab remains inside; Escape restores focus; Cancel is the initial destructive-confirmation action. | [Browser events](evidence/browser-regression.json) |
| Preview/recovery | Invalid directory retained; invalid regex visible in dialog; six named candidates; no registered repos after preview; reopened draft retained; changed filter disables continuation. | [Preview](evidence/preview.png), browser events |
| Worker startup | Six repositories analyzed and graph completed with process PATH `/usr/bin:/bin`, excluding the candidate binary directory. | [Browser result](evidence/browser-result.json) |
| Readability/keyboard | Initial scale 0.75 and 16-pixel text = 12 rendered pixels; keyboard opens service inspection. | [Graph](evidence/graph.png), browser events |
| Agent maintenance | New home reports 15m/startup policy; a shortened disposable 2s policy updates the graph without manual ingestion; local edits remain intact. | [Agent events](evidence/agent-regression.json) |
| Contracts | Four extracted field changes match dashboard/MCP; separate graph comparison reports five changed facts. | [Comparison](evidence/four-contract-changes.png), [agent result](evidence/agent-result.json) |
| Company roles | Viewer refresh denied 403; editor queue accepted 202; editor import denied 403; membership loss removes the open graph. | [Company events](evidence/company-regression.json), [viewer](evidence/viewer-permissions.png), [editor](evidence/editor-permissions.png) |
| Public-source scope | Default: 100 exposures, 25 dependencies, 38 references to examples. Explicit `examples/**` exclusion: 95 exposures, 16 dependencies, zero example references; generated metadata records that scope. | [Scope result](evidence/scope-result.json) |
| Lifecycle | Disconnect stops the owned backend; reconnect preserves explicit manual settings. | Agent events |

No browser JavaScript exceptions occurred in the final trials. Console resource
errors were expected responses to intentionally invalid input, denied actions and
revoked access. Failed early harness assertions were corrected (response envelope,
`completed` ingestion status, hash routing, fleet-job trigger and selector
cardinality); they are not presented as product failures or hidden proof of success.
The browser trial did uncover a real delayed-focus-restoration issue, fixed by
using a layout effect before interaction.

## Evidence limits and reproduction

This is an automated usability regression, not a study with recruited people.
Only Linux/WSL AMD64 and Chromium at 1440×1000 were exercised live. The 15-minute
production default was inspected; timed maintenance was observed using a 2-second
policy in a disposable home. No company's credentials, personal MCP configuration,
provider PR comments or production deployments were changed. PR matching/risk
heuristics, offboarding tokens, lifecycle contention and corpus accuracy are not
claimed fixed by this batch.

The [retained harness sources](harness/) are exact copies of the local trial
scripts with `.txt` suffixes. They refer to the prior study's disposable fixture
paths and this batch's `/tmp` root; they are evidence, not a portable CI harness.
For independent reruns, build the frontend then Go binary, generate disposable
Git repositories from `examples/demo-shop/repositories`, and adapt the paths and
Playwright module location. Restore the checkout baseline before applying
`scripts/apply-demo-shop-change.sh`. For the public-scope check, use an isolated
clone of the recorded source revision and compare default analysis with the
explicit exclusion. Preserve old runs when comparing their metadata.

Logs are under [evidence](evidence/); the manifest makes their retained contents
checkable. Future independent-human, actual-host and release/platform verification
remain separate project gates, particularly T27 and T28.
