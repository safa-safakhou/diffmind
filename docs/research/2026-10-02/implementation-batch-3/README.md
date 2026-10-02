# Implementation batch 3 — reviewed sources and provider continuity

This batch advances T03, T04, T12 and T19 together. It connects approved import settings to later provider queries and makes unavailable evidence distinct from a successfully empty result.

## Committed implementation

- `f3f059b`: explicit import review metadata, persisted GitHub API configuration, provider states and active analyzer-form feedback.
- `f34b8d8`: source/API approval continuity, matching-host credentials, PR response isolation and changed-head refresh.
- `b513293`: consistent unknown PR counts in the main and secondary UI summaries.

The exact compiled candidate, changed inputs and retained observations are identified in [manifest.json](manifest.json). [Linear readback](linear-status.json) confirms the four tasks remain In Progress with their remaining criteria recorded. No changes were pushed or deployed.

## Behavior

**Import review.** Preview now returns project/request scope, source type, effective Git branch and local `analysis_paths`. The browser supports a branch override, shows additional include/exclude rules and explains the default source safety boundaries. MCP receives the same metadata. Managed file configuration is explicitly unknown before checkout. The existing digest detects changed approval scope; preview does not freeze source content or establish complete extraction. Legacy clients can still submit direct imports without a prior digest.

**Provider continuity.** GitHub import `api_base` persists as `git_api_base`. Import, PR listing, PR files/impact and live status use the approved endpoint. Repository add/update also accepts the field. Changing the source URL clears the old custom endpoint unless the update explicitly supplies its replacement. Repeating the same URL preserves the approval. Custom sources without approved API metadata require configuration rather than an inferred enterprise endpoint.

**Credentials and errors.** API endpoints require HTTPS, with loopback HTTP allowed for local integrations. GitHub CLI credentials are selected by the approved API hostname, including hosts whose names do not contain “github”. Redirects to a different origin are rejected before contacting it. Exact repository host parsing rejects lookalike and substring-based public GitHub matches. Authentication, access, rate-limit and transport failures give safe next steps without reflecting provider bodies or redirect URLs. These changes reuse the existing credential sources; no credential arguments or secret-bearing examples were introduced.

**Provider availability.** Local-only, missing remote, unsupported provider, unavailable configuration and failed requests remain distinct from a successful query with zero PRs. The open-PR filter keeps unavailable repositories visible. `checked_count` counts successful queries; `repo_count` includes all registered sources. Unknown-only scopes show an unknown observed PR count. Mixed availability does not prove no PRs exist.

**Async views and form recovery.** Late provider responses cannot replace another workspace after navigation. Impact reloads when the selected PR head or update time changes. Batch and individual analyzer failures reach their active dialogs and preserve drafts; the UI does not silently retry rejected mutations.

## Verification

- Full Go suite passed: 65 packages with tests, plus packages without tests.
- Frontend suite passed: 16 library and 19 component tests; production Vite output was embedded in the candidate.
- Race checks passed for UI, store and agent API. Provider fixtures cover exact hosts, approved custom endpoint persistence, PR/live-status queries, cross-origin redirect rejection, safe auth feedback, API URL validation, source URL changes and credential hostname selection.
- Real Chromium showed local path exclusions, retained a batch draft after an injected 409 and an individual analyzer draft after an injected 403, kept local-only provider context visible, showed a distinct successfully empty result, displayed a fixture PR with changed-file impact, and showed safe authentication failure feedback.
- Real stdio MCP returned the same local path rules and managed branch review metadata as the browser/backend preview.
- Real public GitHub reads queried `safa-safakhou/diffmind`, observed five open PRs at the recorded time and loaded PR 9 impact successfully. No comments, reviews, checks or repository contents were changed. Counts are observations at that time, not permanent facts.

The form 409/403 responses were injected at the browser network boundary to verify error handling; they do not establish new authorization-enforcement evidence. Custom-provider API behavior uses a real local HTTP fixture server; the matching-host CLI test uses a controlled stub rather than company credentials. Cross-origin tests assert that the destination receives zero requests.

Final live logs contain no JavaScript exceptions or assertion failures. Expected 401/403/409 browser resource errors are deliberate cases. An initial observer used the wrong batch route and was corrected before the retained passing run; that timeout is not a product finding.

## Remaining work

| Task | Remaining acceptance / evidence |
| --- | --- |
| T03 / MNI-177 | Legacy-client review enforcement policy; inspection/pinning policy for managed file-analysis boundaries. |
| T04 / MNI-178 | Complete recovery taxonomy across operations, other confirmation/detail errors and full failure matrix. |
| T12 / MNI-185 | Broader partial/loading state matrix and shared-readiness dependencies. |
| T19 / MNI-192 | Target enterprise deployment and actual credential trial. |

The shared readiness contract, broader background contention matrix, company access continuity, provenance, real-company corpus, installed-host adoption and release gates remain in the existing backlog. This batch does not establish independent-human usability, unprompted agent adoption, native-platform readiness, enterprise deployment support, complete dependency coverage or release readiness.

## Evidence and reproduction

The [evidence directory](evidence/) retains final timestamped events, six screenshots, server/stdout logs and test/build output. [Harness files](harness/) have `.txt` suffixes to avoid accidental execution. Copy them into a fresh disposable directory, remove the suffixes, place the compiled `diffmind` binary beside them and run `provider-live.mjs` with Node. The helper's Playwright runtime path and the original disposable Demo Shop fixture path are recorded in the manifest and need adapting on another machine.

The observer creates private disposable homes, copies the local fixture and owns its HTTP fixture, browser and backend processes. It stops those processes on exit. Personal agent settings and existing company workspaces are untouched. Public repository registration retrieves PR metadata without cloning or analyzing that repository. No tokens are retained in these logs.
