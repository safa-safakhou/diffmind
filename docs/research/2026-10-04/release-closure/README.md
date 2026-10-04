# Release closure implementation and verification — 4 October 2026

Three product gaps were repaired: handler-body changes now retain exact PR caller evidence, MCP guidance distinguishes saved company context from local files, and native WebKit controls remain readable in the dark theme. The branch is committed and pushed. **This is a verified candidate within the recorded scope, not an unconditional release sign-off:** repeated installed-host discovery is still inconsistent, and an actual Enterprise/production SSO deployment was unavailable.

Codex executed and evaluated all verification. No human participant or human sign-off is required. [Linear readback](linear-receipts.json): **25 of 28 child tasks Done**. MNI-180 was reopened; MNI-192 and MNI-202 remain In Progress. The three related grouping epics remain In Progress. Earlier successful observations are retained as history and do not erase new failures.

## Candidate identities and changes

| Cohort | Source commit | Purpose |
| --- | --- | --- |
| Handler | `cb2ec52b8f45329552638ff0b448ddf550ec828f` | Handler ownership/evidence and changed-line precedence fix; initial actual upstream/model trials |
| Guidance | `19bf84ecd53af779654761fd3eb2484f3609131a` | Includes `a17e2f7` MCP guidance and reusable verification scripts; full local checks and ten model journeys |
| Final product | `63ca9d31fe988da9f49abeb78fd98393005a4671` | Dark native controls in both UIs, rebuilt embedded assets; all final native/browser/backend checks |

[Final binary identity](evidence/candidate-final.json): Linux AMD64, version `0.0.0-validation`, SHA256 `d53ddebf5ce8cc4ea898e4d367f9bde0217c9c039cb001f7a499b0288b888818`. Both guidance and final products were built from clean source. Native CI builds use version `0.0.0-ci` and their own platform binaries. Subsequent documentation commits refer to this frozen product revision.

The final product differs from the guidance cohort only in the two CSS color-scheme declarations and rebuilt UI assets. Backend, CLI and MCP source is identical between those cohorts. Model results retain their actual binary hashes; they are not relabeled as executions of the final binary.

- **Handler-body evidence — cb2ec52.** Framework routes retain registration/decorator identity plus a separately owned handler range and deterministic handler evidence. That evidence survives protocol serialization and workspace hydration. Exact qualified symbols must belong to the binding's indexed file and its owning function/annotation; ambiguous, unrelated and unsupported bindings are rejected. PR matching chooses changed-line evidence over an earlier weaker file-only location. Tests include real Flask AST ranges, unrelated handlers, ambiguous ownership, stale/dirty sources and missing patches.
- **Saved-context guidance — a17e2f7.** MCP initialization and project discovery explain that company architecture can exist even when the local directory is empty. Dependency/change questions should discover saved context; unrelated arithmetic should remain local. Repeated trials show this improves guidance but does not guarantee host adoption.
- **Live verification tools — 19bf84e.** Added provider, background/proxy, actual-upstream-PR and systemd scripts; both browser harnesses now support Chromium, Firefox and WebKit.
- **Native control readability — 63ca9d3.** WebKit rendered white native selects with pale text. Both UIs now declare dark color-scheme. This was found by screenshot inspection despite zero axe violations: [before](evidence/webkit-before.png), [final workspace](evidence/release/workspace-webkit/phone-graph-selection.png), [final extractor](evidence/release/extractor-webkit/run-form-390.png).

## Verification results

