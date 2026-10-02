# DiffMind proposed work items

These are proposed solutions derived from the October research. T/E keys are stable local references. Their Linear counterparts are recorded in [linear-backlog.md](linear-backlog.md). No task is completed merely because this plan or its issue exists.

## T01 — Define shared readiness freshness and capability states

**Epic:** E1. **Milestone:** M1. **Priority:** High. **Initial state:** Todo. **Depends on:** None. **Evidence:** F11, F12, F13. **Stages:** A06, A07, A09, B05, B09, C04.

Define a single product contract consumed by UI and MCP: runtime available, graph queryable, current job phase, source revision/freshness, known coverage limits and permitted actions are separate dimensions. Reuse existing session, capabilities, ingestion and graph data; completion must never imply comprehensive architecture.

### Acceptance criteria

- Specify empty, queued, running, completed, partial, failed, disconnected and access-denied cases with the saved run and next permitted action.
- Keep live checkout freshness distinct from snapshot validity and PR-head eligibility; represent unavailable information as unknown.
- Existing query-only and viewer connections remain read-only; older snapshots remain usable with explicit age/provenance.

### Validation

Review a cross-interface state table against the retained live evidence and existing API/MCP fixtures. Exercise partial/failed work while an older graph remains queryable.

## T02 — Use one primary refresh action with role-aware advanced operations

**Epic:** E1. **Milestone:** M1. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T01. **Evidence:** F20. **Stages:** A06, B08, B09, C03, C04, C10.

Give the existing outcome a single primary action: Update context. It uses the existing appropriate ingestion/queued-refresh path for the user's authority. Separate a data-view reload from source analysis; move low-level run/build/setup controls to an advanced area and expose only relevant authorized actions.

### Acceptance criteria

- Document the label/action mapping for owner, viewer, editor and administrator; UI and agent guidance use the same meanings.
- A viewer sees exploration and an explanation of who can update; editors queue saved-configuration refresh without receiving import/configuration authority.
- No redundant analyzer runs are launched just because someone reloads a screen; retain advanced troubleshooting routes.

### Validation

Walk the role matrix through the existing six-service fixture and assert route selection plus backend permission denial independently of UI visibility.

## T03 — Define one approved scope contract for discovery and automation

**Epic:** E1. **Milestone:** M1. **Priority:** High. **Initial state:** Todo. **Depends on:** None. **Evidence:** F08, F09, F10. **Stages:** A04, A05, B04, C03.

Make the project, repository candidates, source type, branch, inclusion/exclusion rules and file-analysis boundaries explicit before mutation. Use the same preview semantics for agents and UI. Persist the chosen scope and refresh policy so later automation updates existing approved registrations rather than rediscovering an organization.

### Acceptance criteria

- Preview names/counts, local versus managed behavior and intended effects without cloning, registering or analyzing candidates.
- Changing provider, directory, filters or file scope invalidates the prior preview; the commit step identifies the exact current scope.
- Automatic refresh cannot expand the repository set, pull a local checkout, create ambiguous projects or overwrite dirty managed worktrees.

### Validation

Compare UI/MCP previews for identical inputs and run negative cases for changed scope, excluded sources, ambiguity and an invalid directory.

## T04 — Standardize actionable feedback and recovery across UI and agents

**Epic:** E1. **Milestone:** M1. **Priority:** High. **Initial state:** Todo. **Depends on:** None. **Evidence:** F11, L03. **Stages:** A06, A10, B09, C03, C04, C10.

Use a consistent error/recovery contract for validation, environment failures, permission denial, conflicts, partial work and unavailable providers. Keep feedback at the point of action and preserve drafts. Derive retryability and next actions from the actual operation; successful acceptance must still point to completion tracking.

### Acceptance criteria

- Directory errors are visible inside the active form and associated with the affected field; async progress remains observable after acceptance.
- Distinguish recoverable failure, non-retryable configuration and denied authority; do not blindly retry mutations or conflicts.
- Keep provider/auth details safe, accessibility announcements usable and user-facing language free of internal stack traces; retain diagnostics for authorized operators.

### Validation

Replay the retained invalid-root and missing-executable cases, plus conflict/role-denial cases, through UI and MCP and compare their meaning.

