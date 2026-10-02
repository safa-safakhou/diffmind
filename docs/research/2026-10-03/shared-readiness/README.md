# Shared readiness implementation and validation

This candidate extends `origin/codex/ux-context-batches` at `7e5ab06` on `codex/shared-readiness`. It completes the combined runtime/work/action contract remaining in MNI-175.

## Changes

- One shared query-service projection supplies the HTTP readiness endpoint, workspace metadata, native MCP readiness tool and management inspection catalog.
- Saved graph validity, age/input provenance, current work, analysis freshness, observed revisions, permitted actions and unverified coverage remain independent. PR-head eligibility stays unknown.
- Queued/running jobs suppress conflicting actions. Partial/failed work preserves usable saved context. Query-only and viewer connections retain read-only behavior. Job success preserves correlated partial results.
- Browser connection failures retain the loaded view and pause actions; confirmed denial removes it. Late project/poll responses are discarded. MCP error results carry safe structured recovery without project data or authority.
- Saved artifacts are checked for null/mismatched graphs, and HTTP readiness reuses the validated artifact cache. Workspace metadata points to the same saved run as readiness.

## Validation and limits

The manifest and retained test outputs identify the exact candidate and checks. Contract tests exercise empty, queued, running, completed, partial, failed, cancelled, interrupted, invalid artifact and access-loss states. They check saved-graph fallback, timestamp/input provenance, query-only authority, role ceilings, private-project isolation and project-token revocation. Browser component/library tests cover the state explanations and transport adapter.

The live trial used isolated copies of two real company repositories, a private disposable DiffMind home, the compiled candidate and headless installed Google Chrome. The original repositories and personal agent settings were not modified. Company source, extracted graphs, raw server logs and credentials remain outside this repository. Retained observations contain only aggregate counts, generic workspace/run identifiers and readiness state.

Real operations reached completed context, observed running work with a saved graph, generated partial results when a controlled analyzer wrapper failed one source, and failed after both analyses were deliberately rejected and the disposable analysis-artifact directory was moved aside. The immutable saved graph remained queryable when analysis artifacts were unavailable. These are injected analyzer/artifact failures in a real runtime, not claims about failures inherent to those company repositories.

HTTP and remote MCP readiness matched for viewer, editor and admin roles. Real query-only stdio MCP preserved query access and advertised no mutation authority. The browser displayed partial and failed work independently of saved context. A controlled 503 capability response verified connection continuity; actual scoped membership removal returned 404 and removed the viewer's graph. No JavaScript exceptions occurred in the passing trial.

This is automated observer evidence on macOS arm64. It does not establish independent-human usability, unprompted installed-host adoption, enterprise provider/credential readiness, large-corpus latency, comprehensive dependency extraction or release readiness. Existing MNI-181/MNI-182/MNI-192/MNI-197/MNI-201/MNI-202 acceptance remains separate.

Early verification found old hard-coded tool counts/isolation argument tables, a broken system Node dependency, and an absent bundled Chromium executable. Tests were updated for the additive readiness tool; the bundled Node and installed Chrome completed the final checks. The initial failure trial produced partial context from retained artifacts, so the final failed trial explicitly made analysis artifacts unavailable. Those setup/observer outcomes are not product regressions.

## Reproduction

Run the full Go suite, the dashboard's tests/build, and targeted query/UI/MCP/agent API race checks. Build `./cmd/diffmind` after embedding the dashboard production output.

In a private temporary directory, copy two approved repositories, launch the candidate with a disposable `DIFFMIND_HOME`, scoped project access and a temporary trusted-proxy secret, then preview and approve exactly those sources. Compare HTTP `/readiness`, `workspace.readiness` and MCP `get_readiness` for each role. Force refresh through a controlled analyzer wrapper to observe running/partial outcomes; move only the disposable analysis artifacts aside and reject both analyses to exercise failed-work saved fallback. Check standalone `mcp` against the same home for query-only behavior. Restore successful transport reads after a controlled browser 503, then revoke the disposable viewer membership to verify confirmed denial. Stop all owned processes after the trial.

The retained [live observer](harness/live.mjs.txt) has a `.txt` suffix to prevent accidental execution. Copy it into a private trial directory as `live.mjs`, supply `sources/` with exactly two approved repository copies and an `analyzer` wrapper controlled by `failure-mode` (`none`, `slow`, `partial`, `all`). The wrapper delegates to `/tmp/diffmind-readiness-candidate`, waits three seconds in slow mode, rejects one selected copied source in partial mode and rejects both in all mode. Adapt the bundled Node/Playwright and installed Chrome paths for the target machine. Its logs and DiffMind home stay beneath that private trial directory.
