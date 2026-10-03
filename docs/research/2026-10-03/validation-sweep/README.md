# Validation sweep: public source, access, maintenance and recovery

Candidate **e2e9530b06024202fbb9385a0f72a0c84e13d2eb**, built from a clean tree.
Binary and source identities are in [manifest.json](manifest.json). This report
is a later documentation commit; the tested product revision is unchanged.
This is a Linux AMD64/WSL observation, not certification of every native target.
No release, PR comment, public contribution or remote push was performed.

## Outcomes

Five remaining tasks closed with evidence: **MNI-177, MNI-181, MNI-189,
MNI-193 and MNI-196**. Linear now has **21 Done and 7 In Progress** child tasks.
The shared-model, exploration, evidence and company-access epics have all
children complete; install/agent, PR and final validation epics remain open.
[Linear receipt](linear-receipt.json) records verified updates.

| Area | Actual trial and result |
| --- | --- |
| Public corpus | Pinned Online Boutique, eight service sources plus exact corresponding deployment files copied into explicit derived Git corpus directories. Upstream services were not executed. |
| Correction | Baseline: eight services, two source-extracted edges, zero of 12 labeled declarations. Explicit tested pack: 12 of 12 declarations with no extra pack edges. This is declared-subset precision/recall, not whole-graph accuracy. |
| Real bug | UI/API repository selections use pack storage keys; resolver expects manifest IDs. Dotted IDs silently omitted a selected correction. Graph assembly now resolves both forms and reports missing/ambiguous selections. |
| Review/rollback | Actual private gap state transitions, editor reporting without pack authority, positive/near-match negative fixtures, HTTP/MCP pack provenance parity, baseline graph bytes and corrected snapshot preserved after rollback. |
| Scope | Actual browser and MCP use the same request and return the same candidates/digest. Preview registers nothing. Filter edits invalidate review; added candidate returns 409 before registration. Invalid root and include controls pass. |
| Shared access | Legacy access changes only after explicit process restart into scoped mode. Browser shows actual mode and project-wide evidence boundary. Viewer, editor, ungranted and admin recovery are tested across actual browser/HTTP/MCP routes. |
| Maintenance | Clean unchanged corpus reuses all eight artifacts. Approved two-second test cadence analyzes one changed local source and reuses seven, preserves the local edit and eight registrations. Explicit manual pause stops scheduling. Fifteen-minute cadence remains a product default, not an optimized measured interval. |
| Recovery | Prior implementation binary e5c3505 creates graph and private snapshot; current binary restores at the original path with rollback copy. JSON/SQLite both reproduce old token/membership resurrection and independently reconcile them while preserving service/admin access and graph bytes. |
| Native archive | Same candidate passes Linux AMD64 archive installation, embedded UI, SQLite, installed graph/MCP and managed recovery gate. Artifact version 0.0.0-validation is local only. |

[Live results](live-results.json), [source-first labels](source-labels.json),
[recovery results](recovery-results.json), and actual browser screenshots
[graph](corrected-public-graph.png) / [access setup](access-setup.png) are retained.
The graph screenshot shows the readable-scale tradeoff: larger topology extends
beyond the viewport, so pan/search/focus remains necessary. A label-size assertion
does not mean every node fits or establish human comprehension.

## Verification

Full Go suite passes; runmgr/UI/store/MCP race suites pass. Frontend has 20
library and 26 component tests passing. Four exact pack assertions pass.
Native archive gate runs installed-binary company and agent acceptance; these
include cancellation/retry, queue/backup safeguards, controller contention,
maintenance, disconnect, crash cleanup and saved-history reconnect.
Actual browser trial reports zero JavaScript exceptions.

The recovery observer measures create/restore on its small fixture; exact
durations are in recovery-results.json. These are not production RTO guarantees.
Existing interrupted-job recovery tests cover both JSON and SQLite. Archives,
token secrets, raw private server logs and cloned upstream source are excluded.

## Reproduce

Use a supported Node runtime and Playwright/Chromium. Keep all generated homes,
archives and credentials in a fresh private root. Prepare labels before running
the candidate, then explicitly lint/test the evaluation-only pack:

    python3 scripts/validation/prepare-boutique.py /absolute/new/trial
    diffmind pack lint testdata/workspace/boutique-correction
    diffmind pack test testdata/workspace/boutique-correction
    DIFFMIND_PLAYWRIGHT_MODULE=/absolute/node_modules/playwright/index.mjs \
    DIFFMIND_CHROMIUM=/absolute/chromium \
    node scripts/validation/public-correction.mjs /absolute/diffmind \
      /absolute/new/trial /absolute/repo/testdata/workspace/boutique-correction/pack.json

The preparation script creates an absent root, clones the exact public revision,
copies reviewed scope and creates private derived commits for clean/dirty/cache
tests. Original upstream pin/digests remain separately recorded. The observer
creates only isolated homes, uses synthetic trusted-proxy headers, and shuts
down its processes. It is not a production IdP integration trial.

For prior-to-current recovery, build the prior commit in an isolated checkout
and use an absent private output directory:

    python3 scripts/validation/recovery-upgrade.py /absolute/prior-diffmind \
      /absolute/current-diffmind /absolute/new/recovery-trial \
      /absolute/public-trial/repositories

It never connects ordinary ingress. It preserves rollback copies, probes only
loopback, verifies exact paths, and stops every owned server even on failure.

## Remaining tasks and external evidence

| Task | Evidence still required |
| --- | --- |
| MNI-180 | Unprompted context selection in an actual installed model host under T27. Scripted discovery and schema correctness already verified. |
| MNI-192 | Target Enterprise deployment and real matching-host credential trial; controlled host/redirect tests and public GitHub reads are already retained. |
| MNI-197 | Independent blind labeling/agreement, supported source-extracted graph accuracy and actual upstream PR review outcomes. Labels here are source-first by the implementation agent, not a second blind review. |
| MNI-198 | Independent corpus acceptance of PR revision/deletion results; bounded eligibility and conservative unsupported-output tests are already retained. |
| MNI-199 | Independent corpus and human comprehension of evidence/heuristic ordering. |
| MNI-201 | Consented independent persona observations and unprompted installed-agent adoption. Human observations remain zero; none were simulated. |
| MNI-202 | Resolve the above and run the same current candidate on other native platforms, with final all-finding reconciliation. Candidate remains unreleased. |

External deployment/reviewer/participant information was requested. None of the
missing evidence is replaced by green tests, source scan counts or prepared
fixtures. Public declared-configuration corrections do not teach universal
gRPC semantics or prove traffic/merge safety.

## Observer corrections during development

One real product defect was repaired: the storage-key/manifest-ID mismatch.
Separate observer issues were corrected: the positive fixture initially expected
the stanza line instead of the captured value line; the observer initially used
file/line rather than the shared file_scope schema; an ungranted identity probe
initially supplied an invalid global role. These were test assumptions, not
additional product defects. Final reproduction uses the committed observers.