## T05 — Make installed startup and worker execution independent of PATH accidents

**Epic:** E2. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T01. **Evidence:** F05, F06, L05. **Stages:** A02, A03, A06, C01, C03.

Resolve analyzer execution from the installed/running executable identity rather than assuming a second PATH lookup succeeds. Verify the bootstrap artifact and actual worker path, and explain supported OS/architecture and client-registration steps as separate readiness checks.

### Acceptance criteria

- An absolute-path UI launch can analyze the fixture without an unrelated PATH entry; installed/source/embedded-UI paths identify the executable actually used.
- Doctor/bootstrap clearly separate executable health, host registration and usable completed context; do not silently edit unrelated personal client settings.
- Document the supported macOS/Linux matrix and WSL/shared-server alternatives without claiming native Windows support.

### Validation

Reproduce the live absolute-path failure with a clean environment, then verify ingestion and bootstrap against an installed candidate using existing distribution checks.

## T06 — Guide agent onboarding project selection and evidence use

**Epic:** E2. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T01, T03, T05. **Evidence:** F03, F04, F08, F25. **Stages:** A03, A04, A07, B03, B04, B06.

Align MCP descriptions and the host playbook around the actual connection mode and current workspace. Reuse an unambiguous matching project; ask only for unresolved scope/authority decisions. Teach agents to use focused/paginated evidence, report the answering snapshot and recognize when a viewer or query-only connection cannot onboard or refresh.

### Acceptance criteria

- Empty, sole-project, multi-project and query-only cases have explicit next steps and never silently answer from the wrong company.
- Agent guidance distinguishes read operations from mutation operations, terminal readiness from acceptance and exact evidence from missing context.
- No secrets are embedded in prompts/config examples; ordinary context queries remain read-only and bounded.

### Validation

Use scripted schema checks for correctness, then evaluate unprompted tool choices in an actual installed host under T27; do not equate prescribed calls with adoption.

## T07 — Enable bounded background refresh after approved onboarding

**Epic:** E2. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T01, T03, T06. **Evidence:** F02, F10, F12, L06. **Stages:** A05, A09, B05, B08, B10.

Make maintenance a documented onboarding default for approved projects using the existing scheduler and incremental ingestion. While the owning/shared runtime is active, propose a 15-minute cadence and a non-blocking catch-up after reconnect when overdue. Return the last completed graph promptly and expose refresh progress; allow explicit pause/manual mode and existing resource limits.

### Acceptance criteria

- First onboarding records the effective policy; no repeated prompt is needed for updates already authorized by that policy.
- Reuse unchanged artifacts, coalesce equivalent pending work, honor project/global caps, preserve cancellation and bounded retries, and avoid repeated failed refresh loops.
- Local sources are analyzed in place without automatic Git pull; managed clones use their approved branch and dirty-worktree safeguards.
- Connection closure does not imply continued indexing: display session-owned versus shared availability, last successful update and next eligible refresh.

### Validation

Run a controlled source edit/reconnect schedule, paused policy, busy project, failed update and two-agent contention case; prove reuse and no unauthorized scope/source changes. Cadence is a proposed default, not a measured optimal interval.

## T08 — Clarify local lifecycle and handle shared clients without competing owners

**Epic:** E2. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T01, T05, T07. **Evidence:** F07, F12, F37, L06. **Stages:** A03, A09, A11, B10.

Retain one lifecycle/writer owner per home. Give secondary clients an actionable connection/ownership explanation and make dashboard endpoint changes discoverable. Keep persisted data on disconnect and resumable work on reconnect. Use the existing secured shared deployment when refresh must outlive developer sessions.

### Acceptance criteria

- Two clients cannot create competing controllers or duplicate refresh schedules; lock contention explains the current supported connection route.
- Disconnect, crash, restart and maintenance leave saved graphs intact and advertise the current endpoint and recovery state.
- Never silently install an OS daemon, forward an unauthenticated local backend, bypass locks or claim multi-replica availability.

### Validation

Exercise owner/secondary-client behavior, controller crash cleanup, reconnect and a maintenance restart against disposable homes.

## T09 — Repair shared modal form and confirmation accessibility

