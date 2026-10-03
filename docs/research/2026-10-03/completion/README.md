# DiffMind completion validation — 2026-10-03

Candidate parent: `9b164df37c5976501ac35cdb3a3d7e5ad4e10733`, branch
`codex/diffmind-backlog-batch`. This directory is committed with the follow-up
implementation. Binaries were built with the current uncommitted product changes;
full hashes and build metadata are retained. The clean source installer uses
`-trimpath`, so its binary differs from the ordinary trial build. Both record the
same parent and modified source. No release was published or branch pushed.
PACT and original company checkouts were excluded.

## Completed technical scopes

- **MNI-178:** recovery matrix compares actual HTTP handlers and the MCP management
  adapter for invalid roots/configuration, authentication, denied roles, hidden
  projects, access revision conflicts, queue capacity and safe provider failure.
  A real missing analyzer is accepted (202), then observably fails while retaining
  source registrations and operation diagnostics. Browser confirmation feedback
  retains context for 400/403/409/429/503/transport failure, with one mutation request
  and no automatic replay. Network-injected browser statuses demonstrate feedback;
  the Go matrix separately tests actual enforcement. Previous batch evidence retains
  field-associated import errors and draft preservation.
- **MNI-194:** the actual installed Codex CLI app-server joins the scoped HTTP MCP
  endpoint as viewer, discovers 14 read tools (current catalog, superseding the
  ticket's original 13), and reads a prepared six-repository synthetic graph with
  three source-backed dependencies. A separate replacement credential/client works
  before explicit old-token revocation (old credential then returns 401). Secrets
  are generated privately, supplied via environment, and absent from client config.
  Prior separate browser contexts verify granted/ungranted handoff. This is scripted
  joining, with zero model turns; it does not establish unprompted agent adoption.
- **MNI-200:** a fresh source installation executes current help/doctor commands,
  installed-binary showcase analysis, and installed-client joining. Historical
  cheap/eval instructions are marked obsolete; current validation and persona
  setup routes are linked. CLI subcommand help displays usage but `agent --help`
  and `mcp --help` currently exit 1 and 2. Native execution was macOS ARM64 only;
  the observed release has macOS/Linux AMD64/ARM64 assets, not proof of execution
  on every platform. Release metadata is a dated observation, not a release action.

## Remaining evidence requirements

**MNI-198 and MNI-199 remain In Progress.** Product implementation and bounded
mechanism tests pass, but their required independently labeled corpus under
MNI-197 is incomplete. MNI-199 also requires human comprehension under MNI-201.
Human participants: **0**. Independent blind reviewers: **0**. Recruitment was
requested; no participant response has been recorded. The session protocol in
[human-study.md](human-study.md) is preparation, not participant evidence.

The source-first labels were frozen before candidate scoring (hash in manifest).
Eight pinned Online Boutique services were actually assembled: eight nodes, two
edges, and zero of seven known source gRPC relationships found. Those patterns
were predeclared unsupported/unverified; supported-case precision and recall are
null. Unlabeled emitted facts have unknown truth. This is a coverage limitation,
not a passing graph accuracy result or a claim that missed dependencies are safe.
The OpenTelemetry source was pinned/prepared but not assembled in this follow-up.

Six controlled PRs mutate real pinned public source: addition, removal, transitive
behavior, configuration, contract and a no-impact control. Actual head analyses
and a read-only GitHub-compatible provider exercise HTTP and MCP projections.
All six retain unknown freshness, ineligible scoring and explicit limitations,
with zero exact callers/candidates. These cases verify conservative unsupported
output; they are not positive exact-caller evidence, observed upstream review
outcomes, score calibration, or blind-label agreement. Existing source-positive
unit fixtures separately verify exact/candidate separation and revision gating.
The previous batch retains the real public PR 9 default-branch/head mismatch.
MNI-197 still needs blind agreement, labeled supported-case graph accuracy,
enclosing-repository controls and actual PR outcome evidence.

## Additional fixes found by validation

Exact PR callers now require their own clean saved analysis and a fact revision
matching that analysis; missing, dirty, unknown and older caller revisions do not
qualify. Caller repositories need not share the changed repository's commit.

An explicitly local repository stays local even with PR remote metadata: creation,
manual sync, refresh and missing-path handling cannot silently clone or replace its
checkout. Tests retain a dirty draft, active branch, HEAD and registration, and
verify that no managed worktree is created. Existing inferred Git defaults remain;
previously stored registrations are not migrated automatically.

A separate exploratory full-management stdio registration encountered competing
local-controller startup during installed-client inventory. A scripted thread
query succeeded, but stable global manager adoption was not established. This
belongs to the existing host/controller boundary (MNI-182); the company HTTP
joining result above is independently verified and does not claim to resolve it.

## Verification and reproduction

Full Go suite: 65 tested packages. Targeted race suites: store, UI, agentapi and
mcpserver. Frontend: 20 library and 26 component tests. Installed candidate
showcase: actual six-repository analyzer and source-change checks. Browser recovery:
nine retained observations, zero JavaScript exceptions. Logs and safe aggregates
are committed; raw client/provider logs, credentials, private homes and cloned
source are excluded. These counts describe bounded checks, not universal coverage.

Use [canonical validation](../../../validation.md) and
[agent setup](../../../../AGENT_SETUP.md) for installation. The retained observer
sources have `.mjs.txt` extensions to avoid being mistaken for portable product
scripts. To replay, copy them to a fresh private validation root as `.mjs`, create
`evidence/`, build `bin/diffmind`, install `final-installed-bin/diffmind`, prepare
public sources with `scripts/prepare-public-benchmarks.sh`, and copy the frozen
labels into `evidence/`. Adjust the explicit local Codex, Chrome and Playwright
paths. Observers create only owned private homes/worktrees and stop their processes
in `finally`. They are prescribed checks, not human sessions or autonomous adoption
experiments. Retain separate source and installed binary hashes for each trial.
