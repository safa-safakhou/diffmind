# DiffMind live usability observation

**Study date:** 2 October 2026. **Product source:** `5d2548514bfc5d920dfd40def5fcacd80162e9ae`. **Method:** live AI-observer walkthrough through the compiled application, browser and MCP interfaces. **Participants:** no recruited humans.

This report supplements the [journey audit](journey-audit.md). The earlier report described what source and documentation imply. This record describes what happened when the product was operated. It documents outcomes and gaps without selecting solutions or adding features.

## Principal assessment

DiffMind's existing paths can deliver useful architecture context. The agent route completed onboarding, built the six-service fixture graph, survived reconnect with saved state, refreshed incrementally and compared contract fields. The shared route delivered the same graph to a permitted viewer, enforced role boundaries and accepted queued refresh work. The PR route retrieved actual public GitHub data and distinguished file-based review signals from stale company context.

The experience still depends on an operator establishing repository scope, runtime environment and refresh responsibility. Live interaction exposed gaps that the source audit alone could not demonstrate: an invisible preview result, lost form state, obscured validation feedback, weak modal keyboard behavior, unreadable initial graph labels and example-source facts attributed to a real repository's service. These observations do not establish how frequently human users encounter them.

## Method, environment and evidence

The observer adopted the task goals of the three personas and performed actual UI clicks, typing, search, keyboard navigation, browser reloads and MCP requests. The same source revision was compiled to a disposable location outside the checkout. Tracked embedded dashboard assets were served without rebuilding them. No alternate implementation or mocked backend replaced DiffMind.

The main browser viewport was 1440 × 900. Supplemental screenshots used 1280 × 720 and 1920 × 1080. The browser was headless Chromium 153.0.8010.12, controlled through Playwright, with Node 24.12.0. Tests ran on Linux/WSL. Company identities were separate browser contexts carrying the documented trusted-proxy headers to a loopback server in scoped mode. This exercised DiffMind after the identity boundary; it did not deploy a proxy, OIDC or production TLS.

Three separate disposable homes held individual-agent, browser-first and public-repository trials. The company trial reused the agent fixture's completed snapshots after its owning agent stopped. The six Demo Shop source repositories were copied and initialized using the repository's existing preparation script. Its existing change script supplied the checkout-v2 contract change. A separate anonymous shallow clone of the public DiffMind repository supplied the real PR trial; its commit matched the reviewed source.

All observation timestamps are UTC. UI screenshots display the browser's local date/time formatting. `elapsed_ms` and `duration_ms` belong to the automated driver and requests. They are not human task-completion times or representative performance benchmarks. The synthetic repositories are deliberately small.

The evidence package contains 143 timestamped events and 22 unchanged PNG screenshots. JSON strings replace the disposable root and source checkout with `<study-root>` and `<source-root>`. Authentication headers and token secrets were never recorded. Token IDs shown in screenshots are metadata rather than bearer secrets. The [evidence guide](evidence/README.md) links the logs; the [reproduction record](reproduction.md) describes execution and observer errors.

### Evidence labels

| Label | Meaning |
| --- | --- |
| Observed | Actual browser, MCP, API, process or persisted-output result recorded in this study. |
| Corroborated | Observed result additionally supported by another interface, a saved artifact or source provenance. |
| Interpretation | Consequence inferred from the observed behavior; not a measured human response. |
| Unexamined | Requires another environment, dataset or participant study. |

High confidence below applies to the bounded observation. A UX interpretation can remain uncertain even when the underlying behavior is directly observed. Consequence labels express relevance to a task, not a vulnerability rating or implementation priority.

## A. Individual developer using an agent

**Task goal:** connect DiffMind, obtain useful repository context, return later and understand a source change. The observer acted as the agent's operator through the actual advertised tool interface. A personal Codex/Cursor configuration was not edited.