**Epic:** E3. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T04. **Evidence:** F23, L03, L04. **Stages:** A10, B03, C02, C03.

Fix the shared interaction primitives once: named dialog semantics, appropriate focus entry, inert background, contained Tab navigation, meaningful labels, visible errors, Escape/cancel and focus restoration. Use these primitives throughout first use, import, confirmation and access/token forms.

### Acceptance criteria

- The first-use Name field has a programmatic label; close controls and all dialog titles have accessible names.
- Tab/Shift+Tab remain in an active modal, feedback is announced in context, and closing restores a logical focus target.
- Provide a non-mutating way to leave first setup so required initialization does not trap an explorer; retain deliberate confirmation for destructive actions.

### Validation

Replay the retained first-use Tab/Escape sequence, test nested confirmation and form-error focus, and include independent keyboard/screen-reader tasks in T27. Follow W3C APG modal guidance.

## T10 — Turn import preview into a retained review-and-build flow

**Epic:** E3. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T03, T09, T13. **Evidence:** F09, L01, L02, L03. **Stages:** A05, C03.

Keep discovery results visible with named candidates, counts, exclusions and analysis purpose. Preserve the draft and preview state when returning. Continue from the reviewed scope into the existing import/build pipeline, exposing progress and partial failures in the same flow.

### Acceptance criteria

- Preview displays the six candidates returned under the backend's results field instead of closing into an empty workspace.
- Back/cancel/reopen behavior retains an understandable draft; changed input requires a new preview before continuation.
- Invalid roots stay in the active form, and preview remains non-mutating; build acceptance is followed to a completed/partial/failed outcome.

### Validation

Replay the exact live preview, reopen and invalid-root sequence through real Chromium, then complete the six-service graph and compare the approved scope with MCP.

## T11 — Make graph selection readable keyboard-accessible and fully inspectable

**Epic:** E3. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T02, T09. **Evidence:** F19, F24, F25, F15, L07. **Stages:** C05, C06, A07.

Use a readable initial team/service focus rather than fitting all content below useful text size. Keep overview topology available and make search/focus discoverable. Provide equivalent keyboard selection for services, resources and relationships, and explicit pagination/load-more for bounded evidence.

### Acceptance criteria

- The six-service initial view avoids the observed six-pixel service label; validate a proposed minimum 12 CSS-pixel rendered label at supported desktop sizes.
- Mouse and keyboard can select a service, follow a relationship and inspect its source evidence with a visible focus indicator.
- Evidence totals, omitted items and continuation controls are clear; agents receive equivalent bounded paging information.
- Validate larger rendered and real extracted graphs separately; do not claim the generated 150-service fixture establishes extraction performance.

### Validation

Replay 1440x900 plus 1280x720/1920x1080, keyboard selection and a >80-object/>120-trace case; compare focused/overview readability and measure the candidate rather than screenshots alone.

## T12 — Distinguish empty loading denied failed and partial workspace states

**Epic:** E3. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T01, T02, T04, T09. **Evidence:** F11, F20, F21, F22, L09, L10. **Stages:** B02, B04, B09, C02, C04, C08, C10.

Use explicit page states instead of treating request failure as an empty project collection. Let first-time explorers choose their existing entry route before creating state. Keep primary exploration visible, advanced operations role-appropriate, and show plain next steps for missing graph, partial ingestion, denied access and unsupported PR context.

### Acceptance criteria

- A failed project listing offers retry/diagnostics and cannot automatically launch project creation; a truly empty authorized workspace offers an explicit start action.
- Viewers do not receive a toolbar full of unusable administrative controls; permitted refresh/recovery actions match the capability contract.
- Access-loss pages give generic request-access/return guidance without revealing hidden project existence or cause.
- Existing one-snapshot compare guidance remains clear; provider-not-ready and genuinely zero-PR states are distinct.

### Validation

Use real UI walkthroughs with controlled list failure, no-project, viewer, removed membership, partial graph and local-only PR states.

## T13 — Apply approved file scope consistently to facts contracts and metrics

**Epic:** E4. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T03. **Evidence:** F09, F13, L11. **Stages:** A05, A06, B05, C03, C04, C06.

