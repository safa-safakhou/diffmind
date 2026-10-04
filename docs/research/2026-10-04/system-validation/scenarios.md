# Scenario execution record

Product candidate: 2cd9d328cc24de955184fc8aab863892a6ae91b3. These are the 29 operator scenarios and three persona protocols from [the test pack](../../../testing/README.md). PASS is bounded by the last column. BLOCKED means required evidence remains unavailable; passing subcases are retained without promoting the whole scenario. NOT_RUN means the full case was not executed, even if related regressions passed. None of these statuses means universal correctness.

| Scenario | Overall | Executed evidence | Limit or remaining gate |
| --- | --- | --- | --- |
| A01 | BLOCKED | Native installed bootstrap/management acceptance passes. | Actual model first activation blocked by host login. |
| A02 | BLOCKED | Read-only stdio, managed agent and HTTP MCP mechanisms covered by full/native suites; actual HTTP management regression passes. | Model behavior across all three modes unobserved. |
| A03 | BLOCKED | Scoped hidden-project and selector mechanisms covered by automated suites. | Ambiguous two-project unprompted model choice unobserved. |
| A04 | PASS | Eight-source reuse; one changed source analyzed, seven reused; manual pause; native reconnect/history checks. | Bounded two-second operator trial; no retention or always-on claim. |
| A05 | PASS | Native managed-agent and full race suites cover owner contention, cancel/retry, restart/crash/disconnect and saved-state recovery. | Scripted clients and disposable homes only. |
| A06 | BLOCKED | Installed CLI invoked with an ordinary dependency prompt; login refresh failed before any model/tool turn. | Actual host authentication and unprompted model use. |
| B01 | BLOCKED | Live synthetic viewer/editor/admin/ungranted browser and MCP checks pass. | Actual production identity proxy and deployment. |
| B02 | BLOCKED | Both backends reconcile membership and assigned-token grants independently; unrelated service/admin survives. | Independent joining-developer/support-owner interpretation. |
| B03 | BLOCKED | Live preview has no registration, edits invalidate review, candidate change returns 409, approved cadence works. | Participant understanding of whole-project scope and maintenance. |
| C01 | BLOCKED | Pinned Flask extraction matches 12 explicit runtime operations; structured caller/near-match regression controls pass. | Independent blind graph labels and broader protocol/language corpus. |
| C02 | NOT_RUN | Public extraction uses a real tutorial subdirectory; automated identity/path tests pass. | The complete paired root/subdirectory public corpus case was not executed. |
| C03 | PASS | Full public gap/pack/reuse/provenance/rollback workflow passes; selection regressions included in full suite. | Twelve literal declarations are a subset; no whole-graph accuracy claim. |
| R01 | BLOCKED | Controlled decorator-only PR: two changed entrypoints, one clean exact caller; HTTP/MCP parity passes. | Independent labels and actual upstream positive impact acceptance. |
| R02 | BLOCKED | Controlled stale-head and dirty-caller results reject exactness; existing eligibility regressions pass. | All revision combinations in an independently labeled corpus. |
| R03 | BLOCKED | Controlled deleted-line evidence does not match the unchanged head; deletion regressions pass. | Actual independently labeled removed-endpoint/file corpus and human interpretation. |
| R04 | BLOCKED | Decorator coordinates and existing conservative impact tests pass. | Independent internal/transitive change path labels; no new live helper-change case. |
| R05 | BLOCKED | Public configured-edge provenance passes; compatibility regressions and six-service showcase pass. | Independent real config/contract PR corpus and human interpretation. |
| R06 | BLOCKED | Unrelated filename, method/case/slash/escaped address and protocol/ID regression controls pass. | Complete independent no-impact corpus; no safe-to-merge inference. |
| R07 | BLOCKED | Live missing-patch case rejects exactness; incomplete-list regression coverage passes. | Independent incomplete provider corpus and human score interpretation. |
| R08 | BLOCKED | Actual PR browser shows on-demand, stale/unknown and heuristic context. | Human observations and evidence of any configured posting lifecycle. |
| P01 | PASS | Actual Flask PR 5918 provider, browser and stateful HTTP MCP company parity pass; checkout preserved. | Stale negative case only; head and response recorded. |
| P02 | BLOCKED | Controlled provider tests pass as mechanisms. | Actual Enterprise endpoint/version/authorized credential unavailable. |
| P03 | NOT_RUN | Full provider/component regression suite passes; live public no-token read succeeds. | No new full live invalid-credential/grant/missing-PR recovery matrix. |
| P04 | PASS | Existing controlled provider host/redirect/API-base boundary regressions pass in full suite. | Controlled endpoints and synthetic credentials; not an Enterprise deployment. |
| P05 | PASS | Controlled/public PR observers preserve clean local checkout; background edit stays intact; missing-path regressions pass. | No automatic checkout change or cloning inferred from a PR query. |
| O01 | PASS | Full JSON/SQLite queue/history tests, full races and installed managed recovery gate pass; live unchanged reuse observed. | Small fixtures; no distributed-writer or long-duration load claim. |
| O02 | PASS | Prior e2e9530 to candidate restore at original path passes in JSON/SQLite with rollback and authority reconciliation. | Prior implementation build, not published-release migration certification. |
| O03 | BLOCKED | Same binary passes native Linux AMD64 archive/install/UI/SQLite/company-agent gate. | Linux ARM64 and both macOS architectures; WSL is not Windows release coverage. |
| O04 | BLOCKED | Exact candidate manifest, 53-finding map, scenarios and Linear receipts retained. | External gates plus Ruby/systemd and broader operational acceptance remain open. |
| U01 | BLOCKED | Individual-developer mechanisms exercised through real binaries and browser automation. | Zero independent individual-developer participants; actual model trial blocked. |
| U02 | BLOCKED | Company-joiner role/access/recovery mechanisms exercised on loopback. | Zero independent company-joiner participants; no deployed identity proxy trial. |
| U03 | BLOCKED | Explorer search/graph/inspector/PR/correction mechanisms exercised in Chromium. | Zero independent architecture-explorer participants or screen-reader sessions. |

Evidence is indexed in [the report](README.md) and hashed in [manifest.json](manifest.json). Source tests run on this candidate are discoverable in the full Go/frontend logs; those logs report package/suite outcomes rather than independent human judgments. No new long-duration soak, multi-browser, screen-reader, production load or end-user trial was performed.