| Stage and action | Observed result | What the result means for the journey |
| --- | --- | --- |
| Execute the source bootstrap into a disposable home and binary directory | Bootstrap exited successfully and returned launch configuration containing the installed binary, `agent` argument and explicit home. | A usable connection artifact can be generated. Copying/merging it into an actual agent host remains outside this trial. |
| Connect to the stdio agent and discover tools | Initialization succeeded; 18 tools were advertised. | The product exposes both query and lifecycle/management capabilities. |
| Query an empty installation | Project listing returned empty; a graph request explained that no projects exist and a project must be selected. | Connection alone does not create usable architecture context. |
| Create a project and preview local repositories | Six named candidates were returned; repository listing still contained zero registrations. | Agent preview exposed the scope and did not mutate registration. |
| Start ingestion and poll completion | All six repositories were analyzed; the completed graph had six services and nine relationships. | The existing tool path reached first value without dashboard setup. |
| Query gateway dependencies and checkout's full service evidence | Structured results exposed the graph and source-backed objects. A browser opened the generated dashboard. | Agents and people could inspect the same saved context. |
| Commit the supplied checkout-v2 fixture change, then wait five seconds without refresh | The ingestion run stayed unchanged and contract queries still referred to the old snapshot. | This default local session did not automatically reindex within the observed window. |
| Disconnect the owning MCP session | The old dashboard URL stopped responding. | The observed local backend lifecycle belonged to the agent connection. |
| Reconnect using the same home | Saved graph queries worked; the dashboard had a new loopback port; the old snapshot was still selected. | Persistence succeeded, while reconnect did not imply a fresh analysis under these defaults. |
| Explicitly refresh | One changed repository was analyzed and five were reused. | Incremental reuse worked for this controlled source change. |
| Compare endpoint contracts across the two saved runs | Four changes were returned: three potentially breaking and one compatible. | The contract tool exposed concrete field-level interpretation. |
| Repeat ingestion without another source change | Zero repositories were analyzed and all six were reused. | The unchanged path reused completed analysis. |

Evidence: [agent observations](evidence/agent-observations.json), [correct contract comparison and repeat](evidence/agent-continuation.json), [bootstrap output](evidence/bootstrap-output.json), [first agent graph](screenshots/A-first-agent-graph.png).

The contract comparison reported required additions `deliveryPostalCode` and `market`, removal of `postalCode`, and an optional nullable `deliveryWindow` addition. The first three were potentially breaking; the fourth was compatible. This agrees with the controlled scenario. It does not independently validate every framework's contract extraction or a real organization's compatibility policy.

The change-to-refresh observation is bounded to five seconds and one reconnect with `refresh_interval` empty and `refresh_on_start` false. It supports a lifecycle finding under those defaults, not a claim that configured scheduled refresh cannot work. The source audit's distinction between installed MCP and continuously maintained context remains relevant.

**Perspective assessment:** an agent can complete the path once it knows the appropriate operations and has authority. The developer's trust depends on understanding who initiates refresh and how long the local runtime exists. Engineering correctness was demonstrated for persistence and reuse in a fixture; operational reliability over long sessions remains unexamined. A manager can treat this as evidence of a functioning path, without treating it as evidence of adoption, retention or unattended indexing.

## B. Developer joining an existing company workspace

**Task goal:** enter an existing workspace, read architecture, connect a personal agent with appropriate access, and understand what changes when access is removed. An administrator prepared grants through the live UI; an editor exercised refresh. These were isolated synthetic identities, not actual employees.

| Stage and action | Observed result | What the result means for the journey |
| --- | --- | --- |
| Visit scoped workspace as an ungranted viewer | “No accessible projects. Ask an administrator to grant your user access.” | The initial access boundary had a useful next-action explanation. |
| Administrator adds viewer/editor subjects and saves | The UI saved the membership policy. | Grants could be established using the existing access screen. |
| Viewer reloads and opens the project | The completed six-service graph rendered. | A joining viewer did not need to ingest repositories. |
| Inspect viewer controls and attempt refresh directly | Several operation controls remained visible but disabled; a direct refresh request returned 403, “insufficient role.” | The backend enforced the restriction. The toolbar still displayed unavailable operational vocabulary. |
| Administrator issues a project viewer token | The defaults were viewer role and 30 days. Issuance returned 201; the secret was shown once and then hidden through the live button. | The connection credential had an explicit scope/lifetime and one-time display. No secret appears in the evidence. |
| Connect through HTTP MCP with the viewer token | Initialization and queries succeeded. Tool discovery exposed 13 read-only tools. | An agent could use the existing company graph without local ingestion authority. |
| Editor selects Update graph | The UI request returned 202. Operations displayed the queued work and its running state. | Editor authority reached the durable refresh path. |
| Check the resulting persisted ingestion and refresh job | The same job ID completed successfully, with zero analyzed and six reused; a new graph run was recorded. | The accepted job actually completed. The early Operations screenshot alone shows running state, not completion. |
| Viewer opens the saved before/after comparison | Five modified graph facts were displayed: three checkout objects and two relationships. Repository and pack provenance were visible. | The graph comparison was readable and explicit about saved architectural facts. |
| Administrator removes the verified viewer subject while its graph remains open | After a 4.5-second observation window, zero service nodes remained and the page displayed “Projects / not found.” | Open-page access reacted to grant removal. The message offered less guidance than initial no-access. |
| Use the separately issued viewer token after membership removal | Its graph query still succeeded. | Membership and agent-token grants were independent, as the access screen explains. |
| Explicitly revoke that token | A subsequent MCP tool-list request returned 401. | Explicit token revocation blocked its further use. |