| Area | Executed result | Evidence and limits |
| --- | --- | --- |
| Final product CI | All five jobs passed | [Exact-revision run](https://github.com/safa-safakhou/diffmind/actions/runs/37222725816), [safe job receipts](evidence/ci.json), [test log excerpt](evidence/ci-test-excerpt.txt) |
| Native distribution | Linux AMD64/ARM64 and macOS Intel/ARM64 passed | Actual native installer, embedded UI, SQLite, graph/MCP and managed-recovery checks; host-model trials remain Linux AMD64 |
| Engineering checks | Full Go and full race suites, go vet, 60 frontend tests, both builds, pack/distribution/install/demo/agent-setup checks, Ruby syntax, vulnerability checks passed | [Full local make verify log](evidence/verified/verify.txt), final CI; local full Go/race run is guidance cohort; final CI runs full Go, targeted races, vet and frontend/distribution checks on final product. Zero vulnerabilities reported at execution time |
| Browsers | Chromium 149.0.7827.55, Firefox 151.0, WebKit 26.5; 120 observations, 93 axe scans, zero reported violations or page exceptions | [Browser summary](evidence/browser-summary.json); keyboard/dialog focus and restoration, graph activation, desktop/tablet/phone geometry, short-phone inspector, actual extraction actions |
| Runtime/source corpus | All six explicit route-map and caller expectations pass | [Results](evidence/release/corpus/results.json); pinned public Flask and controlled mutations, one synthetic literal caller per case |
| Actual upstream PR | Flask PR4139 at its actual clean head: one exact caller | [Results](evidence/release/upstream-pr/results.json); GitHub data is real, caller is synthetic |
| Provider behavior | All ten HTTP/stateful-MCP cases pass | [Results](evidence/release/provider/results.json); controlled compatible endpoint, not an actual Enterprise installation |
| Background/access | 62 completed jobs, 7,008 read requests across JSON and SQLite, restart recovery and permission checks pass | [Results](evidence/release/background/results.json); 120 seconds per backend, eight readers, six small repositories |
| Native systemd | Active app resumes active; inactive app remains inactive after backup | [Results](evidence/release/systemd/results.json); isolated units using shipped helper/protections, private verified archives, retention and marker cleanup |
| Upgrade/restore | Prior implementation 66cf19a to final product, JSON and SQLite pass | [Results](evidence/release/recovery/results.json); previous graph bytes, independent grants, reconciliation, nonoverwrite and live-lock safeguards |

The workspace browser fixture contains 150 **synthetic** services from the prior audit. The extractor browser reads historical saved CLI output and launches new analyses. This exercises current UI behavior; it is not a new accuracy measurement for every fixture edge. New final-product extraction accuracy is measured separately by the runtime corpus. Playwright browser engines on Linux do not establish physical iOS/Safari or screen-reader behavior.

The background trial performed unchanged-source reuse and one-changed/five-reused refreshes while readers queried. It retained approved scope, local edits and previous graph bytes across a graceful restart. Peak running and queued jobs per project were each one. Measured p95 read time was 28.966 ms for JSON and 28.715 ms for SQLite; these are observations, not service-level guarantees.

A controlled proxy injected verified subject/role and stripped spoofed identity headers: viewer reads returned 200, forged admin writes 403, direct untrusted identity 401. Removing membership denied both the browser and an existing MCP session. Separate project-token revocation cases preserve the distinction between membership and independent credentials. This does not qualify a real IdP/SSO installation.

## PR evidence and runtime controls

Public Flask 3.1.2 is pinned at `2c1b30d0503cfb064f1cb252e6614a06915a362a`. Source/runtime oracles were frozen before analysis. Explicit URL-map agreement is 12 operations for direct rename, helper, configuration, form-field and documentation cases; endpoint removal has 11. Implicit HEAD/OPTIONS and the framework static route are outside that labeled subset.

The earlier form-field failure is repaired **without changing the expected exact-caller count**. Changing the login handler from username to account_name leaves the route registration unchanged. The same runtime form returns 200 at baseline and 400 at the changed head. DiffMind now reports the two methods sharing that changed handler and the one matching POST caller. This is handler-body changed-surface evidence, not a complete form-schema compatibility analysis or proof that every caller fails.

Removal, indirect helper and configuration changes remain conservative: no exact caller is invented when supporting changed-surface evidence is missing. Deleted baselines and arbitrary cross-file call propagation remain explicit limitations. The additional [eligibility controls](evidence/verified/pr/controls.json) cover missing patch, unrelated filename, deleted line, stale head and dirty caller.

The actual upstream [Flask PR4139](https://github.com/pallets/flask/pull/4139), “Avoid race condition in example app,” was checked at `51196575479e34c3ea43612e1eb770db3aa5d114`. Actual GitHub PR files and the clean matching checkout identify register's changed body and one exact synthetic POST caller. Local source was preserved; no comment, check, review or merge action was posted.

Two host-generated code-review claims were independently executed after the model answers:

- **PR4139: false finding.** The handler-cohort host claimed that `sqlite3.Connection.IntegrityError` does not exist and recommended requesting changes. It does exist. Actual tutorial execution returned 302 for first registration and 200 with the duplicate-user validation message for the second. [Disproving runtime control](evidence/upstream-pr/runtime-control.json). The [retained model answer](evidence/upstream-pr/model/answer.txt) contains that disproved claim and must not be treated as a valid review. DiffMind's deterministic caller evidence was correct; this additional host inference was not.
- **PR5918: confirmed bounded finding.** The guidance-cohort host correctly classified the saved company graph as stale relative to actual head `a82e942870b6472bb40017349cca772c242eb1ad`, and did not promote its dirty caller to exact impact. Its separate claim that the routes CLI mutates the method set was confirmed by executing that actual head: the GET rule changed from GET/HEAD to GET/HEAD/OPTIONS after the command. [Runtime control](evidence/verified/pr/public-runtime-control.json). The first exploratory probe accidentally selected the distinct automatic-OPTIONS rule; the retained check explicitly selects the GET rule.

These are post-hoc controls of host assertions, not pre-frozen corpus oracles. Neither review is a merge-safety guarantee.

## Installed agent results

Actual installed Codex CLI **0.160.0**, user's **gpt-6-luna/high** configuration, Linux AMD64. The runner disables unrelated app/MCP configuration, keeps shell commands read-only and authorizes only the required disposable-workspace management tools. No model agents were delegated to implement work; Codex CLI is the client under test.

The guidance cohort's ten final answers met their bounded DiffMind expectations:

| Journey/case | Observed behavior | Seconds |
| --- | --- | ---: |
| Individual: fresh managed setup | Six approved local repositories, complete graph, accurate HTTP/Kafka evidence, requested manual schedule retained | 135.610 |
| Individual: returning | One analysis/five reuses; edit, scope and manual schedule preserved | 164.513 |
| Individual: query-only setup | Correctly explains missing management capability; no invented import | 39.321 |
| Individual: local arithmetic | 391, no DiffMind calls | 3.698 |
| Joiner: viewer | Correct evidence, no hidden project, no refresh/configuration authority | 66.965 |
| Joiner: editor | Permitted queue completes one/five; configuration still denied | 159.216 |
| Joiner: revoked token | No fabricated architecture; existing MCP session denied, unrelated editor remains valid | 24.389 |
| Explorer: ambiguous | Asks East or West Commerce | 47.992 |
| Explorer: exact controlled PR | One exact GET caller, clean matching head, uncertainty retained, no posting | 68.822 |
| Explorer: stale actual PR5918 | Stale/dirty limitations retained; additional code claim independently checked above | 111.498 |

[Session index](evidence/sessions.json) and [per-session assessments](evidence/session-assessments.json) retain prompts, tool choices, shell commands, binary identity, timing and final answers. Exit zero alone is not a pass. Time is the entire model turn, not first useful result or an SLA.

**Repeatability failed.** Across the main ambiguous case and three additional identical vague prompts (“What services could be affected by changes to checkout? Explain the evidence.”), two asked for project selection, one gave clearly separated East/West results without asking first, and one used no MCP tools and incorrectly concluded that an empty local folder meant no evidence. The separated answer did not silently select the wrong company, but missed the strict ask-first expectation. The earlier handler cohort also had an empty-folder discovery miss.

Two final-product variants explicitly saying “Using the saved company architecture” both asked for project selection (30.991 and 62.471 seconds). That is a different prompt, not evidence that the vague prompt now always works.

**Evaluation boundary:** these are executable workflow observations, not blind filesystem-isolated evaluations. Some exploratory sessions inspected adjacent harness filenames; repeated ambiguity case 2 read neighboring specs, and explicit-context case 1 read its own spec/transcript. Their answers therefore cannot independently establish uncoached adoption. These shell actions are retained. The discovery failure and the false upstream review remain visible rather than being discarded as outliers. MNI-180 is reopened for this measured reliability gap.

## Current disposition and reproduction

[Scenario deltas](scenarios.md) and the [complete 53-finding map](finding-map.json) distinguish new evidence from inherited historical checks. No claim of exhaustive language, protocol, dynamic-runtime or production-scale correctness is made.

- **MNI-180 — In Progress:** intermittent discovery of saved MCP context under ordinary dependency prompts.
- **MNI-192 — In Progress:** controlled/public provider behavior verified; actual Enterprise endpoint and production SSO qualification unavailable.
- **MNI-202 — In Progress:** final candidate has broad executable evidence, but the above gaps and host review reliability prevent an unconditional all-clear.
- Long production soak, physical devices/screen readers and exhaustive extraction coverage remain unexecuted boundaries. Automatic upstream PR delivery remains outside this existing on-demand product behavior. No release is published.

Reproduction entry points are [source-corpus.py](../../../../scripts/validation/source-corpus.py), [upstream-pr.py](../../../../scripts/validation/upstream-pr.py), [provider-matrix.py](../../../../scripts/validation/provider-matrix.py), [background-matrix.py](../../../../scripts/validation/background-matrix.py), [systemd-native.py](../../../../scripts/validation/systemd-native.py), [installed-host.py](../../../../scripts/validation/installed-host.py), [company-host.py](../../../../scripts/validation/company-host.py), and the two browser scripts in the same directory. Each script documents its arguments. Browser selection uses `DIFFMIND_BROWSER=chromium|firefox|webkit`; `DIFFMIND_CHROMIUM` only overrides Chromium's executable. Use disposable output directories, authorized test identities and pinned public source.

Exploratory harness errors were repaired before the retained passes: pagination initially matched page=1 inside per_page=100; the background harness initially omitted the MCP initialized notification and used an admin-only viewer probe; systemd initially looked for the archive at the wrong depth. The first handler draft retained locations but lost them during evidence hydration; the committed fix retains explicit handler evidence and the corpus was rerun. Those exploratory failures are not presented as successful product executions.

The [manifest](manifest.json) hashes retained files. Private raw transcripts, credentials, workspace databases, archives, source clones and binaries are excluded; text paths use trial placeholders and trailing whitespace is normalized. Original oracle digests identify pre-redaction trial files; the manifest separately hashes the retained sanitized bytes. Screenshots show disposable data. Native systemd trial units were removed. Feature-branch pushes triggered verification, not release publication.
