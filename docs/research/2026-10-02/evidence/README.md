# Live evidence index

This package supports [the live usability report](../live-usability-study.md). The [reproduction record](../reproduction.md) explains task order, isolation and observer errors. JSON records include timestamps, inputs and returned data; screenshots capture actual browser states.

| File | Evidence supported |
| --- | --- |
| [study-metadata.json](study-metadata.json) | Environment, source/clone revisions, limits, transformations and known observer mistakes. |
| [bootstrap-output.json](bootstrap-output.json) | Successful source bootstrap's emitted agent launch configuration. |
| [agent-observations.json](agent-observations.json) | Real stdio discovery, empty state, candidate preview, initial graph, queries, source change, persistence/lifecycle and incremental refresh. Ends after the observer's wrong tool name. |
| [agent-continuation.json](agent-continuation.json) | Advertised contract schema, successful field comparison and unchanged analysis reuse. |
| [explorer-observations.json](explorer-observations.json) | First-use focus/labels, preview/continuation, invalid-root validation and graph wait failure. The candidate-name logger key was incorrect; see reproduction exclusions. |
| [explorer-continuation.json](explorer-continuation.json) | Confirmed analysis failure, process PATH recovery, readable focused graph, evidence views and comparison/PR empty states. |
| [company-observations.json](company-observations.json) | Corrected final scoped company trial: grants, role denial, HTTP MCP, editor queue, comparison, verified membership removal and token revocation. |
| [company-completion.json](company-completion.json) | Successful later MCP readback of the final completed ingestion. Its unfinished route probe is not product evidence. |
| [company-refresh-job.json](company-refresh-job.json) | Persisted succeeded job associated with that completed ingestion and graph run. |
| [public-pr-observations.json](public-pr-observations.json) | Anonymous public source ingestion and live five-PR listing. Initial PR impact operation was routed incorrectly by the observer. |
| [public-pr-continuation.json](public-pr-continuation.json) | Correct PR #9 impact response, actual displayed text and browser list refresh. |
| [public-graph-provenance.json](public-graph-provenance.json) | Full-service and dependency results tracing queue/example facts to their source locations. |
| [manifest.json](manifest.json) | SHA-256 and size of retained evidence JSON and screenshots; the manifest does not include itself. |
| [harness/](harness/README.md) | Observer script copies documenting the executed calls and selectors. |

## Screenshot navigation

All 22 PNG files in [screenshots](../screenshots/) were copied unchanged. Most informative pairs/states:

- Preview: [after dry run](../screenshots/C-after-dry-run-preview.png) and [invalid-root modal](../screenshots/C-invalid-root-error.png).
- Onboarding and recovery: [first-use modal](../screenshots/C-empty-first-use.png) and [analysis failure](../screenshots/C-ingestion-path-failure.png).
- Readability: [initial fit](../screenshots/C-initial-fit-1440.png) and [gateway search/focus](../screenshots/C-gateway-search-focus.png).
- Evidence: [gateway details](../screenshots/C-gateway-evidence.png).
- Company entry and access loss: [no grant](../screenshots/B-no-project-grant.png) and [open page after removal](../screenshots/B-open-page-after-revocation.png).
- Company credentials: [secret hidden and token metadata](../screenshots/B-access-and-token-metadata.png) and [explicit token revoked](../screenshots/B-token-revoked.png).
- Change interpretation: [five modified graph facts](../screenshots/B-contract-snapshot-comparison.png); field-level compatibility is in the agent JSON, rather than this screenshot.
- Public PR: [selected PR and file evidence](../screenshots/P-real-pr-9-impact.png) and [company context after scrolling](../screenshots/P-real-pr-company-context.png).

## Reading rules

`<study-root>` replaces the disposable absolute directory in JSON. `<source-root>` replaces the working checkout path. Screenshots are unmodified and therefore can display the original disposable path. Only synthetic/local-test identities and public repository material were recorded. Token IDs and lifetime metadata are not token secrets.

Do not treat `elapsed_ms` as human completion time; do not treat observer `study_driver_failure` records as application defects without the report's explanation. The metadata and reproduction file list exclusions. Expected access/validation status codes are part of the test protocol.
