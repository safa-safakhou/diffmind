# Autonomous verification follow-up — 4 October 2026

Codex owns all execution, evaluation and repair. The owner's instruction removes human recruitment, consented sessions and independent human sign-off from acceptance. No human verification is requested or needed. This report supplements the [earlier system validation](../system-validation/README.md); it does not rewrite historical observations.

Three additional product defects were repaired and committed. Ten installed-model verification cases completed successfully, with real tool selection and answers checked against source/runtime expectations. An earlier public-PR attempt timed out and passed on retry. The expanded corpus also retains a genuine coverage limit: a changed Flask form field yields a potential caller instead of exact changed-surface evidence. These results do not establish universal correctness, human satisfaction, adoption or continuous availability.

## Candidate identities

| Evidence | Source | Binary SHA256 |
| --- | --- | --- |
| Earlier broad system validation | `2cd9d328cc24de955184fc8aab863892a6ae91b3` | `6f942fc23e44a7b7c7e79eb668852fc8d6f406e9f3b834133bd29341746fb3de` |
| Repeated browser checks; editor guidance repair | `a7c83692a1c549ae0f50bcc7e440617f34a5165e` | `908f15ee16f3548d89d76a3e8985f40fcaf6780c0ae9967cc81230866aad37d7` |
| Ten final model cases; PR discovery repair | `36db2ab3f800cfc5baba0383d2b1d7ad5c76d958` | `8c7f1f13124dfcb8b731b12204462010d003c98485be1251249d9a5dec6191fe` |
| Latest product; CLI typo repair, full Go/native gates and six-source-case corpus | `66cf19a1e7c86223979bf284b06336117ff3ded6` | `7711abc0bf0ae2ea277e8d957234d4d7f576cfb29824346fd8abcf6a4818f623` |

The last product change only rejects unknown CLI command names before opening a workspace. Agent/server/extraction code is unchanged from the successful model cohort; the UI assets are unchanged from the browser cohort. Subsequent commits contain harnesses or documentation. Each cohort keeps its actual identity; model sessions are not falsely attributed to the later binary.

## Repairs found by actual use

1. **Scoped editor refresh guidance — 011e4ee.** The model selected direct analysis and ingestion, received 403, and reported a capabilities conflict. Editors were correctly authorized for the refresh queue, but the catalog's generic workflow did not make that distinction clear. Catalog, operation descriptions, documentation and denial recovery now name the permitted queue path. An actual HTTP MCP regression follows denial to accepted queue work while preserving configuration denial. The repeated model task selected that path, completed one analysis/five reuses and preserved the source edit.
2. **Read-only PR discovery — 36db2ab.** A model reviewed graph/source information without finding the dedicated PR report, leaving available head eligibility unknown. Initialization and tool descriptions now explain where read-only PR inspection lives. Repeated exact and stale reviews autonomously invoked the PR operation and correctly described revision, caller, uncertainty and delivery.
3. **Unknown CLI command — 66cf19a.** An exploratory `analyze --help` typo followed the default UI path and started a server. The owned process was stopped. Unknown commands now exit 2 before acquiring a workspace lease; documented UI flag shorthand still works. A subprocess regression and the built binary verify that no workspace directory is created.

## Actual installed-model journeys

Final successful cohort: installed Codex CLI **0.160.0**, configured **gpt-6-luna / high**, WSL Linux AMD64. The earlier exploratory record includes the host state observed then. Login now works. No reviewer agents were delegated: Codex CLI was exercised as the DiffMind client under test.