Evidence: [company observations](evidence/company-observations.json), [completed ingestion readback](evidence/company-completion.json), [persisted succeeded job](evidence/company-refresh-job.json), [access metadata](screenshots/B-access-and-token-metadata.png), [viewer graph](screenshots/B-viewer-ready-graph.png), [Operations while running](screenshots/B-editor-operations.png), [comparison](screenshots/B-contract-snapshot-comparison.png), [open page after revocation](screenshots/B-open-page-after-revocation.png).

The viewer and editor toolbars were inspected in the real DOM as well as screenshots. The current experience does not simply remove all unavailable operations. It retains several controls in a disabled state; access administration is not available to the viewer. The report therefore uses “disabled” rather than treating every restriction as a hidden control.

The editor's final job was `refresh-4615f0ec8f9e6df5337ada4f9ae29f84`, associated with graph run `20261002T163243Z`. The retained Operations screenshot was captured before its displayed state updated. Completion is supported by the later ingestion query and saved job record, rather than inferred from that screenshot.

Membership removal was verified against the subject value before the removal button was clicked. An earlier observer attempt removed the wrong sorted row; that attempt was discarded and is not evidence of a revocation defect. The retained final trial shows the correct subject removal and graph disappearance.

**Perspective assessment:** joining works when an administrator has prepared both project access and useful saved data. Security behavior passed the exercised role and revocation cases. Token grants need to be understood separately from person grants, which is already explained in the access screen. The ongoing UI does not explain the loss of access as clearly as first entry. Production identity onboarding, account support and real-team handoff remain unexamined.

## C. Person exploring architecture

**Task goal:** start through the dashboard, choose repository scope, obtain a readable map, inspect supporting evidence and understand what other views mean. This persona traveled through the browser-first route as well as the prepared graph.

| Stage and action | Observed result | What the result means for the journey |
| --- | --- | --- |
| Launch the UI with an empty home | A required first-project modal appeared. Initial focus stayed on BODY; no dialog role was exposed. | The page visibly guided first setup but did not provide the expected modal focus semantics in this test. |
| Press Escape and navigate with Tab | Escape left the required modal open. Tab reached a background New Project button, later the modal fields, then BODY/background again. | Focus was not contained in the first-use modal. Its required nature alone does not establish that Escape should close it. |
| Find the visible Name field by accessible label | No textbox named Name was found, despite visible Name text. | The displayed label did not provide that accessible name. No screen-reader session was performed. |
| Create a project through the form | The project opened successfully. | The browser-first path reached repository selection. |
| Choose Local directory, enter the fixture root, select Dry run only and Preview import | The request returned 200 and count six, then the modal closed. The underlying page still showed zero repositories/no graph and no candidates. | Successful discovery was not exposed as a usable preview result. |
| Reopen import | The root was empty and Dry run only was unchecked. | Continuing from preview required re-entering the scope and preview state. |
| Preview an invalid directory | The API returned 400 with the exact missing-path error. The modal remained open; the error was outside it and obscured by the overlay. | Input validation existed but its feedback was detached from the active task. |
| Import/build while the UI executable was launched by absolute path outside PATH | Registration occurred, but all six analyses failed with `exec: "diffmind": executable file not found in $PATH`. | Launching the dashboard successfully did not establish the worker's executable environment. |
| Restart the same disposable home with the study binary directory on that process's PATH and select Update graph | Ingestion completed and all six service nodes rendered. No product code changed. | The failure was recoverable by satisfying the process environment requirement. |
| Inspect initial fit at 1440 × 900 | The measured gateway label was six pixels high on screen, with scale approximately 0.298. | The initial overview was difficult to read at the captured scale. |
| Search gateway and press Enter | Graph scale became approximately 0.9 and the measured service label became 16 pixels high. | Search/focus provided a successful path to legible detail. |
| Click the gateway node | Details showed service facts and supporting objects. | Architecture exploration reached evidence rather than only decorative nodes. |
| Inspect graph keyboard affordances | The compact checkout service group had no `tabindex` or role; four gateway connection summary rows were focusable. | Keyboard affordances differed between graph elements. This is a bounded DOM/focus observation, not a complete accessibility conformance assessment. |
| Open PR impact for local synthetic repositories | Zero PRs were shown with six repositories checked and “Only repositories with open PRs” selected; the page said no repositories matched the scope. | The empty result did not clearly distinguish local-only repository context from an ordinary filtered PR result. |
| Open comparison with only one graph | The page explained that two completed graphs were needed and that self-comparison was possible. | This was a clear prerequisite message. |
| Capture other viewport sizes | 1280 × 720 and 1920 × 1080 screenshots were saved. | They support inspection of these states, not a mobile/reflow certification. |

