# DiffMind system validation

The live tests found defects beyond the existing passing suite. Six product commits now repair graph usability, Flask extraction, PR evidence matching, local HTTP MCP management, and the separate extraction viewer. The final checks passed on clean product candidate **2cd9d328cc24de955184fc8aab863892a6ae91b3**.

This is an operator-led assessment on Linux AMD64 inside WSL, performed on 4 October 2026. It supports the individual developer, company joiner and architecture explorer journeys within the tested conditions. It does not establish universal repository coverage, continuous availability, human comprehension or performance at production load. The candidate remains unpublished; no remote push or PR comment was made.

## What was exercised

| Surface | Real action and result |
| --- | --- |
| Installation and agent mechanisms | Packaged and installed the clean binary through the native archive gate. Embedded UI, SQLite, installed company/MCP acceptance and managed recovery passed. The underlying acceptance suite covers controller contention, reconnect, crash cleanup, cancellation, retry and maintenance. |
| Company dashboard | Opened six routes at 1440, 768 and 390 pixels with a synthetic 150-service, 10-team graph. Selected a service using search and keyboard, opened its inspector, exercised dialog focus containment, Escape and focus restoration, then tested a short 390 by 640 phone. Graph space, readiness visibility and control bounds are asserted. |
| Extraction viewer | Opened repository, add-dialog, saved CLI artifact, run-form, completed-run and outcome-graph states at desktop and phone widths. Launched real extraction twice, including after storing an unrelated repository as the previous target. Both runs completed with replayable events and no document overflow. |
| Accessibility | 19 dashboard scans, 12 extraction-viewer scans and one actual public PR scan: no axe WCAG A/AA violations in the inspected states after fixes. Screenshots and keyboard behavior were also inspected. This is not full assistive-technology certification. |
| Public extraction | Analyzed the pinned Flask tutorial and compared its HTTP operations with the actual Flask 3.1.2 runtime URL map. Final extraction matches all 12 unique explicit method/path pairs in scope. |
| Company relationships | Analyzed eight pinned Online Boutique services, applied the existing reviewed correction pack, checked HTTP/MCP provenance parity, reused unchanged artifacts, edited one source under approved background maintenance, paused maintenance and rolled back the correction. All 11 recorded workflow stages passed. |
| PR investigation | A controlled change to copied Flask source with an explicitly synthetic caller produces two changed entrypoints and one eligible exact caller. Missing patch, unrelated filename, deleted-line, stale-head and dirty-caller controls reject exact eligibility. Actual HTTP MCP initialization/session/tool use returns the same company evidence as HTTP. |
| Actual upstream PR | Read [pallets/flask PR 5918](https://github.com/pallets/flask/pull/5918), observed head a82e942870b6472bb40017349cca772c242eb1ad, through the provider, browser and HTTP MCP, with matching company evidence. The saved tutorial graph is stale for this PR; exact eligibility is correctly false. The UI identifies on-demand evidence and the heuristic limit. No positive exact-impact claim is made for this upstream PR. |
| Access and recovery | Exercised scoped viewer/editor/admin/ungranted identities on loopback, rejected stale scope approval, and verified immutable history. Upgraded/restored from prior implementation e2e9530 on JSON and SQLite, independently reconciled resurrected membership and assigned-token grants, and preserved unrelated service/admin access. |

Source and binary identities, package versions and evidence hashes are in [manifest.json](manifest.json). Detailed coverage is in [scenarios.md](scenarios.md) and [findings.md](findings.md). The [Linear snapshot](linear-status.json) separates task state from this run's evidence.

## Defects found and repaired

| Finding | Before | Verified result | Commit |
| --- | --- | --- | --- |
| Dashboard accessibility | Primary buttons and dim labels failed contrast; graph selectors lacked accessible names; mobile status scrolling was not keyboard reachable. | Contrast, names and keyboard focus now pass browser checks. | 20b21ba |
| Narrow dashboard layout | Header, readiness and inspector consumed the phone viewport; the graph was nearly invisible. Tablet header text clipped. | Wrapped controls, bounded readable panels, usable graph height and scroll access to the inspector. A geometry assertion supplements axe. | 4aea079 |
| Flask operation extraction | Blueprint prefix was dropped and a multi-method decorator emitted only one method. Against the runtime oracle: 8 emitted operations, 5 correct, 7 missing, 3 unexpected. | 12 of 12 unique in-scope operations, no extras. Tests cover reused receiver names, literal tuples/lists, registration overrides and conservative handling of dynamic/conflicting prefixes. | 4e3c0d5 |
| PR changed-line evidence | Route source location pointed at the function body, so a decorator-only edit was missed. | Each route uses its own decorator coordinates, including stacked decorators. | 5d1dadb |
| PR HTTP caller matching | Real extracted caller display names contain service/URL information and did not match a route operation. | Matching uses structured method/path/URL evidence; negatives preserve method, case, trailing slash, escaping and protocol distinctions. | 5d1dadb |
| Local HTTP MCP management | Internal dispatch reused a synthetic host and rejected a valid local request with 403. | The adapter carries the already-validated listener host while rechecking current credentials/permissions. A real HTTP regression still rejects a forged external Host. | 18b32c3 |
| Extraction run target | Clicking Run on a selected repository opened a blank path; saved defaults could target a previous repository. | The selected repository overrides only the saved target, preserving worker/confidence settings. Real UI runs and three unit regressions cover first use and remembered settings. | 2cd9d32 |
| Extraction accessibility and phone layout | Dialog had no role/focus boundary, fields lacked associated labels, contrast failed, activity/graph scrolling lacked keyboard focus, and content overflowed the phone. | Shared dialog focus behavior, labels, contrast, scroll focus, wrapping and bounded columns pass desktop/phone checks. | 2cd9d32 |
| Saved extraction history | A completed CLI artifact with no events showed pending stages, waited for events and displayed an enormous elapsed duration from a zero timestamp. | Missing history is explicit; absent terminal stages say not recorded; missing dates show unavailable and terminal durations use finish time. Live UI runs still show recorded stage completion. | 2cd9d32 |

Before evidence includes [Flask comparison](evidence/flask-before-after.json), [dashboard audit](evidence/workspace-before.json), [extraction audit](evidence/extractor-before.json), and [run-target failure](evidence/logs/extractor-target-before.txt). Regression failure logs for Flask, caller matching and local MCP are retained under evidence/logs.

The [phone graph before](evidence/screenshots/workspace-phone-before.png) and [after](evidence/screenshots/workspace-phone.png) illustrate why passing automated accessibility checks did not by themselves establish usability. The final [extraction phone view](evidence/screenshots/extractor-phone.png), [keyboard-selected graph](evidence/screenshots/workspace-keyboard.png), [company access screen](evidence/screenshots/company-access.png) and [public PR screen](evidence/screenshots/public-pr.png) come from the tested binary.

## Verification and interpretation

Full Go tests and the full race suite each passed across 65 packages with tests; go vet passed. Frontend tests passed: 20 workspace-library, 26 workspace-component and 14 extractor tests. Both production bundles were rebuilt and embedded. The three bundled knowledge packs and the evaluation-only Boutique correction passed lint and their positive/negative assertions. Installer, demo and six-service analyzer showcase checks passed. govulncheck 1.7.0 and both npm audits found no reported vulnerabilities at the time of this run. These are bounded automated checks, not a security audit.

The public [Flask source](https://github.com/pallets/flask/tree/2c1b30d0503cfb064f1cb252e6614a06915a362a/examples/tutorial) is pinned to 2c1b30d0503cfb064f1cb252e6614a06915a362a. The reproducible [runtime oracle](harness/flask-runtime-oracle.py) instantiated the tutorial app and inspected its URL map; it did not send application traffic. It excludes static files and implicit HEAD/OPTIONS, and deduplicates an explicit GET / alias. A match on these 12 operations is not a general Flask or other-language accuracy estimate. Dynamic imports, runtime mounts and unsupported patterns remain outside this evidence.

[Online Boutique](https://github.com/GoogleCloudPlatform/microservices-demo/tree/b9a978db9e01f4ad3dca9494a22cb9edc17548fe) is pinned to b9a978db9e01f4ad3dca9494a22cb9edc17548fe. Baseline found eight services and two source-extracted edges, with none of the 12 labeled configuration declarations. Explicitly installing the tested pack produced exactly those 12 declarations. This measures the declared subset; source-extracted extras remain unlabeled, and no complete graph precision or runtime reachability claim follows. Upstream services were not executed.

The controlled PR mutates a private copy of public tutorial code and adds a synthetic requests caller. It is not an actual upstream PR. Conversely, the actual upstream read proves provider/browser handling and conservative stale evidence, not independently verified positive PR impact. Independent human reviewers: zero.

Background maintenance analyzed one changed source and reused seven; unchanged ingestion reused all eight. The trial used a two-second cadence to observe behavior quickly. It does not validate an optimal production cadence, unattended indexing after the owner exits, or a long-running availability objective. Browser timing is one observation per route/viewport, not a latency SLA.

Recovery preserved graph bytes, private archive permissions and rollback copies on both backends. Small-fixture create/restore timings are retained in [recovery.json](evidence/recovery.json); they are not production RTO. Synthetic trusted-proxy headers test application authorization, not a deployed identity provider.

## Agent and external gates

The installed Codex CLI 0.144.4 was invoked as a test subject with only DiffMind connected, the user's existing model setting, an ephemeral read-only session, and an ordinary dependency question that did not name a tool. Login refresh failed with HTTP 401 before a model turn or DiffMind tool call. [The sanitized result](evidence/model-host.json) records BLOCKED. The user was asked to renew that login. Scripted installed-agent acceptance passes, but it cannot substitute for this actual model-driven trial.

The existing seven open child tasks remain open:

| Task | Remaining evidence |
| --- | --- |
| MNI-180 | Successful unprompted installed-host turns and context selection, after host login is restored. |
| MNI-192 | Authorized actual Enterprise endpoint/version and matching-host credential trial. |
| MNI-197 | Independent blind graph/PR labels and reviewer agreement across a broader supported corpus. |
| MNI-198 | Independently accepted revision/deletion corpus, beyond the controlled cases and regressions. |
| MNI-199 | Human understanding of exact/candidate/unknown evidence and heuristic ordering. |
| MNI-201 | At least three independent observations for each of the original three personas; actual agent adoption. Human participants remain zero. |
| MNI-202 | Those gates plus the exact candidate on Linux ARM64 and both macOS architectures, final acceptance and sustained operational evidence appropriate to deployment. |

Other checks not run here: real screen-reader sessions, additional browser engines, production identity proxy, long-duration load/soak, Ruby formula syntax and systemd unit installation/verification. Ruby is absent and the expected installed service helper is absent. Existing Go maintenance/formula tests passed, but they do not replace those native checks. WSL is a Linux test environment, not a Windows release.

## Reproduce

Use fresh private output directories and the recorded product revision. The [test runbook](../../../testing/runbook.md) and [scenario definitions](../../../testing/scenarios.md) remain the acceptance source. Reusable browser observers require the installed Playwright and axe module paths and Chromium executable through DIFFMIND_PLAYWRIGHT_MODULE, DIFFMIND_AXE_MODULE and DIFFMIND_CHROMIUM.

    go build -trimpath -ldflags '-X main.version=0.0.0-validation -X main.commit=2cd9d328cc24de955184fc8aab863892a6ae91b3' -o /private/diffmind ./cmd/diffmind
    go run ./scripts/enterprise-showcase --home /private/new-scale-home
    node scripts/validation/browser-audit.mjs /private/diffmind /private/new-scale-home /private/new-browser-evidence

    python3 scripts/validation/prepare-boutique.py /private/new-public-trial
    /private/diffmind pack lint testdata/workspace/boutique-correction
    /private/diffmind pack test testdata/workspace/boutique-correction
    node scripts/validation/public-correction.mjs /private/diffmind /private/new-public-trial testdata/workspace/boutique-correction/pack.json

    python3 scripts/validation/flask-pr.py /private/diffmind /private/pinned-flask/examples/tutorial /private/new-pr-trial
    DIFFMIND_HOME=/private/analysis-home /private/diffmind run --repo /private/pinned-flask/examples/tutorial --out /private/flask-artifacts
    node scripts/validation/extractor-browser.mjs /private/diffmind /private/flask-artifacts /private/new-extractor-evidence

    python3 scripts/validation/recovery-upgrade.py /private/prior-diffmind /private/diffmind /private/new-recovery-trial /private/public-trial/repositories

Build the prior e2e9530 implementation in a separate checkout. The extraction observer copies input artifacts into its private output and launches new runs there. The actual upstream PR observer is preserved in harness with this execution's private paths; adjust them for a fresh controlled-PR home before replay. Observe and record the current upstream head if it changes.

Private homes, auth material, backups, raw server logs and cloned source are excluded from the checked-in evidence. Report/Linear commits follow the frozen product commit and do not change the tested implementation.