Make analysis purpose and existing paths.include/paths.exclude configuration visible during scope review. Use that same effective file scope for detector facts, OpenAPI enrichment, packs and metrics; distinguish potential examples/tests/scenarios before attributing them to the enclosing service. Do not globally discard directories solely because their names resemble fixtures.

### Acceptance criteria

- Reproduce default DiffMind-repository example attribution as a regression, then demonstrate an explicitly approved production scope excludes those example/scenario facts and contract fields.
- All extraction/enrichment stages agree on the effective scope; excluded alternate OpenAPI scenarios cannot enrich an included route.
- Scope changes invalidate reuse and appear in provenance; historical snapshots retain their original meaning.
- Multiple applications in one repository remain an explicit limitation unless existing service configuration separates them; do not invent a general monorepo model.

### Validation

Use the retained public-clone provenance plus positive examples-as-production and negative alternate-scenario cases. Inspect assembled UI/MCP graphs and contracts, not only detector counts.

## T14 — Explain evidence origin uncertainty and coverage consistently

**Epic:** E4. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T01, T13. **Evidence:** F10, F13, F16, F18, L11. **Stages:** A07, A10, B05, B06, C04, C06, C09.

Present source-extracted, identity-resolved and pack-declared evidence distinctly with revision/file scope and limitations. Keep processing success, freshness, snapshot validity and known coverage separate. Use explicit unknown/unsupported explanations instead of interpreting zero edges as proof that systems are disconnected.

### Acceptance criteria

- UI and MCP report the same evidence class, answering run and applicable source/coverage limits for a relationship.
- Declared edges and static calls cannot be relabeled as observed runtime traffic; stale/dirty/unknown context never becomes exact PR evidence.
- Absent findings and unsupported patterns are described as incomplete evidence, with a supported existing correction path.

### Validation

Inspect a declared relationship, unresolved destination, excluded example, stale source and a successfully analyzed unsupported convention through both interfaces.

## T15 — Use the existing gap and pack workflow for reviewed corrections

**Epic:** E4. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T03, T14. **Evidence:** F14, F17, F34. **Stages:** A10, B08, C10.

Route missing-convention reports into the existing private improvement-gap workflow, with an owner/support path. Require a synthetic positive and near-match negative example before changing existing detector/config/pack behavior. Keep activation reviewed, auditable and reversible; measure correction effort rather than assuming indexing teaches conventions.

### Acceptance criteria

- Viewers/editors receive an appropriate existing report-or-admin handoff without gaining host/path/pack authority.
- An accepted correction includes expected graph evidence, negative fixtures, provenance and rollback validation.
- No proprietary snippets or public contribution are emitted automatically; unsupported frameworks remain explicitly scoped.

### Validation

Take one labeled coverage gap from T26 through report, review, tested correction and rollback using existing pack/detector mechanisms; verify previous snapshots stay immutable.

## T16 — Explain snapshot differences through input provenance and evidence

**Epic:** E4. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T14. **Evidence:** F39, L08. **Stages:** C07, C09, A08, B07.

Keep graph comparison about saved facts while making changed source artifacts, file scope, analyzer/config and pack inputs discoverable. Distinguish topology/fact changes from evidence-only changes and avoid claiming a cause solely from a diff. Link the existing contract view when a reviewer asks a compatibility question.

### Acceptance criteria

- The known checkout before/after result retains five modified facts, with clear evidence/context and stable pinned run IDs.
- Changes to scope or packs are visible as possible inputs; historical missing provenance is unknown rather than guessed.
- An evidence-only change can be identified without implying a runtime regression or silently dropping meaningful source changes.

### Validation

Compare source-change, scope-change, pack-change and evidence-only pairs; verify self-comparison, pagination and immutable historical interpretation.

## T17 — Expose existing contract compatibility alongside graph changes

**Epic:** E5. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T12, T14, T16. **Evidence:** F39, L08. **Stages:** A08, B07, C07, C08.

Bring the existing compare_contracts result into the dashboard's change-investigation path using the same pinned before/after snapshots and service selection. Clearly separate architectural fact differences from request-field compatibility and show the supported OpenAPI limits.

### Acceptance criteria

