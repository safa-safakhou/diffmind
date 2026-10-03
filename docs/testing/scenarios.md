# Scenarios and expected results

Use frozen artifacts, private test homes and public/synthetic code. Preserve branches/dirty edits. Each case records actual actions, expected/actual behavior, candidate identity and evidence in result-template.json. Keep facilitator expectations out of participant cards.

## Individual developer and installed agent

### A01 - First activation
Setup: empty home, installed candidate, actual model host, labeled public repo.
Action: register documented full-management connection; give U01 without prescribed tools.
Expected: discovery, approved scope, saved usable graph, source-backed answer; accepted work distinguished from finished analysis.
Evidence: redacted registration, host/settings, model turns/tool calls, project/job/run IDs, citations, first-useful-result time, assistance.

### A02 - Connection capabilities
Setup: query-only stdio MCP, management agent and company HTTP MCP in separate sessions.
Action: dependency question then import/refresh request.
Expected: actual capabilities/role honored; denied management actionable; no invented success/global-token fallback.
Evidence: discovery/calls/results/final answer per mode.

### A03 - Context selection
Setup: two similarly named accessible projects with different dependencies; hidden project.
Action: named-project question then ambiguous question.
Expected: correct context; ambiguity clarified or justified choice explicit; hidden evidence never leaks.
Evidence: selected IDs/citations, wrong-context answers, missed/unnecessary calls.

### A04 - Return and background work
Setup: saved graph, approved cadence, clean repo.
Action: reconnect, edit one repo, wait cadence/query; pause schedule/repeat edit.
Expected: state survives, changed input refreshes/unchanged reuses, edits preserved, revision/staleness clear, pause stops scheduling.
Evidence: HEAD/diff hashes, cadence/times, registrations, run IDs/reuse counts.

### A05 - Failure/reconnect/owner contention
Action: controlled failure, correction/retry, disconnect/reconnect, second owner in disposable home.
Expected: durable distinct job states/history, no competing writes or silent unauthenticated forwarding; separate-home/scoped HTTP handoff works.
Evidence: error/recovery, lease/queue/history/process cleanup, subsequent successful query.

### A06 - Unprompted ordinary work
Setup: actual installed host in fresh and returning sessions; known dependency answer.
Action: ask 'What other services could this change affect, and what evidence supports that?' without naming DiffMind/tools; also local-only task.
Expected: appropriate context, source-backed answer/uncertainty. Actual model turns required.
Evidence: conversation/tool sequence/settings/expected vs actual answer/missed and unnecessary calls. Scripted discovery and brief sessions prove neither adoption nor retention.

## Company entry

### B01 - Role parity
Setup: actual identity proxy/scoped mode; viewer/editor/ungranted/admin; two differently granted projects.
Action: browser/HTTP MCP read, refresh, import/configure/manage access; guessed inaccessible URL.
Expected: viewer reads only; editor queues refresh but cannot import/configure/grant; ungranted no project/same 404 as missing; admin recovery; interface parity.
Evidence: actual proxy subject/capabilities, redacted responses/browser actions/filtered totals. Synthetic headers alone do not validate production proxy.

### B02 - Offboarding
Setup: membership plus independently issued token; unrelated service grant.
Action: revoke membership/test browser-token; revoke token/test existing/new connections.
Expected: independent lifecycle clear; both removals deny access, unrelated service/admin survives.
Evidence: operations/timestamps/statuses and participant support-owner explanation.

### B03 - Scope approval
Action: preview, edit filters, add discovery candidate before submission; approve valid scope/cadence.
Expected: preview registers nothing; edits invalidate review; stale digest rejected before registration; whole-project exposure/maintenance scope understood.
Evidence: digests/registrations/rejection/cadence and participant words.

## Independent graph

### C01 - Supported extracted positive and negative
Setup: pinned supported public exposures/callers, blind labels; declarations/unknowns separate.
Action: analyze; compare browser/MCP graph evidence.
Expected: correct direction/protocol/operation/location/revision/class, near-match negatives not exact.
Evidence: frozen reviewer receipts/pins/graph/facts/denominators/mismatches.

### C02 - Enclosing repo/subdirectory
Setup: same public monorepo at root and explicit subdirectory in separate projects; duplicate basenames.
Action: compare source locations/PR paths; change unrelated same-basename file.
Expected: correct identity; no basename/suffix guessing to exact; unsupported mapping candidate/unknown.
Evidence: roots/upstream/derived revisions/locations/filenames/tiers.

### C03 - Reviewed correction/rollback
Setup: pinned Boutique/evaluation pack/frozen declared labels.
Action: baseline/gap review/selected pack ingest/unchanged repeat/detach-delete pack/rollback.
Expected: declared provenance, negative controls absent; key/manifest selection works; missing/ambiguous selection visible; baseline and corrected history preserved.
Evidence: graph hashes/assertions/review states/HTTP-MCP provenance/reuse. Twelve declared edges are subset, not total accuracy.

## PR cases

Use independent labels and browser/API/HTTP MCP parity. Pin base/head/merge-base/file-list completeness. Controlled public mutations cover mechanisms; actual upstream PR cases are additionally required, and cannot be replaced by local mutations.