The runner isolates user configuration, disables unrelated apps, keeps shell access read-only and supplies test credentials through an environment variable. Only explicitly authorized disposable-workspace tools are preapproved for mutation trials. The supported per-tool setting is documented in [Codex MCP configuration](https://learn.chatgpt.com/docs/extend/mcp). No user configuration or login credentials were edited.

| Case | Observed result | Session seconds |
| --- | --- | ---: |
| Individual: fresh managed setup | Exactly six reviewed repositories; completed saved graph; manual schedule verified; correct HTTP/Kafka answer | 296.918 |
| Individual: returning, same binary | One changed repository analyzed, five reused; edit and scope preserved; dirty-source limit explained | 86.223 |
| Individual: query-only setup request | No invented project or import; accurate management-capability guidance | 35.904 |
| Individual: local arithmetic control | 391; no DiffMind calls | 5.069 |
| Joiner: viewer | Correct saved relationships; no hidden-project disclosure; refresh/configuration unavailable | 79.214 |
| Joiner: editor | Permitted refresh queue; completed one/five analysis/reuse; configuration still forbidden | 92.008 |
| Joiner: revoked credential | Connection unavailable; no fabricated architecture; existing session also returned 401 | 17.391 |
| Explorer: ambiguous workspace | Asked which of two accessible companies to use | 13.149 |
| Explorer: exact controlled PR | Matching clean head, one GET caller, removed-route uncertainty, no posting claim | 89.355 |
| Explorer: actual stale public PR | Exact head mismatch and dirty caller; candidates distinguished from proven impact; no posting claim | 105.598 |

[Session index](evidence/sessions.json) links the retained result files by case. Each case directory contains prompt/settings, tool calls, critical responses and full final answer. Codex checked the answers against the source/runtime oracles as well as running structural assertions. Session duration is total model-turn time, not a first-useful-result or latency SLA.

Exploratory outcomes are retained, including the initial host profile blocking mutation, a second profile blocking runtime settings, editor failure, PR discovery miss and a **300-second public-PR timeout with no answer**. That identical public task passed with a 600-second timeout on retry. Fresh setup also recovered from an invalid `0s` interval using the documented `0`. This is not a claim of perfect first-attempt model behavior.

Membership and project tokens remain independent grants: removing the membership did not revoke the token; explicit token revocation denied an existing MCP session and a new model connection. The unrelated editor token remained valid. Test identities used real project-token enforcement on loopback, not a production identity provider.

## Source/runtime corpus

Public source: pallets/flask 3.1.2, pinned `2c1b30d0503cfb064f1cb252e6614a06915a362a`. Oracles were written before DiffMind analysis. Flask 3.1.2's actual URL map provides the explicit method/path expectations; implicit HEAD/OPTIONS and Flask's static route are excluded.

- Six actual clean derived heads exercise direct route rename, endpoint removal, internal helper change, application configuration change, form-field change and documentation-only change.
- All six extracted explicit route sets match the runtime URL maps: 12 operations per case, except the removed endpoint case with 11.
- Direct rename yields one exact clean synthetic caller. Removal/helper/configuration/docs controls do not invent exact callers or merge safety.
- **Exactness expectations met: 5 of 6.** The form-field case expected one exact caller but received zero, with that caller retained as a service-level candidate. A runtime test submits the same username/password form: baseline returns 200, changed head returns 400. The unchanged route location does not overlap the changed function-body line. This is a measured contract-analysis limit, not a passed exactness expectation. The original failing receipt and the retained false result remain in evidence; the expectation was not relabeled to manufacture a green corpus.
- Separately, matching positive plus missing patch, unrelated file, deleted line, stale head and dirty caller controls verify HTTP/MCP behavior. Controlled provider responses and synthetic callers are labeled explicitly.
- The real upstream [Flask PR 5918](https://github.com/pallets/flask/pull/5918) was independently fetched. Its observed head is `a82e942870b6472bb40017349cca772c242eb1ad`; the saved tutorial head is different and correctly ineligible.
- Paired enclosing-repository and tutorial-subdirectory extraction each recover the same 12 tutorial operations, with correct respective `examples/tutorial/flaskr/...` and `flaskr/...` source pointers. This checks the labeled tutorial subset, not every test/example in the enclosing repository.
- The synthetic 150-service overview was also exercised with an installed model. Its static/synthetic provenance and unknown current freshness remain explicit. This is separate from real-source graph accuracy.

See [six-case results](evidence/source-corpus-final/results.json), [contract runtime oracle](evidence/source-corpus-final/contract/oracle.json), [scope comparison](evidence/scope-pair/results.json) and [PR controls](evidence/pr-verified/controls.json). Earlier Boutique correction, rollback, source declarations and prior/current recovery evidence remain in the system report.

## Engineering and browser verification

Latest product `66cf19a`: full Go suite (65 tested packages), targeted race suites for CLI/agent API/workspace UI/MCP, and `go vet` pass. The same latest binary passes the native Linux AMD64 archive installer, embedded UI, SQLite, graph/MCP and managed-recovery gate. The invalid-command subprocess test and actual binary both fail promptly without touching a workspace.

Browser cohort `a7c8369`: **22 workspace observations and 18 extractor observations**, including **31 axe scans**, keyboard interaction, desktop/tablet/phone geometry, a real UI-launched analysis, saved CLI history and selected-repository preservation. No browser exceptions or axe violations were reported. Selected phone screenshots were visually inspected. Browser assets did not change in subsequent product commits.

The earlier full race suite (65 packages), 60 frontend tests, dependency audits, public correction/maintenance/rollback and JSON/SQLite restore trials are retained with their original product identity. They were not silently relabeled as fresh executions.

## Current disposition

Linear readback: **26 of 28 child tasks Done**; MNI-192 and MNI-202 remain In Progress. See [scenario deltas](scenarios.md), [53-finding reconciliation](finding-map.json) and [verified issue receipts](linear-receipts.json). Grouping epics retain their separate integrated-candidate requirements.

All verification ownership is automated. The former login and human-participant blockers are removed. The model behavior, runtime-oracle corpus, revision handling and interpretation tasks now have executable evidence. Completion of the corpus task means the corpus and its misses are documented; it does not mean complete framework/contract support.

Technical limits remain: actual Enterprise deployment/SSO, other native architectures/platforms, production identity proxy, other browser engines/screen readers, Ruby/systemd installation, sustained production load/soak, broader language/protocol accuracy and an actual upstream positive exact-caller case. No release has been published or changes pushed. The final candidate remains unreleased.

## Reproduction and evidence handling

Use [installed-host.py](../../../../scripts/validation/installed-host.py) with a JSON spec and an absent output directory; supply the binary, mode, private home or loopback URL, working directory and ordinary prompt. The host chooses tools. For authorized mutation trials use only the needed `approved_tools`; project tokens use `token_env`, never literal values in the spec.

[company-host.py](../../../../scripts/validation/company-host.py) provisions scoped disposable roles and runs viewer/editor/revocation cases. [flask-pr.py](../../../../scripts/validation/flask-pr.py) accepts an optional installed-host runner for exact/stale PR cases. [source-corpus.py](../../../../scripts/validation/source-corpus.py) accepts binary, pinned Flask clone, absent output directory and isolated runtime-package directory. It intentionally retains `exact_expectations_met: false` for the known contract limit.

The [manifest](manifest.json) hashes retained evidence. Private homes, tokens, authentication files, raw unrelated-app transcripts, cloned source, databases and binaries are excluded. Paths in text evidence are replaced by trial placeholders; screenshots show only disposable test data. Read-only GitHub retrieval did not post comments or alter a PR.