- The checkout scenario displays the same four field changes and labels in dashboard and MCP, alongside its distinct five-fact graph diff.
- Required additions/removals/type changes remain potentially breaking; optional nullable additions retain their existing compatibility interpretation.
- Unsupported/absent contracts cannot be summarized as no breaking changes; a missing baseline is an explicit prerequisite.

### Validation

Use the existing checkout-v2 comparison, absent/unsupported-contract and same-run cases across UI/HTTP/MCP without inventing runtime compatibility guarantees.

## T18 — Make PR revision eligibility and deletion limitations actionable

**Epic:** E5. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T14, T16, T26. **Evidence:** F27, F29. **Stages:** A08, B07, C08.

Explain why current source freshness can differ from PR-head eligibility. Use matching clean baseline/head saved evidence where available; retain candidate/unknown output for removed entrypoints and indirect changes that cannot be proven. Guide authorized agents to prepare existing temporary checkout/analysis workflows when precise evidence is needed, without switching a user's active branch.

### Acceptance criteria

- The observed public PR explains its default-branch/head mismatch; refreshing the same default branch is not presented as a guaranteed cure.
- Deletions, unchanged-route internal changes and configuration changes are explicitly classified as supported evidence or unproven limitations.
- Exact callers require matching revision evidence; no automatic clone, branch switch or new PR-head capture service is introduced.
- An absent exact caller result is never described as safe-to-merge.

### Validation

Use independently labeled add/remove/configuration/transitive PR cases from T26 plus the retained stale public PR response; verify exact and candidate separation.

## T19 — Align existing repository provider settings with PR retrieval states

**Epic:** E5. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T03. **Evidence:** F30, L09. **Stages:** A05, B07, C08.

Use the repository's approved provider/API-host configuration consistently for import and existing PR retrieval. Explicitly identify local-only, missing remote metadata, unsupported host, credentials failure and a successful zero-PR result. Keep the provider contract limited to currently intended GitHub-compatible deployments.

### Acceptance criteria

- Local Demo Shop repositories show provider-unavailable context rather than pretending a filtered zero list proves their GitHub PR state.
- An approved custom GitHub API host follows the same credential/hostname policy for import and PR queries; test redirects and prevent credential forwarding to unrelated hosts.
- Actual public PR listing/refresh still works, with safe actionable auth/rate-limit feedback.

### Validation

Replay the live local/public cases and a controlled GitHub-compatible API fixture; verify host/credential boundaries. Enterprise behavior remains unverified until its target deployment trial passes.

## T20 — Present PR evidence before heuristic scores and clarify delivery

**Epic:** E5. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T14, T17, T18, T19. **Evidence:** F26, F28. **Stages:** A08, B07, C08, C09.

Lead the existing PR view and agent output with changed surfaces, compatibility findings, exact callers, candidates and missing evidence. Keep the numeric heuristic secondary with explicit reasons and eligibility. State that the current workflow is on-demand inspection; an agent can use its findings within the user's separately authorized review workflow.

### Acceptance criteria

- Low/no score is not presented as a probability, merge recommendation or proof of safety; stale company context remains excluded from scoring.
- Warnings distinguish potentially breaking evidence, attention prompts and unsupported questions; UI and agent language agree.
- Document the actual on-demand delivery path and record automatic code-host posting/checks as a separate scope decision, not an implemented behavior.
- Preserve useful explanations observed on public PR #9 while reducing unsupported certainty.

### Validation

Review outputs against the labeled PR set from T26 and user comprehension in T27; record false alarms/misses without claiming calibrated scores.

## T21 — Make shared access scope explicit and safely adopt scoped defaults

**Epic:** E6. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T01, T03. **Evidence:** F32, F33. **Stages:** B01, B02, B04, B05.

Use explicit scoped access for new company setup and show the active mode/recovery identity requirements. Provide a reviewed migration path for existing legacy deployments rather than silently changing authority. Keep project creation/import aligned with the rule that every member can read all evidence in that project.

### Acceptance criteria

- New shared setup validates ordinary viewer/editor identities and retains admin recovery access before it is considered ready.
- Existing legacy instances are not silently migrated; mode changes and prepared grants are reviewed and clearly reported.
- Project-wide evidence visibility and lack of repository-permission mirroring are explicit during scope/access setup.
- No per-repository ACL or automatic SSO/group provisioning is invented in this batch.