### R01 - Direct addition/change
Action: analyze supported endpoint and clean caller at changed head; inspect PR.
Expected: exact only with extracted changed-head lines and clean matching caller-own revision; unrelated caller candidate.
Evidence: diff/head/caller SHA/source lines/eligible set/parity.

### R02 - Revision controls
Action: default-branch, dirty/missing/mixed endpoint revisions; old/missing/dirty caller facts.
Expected: none exact; actionable limitations; same-branch refresh not promised to capture head; no query clone/branch switch.
Evidence: separate freshness/eligibility/next-step per subcase.

### R03 - Removal
Action: remove baseline endpoint/file, analyze head/inspect deletion.
Expected: missing entity not no-impact proof; removed callers/missing baseline limits visible; deleted old coordinates cannot match unchanged head.
Evidence: baseline/head/diff/facts/limits and participant interpretation.

### R04 - Internal/transitive
Action: helper change below unchanged route with known dependency.
Expected: supported paths distinct from file/service candidates; unsupported transitive conclusion not exact.
Evidence: independent source path/diff/connections/tiers.

### R05 - Config and contract
Action: address/config edit; separate contract field removal/type change.
Expected: declarations not runtime proof; compatibility separate from graph/heuristics; absent contract evidence unknown.
Evidence: before/after config/contracts/labels/provenance/compatibility.

### R06 - No-impact/identity controls
Action: docs/test/unrelated edits; different directory with same basename, method/case/slash/protocol; colliding IDs; unknown kind.
Expected: no false exact; HTTP fallback not RPC/unknown, IDs not cross-protocol, unknown kinds ineligible. No exact callers is not safe-to-merge.
Evidence: every subcase/heuristic. Existing eight-case regression is only a subset.

### R07 - Incomplete provider evidence
Action: missing patch/lines and controlled truncated file list.
Expected: completeness visible; file-scope not changed-line proof; score not calibrated probability/runtime safety/merge permission.
Evidence: flags/coordinates/output ordering.

### R08 - Human interpretation/delivery
Action: participants explain one exact and one uncertain PR result and report destination.
Expected: evidence/freshness/limits discoverable before score reliance; heuristic distinct from proof; on-demand read distinct from automatic PR posting.
Evidence: first viewed information/words/errors/prompts. No posting claim without actual configured integration.

## Provider

### P01 - Public actual PR
Action: read pinned public PR.
Expected: host/repo/revision limits and parity correct; no mutation/comment.
Evidence: endpoint/PR/redacted responses. Not Enterprise validation.

### P02 - Enterprise
Setup: authorized real endpoint/version/known PR/matching-host read credential/SSO.
Action: existing provider/API-base configuration; browser/MCP inspect.
Expected: configured host/credential used/evidence retrieved.
Evidence: host/version/PR/status/redacted destination and auth-presence, never values.

### P03 - Credential/access failures
Action: no/invalid credential, missing repo grant, inaccessible/missing PR; correct access.
Expected: actionable provider states where distinguishable; no fictitious no-impact/secret leak.
Evidence: per-subcase status/interface/recovery.

### P04 - Host/redirect boundary
Setup: controlled endpoints/synthetic credential; P02 remains required.
Action: cross-host redirect/similar host/rejected API-base.
Expected: no credential to unmatched destination; visible failure/no public fallback.
Evidence: destination assertions. Never real token to external echo service.

### P05 - Local remains local
Action: dirty local repo with remote PR metadata; inspect/sync/refresh; missing path control.
Expected: no switch/overwrite/clone/substitution; path failure actionable.
Evidence: branch/HEAD/diff hashes/source kind/path before-after.

## Operations/final candidate

### O01 - Queue/history
Action: JSON/SQLite each queue/cancel/retry, interrupt owned process/restart, unchanged rerun.
Expected: durable distinct recovery/history, immutable graph, reuse/no duplicate writers.
Evidence: backend/jobs/runs/hashes/lease/reuse.

### O02 - Upgrade/recovery/grants
Action: prior snapshot, revoke test grants, stop/preserve rollback; exact-path current restore with isolated ingress; independent offboarding reconciliation on both backends.
Expected: reject damage/existing-home overwrite/live-writer backup; preserve graph; restored old grants removed/unrelated service-admin retained.
Evidence: versions/hashes/verify/rollback/per-identity receipt. Small timings not production RTO.

### O03 - Native matrix
Action: same product revision native Linux/macOS AMD64/ARM64 archive gate.
Expected: version/installer/doctor/UI/SQLite/company-agent acceptance pass; compilation alone insufficient.
Evidence: source/binary/archive hashes/platform/logs. WSL Linux is not Windows release.

### O04 - Final reconciliation
Action: map scenarios/platforms and 42 audit/11 live findings to verified resolution/tested boundary.
Expected: no critical activation/access/evidence defects, required gates complete/external gaps explicit; fix triggers affected reruns.
Evidence: candidate/issue/finding matrix/blocked checks. Publishing separate.

U01-U03 are defined in usability-study.md; three observations per original persona minimum, including fresh/returning, supplement operator tests.