Evidence: [first-use and import observations](evidence/explorer-observations.json), [recovery and graph exploration](evidence/explorer-continuation.json), [initial graph](screenshots/C-initial-fit-1440.png), [focused graph](screenshots/C-gateway-search-focus.png), [source evidence view](screenshots/C-gateway-evidence.png).

The PATH condition was intentional: the UI was started using a valid absolute executable path without making that directory available to subprocess lookup. Manual installation instructions can establish PATH, so this is not evidence that every documented installed launch fails. It is evidence that a successful direct launch does not guarantee analysis readiness. The failure and recovery are both retained.

**Perspective assessment:** the explorer can reach supporting evidence and readable focused views. The initial import/preview experience and first fit create friction before that value is visible. Keyboard observations raise accessibility questions that remain relevant even when mouse operation works. Architecture confidence requires understanding both the repository scope and the type of comparison being shown.

## Actual public repository and PR trial

The fixture exercises known outcomes but cannot establish behavior on ordinary public repositories. The observer therefore also cloned [the public DiffMind repository](https://github.com/safa-safakhou/diffmind) anonymously, ingested it through the real agent, fetched the repository's live PR list and opened [PR #9](https://github.com/safa-safakhou/diffmind/pull/9) in both MCP and Chromium. No PR comments, checks or reviews were posted.

At observation time the PR list returned five open PRs. The selected PR was Dependabot's Go dependency update, changing `go.mod` and `go.sum`. DiffMind returned two changed files, 32 additions, 28 deletions, one commit and a low risk score of 19. The browser explained the dependency-file bucket, showed the changed files and described review signals as attention prompts rather than proof of a vulnerability or breaking change.

The company projection was explicitly **stale**, with confidence `stale_graph_estimate` and `score_eligible: false`. It reported zero exact direct and indirect service matches. The graph's analyzed revision was the cloned default-branch commit, rather than the PR head. That is an important successfully exposed limitation: a newly completed graph can be fresh for its checkout while stale for the specific PR being reviewed. The score is an observed heuristic output, not a validated estimate of real operational risk.

The PR's technical graph also included a publish/consume cycle through `orders.created` attributed to `diffmind`. A subsequent full-service and dependency query traced the facts to these source locations:

| Observed fact attributed to `diffmind` | Returned source evidence |
| --- | --- |
| `orders.created` consumer | `examples/demo-shop/repositories/notification/src/main/java/example/OrderNotificationConsumer.java` |
| `orders.created` publisher | `examples/demo-shop/repositories/checkout/src/main/resources/application.yml` |
| Demo Shop HTTP routes and calls | Files under `examples/demo-shop/repositories/{catalog,gateway,checkout,payment,storefront}/` |
| Additional checkout request fields | `examples/demo-shop/scenarios/checkout-contract-v2/openapi.yaml` as well as the example repository's original `openapi.yaml` |

These facts really exist in the source tree. The problem exposed by this trial is their **scope and attribution when the enclosing repository is registered as one service**. Their appearance is not proof that the deployed DiffMind service publishes orders or implements the Demo Shop APIs. The returned source evidence made that distinction investigable; the resulting service and PR context still required interpretation. This is a confirmed case of example/scenario content entering the default analysis, not a measured population-wide false-positive rate.

The enclosing service response contained 97 HTTP route objects and one queue consumer. That count must not be presented as 97 production APIs. The mix includes application source and example content. The company PR projection kept this repository-wide technical context separate from exact changed-surface callers and did not count the stale context toward this PR's score.

Evidence: [public repository/list trial](evidence/public-pr-observations.json), [actual PR response and displayed text](evidence/public-pr-continuation.json), [full service/dependency provenance](evidence/public-graph-provenance.json), [PR file evidence screenshot](screenshots/P-real-pr-9-impact.png), [company context screenshot](screenshots/P-real-pr-company-context.png).

This tested PR is a dependency change. It does not validate breaking API/queue changes across independent real services, a company's full repository estate, automatic PR posting, review accuracy or private GitHub authorization. Fixture contract comparison provides different, narrower evidence.

## Live findings register

These are current-state findings. “Outcome to assess” describes what a user needs to understand or accomplish; it does not select an implementation.

| ID | Observation and confidence | Consequence to investigate | Outcome to assess |
| --- | --- | --- | --- |
| L01 | Dry-run discovery returned six candidates but the dashboard closed without displaying them. **High, observed;** candidate structure corroborated through MCP. | An explorer cannot inspect the discovered scope before committing through that preview flow. | Whether a person can tell which repositories were discovered and make an informed continuation decision. |
| L02 | Reopening import after preview reset the directory and dry-run state. **High, observed.** | Re-entry creates work and may change the scope/state the person intended to continue with. | Whether the continuation preserves an understandable relationship to the previewed scope. |
| L03 | Invalid-root validation returned an exact error outside the still-open modal, visually behind its overlay. **High, observed.** | The person may not see why the action failed or which input requires correction. | Whether validation feedback is visible and connected to the active task. |
| L04 | First-use modal had no dialog role, initial field focus or focus containment; Name lacked the tested accessible name. **High for observed behavior;** broader assistive-technology effects unexamined. | Keyboard and screen-reader users may have a different onboarding path from mouse users. | Whether those users can orient, enter data and remain in the active task. |
| L05 | Absolute-path dashboard launch succeeded while its analysis subprocess failed without `diffmind` on PATH; process-environment recovery succeeded. **High, corroborated.** | A visible running dashboard can be mistaken for full analysis readiness. | Whether the operator knows the actual prerequisites and can recognize the failure boundary. |
| L06 | Default local session preserved the old graph after five seconds and reconnect; disconnect stopped its dashboard, while saved state survived. **High within this window/settings.** | An installation can be mistaken for unattended continuous indexing. | Whether developer and agent understand refresh ownership and local runtime lifetime. |
| L07 | Initial gateway label measured six pixels high; search/focus raised it to 16 pixels. **High for this graph/viewport.** | The overview is harder to read before the useful focus path is discovered. | Whether a person can orient and find a service at first presentation. |
| L08 | The tested UI graph comparison displayed five modified facts; the agent contract tool displayed four field changes with compatibility labels. **High, corroborated across interfaces.** | A person asking “is this API change breaking?” receives different levels of interpretation depending on the entry path. | Whether the person understands which view answers graph change versus contract compatibility. |
| L09 | Local synthetic repositories produced a zero-PR filtered view with a generic scope empty state. **High for this state;** confusion is inferred. | Local-only or unconfigured provider context can resemble a correctly empty GitHub result. | Whether a person can distinguish “no matching PRs” from “this repository context cannot supply PRs.” |
| L10 | Initial no-access guidance explained the need for a grant; live grant removal left a generic not-found page. Access itself was removed. **High, observed.** | Security behavior is correct in the test, but a returning person has less explanation and recovery context. | Whether they understand the change and the appropriate person to contact. |
| L11 | Default public-repository analysis attributed Demo Shop examples/scenarios to `diffmind`, including queue facts in PR technical context. **High, corroborated through source provenance.** | Architecture consumers can read example source as evidence of the deployed enclosing service. | Whether a person or agent can recognize and interpret the represented source scope before making a company or PR decision. |

L01–L05 mainly affect activation and recovery; L06 affects ongoing use; L07–L08 and L11 affect architecture interpretation; L09 affects PR entry; L10 affects company lifecycle. Several original audit concerns now have direct runtime evidence, while others remain documentary or inferential. The two reports should be read together rather than treating every original hypothesis as confirmed.

### Cross-disciplinary reading

| Perspective | Directly supported assessment | What remains uncertain |
| --- | --- | --- |
| UX/product | The existing journeys reach useful outcomes, but preview feedback, form continuation and access-loss explanations have concrete gaps. | Frequency of abandonment, recovery attempts and actual person confusion. |
| Design/information hierarchy | The toolbar contains many operational choices, including disabled controls for viewers. Initial graph labels are much smaller than their focused state. | Which vocabulary and hierarchy different users understand without help. |
| Accessibility | Keyboard focus escaped the initial modal, Name was not available through its tested accessible label, and compact service nodes lacked the inspected focus/role attributes. | Screen-reader behavior, other dialogs, full keyboard task completion, contrast and WCAG conformance. |
| Agent experience | Tool discovery, onboarding, querying, persistence, incremental refresh and contract comparison worked through real stdio MCP. A shared viewer exposed read-only tools. | Behavior of an actual installed client and independently prompted agents without observer task knowledge. |
| Architecture/engineering | Saved context and controlled graph changes were inspectable. Examples/scenarios entered the enclosing public service's analysis. Graph diff and field compatibility were distinct outputs. | Coverage and error rates over diverse, independently labeled repositories and distributed systems. |
| Operations/support | The observed launch-environment failure was diagnosable and recoverable. Queued refresh succeeded and local lifecycle matched the displayed ownership model. | Long-duration availability, failed-job recovery at scale, deployment upgrades, restore and on-call support. |
| Security/governance | Scoped viewer denial, independent token grants, open-page access removal and explicit token revocation behaved as exercised. | Production identity/proxy configuration, adversarial security testing, larger access policies and organizational offboarding practice. |
| Senior management | There is evidence of functioning end-to-end paths and specific decision-sensitive interpretation gaps. | Adoption, willingness to trust output, support cost, retention and business value. These cannot be inferred from a successful automated run. |

## Strengths established live

The study should not be read only as a defect inventory. These useful behaviors were observed:

- The advertised MCP surfaces were usable: local agent discovery provided 18 tools; shared viewer discovery provided 13 read-only tools.
- Agent dry-run candidates were named and non-mutating. The controlled multi-repository graph completed and supported source-backed queries.
- Saved context survived agent disconnect/reconnect. Incremental refresh correctly reused five unchanged repositories after one source change, then six on an unchanged repeat.
- Contract comparison matched the supplied scenario's required additions, removal and optional addition.
- The UI recovered from the documented process-environment condition without a product change, and search/focus made graph detail readable.
- Scoped company access enforced the exercised role boundaries. An editor's queued job actually succeeded; person membership and token revocation behaved distinctly and predictably.
- Snapshot comparison exposed its saved-fact nature and provenance. Its one-snapshot state explained the prerequisite.
- A real public PR loaded in both interfaces, with changed-file evidence, score reasons, stale graph status and explicit cautions about unproven downstream impact.
- Full-service and dependency evidence allowed an unexpected architectural fact to be traced back to example source files.

## Interpretation boundaries and remaining research

No recruited person performed a task in this study. There are no satisfaction scores, task-success percentages, cognitive-load measurements, interview quotes or measured abandonment rates. The observer knew the repository's architecture and could diagnose failures with tooling; ordinary users may behave differently. Headless desktop Chromium is not proof of mobile, native desktop, screen-reader or multi-browser behavior.

The installation trial generated and verified a connection artifact but did not modify a real agent host's settings. HTTP MCP was exercised by a direct protocol client rather than a installed personal agent integration. The company test simulated trusted identity headers after the authentication boundary and did not test production OIDC/TLS. Existing company data was not imported. Six fixture repositories and one public repository/PR are too narrow for universal architecture or compatibility claims.

The public PR data was live at the recorded time and can change. The default branch was ingested; the PR head was not checked out and analyzed. The resulting stale status is reported as observed, rather than hidden behind a “fresh graph” claim. No score calibration or actual dependency-upgrade validation was performed.

Research questions remaining after this work:

1. Can a developer with no DiffMind knowledge reach first usable context through their actual installed agent and explain who maintains it afterward?
2. Can an explorer understand candidate scope, feedback and graph meaning through the UI without access to raw API responses?
3. Can keyboard and screen-reader users complete the same onboarding and evidence-inspection tasks?
4. Can a joining employee connect both browser and agent access without conflating the two grants, and understand later access loss?
5. Across independent real repositories, how often do examples, fixtures, alternate scenarios or multiple applications affect the interpretation of a service's architecture?
6. For known real cross-service API/queue changes, do the warnings agree with independently established affected consumers, and do users distinguish exact evidence from candidates?
7. In a representative company deployment over time, do scheduled freshness, failure recovery and support expectations match what developers and managers think the service provides?

These questions identify missing evidence. They are not a proposed feature backlog or an engineering repair plan.