### Validation

Test new scoped setup and an explicit legacy-to-scoped migration with no-grant/viewer/editor/admin identities and guessed project URLs.

## T22 — Make joining browser and agent setup a clear role-aware handoff

**Epic:** E6. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T05, T06, T21. **Evidence:** F31, F34. **Stages:** B01, B02, B03, B04, B08.

Give joining developers a single clear explanation of project access, browser identity, MCP endpoint and the supported project-token/proxy route. Use the existing admin token issuance UI with least-privilege defaults and explicit expiry/renewal/support responsibility. Avoid making ordinary users configure local ingestion for a prepared company graph.

### Acceptance criteria

- An ungranted user gets a generic admin-help route; a granted viewer can open existing context and connect a viewer agent with 13 read tools.
- The handoff explains that browser cookies do not configure agent identity and identifies the token renewal/support owner.
- One-time token secret handling remains explicit and private; do not distribute admin recovery credentials or silently broaden grants.

### Validation

Repeat company entry with separate browser/agent contexts, verify viewer tool discovery and token expiry/renewal guidance, then observe actual installed-client joining under T27.

## T23 — Make access loss and offboarding understandable across independent grants

**Epic:** E6. **Milestone:** M2. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T04, T21, T22. **Evidence:** F35, L10. **Stages:** B02, B10, B11.

Keep prompt access removal while showing generic return/request-access guidance. Add an operator offboarding sequence that separately reviews person memberships and explicitly associated project agent tokens using existing metadata; never infer token ownership just from similar labels or assume membership removal revoked a service grant.

### Acceptance criteria

- An already-open viewer graph loses protected content after the existing poll; its explanation does not reveal whether an inaccessible project exists.
- Admin workflow makes independent token state and explicit selection visible; revocation targets the verified subject/token rather than row position.
- Membership-only removal leaves independent tokens unchanged by design, while explicit token revocation returns 401; the distinction is clear to the operator.

### Validation

Replay the corrected live revocation trial and reorder membership rows; test ordinary/access-lost responses and explicit token revocation.

## T24 — Verify upgrade backup and restore with access-policy reconciliation

**Epic:** E6. **Milestone:** M2. **Priority:** High. **Initial state:** Backlog. **Depends on:** T08, T21, T23. **Evidence:** F36, F37. **Stages:** A11, B09, B11.

Provide one operational continuity runbook around existing offline backup, rotation, verification and lifecycle controls. Before a restored shared instance accepts ordinary traffic, review memberships/tokens against the current offboarding record because older backups can restore formerly revoked grants. Keep single-server/path/schema boundaries explicit.

### Acceptance criteria

- Upgrade/maintenance preserves valid context or gives a clear rollback path; no implicit schema/path relocation or distributed availability claim.
- A restore drill demonstrates the older-token resurrection risk and a verified reconciliation/revocation step before reopening ordinary access.
- Backups, credential-bearing artifacts and operator logs remain private; verify integrity and existing non-overwrite/lock safeguards.
- Document recovery targets observed in the drill rather than promising an unmeasured recovery time.

### Validation

Run the existing backup/restore gates plus a disposable revoke-after-backup scenario and queued-work interruption/recovery against the supported storage modes.

## T25 — Align installation operation and evaluation docs with the implemented journey

**Epic:** E7. **Milestone:** M3. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T05, T06, T20, T21. **Evidence:** F01, F40. **Stages:** A02, A03, A10, A11, C10.

Consolidate canonical entry guidance for the three personas, update obsolete release and evaluation instructions, and cross-link host setup, maintenance, shared joining and PR boundaries. Mark historical documents as historical while preserving the audit records and existing roadmap history.

### Acceptance criteria

- AGENT_SETUP and entry docs agree on actual releases, supported platforms, management/query-only modes and the tested installed-client route.
- Discoverable evaluation commands match the current CLI; obsolete cheap/eval instructions cannot be mistaken for a working validator.
- The final docs explain the actual background policy, access setup and on-demand PR delivery; exact candidate and evidence limits are recorded.

### Validation

Execute the canonical commands from a clean candidate installation and check links/examples against the CLI/tool catalog; do not publish a release as part of this issue.

## T26 — Establish an independently labeled real-source graph and PR corpus

**Epic:** E7. **Milestone:** M1. **Priority:** High. **Initial state:** Todo. **Depends on:** None. **Evidence:** F14, F15, F28, F29, F41, L11. **Stages:** A08, B05, B07, C01, C04, C08.

Build on the pinned OpenTelemetry/Boutique public scans and the existing six-repository fixture, but validate full assembled graphs against independently labeled relationships. Include the public enclosing-repository example case and PR additions/removals/contracts/configuration/transitive/no-impact controls. Label unknown ground truth rather than counting it as an error.

### Acceptance criteria

- Record pinned commits, source scope, expected known edges, relevant callers, contract changes and negative cases independently of DiffMind output.
- Report assembled-graph recall/precision only for labeled supported cases, plus missed supported conventions and explicit unsupported patterns.
- Separate synthetic rendering/load checks, detector counts, full graph correctness and actual PR review outcomes; record false alarms/misses without calibration claims.
- Retain safe public/synthetic inputs; no private company code is published automatically.

### Validation

Reproduce the corpus and blind-label agreement, then use it as the acceptance basis for T13/T15/T18/T20/T28. Targets are reviewed before candidate scoring; this issue is evidence work, not universal framework expansion.

## T27 — Observe independent developers and actual installed-agent use

**Epic:** E7. **Milestone:** M3. **Priority:** Medium. **Initial state:** Backlog. **Depends on:** T06, T10, T11, T12, T17, T20, T22, T23. **Evidence:** F04, F05, F38, F42. **Stages:** A01, A03, A07, A09, B01, B03, B06, B11, C01, C05, C06, C09.

Run independent task sessions for the three personas using the improved candidate, including at least three observations per path and a keyboard/screen-reader participant where possible. In actual installed hosts, evaluate whether agents seek relevant context without a prescribed tool script and preserve uncertainty. Measure first useful result, recovery and interpretation.

### Acceptance criteria

- Retain consented task evidence and clearly separate human responses from automated checks; document recruitment shortfalls rather than simulate participants.
- Measure task success/blockers, setup time, source/snapshot understanding, wrong-context answers and missed/unnecessary agent calls; establish baseline before numerical performance targets.
- Participants can identify approved scope, maintainability, evidence limits and PR uncertainty; critical task blockers return to their owning epic.
- Avoid compulsory telemetry: consented session records and existing operation logs are sufficient for this batch.

### Validation

Compare the candidate with the retained observer baseline; use fresh and returning sessions plus ordinary source edits and access loss. Results do not establish long-term retention without a separate longitudinal observation.

## T28 — Verify the exact candidate across all three journeys and safeguards

**Epic:** E7. **Milestone:** M3. **Priority:** High. **Initial state:** Backlog. **Depends on:** T07, T08, T11, T15, T20, T24, T25, T26, T27. **Evidence:** F38, F41. **Stages:** A06, A09, A11, B05, B09, B11, C03, C06, C08.

Re-run live cross-journey acceptance against the same installed candidate used for the native/distribution gates. Cover approved onboarding, background refresh, readable evidence, contract/PR interpretation, company restrictions and recovery together. Preserve existing durable queue, immutable history, credentials and writer-lock safeguards.

### Acceptance criteria

- No critical activation, evidence-interpretation or access-control blocker remains in the agreed tasks; all 42 audit and 11 live findings have a verified resolution or explicit boundary/measurement result.
- Real-binary/browser/HTTP-MCP/stdio/native gates are recorded for the exact revision with corpus and configuration identities; required checks pass or the candidate stays unreleased.
- Candidate results preserve previous useful behavior: six-service graph, incremental reuse, field changes, permissions, token revocation and honest stale-PR output.
- Publish a concise outcome/limitations record; do not mark success from accepted work, green scan status or a prepared fixture alone.

### Validation

Use the existing targeted/full readiness and native archive gates as appropriate, plus the real browser protocol and independent trials. Publication remains a separate maintainer action.
