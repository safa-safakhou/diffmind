> This is the source and documentation audit completed before live interaction. See [the live usability study](live-usability-study.md) for subsequently observed behavior. Product source revision remains the same; documentation was added afterward.

# Research scope and principal assessment

Research date: **2 October 2026**. Repository reviewed: **DiffMind**, source revision **5d2548514bfc5d920dfd40def5fcacd80162e9ae**. This report maps the existing product and its experience boundaries. It contains findings, uncertainty and desired outcomes; it does not select solutions, add features or prescribe engineering work.

**Principal assessment:** DiffMind already has a substantial foundation for an evidence-based, persistent architecture workspace. It can be operated entirely through an agent after bootstrap. The experience is nevertheless conditional on several responsibilities being fulfilled: someone selects repository scope, initiates ingestion, checks extraction coverage, establishes a refresh policy, and maintains the relevant runtime. Installing an MCP connection alone is not evidence that those responsibilities have been fulfilled.

The three journeys share the same graph but have different readiness requirements:

| Person and purpose | What makes the experience usable | Main unresolved experience boundary |
| --- | --- | --- |
| Individual developer using an agent | Working client connection, authorized repositories, completed graph, understood freshness policy | What the agent does automatically after installation and after later code changes |
| Developer joining a company workspace | Existing service, separate browser and agent credentials where needed, project grants, administrator-prepared data | Whether “I can access the workspace” also means “I have the context and authority needed for my task” |
| Person exploring architecture | Accessible graph, readable scope, understandable evidence and coverage | Whether the visible map represents enough of the real system to support the intended decision |

These are overlapping roles. An individual may explore the dashboard; a company viewer may use an agent; an administrator may travel through all three. Installation, administration and architecture exploration should therefore not be mistaken for equivalent tasks.

The largest findings concern **responsibility and lifecycle clarity, interpretation of incomplete architecture, company access boundaries, and the distinction between on-demand PR impact and automatic PR review**. UX and accessibility findings add concrete concerns around first-use control, crowded operations vocabulary and dialog interaction. None of these conclusions establishes that all users experience the same friction.

Evidence for the core product and the separate connection modes: [R01 Product promise and agent first entry](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91); [R04 Connection types and lifecycle](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1); [R12 Tested analysis coverage and limits](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R16 Access modes permissions and trust boundaries](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3).

# Method and evidence discipline

The review covered current entry documentation, installation scripts, local lifecycle management, MCP schemas and management discovery, project/import/dashboard components, graph and query behavior, PR projection and scoring, refresh events, authorization, credentials, recovery, knowledge packs, controlled acceptance fixtures and public compatibility records. Two checked-in dashboard screenshots were inspected for visual hierarchy. Public research used official product documentation, published research and W3C/MCP guidance.

**Audit-phase scope:** No source, configuration, client registration, credentials or existing workspace state was modified. The report was originally saved separately from the repository and subsequently copied here at the user's request. No real company was imported. The later live study used disposable runtime state and is documented separately.

Evidence labels used throughout:

| Label | Meaning | What it does not establish |
| --- | --- | --- |
| V | Directly inspected code, file, screenshot, or automated behavior exercised in this conversation | Real-user success, production reliability or complete framework coverage |
| D | Behavior described by repository documentation or a vendor’s public documentation | Independent verification of every claim or every deployment |
| I | Inference from an inspected fact or a usability principle | A measured rate of confusion, failure or business impact |
| H | Hypothesis or important unanswered question | A confirmed defect or an observed user complaint |

Confidence applies to the stated observation. Consequences may remain inferred even when the underlying observation has high confidence. “High consequence” below identifies a decision-sensitive area, not a numerical risk estimate or vulnerability classification. A documented boundary is not automatically a bug.

**Automated evidence already collected at this same revision during this conversation:** targeted Go suites for workspace UI/backend, graph resolution, graph assembly, queries, discovery, connections and protocol passed. The real stdio `TestAgentAcceptance` passed and reported onboarding, graph queries, incremental reuse, pack installation, recovery, backup, SQLite queue migration, scheduling and reconnect coverage. These checks used disposable test state. They were not a test of an actual installed desktop client, a production deployment or user behavior. Test definitions: [R24 Real stdio agent acceptance](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/cmd/diffmind/agent_test.go#L27); [R25 Real binary company fixture acceptance](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/company_acceptance_test.go#L30).

The repository also records broader earlier validation, including full/race/vet/frontend and native checks, with explicit exclusions. Those are **D evidence from the historical checkpoint**, not a claim that all of them were rerun in this research audit. [R23 Existing validation record and exclusions](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1).

**Audit-phase limits:** no participant interviews or usability sessions; no live browser interaction or screen-reader exercise; no actual client registration; no new deployment/restore drill; no independent competitor trials; no measured retention, setup time, false-positive rate, false-negative rate or production scale. Source inspection can demonstrate a missing dialog behavior or a default setting. It cannot demonstrate how often that affects people. Public documentation was checked on the research date and may change. The subsequent [live study](live-usability-study.md) adds actual browser, MCP and PR observations while retaining explicit limits of its own.

Internal references R01–R38 link to the inspected commit, with entry line anchors where useful. External citations link to primary sources. A reference anchor is a starting point for inspection, not a claim that one line contains the whole behavior.

# The existing experience model

The practical lifecycle is:

**Authorize scope → install and connect → identify workspace and project → preview or import repositories → synchronize and analyze → assemble and persist a graph → inspect or query evidence → compare changes or inspect PR impact → maintain freshness → recover, hand off or leave.**

Several distinct states matter:

- **Installed:** the executable is present.

- **Connected:** the host client can discover and call the intended tools.

- **Registered:** repositories are known to the workspace.

- **Accepted:** an asynchronous job has been admitted.

- **Analyzed:** repository artifacts have been produced, possibly with omissions or failures.

- **Queryable:** a completed graph snapshot exists.

- **Useful for this question:** that snapshot includes relevant relationships at an appropriate revision.

- **Maintained:** someone or something refreshes it under an understood operating policy.

DiffMind already represents many of these states. The experience gap is that they do not collapse into one guarantee simply because the MCP server appears connected. Saved query tools read completed snapshots. They do not implicitly build a new graph from current edits. [R11 Project resolution and completed snapshots](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R27 Management discovery and async workflow](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agentapi/api.go#L132); [R38 Ingestion completion contract](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/workspace/README.md#L35).

| Connection | Practical authority and lifecycle | Experience implication |
| --- | --- | --- |
| Local `agent` | Full local management; client starts DiffMind and DiffMind owns a backend while its controller lives | The host agent can complete onboarding, but must actually perform the workflow |
| Local `mcp` | Query-only access to saved graphs, trusting the OS user | A working connection cannot import or maintain the graph |
| Shared HTTP MCP | Rights of authenticated identity; viewers query, eligible editors/admins manage within policy | The organization’s service and access setup become prerequisites |
| Dashboard | Human inspection and permitted operations on the same workspace | It can be optional for the solo agent path and primary for an explorer |

Local agent refresh is disabled by default: the default settings do not supply an interval or refresh-on-start. The playbook configures refresh when requested. The backend stops when the owning local connection ends; data and history persist. Shared Compose defaults instead request refresh on startup and every 15 minutes. These are **verified different operating contracts**, not inconsistent execution of one default. [R03 Local runtime defaults](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agenthost/host.go#L27); [R02 Agent bootstrap scope and onboarding](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R15 Shared deployment defaults](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/compose.yaml#L19).

The architecture is deterministic static analysis. A relationship may be source-extracted, resolved through identity conventions, or declared through a knowledge pack. These are useful forms of evidence with different meanings. The graph is not a runtime traffic inventory, and a missing edge does not prove that two services never communicate. [R12 Tested analysis coverage and limits](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R22 Convention teaching provenance and integrity](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91).

# Journey A Individual developer using an agent

**Person:** a developer who wants their usual coding agent to understand relationships across their authorized repositories. Their preferred routine is to work normally and receive useful context without repeatedly operating another tool.

**Observable end state:** the agent can answer a real cross-repository question with source evidence, explain the relevant snapshot and uncertainties, and remain useful after changes and reconnects.

Stages A01–A11 cover discovery through exit. The journey is possible today; the degree to which it happens automatically in a real client is unmeasured.

## A01 Understand the promise

**Current path:** README positions DiffMind as an architecture graph for humans and coding agents, with an agent-first request and optional dashboard.

- **UX:** The developer must recognize whether this complements their code search/memory tools and what information it adds.

- **Design and information:** The leading story is clear; depth, supported patterns and the PR experience require following additional material.

- **Agent experience:** The host agent is an operating participant, not merely a query consumer.

- **Architecture and engineering:** The graph is deterministic and does not require an LLM key for extraction; the host agent may still use a model.

- **Management and product:** A broad company overview is an aspiration bounded by actual detector and identity coverage.

- **Operations and trust:** Expectations about source access, local storage and downstream evidence sharing begin before installation.

**Evidence:** [R01 Product promise and agent first entry](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91); [R12 Tested analysis coverage and limits](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3) V/D for the product description; I for possible expectation mismatch.

## A02 Choose an installation route

**Current path:** The source bootstrap is the documented agent-first route. A checksummed release installer and native macOS/Linux releases also exist. AGENT_SETUP still says no release exists.

- **UX:** The person or agent can encounter incompatible answers to “what do I install?”

- **Design and information:** Two legitimate routes become confusing when the playbook contains obsolete availability advice.

- **Agent experience:** The default playbook requires the host to build source and handle prerequisites, despite an available release route.

- **Architecture and engineering:** Source builds need Go, Git and C/CGO; release binaries reduce those prerequisites. Native Windows binaries are not offered.

- **Management and product:** Actual acquisition friction and unsupported-platform demand are unknown; a documentation contradiction is confirmed.

- **Operations and trust:** Executable installation alone neither grants repository access nor configures the client or scope.

**Evidence:** [R01 Product promise and agent first entry](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91); [R02 Agent bootstrap scope and onboarding](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R32 Release installation behavior](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/install.sh#L1); [R33 Source bootstrap implementation](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/scripts/agent-setup/main.go#L1); [published v0.1.0 release](https://github.com/safa-safakhou/diffmind/releases/tag/v0.1.0). V/D high-confidence contradiction.

## A03 Connect the host client

**Current path:** The bootstrap emits absolute paths and MCP configuration. The host registers it, reconnects and verifies discovery. Full management uses agent; the manual guide also presents query-only mcp.

- **UX:** The meaningful milestone is usable tools, rather than a copied configuration snippet.

- **Design and information:** Commands that look similar conceal different lifecycle and management capabilities.

- **Agent experience:** The host must complete client-specific registration; initial permission/reload remains outside DiffMind’s control.

- **Architecture and engineering:** Tool count and runtime status provide a checkable connection contract. Actual client variations were not exercised here.

- **Management and product:** Time to a working connection and host-specific success rates are not established by protocol acceptance tests.

- **Operations and trust:** One lifecycle controller owns a home; a second agent needs a compatible sharing or separate-workspace arrangement.

**Evidence:** [R02 Agent bootstrap scope and onboarding](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R04 Connection types and lifecycle](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1); [R05 Manual query only agent setup](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/personal-setup.md#L176); [R24 Real stdio agent acceptance](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/cmd/diffmind/agent_test.go#L27) V/D for modes; H for client-level completion rates.

## A04 Identify workspace and project

**Current path:** The full agent can list projects and create one if needed. Graph calls can use a configured or sole accessible project; multiple projects require selection.

- **UX:** The user may think in company, repository collection or current checkout rather than project ID.

- **Design and information:** Workspace home, project name, project ID, repository and run are distinct concepts the agent must translate.

- **Agent experience:** Discovering an existing project is necessary before creating or choosing one; a lost create response cannot be retried blindly.

- **Architecture and engineering:** Project selection is explicit when ambiguous. A chosen default is a routing convenience, not a local authority boundary.

- **Management and product:** Accumulating duplicate or unused projects is a plausible lifecycle problem, not a measured one.

- **Operations and trust:** Personal and company state can be separated by home; ownership and future handoff still need to be understood.

**Evidence:** [R02 Agent bootstrap scope and onboarding](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R11 Project resolution and completed snapshots](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R27 Management discovery and async workflow](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agentapi/api.go#L132) V for resolution and workflow; I for mental-model burden.

## A05 Authorize and discover repositories

**Current path:** The playbook asks the agent to preview authorized local roots or GitHub scope, inspect the result and then ingest. Include/exclude filters and limits exist.

- **UX:** A user must be able to understand what will be included without becoming an import operator.

- **Design and information:** Repository scope, clone choice and credentials are different decisions; a preview gives concrete scope evidence.

- **Agent experience:** The agent must resolve intended paths/organization, review the preview and avoid guessing authorization.

- **Architecture and engineering:** Existing local checkouts and managed clones have different synchronization behavior.

- **Management and product:** The cost and usefulness of a company overview depend on scope quality, not just repository count.

- **Operations and trust:** Credentials belong in approved handling and the backend’s environment; authorized access remains a prerequisite.

**Evidence:** [R02 Agent bootstrap scope and onboarding](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R04 Connection types and lifecycle](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1); [R14 Company deployment and operations](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1) V/D; H for accidental overscope frequency.

## A06 Wait for the first usable graph

**Current path:** Ingestion is asynchronous and may complete, partially complete or fail. The agent polls persisted status and then queries known dependencies.

- **UX:** The person needs a clear answer about readiness without mistaking acceptance for success.

- **Design and information:** Progress, failure details and graph availability convey different information; completed processing does not measure semantic coverage.

- **Agent experience:** The management workflow explicitly tells the agent to poll, inspect failures and verify evidence.

- **Architecture and engineering:** Incremental reuse, cancellation and durable state support repeatable work; query tools require a completed graph.

- **Management and product:** The first meaningful answer is a more relevant activation milestone than installation or imported repository count.

- **Operations and trust:** A partial run can leave useful data, but its omissions need to accompany any conclusion based on that data.

**Evidence:** [R27 Management discovery and async workflow](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agentapi/api.go#L132); [R38 Ingestion completion contract](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/workspace/README.md#L35); [R24 Real stdio agent acceptance](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/cmd/diffmind/agent_test.go#L27); [R25 Real binary company fixture acceptance](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/company_acceptance_test.go#L30) V for tested lifecycle; H for actual time to value.

## A07 Use context in ordinary coding

**Current path:** The agent can list services, inspect dependencies, query evidence, trace and compare contracts. Compact results are default; richer evidence is requested explicitly.

- **UX:** The hoped-for experience is helpful context during normal tasks, with minimal prompting about DiffMind.

- **Design and information:** Users need understandable explanations rather than unexplained tool names or raw IDs.

- **Agent experience:** Tool availability does not guarantee the host chooses the right tool, follows pagination or requests full evidence before making a claim.

- **Architecture and engineering:** Schemas expose scope/run and bounds; results describe the extracted architecture rather than arbitrary current-file search.

- **Management and product:** No observed data shows how often agents invoke the tools, or whether answers improve compared with an agent alone.

- **Operations and trust:** Source evidence sent to the host can enter the host’s model context; local extraction alone does not establish the whole data path.

**Evidence:** [R26 Agent schemas descriptions and graph tools](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35); [R11 Project resolution and completed snapshots](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R12 Tested analysis coverage and limits](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3) V for tool contracts; H for natural adoption and answer quality.

## A08 Understand a planned change or PR

**Current path:** Saved-run comparisons and on-demand GitHub PR impact use graph context and limited changed-surface evidence. They are separate surfaces.

- **UX:** The developer must know whether they are asking about a saved change, a working edit or a particular PR revision.

- **Design and information:** Risk labels can appear more definitive than a heuristic score and partial coverage justify.

- **Agent experience:** The host can interpret supplied evidence, but is not guaranteed to prepare the matching snapshots or publish a PR review.

- **Architecture and engineering:** Exact PR attribution is revision-sensitive; general dependency candidates are kept separate from proven changed-endpoint callers.

- **Management and product:** The company needs clarity about whether this aids investigation, informs review or is intended as a gate; current evidence does not validate gate accuracy.

- **Operations and trust:** The inspected webhook triggers push refresh, not automatic PR comments/checks or a persisted review conversation.

**Evidence:** [R19 Live PR projection scoring and revision checks](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24); [R20 Push refresh event contract](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/webhook.go#L67); [R21 Bounded contract comparison](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/contracts.go#L95) V high confidence for current boundaries; I for interpretation risk.

## A09 Continue after edits and reconnects

**Current path:** Graphs persist. The local backend is session-owned; refresh has no interval by default. Reconnect restores a backend, not an implicit promise to reindex all edits.

- **UX:** The developer’s “it just works in the background” expectation is conditional on refresh and ownership policy.

- **Design and information:** Persistent data, a running process and current context must be distinguishable in explanations.

- **Agent experience:** The host may need to notice stale context, initiate permitted work and wait before answering a revision-sensitive question.

- **Architecture and engineering:** Saved snapshots protect history; managed clones track configured branches, while local checkouts are analyzed in place.

- **Management and product:** Repeat usefulness after the first day is unmeasured, including maintenance effort and silent staleness.

- **Operations and trust:** Closing the owning client changes service availability; shared deployment has a different continuity contract.

**Evidence:** [R03 Local runtime defaults](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agenthost/host.go#L27); [R04 Connection types and lifecycle](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1); [R11 Project resolution and completed snapshots](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R14 Company deployment and operations](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1); [R15 Shared deployment defaults](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/compose.yaml#L19) V/D; H for perceived background behavior.

## A10 Encounter omissions or failures

**Current path:** The product offers errors, status/history, retries, doctor, backup and knowledge-pack teaching. Pattern limits are documented.

- **UX:** A missing relationship can mean unsupported detection, unresolved identity, excluded scope, failed work or genuinely absent evidence.

- **Design and information:** Recovery language must convey which state occurred and who can act without suggesting that empty results are proof of safety.

- **Agent experience:** The playbook instructs the host to inspect evidence and avoid inventing edges; teaching requires constrained examples and validation.

- **Architecture and engineering:** Pack provenance/integrity and immutable history support correction while preserving earlier results.

- **Management and product:** The effort required to adapt an unfamiliar company convention is not measured.

- **Operations and trust:** Recovery should preserve valid history and source checkouts; existing safeguards help, but real-user comprehension remains unknown.

**Evidence:** [R02 Agent bootstrap scope and onboarding](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R12 Tested analysis coverage and limits](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R22 Convention teaching provenance and integrity](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91); [R23 Existing validation record and exclusions](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1) V/D for mechanisms; H for correction burden.

## A11 Upgrade move or stop using DiffMind

**Current path:** Release/source installation, offline recovery and separate homes exist. No unified end-user journey for retirement or migration was demonstrated.

- **UX:** The person needs to understand which connection, process, executable and persisted data are part of their installation.

- **Design and information:** Hidden background state and visible client registration can outlive different portions of use.

- **Agent experience:** The host must distinguish reconnect, upgrade, stop, backup and restore; these have different effects.

- **Architecture and engineering:** Restore is bounded to workspace files and paths; it is not general application-schema or checkout migration.

- **Management and product:** Lifecycle ownership and retention decisions matter when a personal experiment becomes shared company infrastructure.

- **Operations and trust:** Restored credentials reflect snapshot-era authority; external local source layouts may need separate protection.

**Evidence:** [R18 Recovery scope and restored authority](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/backup-recovery.md#L35); [R23 Existing validation record and exclusions](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1); [R04 Connection types and lifecycle](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1) D/V for recovery boundaries; H for migration and exit experience.

# Journey B Developer joining an existing company workspace

**Person:** a developer who did not install or curate the company service. They expect to receive access and use company context in their normal work.

**Observable end state:** their browser and agent reach the intended workspace with appropriate rights, query trustworthy context, understand their limits and know who owns recovery or missing coverage.

The company operator is a backstage participant in this journey. The joining developer’s convenience depends on deployment, credentials, access policy, source availability and graph preparation already being handled.

## B01 Receive an entry point

**Current path:** The organization provides a service URL and an authentication arrangement. Deployment and reverse-proxy identity are operator responsibilities.

- **UX:** The joining user needs to know what service this is and whether it is already ready for use.

- **Design and information:** Entry information needs to distinguish service location, browser login and agent connection.

- **Agent experience:** The host cannot infer a company endpoint or identity from a local MCP installation.

- **Architecture and engineering:** DiffMind supports one continuously refreshed server, rather than providing an organization signup/control plane.

- **Management and product:** Access handoff and ownership of the service are organizational dependencies.

- **Operations and trust:** TLS, trusted-proxy configuration and protecting direct backend access belong to the deployment contract.

**Evidence:** [R14](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3) D/V for deployment responsibilities; H for handoff quality.

## B02 Authenticate and reach the right projects

**Current path:** Scoped access is opt-in. In scoped mode explicit grants and a global role ceiling determine accessible projects. Legacy mode retains global-role behavior.

- **UX:** A successfully authenticated user may still see no projects.

- **Design and information:** The UI’s no-access state explicitly asks for an administrator; it is different from an administrator’s empty workspace.

- **Agent experience:** The agent must list what this identity can access rather than reuse an unrelated user’s project assumptions.

- **Architecture and engineering:** Policies filter listings before pagination; inaccessible IDs return a non-disclosing response and invalid policies fail closed.

- **Management and product:** Service rollout is not complete until ordinary identities can actually use intended projects.

- **Operations and trust:** Default legacy mode and shared administrator credentials have broader authority than project-scoped membership.

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R06](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/Projects.jsx#L16); [R15](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/compose.yaml#L19) V/D high confidence for access behavior; I for onboarding dependency.

## B03 Connect a company agent

**Current path:** Remote agents can use proxy-authenticated requests or separately issued project viewer/editor tokens. Browser cookies alone do not configure an agent; DiffMind does not implement OAuth login.

- **UX:** The browser and agent may each require successful authentication.

- **Design and information:** The distinction between a personal identity, independent service grant and recovery token is essential to comprehension.

- **Agent experience:** Tool discovery and capability depend on current request credentials; a viewer connection remains query-only.

- **Architecture and engineering:** Scoped MCP is stateless per request and rechecks authorization, avoiding retained session privilege across identities.

- **Management and product:** Agent onboarding may create ongoing administrator work; the volume and delay are unknown.

- **Operations and trust:** A server administrator token is not an ordinary joining credential. Tokens expire and are independently revocable.

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63); [R30](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectTokens.jsx#L1) V/D for credential behavior; H for onboarding effort.

## B04 Orient within company scope

**Current path:** The developer selects among authorized projects, teams, services and repositories. A sole accessible project can be the MCP default.

- **UX:** The user must recognize whether the available map covers their work or just a subset.

- **Design and information:** Project labels, service names and team groupings become the navigation vocabulary.

- **Agent experience:** The host should resolve the intended project explicitly when multiple options exist.

- **Architecture and engineering:** A project is an authorization and graph boundary; team filtering is a view boundary inside it.

- **Management and product:** Company structure, software ownership and imported repository organization are not automatically the same model.

- **Operations and trust:** A project grant exposes all project data, including source evidence and configuration; it does not inherit per-repository code-host permissions.

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R11](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R09](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCanvas.jsx#L1218) V/D for boundaries; I for orientation burden.

## B05 Establish readiness and coverage

**Current path:** A prepared graph can be queried immediately. If none exists, viewers depend on an administrator or refresh worker.

- **UX:** The developer should know whether empty information is missing access, missing ingestion or limited extraction.

- **Design and information:** Completion/freshness indicators are useful but cannot communicate semantic completeness on their own.

- **Agent experience:** The host can inspect summaries and known relationships before trusting the workspace for a task.

- **Architecture and engineering:** Known supported patterns are bounded. Imported language inventory does not imply relationship coverage for every framework.

- **Management and product:** A company-wide inventory claim requires evidence against the company’s known topology, beyond successful imports.

- **Operations and trust:** Source credentials, branch configuration and failed repositories affect the available snapshot.

**Evidence:** [R11](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R13](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63) V/D; H for readiness of an actual company workspace.

## B06 Use context while developing

**Current path:** The developer or host queries the same saved graph and evidence as other authorized users.

- **UX:** Cross-team answers can reduce manual investigation if the graph contains the relationships needed.

- **Design and information:** Teams, evidence and graph direction need understandable meaning; a dependency path is not automatically a required notification list.

- **Agent experience:** The host can combine graph context with current code, but the integration does not prove it will do so correctly.

- **Architecture and engineering:** Completed-run selection permits reproducible answers, while default latest selection can change between calls.

- **Management and product:** Actual reduction in investigation time, interruptions or incorrect assumptions has not been measured.

- **Operations and trust:** Project membership permits reading stored evidence even if the user lacks corresponding upstream repository permission.

**Evidence:** [R26](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35); [R11](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3) V for contracts; H for productivity outcomes.

## B07 Investigate a change or review a PR

**Current path:** PR impact overlays live GitHub details on saved graph context. Saved graph/contract comparison is available separately.

- **UX:** The developer needs to know which PR revision and graph revision the result describes.

- **Design and information:** Proven changed-surface callers and potential service-level candidates have different decision weight.

- **Agent experience:** The host should preserve that distinction when explaining a result; tool availability alone cannot guarantee it.

- **Architecture and engineering:** Public GitHub endpoints are hardcoded in the inspected PR handler, despite configurable API bases in import workflows.

- **Management and product:** A company may expect a review inside its code host; an on-demand workspace analysis has a different adoption path.

- **Operations and trust:** PR access depends on provider credentials and availability. An absent finding is not a merge-safety certificate.

**Evidence:** [R19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24); [R14](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1); [R21](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/contracts.go#L95) V for PR contract; I for expectation mismatch and enterprise-provider concern.

## B08 Need newer or more complete context

**Current path:** Scoped editors can queue refresh using saved configuration. Viewers cannot refresh; imports, paths, packs and configuration require a global administrator.

- **UX:** Read access does not mean the user can repair missing context themselves.

- **Design and information:** Role-aware controls hide unavailable operations, while no-data states identify administrator dependency.

- **Agent experience:** The host must discover actual capability and avoid trying to solve every omission through a forbidden management action.

- **Architecture and engineering:** Scheduled work and user-triggered queued work share durable state and controlled concurrency.

- **Management and product:** Support responsiveness becomes part of the user experience when only an administrator can alter scope or conventions.

- **Operations and trust:** Already-admitted jobs and scheduled refresh have independent service authorization.

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R28](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/operations.md#L19); [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246) V/D for capability boundaries; H for support delays.

## B09 Observe progress and recover from failure

**Current path:** Operations exposes jobs, attempts, queue/cancel/retry behavior and project resource limits. Per-repository outcomes and historical records persist.

- **UX:** The user should distinguish waiting, active work, partial completion and terminal failure.

- **Design and information:** Operational detail exists; the amount a non-operator must understand remains a usability question.

- **Agent experience:** The host can inspect persisted status rather than assuming acceptance equals readiness.

- **Architecture and engineering:** One project’s work is serialized; global budgets and queue capacity constrain progress across projects.

- **Management and product:** Reliable status is a strength, while actual queue wait, resource contention and support workload are unmeasured.

- **Operations and trust:** Metrics and fleet-level administration have broader authority requirements than reading one project’s status.

**Evidence:** [R28](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/operations.md#L19); [R27](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agentapi/api.go#L132); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3) V/D for operations; H for service-level experience.

## B10 Return over weeks and organizational changes

**Current path:** Shared Compose schedules refresh by default. Repository and identity policies still require operator maintenance; token expiry can interrupt agent access.

- **UX:** The ordinary experience should be predictable across repeated visits, branch changes and role changes.

- **Design and information:** Freshness, access loss and service failure are different states that should not be explained as one generic outage.

- **Agent experience:** The host needs to re-resolve credentials and capability after changes rather than retain earlier authority assumptions.

- **Architecture and engineering:** GitHub push webhooks refresh tracked configuration; they do not subscribe the user to automatic PR reviews.

- **Management and product:** Repository churn, ownership changes and maintenance cost are unknown in a real company pilot.

- **Operations and trust:** User memberships and independently issued agent tokens have separate lifecycles.

**Evidence:** [R15](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/compose.yaml#L19); [R20](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/webhook.go#L67); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63) V/D; H for long-term continuity.

## B11 Leave the company or survive service recovery

**Current path:** Membership changes stop new user access. Independent tokens require separate revocation; already-admitted work continues. Offline restore reinstates snapshot-era access state.

- **UX:** Offboarding should have a clear end state for both person and agent.

- **Design and information:** Token UI documentation already explains revocation limitations; comprehension and operator execution are unmeasured.

- **Agent experience:** A client holding an independent token can retain its grant after browser/user membership changes unless that token is revoked.

- **Architecture and engineering:** Backup restores data and authorization records together; it cannot know later revocations absent from the snapshot.

- **Management and product:** Offboarding and disaster recovery are organizational responsibilities, not only storage events.

- **Operations and trust:** Downloaded evidence cannot be retracted. These documented boundaries are not evidence of an unauthorized-access defect.

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63); [R18](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/backup-recovery.md#L35); [R30](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectTokens.jsx#L1) V/D high confidence for authority lifecycle; H for operational follow-through.

# Journey C Person exploring architecture

**Person:** an architect, senior engineer, engineering manager or developer who wants to understand a system before changing it. They may never use an MCP client or run an import.

**Observable end state:** they can find an area, follow relevant relationships, inspect why those relationships exist, recognize missing/stale information and communicate a defensible conclusion.

This journey evaluates exploration and interpretation. Administrative import is included only where the explorer is also the workspace owner.

## C01 Choose a way to explore

**Current path:** The README presents an optional local dashboard, public fixtures and a shared deployment route.

- **UX:** The person needs to distinguish a ready company map from a demonstration or an empty installation.

- **Design and information:** Demo Shop, real compatibility scans and the synthetic enterprise map are explicitly separated.

- **Agent experience:** An agent can prepare context and provide a dashboard URL, but exploration can also be entirely human-led.

- **Architecture and engineering:** Generated navigation data, source extraction and runtime observations establish different kinds of evidence.

- **Management and product:** The relevant success is answering an architecture question, not simply seeing a rendered graph.

- **Operations and trust:** A safe public fixture proves presentation or controlled behavior without exposing company data.

**Evidence:** [R01](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91); [R34](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/enterprise-showcase/README.md#L1); [R13](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30) V/D; I for expectations at entry.

## C02 Enter a project

**Current path:** A ready user selects a project. An authorized empty local workspace automatically opens a first-project dialog, defaults its name to DEFAULT and prevents closing until a project exists.

- **UX:** The explorer can be required to define a container before learning the product’s concepts.

- **Design and information:** The first form includes project root/instruction options and uses internal project vocabulary.

- **Agent experience:** A host could perform this preparation; the dashboard alone makes it a human responsibility.

- **Architecture and engineering:** Read-only company users instead see a grant-related empty state and cannot create projects.

- **Management and product:** A forced creation path may reduce abandonment by removing choices or increase it by reducing control; neither effect has been measured.

- **Operations and trust:** The no-project error path also sets an empty project list, although an error is shown; empty and failed discovery can overlap visually.

**Evidence:** [R06](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/Projects.jsx#L16) V high confidence for control flow; I/H for usability consequence.

## C03 Prepare scope when no graph exists

**Current path:** Administrators can add repositories or import an organization/local root. The import dialog defaults to GitHub, no dry run, and sync/analyze/build enabled.

- **UX:** Repository discovery and interpretation require different skill levels.

- **Design and information:** Filters, team, limit, depth, API base and execution options are exposed during preparation.

- **Agent experience:** The agent playbook recommends preview before import, whereas the UI does not make preview the initial default.

- **Architecture and engineering:** Import registration and actual graph construction can be combined or separate.

- **Management and product:** First useful architecture and setup effort are coupled when the explorer is also the administrator.

- **Operations and trust:** An explicit limit and clear scope are especially consequential for large organizations; actual accidental scope errors are unknown.

**Evidence:** [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246); [R02](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34) V for defaults; I for cognitive/informed-scope burden.

## C04 Read readiness and quality

**Current path:** The workspace shows ingestion outcome, repository freshness, processing state and graph completion, with expandable errors and reuse details.

- **UX:** The explorer needs to know whether the displayed information is suitable for their question.

- **Design and information:** Green completion communicates process success; semantic quality requires additional interpretation.

- **Agent experience:** An agent can supplement the display with known-relationship verification, but that is not an observed dashboard user behavior.

- **Architecture and engineering:** Failed, stale, dirty and unknown repository states have distinct implications.

- **Management and product:** Successful jobs do not establish accuracy or completeness of the company overview.

- **Operations and trust:** Existing aria-live/status surfaces are a positive accessibility practice; their actual announcements were not exercised.

**Evidence:** [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R25](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/company_acceptance_test.go#L30) V/D; I for interpretation risk.

## C05 Orient on the graph

**Current path:** Overview/full detail, service search, zoom and team filtering are available. Large workspaces start in a team scope and state visible versus total service count.

- **UX:** The explorer can navigate a bounded area before expanding to the portfolio.

- **Design and information:** Checked-in screenshots show consistent dark hierarchy and resource shapes, but small card text at whole-view fit and a crowded upper command strip.

- **Agent experience:** Queries provide an alternative route to the same concepts, with different information density.

- **Architecture and engineering:** Team filtering excludes nodes from view; it is not a declaration that cross-team dependencies are absent.

- **Management and product:** Generated 150-service navigation is useful evidence of presentation, not a real company extraction or latency benchmark.

- **Operations and trust:** Browser size, zoom, long names, low vision and keyboard exploration remain untested.

**Evidence:** [R09](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCanvas.jsx#L1218); [R36](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/demo-shop-graph.jpg); [R37](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/enterprise-overview.png); [R34](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/enterprise-showcase/README.md#L1) V for screenshot/source observations; H for task success.

## C06 Follow a relationship and inspect evidence

**Current path:** Selecting services, resources, edges and traces opens details including extracted objects and evidence. Some lists have explicit display caps.

- **UX:** The person must understand direction, relationship type and whether it is extracted or declared.

- **Design and information:** Detailed provenance is available; interpreting source facts still requires domain understanding.

- **Agent experience:** The host can request full evidence and pagination, whereas the dashboard has its own inspection limits.

- **Architecture and engineering:** Local static traces and service graph paths are not complete end-to-end runtime traces.

- **Management and product:** Evidence enables reviewable conclusions; the ease with which non-experts can inspect it is unmeasured.

- **Operations and trust:** Some graph connection rows support keyboard activation; inspected compact service nodes do not show equivalent tabindex/key handlers.

**Evidence:** [R10](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphDetails.jsx#L25); [R09](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCanvas.jsx#L1218); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R26](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35) V for implementation; I for comprehensibility and reachability gaps.

## C07 Compare architecture over time

**Current path:** The comparison UI selects completed runs, shows changes and before/after evidence, and supports run-pair links. Contract differences have bounded compatibility labels.

- **UX:** The explorer can investigate what changed between snapshots.

- **Design and information:** Run IDs, dates, content provenance and semantic changes must be read together.

- **Agent experience:** The host can compare pinned runs; it must not equate a graph difference with a proved behavioral regression.

- **Architecture and engineering:** Analyzer, pack, scope and identity changes can alter a graph even when runtime behavior is unchanged.

- **Management and product:** Historical comparison is a strength for explaining decisions, but stakeholder comprehension is unmeasured.

- **Operations and trust:** A selected snapshot is reproducible; a link interpreted as an always-current state would have a different meaning.

**Evidence:** [R29](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCompare.jsx#L1); [R21](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/contracts.go#L95); [R22](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91); [R11](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117) V/D; I for causal interpretation risk.

## C08 Inspect a PR overlay

**Current path:** The explorer opens PR impact and obtains live PR metadata, code heuristics, graph candidates, changed surfaces, revision context and notes.

- **UX:** The person must see whether the result is based on the PR head and what is merely a candidate.

- **Design and information:** Scores and confidence terms can influence interpretation beyond the evidence supplied.

- **Agent experience:** An agent explanation should preserve the notes and uncertainty instead of compressing everything into safe/unsafe.

- **Architecture and engineering:** The current projection does not itself persist review decisions, issue code-host checks or automatically construct base/head analyses.

- **Management and product:** The intended decision weight of a risk score is unresolved without calibration and outcome evidence.

- **Operations and trust:** Provider errors, truncated file retrieval and a stale graph affect what the result can support.

**Evidence:** [R19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24); [R20](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/webhook.go#L67) V for boundary; H for how managers/reviewers use scores.

## C09 Form and communicate a conclusion

**Current path:** Evidence details and saved-run comparison can support a discussion. This audit did not verify a complete share/export/read-receipt workflow for arbitrary graph selections.

- **UX:** The explorer needs to state what they learned with scope and caveats intact.

- **Design and information:** A diagram alone may lose provenance when detached from its project, run and filters.

- **Agent experience:** The host can produce a narrative, but correctness depends on preserving direction, limits and evidence.

- **Architecture and engineering:** A graph expresses observed/declarative architecture; it does not independently prove ownership, business criticality or runtime usage.

- **Management and product:** Management decisions need an explicit basis and confidence, rather than an attractive portfolio diagram.

- **Operations and trust:** Evidence sharing remains subject to project/data permissions and the receiving context.

**Evidence:** [R10](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphDetails.jsx#L25); [R29](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCompare.jsx#L1); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3) V/D for evidence; H for sharing and decision quality.

## C10 Return correct or request help

**Current path:** The explorer may inspect newer graphs, ask an editor to refresh, ask an administrator to repair scope, or use packs for declared conventions.

- **UX:** The person needs to know who can act on an apparent wrong or missing relationship.

- **Design and information:** Operational controls, content correction and read-only exploration currently coexist in one workspace interface.

- **Agent experience:** The host can help distinguish unsupported extraction from service failure, when its access and instructions permit.

- **Architecture and engineering:** Teaching can change evidence/identity interpretation and is therefore part of graph governance.

- **Management and product:** Ongoing trust depends on successful correction and repeated useful conclusions; neither has a measured adoption baseline here.

- **Operations and trust:** Recovery, access changes and service ownership must remain understandable even when the explorer is not the operator.

**Evidence:** [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R22](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91); [R28](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/operations.md#L19) V/D for mechanisms; H for feedback and long-term trust.

# Findings register

This register separates confirmed inconsistencies, explicit product boundaries, interpretation risks and missing validation. It is an inventory for understanding the present state, **not a ranked implementation backlog**. Journey references point to the stage maps above.

A finding can be important because a decision depends on it, even when it is a deliberate boundary rather than a defect. Questions identify missing evidence; they do not authorize product changes.

## Entry and agent lifecycle

### F01 Agent setup has obsolete release advice

**Stages:** A02. **Assessment:** Confirmed documentation contradiction · V/D · high.

**Observation:** AGENT_SETUP says no published binary exists; README and the public release page identify v0.1.0.

**Consequence:** An agent following the authoritative playbook can choose an unnecessarily demanding acquisition path or explain release availability incorrectly.

**Unanswered question:** How often do host agents rely on the playbook rather than reconcile other entry documentation?

**Evidence:** [R01](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91); [R02](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R32](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/install.sh#L1).

### F02 Background maintenance is conditional in local mode

**Stages:** A09 B10. **Assessment:** Documented operating boundary · V/D · high.

**Observation:** Default local settings supply no refresh interval or refresh-on-start. The playbook configures refresh if requested; the owned backend stops at disconnect.

**Consequence:** The user's install-once background expectation is not established by connection alone.

**Unanswered question:** What freshness/availability does the person believe they have after closing or reconnecting their agent?

**Evidence:** [R03](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agenthost/host.go#L27); [R02](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R04](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1).

### F03 One product exposes distinct agent contracts

**Stages:** A03 B03. **Assessment:** Experience hypothesis grounded in verified modes · V/I · medium.

**Observation:** agent, local mcp and shared HTTP differ in management authority and runtime ownership. Manual setup documents query-only integration.

**Consequence:** A connected server can be interpreted as capable of onboarding or updating when it is only capable of querying saved data.

**Unanswered question:** Do users and hosts identify their actual mode before deciding what the integration will handle?

**Evidence:** [R04](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1); [R05](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/personal-setup.md#L176); [R26](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35).

### F04 Tool discovery does not establish adoption

**Stages:** A07 B06. **Assessment:** Validation gap · H · unresolved.

**Observation:** Tool descriptions and schemas are present; the MCP server construction does not add general server instructions. Agent acceptance exercises intended calls.

**Consequence:** These facts do not demonstrate that an unprompted host uses the right tools or preserves uncertainty during ordinary work.

**Unanswered question:** What proportion of relevant real tasks receives useful graph context, and what causes missed or unnecessary calls?

**Evidence:** [R26](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35); [R24](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/cmd/diffmind/agent_test.go#L27).

### F05 Bootstrap success remains a host boundary

**Stages:** A02 A03. **Assessment:** Validation gap and external dependency · V/D/H.

**Observation:** The installer emits configuration; the host must register/reconnect. Current acceptance does not execute a real installed client’s configuration journey.

**Consequence:** A complete backend workflow can coexist with incomplete client setup.

**Unanswered question:** Where do real client onboarding attempts fail, and can the person recognize successful discovery?

**Evidence:** [R02](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R33](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/scripts/agent-setup/main.go#L1); [R23](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1).

### F06 Platform availability is narrower than generic agent availability

**Stages:** A02. **Assessment:** Explicit compatibility boundary · V/D · high.

**Observation:** Native binaries support macOS/Linux amd64/arm64; source requires Go/Git/C. Native Windows binaries are not provided.

**Consequence:** Some developers cannot follow the same acquisition path as colleagues. WSL is a Linux environment, not evidence of native Windows support.

**Unanswered question:** How much of the intended audience encounters unsupported OS, toolchain or policy constraints?

**Evidence:** [R01](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91); [R32](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/install.sh#L1).

### F07 Multiple agents share a lifecycle constraint

**Stages:** A03 A09. **Assessment:** Explicit architecture boundary · V/D · high.

**Observation:** Only one local lifecycle controller owns a home. Additional agents can share its HTTP endpoint while it lives or use different state.

**Consequence:** Parallel sessions or changing clients can create ownership and availability questions for the user.

**Unanswered question:** Do multi-client users understand which connection owns availability and whether their state is shared?

**Evidence:** [R02](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R04](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1).

### F08 Project selection requires a domain translation

**Stages:** A04 B04 C02. **Assessment:** Experience hypothesis · V/I · medium.

**Observation:** Projects, repository collections, service identities, homes and run IDs are separate. Ambiguous graph queries require a project ID.

**Consequence:** Users working in a current checkout may not immediately understand the intended architecture scope.

**Unanswered question:** Can each persona identify the correct project without importing duplicates or answering from the wrong company context?

**Evidence:** [R06](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/Projects.jsx#L16); [R11](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R26](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35).

### F09 Import scope has different defaults by surface

**Stages:** A05 C03. **Assessment:** Verified default difference with inferred consequence · V/I · medium.

**Observation:** The playbook previews authorized scope first; the UI initializes dry-run false and pipeline execution true.

**Consequence:** Administrators may act on broader scope before inspecting it. No accidental import was observed.

**Unanswered question:** Can a person accurately predict repository count, exclusions and cloning/analysis effects before committing the import?

**Evidence:** [R02](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34); [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246).

### F10 Local and managed sources have different freshness meanings

**Stages:** A05 A09 B05. **Assessment:** Documented operating boundary · D/V · high.

**Observation:** Existing local checkouts are analyzed in place. Managed clones synchronize configured branches and refuse dirty overwrites.

**Consequence:** Refreshing a local source and refreshing a managed company clone do not imply the same source update.

**Unanswered question:** Does a freshness explanation identify the source layout, branch and revision relevant to the question?

**Evidence:** [R04](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1); [R14](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3).

### F11 Accepted work is not usable context

**Stages:** A06 B09 C04. **Assessment:** Verified state distinction · V · high.

**Observation:** Management returns asynchronous acceptance; ingestion can be completed, partial or failed. Queries require completed graphs.

**Consequence:** Premature readiness claims can conceal missing work, although the playbook explicitly warns against them.

**Unanswered question:** Do real agents and human users wait for terminal state and inspect partial outcomes?

**Evidence:** [R27](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agentapi/api.go#L132); [R38](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/workspace/README.md#L35); [R11](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117).

### F12 Latest completed graph is a snapshot fallback

**Stages:** A07 A09 B06. **Assessment:** Verified query boundary · V · high.

**Observation:** Default graph loading selects a completed persisted run, skipping incomplete runs. It does not rebuild current source on demand.

**Consequence:** A useful older map can remain available during failed/new work, but must not be interpreted as the latest source state.

**Unanswered question:** Can a person tell which snapshot answered the question and whether newer incomplete work matters?

**Evidence:** [R11](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117); [R26](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35).

## Coverage and architecture interpretation

### F13 Process completion does not establish semantic completeness

**Stages:** A06 B05 C04. **Assessment:** High consequence interpretation gap · V/D/I.

**Observation:** A repository can scan successfully while detectors miss wrappers, dynamic configuration or unsupported conventions.

**Consequence:** A green completion badge or zero findings can be overread as evidence of complete architecture or low change risk.

**Unanswered question:** Do users distinguish processing health, source freshness and relationship coverage?

**Evidence:** [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R13](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30).

### F14 Public compatibility scans expose limited detection coverage

**Stages:** B05 C01. **Assessment:** Recorded benchmark boundary · D · high for recorded counts.

**Observation:** The September 10 record reports architecture entities in 4/32 OpenTelemetry targets and 6/12 Boutique targets, with zero extractor connections in both.

**Consequence:** This is meaningful missing-coverage evidence for those pinned scans, not an error-rate estimate or full company graph result.

**Unanswered question:** Which known relationships are absent after full graph assembly, and what is their decision consequence?

**Evidence:** [R13](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30).

### F15 Synthetic scale is separate from real extraction scale

**Stages:** C05 B09. **Assessment:** Evidence boundary · V/D · high.

**Observation:** The 150-service enterprise fixture supplies generated topology and metrics for navigation.

**Consequence:** It demonstrates a large rendered model, not measured precision, production load or extraction of a 150-service company.

**Unanswered question:** How does an equally large real extracted graph differ in readability, latency and completeness?

**Evidence:** [R34](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/enterprise-showcase/README.md#L1); [R37](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/enterprise-overview.png); [R13](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30).

### F16 Declared and extracted relationships have different meanings

**Stages:** A10 C06 C07. **Assessment:** Evidence interpretation boundary · D/V · high.

**Observation:** Packs can add identities and declared relationships. Their declarations do not prove a runtime call or entrypoint reachability.

**Consequence:** Treating all graph edges as equivalent proof can overstate what was observed.

**Unanswered question:** Can people and hosts explain the provenance and practical uncertainty of an individual edge?

**Evidence:** [R22](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3).

### F17 Teaching is an ongoing governance activity

**Stages:** A10 B08 C10. **Assessment:** Operational burden hypothesis · D/I · medium.

**Observation:** Pack rules, tests, integrity and provenance exist; custom conventions are not automatically understood merely because the repository is indexed.

**Consequence:** A company may need an owner for convention quality and correction, beyond an initial installation owner.

**Unanswered question:** How much effort, expertise and repeated maintenance are required for the intended repository mix?

**Evidence:** [R22](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91); [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3).

### F18 Static source is not runtime truth

**Stages:** A07 C06 C09. **Assessment:** Explicit product boundary · D/V · high.

**Observation:** The supported-pattern contract defines deterministic static analysis; configuration does not necessarily establish the effective runtime destination.

**Consequence:** Architecture conclusions need to remain bounded when traffic, deployment profiles or dynamically resolved services are involved.

**Unanswered question:** Which decisions are safe to make from available source evidence alone?

**Evidence:** [R12](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3); [R22](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91).

## Dashboard interaction and accessibility

### F19 Whole-view graph labels can become visually small

**Stages:** C05 C06. **Assessment:** Screenshot observation and usability hypothesis · V/I · medium.

**Observation:** The inspected six- and 150-service screenshots fit cards into a large canvas with small internal text. Zoom/search and team scope exist.

**Consequence:** Initial overview may make topology easier to see than individual semantics; task readability was not tested.

**Unanswered question:** Can intended users locate and interpret a target relationship at their actual display size and zoom?

**Evidence:** [R36](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/demo-shop-graph.jpg); [R37](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/enterprise-overview.png); [R09](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCanvas.jsx#L1218).

### F20 Workspace operations vocabulary is crowded

**Stages:** C03 C04 C10. **Assessment:** Verified UI content with inferred confusion · V/I · medium.

**Observation:** The toolbar places Refresh, Update graph, Run DiffMind all and Build graph near other exploration/administration commands.

**Consequence:** These distinctions reflect real operations but may require users to learn internal stages before obtaining the intended result.

**Unanswered question:** Can a viewer, editor and administrator each predict what their available commands do?

**Evidence:** [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246); [R36](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/demo-shop-graph.jpg).

### F21 First project creation limits exploratory control

**Stages:** C02. **Assessment:** Verified UI behavior with inferred consequence · V/I · medium.

**Observation:** Authorized users with no projects enter a forced creation dialog, cannot close it, and receive the default name DEFAULT.

**Consequence:** A person exploring the product may have to create state before deciding how they want to use it.

**Unanswered question:** Does this clarify first use or interrupt orientation for the explorer persona?

**Evidence:** [R06](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/Projects.jsx#L16).

### F22 Project-list failure also produces an empty list

**Stages:** C02 B02. **Assessment:** Verified code path with uncertain consequence · V/I.

**Observation:** The project-load catch records an error and sets projects to an empty array; the empty-list effect may open creation for an eligible user.

**Consequence:** An unavailable listing and a truly empty workspace can coexist with similar UI state. The error is not silently discarded.

**Unanswered question:** Can users distinguish request failure from no projects and avoid unnecessary creation?

**Evidence:** [R06](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/Projects.jsx#L16).

### F23 Shared dialogs lack several expected interaction behaviors

**Stages:** A10 B03 C02 C03. **Assessment:** Confirmed component gap · V · high; impact untested.

**Observation:** Modal/ConfirmDialog code does not supply dialog role, aria-modal, focus entry/trap/restore or Escape handling; the close icon lacks a descriptive accessible label.

**Consequence:** The component does not establish the modal interaction contract described by [W3C APG](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/).

**Unanswered question:** What happens in actual keyboard and assistive-technology sessions across all uses of these dialogs?

**Evidence:** [R08](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/components/Modal.jsx#L4).

### F24 Graph keyboard interaction is uneven

**Stages:** C05 C06. **Assessment:** Confirmed local implementation difference · V · high.

**Observation:** Connection summary rows have tabindex and Enter/Space handling. Inspected compact service nodes attach click without equivalent keyboard handlers.

**Consequence:** Keyboard reachability of core selections is uncertain and may differ by graph surface.

**Unanswered question:** Can a keyboard-only explorer accomplish equivalent service, resource, edge and evidence tasks?

**Evidence:** [R09](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCanvas.jsx#L1218).

### F25 Detail views bound how much evidence is visible

**Stages:** A07 C06. **Assessment:** Verified inspection limits · V · high.

**Observation:** Service details show the first 80 objects and 120 traces with notices; resource facts use a 40-item slice. MCP has separate response bounds.

**Consequence:** Visible detail is not necessarily all detail; users may stop investigation prematurely.

**Unanswered question:** Are omitted items discoverable for the task through another current route, and do hosts follow all required pages?

**Evidence:** [R10](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphDetails.jsx#L25); [R26](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35).

## PR understanding and review expectations

### F26 PR impact is not an automatic code-host review lifecycle

**Stages:** A08 B07 C08. **Assessment:** Confirmed scope boundary · V · high.

**Observation:** PR handlers provide live GET projections. The signed webhook ignores non-push events. No automatic PR posting/check lifecycle was established in inspected routes.

**Consequence:** An expectation of attached comments, checks, warnings and recurring review exceeds the demonstrated on-demand path.

**Unanswered question:** What does a target user currently expect by “check my PR,” and which parts already occur through their own host workflow?

**Evidence:** [R19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24); [R20](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/webhook.go#L67).

### F27 Exact PR evidence requires matching clean revision

**Stages:** A08 B07 C08. **Assessment:** Verified correctness boundary · V · high.

**Observation:** Changed-line attribution checks graph/entity revision against PR head. Stale/dirty/unknown graph context is labeled; candidate paths remain separate.

**Consequence:** These safeguards prevent overclaiming, while a routine default-branch snapshot may not support exact PR-head attribution.

**Unanswered question:** Do reviewers understand when a useful company graph is still unsuitable for precise PR evidence?

**Evidence:** [R19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24).

### F28 PR scores are heuristics with no demonstrated calibration

**Stages:** A08 B07 C08. **Assessment:** Verified formula and validation gap · V/H.

**Observation:** Code/file signals feed a score; eligible company context contributes through a fixed 62/38 weighting. No empirical probability or outcome calibration was found in reviewed evidence.

**Consequence:** A risk label can be mistaken for a likelihood of failure or a merge recommendation.

**Unanswered question:** How do scores relate to observed serious issues, false alarms and missed impacts in representative PRs?

**Evidence:** [R19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24); [R23](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1).

### F29 Removed and indirectly changed behavior remain a PR evidence limit

**Stages:** A08 B07 C08. **Assessment:** Source-grounded analysis boundary · V/I · medium.

**Observation:** The precise changed-line path uses head coordinates and added lines. A head graph alone cannot prove a removed entrypoint or arbitrary changed behavior behind an unchanged route.

**Consequence:** Important change types can require context beyond the demonstrated direct match; a lack of exact callers does not imply safety.

**Unanswered question:** Which deletion, configuration and transitive-change cases are unsupported or merely candidate-level in current output?

**Evidence:** [R19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24); [R21](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/contracts.go#L95).

### F30 Import and PR provider scope may diverge

**Stages:** B07. **Assessment:** Potential interoperability gap · V/I · medium.

**Observation:** Import supports configurable provider API settings, while the inspected PR handler constructs public api.github.com URLs.

**Consequence:** A workspace usable for enterprise import may not have an equivalent PR path. No enterprise runtime failure was reproduced.

**Unanswered question:** What provider/hostname combinations succeed end to end, including authentication and PR retrieval?

**Evidence:** [R19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24); [R07](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246); [R14](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1).

## Company access and continuity

### F31 Browser identity does not automatically connect the agent

**Stages:** B03. **Assessment:** Documented credential boundary · D/V · high.

**Observation:** Browser cookies alone are insufficient; remote agents need a proxy-supported identity or independently issued project token.

**Consequence:** Joining may be a two-channel onboarding process rather than one successful sign-in.

**Unanswered question:** Where do joining developers discover the correct agent credential and renewal owner?

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63).

### F32 Project access does not mirror code-host repository permissions

**Stages:** B04 B06 C09. **Assessment:** High consequence authorization boundary · D/V · high.

**Observation:** A project grant exposes the whole project including stored evidence, paths and configuration; no per-repository/field privacy filter is supplied.

**Consequence:** Different repository access groups cannot be assumed to retain their distinctions within one shared project.

**Unanswered question:** Does intended project membership fit the data imported into it?

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63).

### F33 Scoped company access must be deliberately enabled

**Stages:** B02 B04. **Assessment:** Documented deployment default · V/D · high.

**Observation:** Legacy is the default in the documented service settings and Compose. Scoped access requires explicit configuration and grants.

**Consequence:** Ordinary users’ effective authority depends on operator setup rather than a universal project-isolated default.

**Unanswered question:** Do real operators recognize the active access mode and validate it with ordinary identities?

**Evidence:** [R15](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/compose.yaml#L19); [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3).

### F34 Correction authority is centralized in scoped mode

**Stages:** B08 C10. **Assessment:** Explicit governance boundary · D/V · high.

**Observation:** Editors can queue saved-configuration refresh; scope, paths, imports, packs and analyzer configuration are global-admin tasks.

**Consequence:** A developer can have useful read/update rights while remaining dependent on centralized help for missing architecture.

**Unanswered question:** What correction/support response is available when context needed for a task is incomplete?

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R28](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/operations.md#L19).

### F35 Offboarding spans independent grants

**Stages:** B11. **Assessment:** Documented lifecycle boundary · D/V · high.

**Observation:** User memberships and project agent tokens are independent; existing admitted work and downloaded information have further limits.

**Consequence:** Ending browser access alone does not demonstrate that all associated agent access has ended.

**Unanswered question:** Can operators identify and revoke the grants belonging to the departing person?

**Evidence:** [R16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63); [R30](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectTokens.jsx#L1).

### F36 Recovery restores authority as well as architecture

**Stages:** A11 B11. **Assessment:** Documented recovery boundary · D/V · high.

**Observation:** Older snapshots can restore tokens revoked after their creation. Original expiry remains absolute; later revocation is not reconstructed.

**Consequence:** Recovering a service has access-policy consequences alongside data continuity.

**Unanswered question:** Is post-restore authority understood and reviewed by the actual service owner?

**Evidence:** [R18](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/backup-recovery.md#L35); [R17](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63).

### F37 Persistence is a single-server operating contract

**Stages:** A11 B09 B11. **Assessment:** Explicit architecture boundary · D/V · high.

**Observation:** Workspace persistence, exclusive writer/lifecycle ownership and offline backup exist. Opt-in SQLite applies to queue storage, not a distributed platform.

**Consequence:** Durability is a strength but does not establish high availability, arbitrary path/schema migration or disaster recovery in a particular deployment.

**Unanswered question:** What continuity and recovery requirements apply to the company that will depend on the service?

**Evidence:** [R14](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1); [R18](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/backup-recovery.md#L35); [R23](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1).

## Product evidence and validation

### F38 Operational measures are not user-outcome measures

**Stages:** A06 A07 B06 B09. **Assessment:** Evidence gap · V/D/H.

**Observation:** Queue, attempts, durations and repository operation metrics exist. Reviewed evidence contains no measured user activation, retention, satisfaction or decision accuracy.

**Consequence:** Healthy processing cannot establish that people or agents receive useful architecture context.

**Unanswered question:** What is the first meaningful result for each persona, and how often is it achieved and repeated?

**Evidence:** [R28](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/operations.md#L19); [R23](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1).

### F39 Graph differences can have multiple causes

**Stages:** C07. **Assessment:** Inference grounded in provenance · V/D/I.

**Observation:** Snapshots preserve graph history and pack provenance. Scope, analyzers, declarations and identity resolution can affect derived facts.

**Consequence:** A difference between maps does not by itself identify a source behavior change or production regression.

**Unanswered question:** Can a reviewer explain which inputs changed and what the graph difference actually supports?

**Evidence:** [R29](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCompare.jsx#L1); [R22](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91); [R21](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/contracts.go#L95).

### F40 Old evaluation instructions remain discoverable

**Stages:** A10 C10. **Assessment:** Confirmed documentation drift · V · high.

**Observation:** The eval fixture README refers to internal/eval and diffmind eval --mode cheap, while the inspected current implementation lacks that subsystem/command.

**Consequence:** A person seeking to validate coverage can follow a historical validation route that no longer exists.

**Unanswered question:** Which other discoverable documents belong to an earlier product model?

**Evidence:** [R31](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/testdata/extractor/eval/README.md#L4); [R01](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91).

### F41 Existing tests prove important mechanisms within bounded fixtures

**Stages:** A06 B05 C01. **Assessment:** Strength plus generalization limit · V/D · high.

**Observation:** Real-binary company and agent acceptance cover controlled topology, evidence, HTTP/MCP, history and recovery. Showcase checks verify generated inputs and individual scans.

**Consequence:** The product has meaningful automated evidence, but it is not equivalent to a ground-truth real company/PR accuracy study.

**Unanswered question:** How well do controlled successes generalize to representative conventions and real PR decisions?

**Evidence:** [R24](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/cmd/diffmind/agent_test.go#L27); [R25](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/company_acceptance_test.go#L30); [R35](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/scripts/test-showcase.sh#L20); [R13](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30).

### F42 Actual behavior of target personas remains unmeasured

**Stages:** All stages. **Assessment:** Research limitation · H · unresolved.

**Observation:** No participants, installed-client observation, longitudinal usage or independent competitor trials were part of this audit.

**Consequence:** Implementation and heuristic findings must not be presented as confirmed user pain or comparative superiority.

**Unanswered question:** Which hypothesized frictions occur most often, with what consequences for each person?

**Evidence:** [R23](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1).

# Comparative research

These comparisons examine **documented experience contracts**, not product quality rankings. Vendor claims were not tested. No competitor’s setup time, accuracy, security certification or business-performance claim is treated as independently verified. The products solve overlapping but different jobs, so a documented capability is not automatically a missing feature DiffMind should add.

## Codebase Memory MCP

The likely public reference for the user’s example is [DeusData Codebase Memory MCP](https://github.com/DeusData/codebase-memory-mcp); the user’s exact repository, installed version and settings were not inspected.

Its current [configuration documentation](https://github.com/DeusData/codebase-memory-mcp/blob/main/docs/CONFIGURATION.md) distinguishes initial automatic indexing from ongoing watching: auto_index defaults false, auto_watch and watcher_enabled default true, and the watcher lives in a daemon that outlasts individual MCP sessions. This helps explain how the experience can feel continuously maintained after setup without assuming all installations initially index automatically.

**Relevance:** the benchmark is the user’s low-maintenance lived experience. DiffMind’s company graph needs explicit scope and completed ingestion, and its local controller lifetime differs. Those differences need an intelligible experience contract; this is not evidence that the same internals or defaults are appropriate.

## Cursor

The formerly indexed codebase-indexing documentation URL currently redirects to [Search](https://cursor.com/docs/agent/tools/search). That page describes automatic local search-index use, automatic agent grep and multi-root workspace context. It also distinguishes keeping a search index local from sending an opened file’s contents in a model request.

**Relevance:** documented background behavior and an explicit downstream data distinction can reduce mental-model ambiguity. Search indexing and cross-company architecture resolution have different coverage requirements. Earlier descriptions of Cursor’s indexing should not be assumed current after this documentation change.

## Augment Context Engine MCP

[Augment’s overview](https://docs.augmentcode.com/context-services/mcp/overview) explicitly separates local working-directory indexing from remote indexing of selected repositories’ default branches. It describes local real-time updates, remote updates on default-branch pushes, and separate setup/authentication routes.

**Relevance:** source scope, branch and update policy are presented as part of the product contract. This is a useful comparison for A05/A09 and B03/B05. It does not demonstrate that the service identifies every architectural relationship, nor that “live understanding” in its positioning means runtime truth.

## Sourcegraph

[Sourcegraph MCP documentation](https://sourcegraph.com/docs/api/mcp) describes an instance endpoint, OAuth or token authentication, selectable tool suites and explicit MCP access controls. Repository permissions remain a separate boundary. Its [permission documentation](https://sourcegraph.com/docs/admin/permissions) describes configured synchronization of source-host access.

**Relevance:** the joining-user path distinguishes service access, MCP authority and repository visibility. DiffMind currently uses project-wide grants and operator-supplied identity or independent tokens; upstream repository permissions are not inherited automatically. Sourcegraph’s documented integration is an adjacent operating model, not proof that DiffMind’s project model is wrong. Actual administrative complexity and enforcement were not independently tested.

## CodeRabbit

The [quickstart](https://docs.coderabbit.ai/getting-started/quickstart) centers on connecting the code host, selecting repositories and opening a PR. [Automatic review controls](https://docs.coderabbit.ai/configuration/auto-review) describe eligible PRs, draft/default-branch behavior and manual review commands. [Multi-repository analysis](https://docs.coderabbit.ai/knowledge-base/multi-repo-analysis) documents linked repository scope and the revision used for cross-repository context.

**Relevance:** PR review is presented as a trigger-to-feedback lifecycle in the place developers review code. DiffMind’s current PR projection is on demand and revision-sensitive. This establishes a category distinction and an expectation to investigate, not a requirement to reproduce CodeRabbit’s feature set. No cross-repository accuracy comparison was performed.

## Greptile

[Greptile’s introduction](https://www.greptile.com/docs/introduction) describes connecting selected repositories through a code-host app, automatic PR comments and reviewer feedback through reactions/replies. The inspected page also makes broad completeness, speed and business-impact claims; those claims were not independently validated.

**Relevance:** scope acquisition, routine delivery and user correction are visible parts of the stated review experience. DiffMind already exposes inspectable evidence but has not demonstrated a corresponding automatic code-host review and feedback lifecycle. “Complete graph” marketing is not a reliable semantic-coverage baseline for either product.

## Backstage

[Backstage’s catalog overview](https://backstage.io/docs/features/software-catalog/) centers on software discoverability and ownership metadata harvested from version-controlled declarations. It presents team-focused browsing, filtering and metadata maintenance. Its [relation model](https://backstage.io/docs/features/software-catalog/well-known-relations/) distinguishes generic dependency, provided/consumed APIs and ownership; ownership does not grant runtime access.

**Relevance:** a manager’s catalog question, an engineer’s dependency question and an authorization question are different. DiffMind’s team/service graph can support orientation, but inferred source relationships do not alone establish accountable ownership or access policy. Backstage’s declarations also show why catalog evidence is not equivalent to observed runtime dependencies.

## What the comparisons establish

Across these sources, the useful comparison dimensions are: **how scope is chosen, what becomes automatic after connection, which source revision is represented, how identity affects visibility, where routine value appears, and how a user understands or corrects uncertainty**.

Public documentation offers examples of those dimensions. It does not establish user satisfaction, implementation equivalence or a reliable performance ranking. The report therefore uses the comparisons to sharpen the unanswered questions, rather than produce a feature parity checklist.

# Research principles used in the assessment

## Predictable state and responsibility

[Nielsen Norman Group’s usability heuristics](https://www.nngroup.com/articles/ten-usability-heuristics/) support inspecting status visibility, real-world language, control, consistency, recognition and error recovery. Applied here, the question is whether the person can predict the transition from installed to usable, distinguish refresh/build stages and recover from incomplete work. These are heuristic assessments, not observations of user failure.

[Progressive disclosure research](https://www.nngroup.com/articles/progressive-disclosure/) distinguishes primary needs from specialized information while cautioning against hiding information needed routinely. The import knobs and workspace toolbar therefore create a task-comprehension question. Their mere presence is not proof of poor design; relevance to each persona matters.

## Agent uncertainty and control

[Microsoft’s Human-AI Interaction guidelines](https://www.microsoft.com/en-us/research/articles/guidelines-for-human-ai-interaction-eighteen-best-practices-for-human-centered-ai-design/) emphasize setting capability expectations, explaining uncertain behavior, supporting correction and communicating consequences over time. They apply here to the **host agent’s interaction and explanation**, even though DiffMind’s analyzer is deterministic. The important issue is whether the host faithfully communicates source limits rather than turns partial evidence into a confident global claim.

## MCP interoperability

The [MCP tools specification](https://modelcontextprotocol.io/specification/2025-11-25/server/tools) makes tools available for model-directed use and supports structured results. This does not promise that any host will invoke a particular tool at the right moment. Tool annotations are not a substitute for enforced authorization.

The [roots specification](https://modelcontextprotocol.io/specification/2025-11-25/client/roots) describes optional client-exposed filesystem scope and consent-related behavior. MCP itself does not define DiffMind’s company/project onboarding policy or guarantee automatic indexing. Interoperable transport and a low-maintenance product experience are separate questions.

## Accessibility

The [W3C modal dialog pattern](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/) describes dialog semantics, focus management, keyboard containment and dismissal. F23 is grounded in missing behaviors in the inspected component, not a formal accessibility certification result.

[WCAG status-message guidance](https://www.w3.org/WAI/WCAG22/Understanding/status-messages.html) explains announcing meaningful status without taking focus. DiffMind already has aria-live/status surfaces for work progress. F24 concerns a separate issue: equivalent access to graph selections. Positive status semantics do not establish that all exploration is accessible.

## Management and developer experience

The [DevEx research by Noda, Storey, Forsgren and Greiler](https://www.michaelagreiler.com/wp-content/uploads/2024/06/DevEx-WhatDrivesProductivity.pdf) examines feedback loops, cognitive load and flow, combining developer feedback with system evidence. This supports looking beyond scans/second or job success to the interruptions and responsibility handoffs in the three journeys.

[Google’s HEART research](https://research.google/pubs/measuring-the-user-experience-on-a-large-scale-user-centered-metrics-for-web-applications/) links user goals to evidence about experience. Applied here, the relevant unknowns are activation, repeated usefulness, task success and confidence calibration. These are missing research measurements, not prescribed telemetry features.

# Strengths supported by the inspected evidence

The audit identifies substantial working mechanisms rather than treating the product as an empty prototype.

| Strength | Evidence | Practical value and limit |
| --- | --- | --- |
| Agent can operate the product | Real stdio acceptance and bounded management catalog | Eliminates mandatory dashboard chores after bootstrap; real-client completion remains unverified |
| Architecture can be inspected | Source evidence, graph details, full-evidence query mode | Supports reviewable explanations; evidence is bounded by extraction and declarations |
| Snapshots persist | Completed graph selection, saved comparisons, historical pack context | Supports reproducibility; latest completed is not necessarily current source |
| Failure is represented | Partial/failed statuses, per-repository outcomes, job attempts and recovery | Supports honest readiness; user comprehension remains unmeasured |
| Incremental work and integrity are considered | Acceptance coverage and documented fingerprint/pack verification | Supports repeated use; does not establish production scale or low maintenance |
| Shared permissions have explicit contracts | Scoped membership filtering, role ceiling, token expiry/revocation | Supports bounded company access; mode is opt-in and project grants are broad |
| PR context avoids some overclaims | Revision checks, candidate/proven distinction, eligibility notes | Supports cautious investigation; does not establish full review coverage or calibrated risk |
| Public proof distinguishes evidence types | Curated company, real pinned scans, generated enterprise fixture | Makes limitations visible; results still require careful interpretation |

Evidence: [agent acceptance](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/cmd/diffmind/agent_test.go#L27), [company acceptance](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/company_acceptance_test.go#L30), [query contract](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L181), [permissions](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L34), [PR attribution](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L641) and [public benchmark record](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30).

# Current state by outcome

This table describes how far the evidence supports each intended purpose. It avoids maturity scores because the audit has no validated scoring model.

| Intended outcome | What is supported today | What remains unestablished |
| --- | --- | --- |
| Agent handles first-use operations | Full management and disposable end-to-end acceptance | Reliable first-use success in actual client environments |
| Ordinary work gets cross-repository context | Graph tools with scoped snapshots and source evidence | Natural host adoption and measurable answer/task improvement |
| Context remains useful after changes | Persisted history, incremental work and configurable refresh | User-understood local continuity and repeated real-world usefulness |
| Company developer joins easily | Shared access model, project grants, separate agent credentials | Actual joining effort, administrative delay and credential-renewal experience |
| Company architecture is broadly understood | Bounded deterministic patterns, declarations and controlled exact relationships | Ground-truth coverage for an arbitrary company’s stack and conventions |
| Explorer forms a defensible conclusion | Scope controls, evidence details, snapshots and comparison | Accessibility/task success, interpretation accuracy and stakeholder comprehension |
| PR changes are understood | On-demand PR projection, heuristics, revision checks and contract comparison | Automatic attached review lifecycle and calibrated impact accuracy |
| Workspace is durably operated | Single-server persistence, queues, locks and offline recovery | Production availability, recovery outcomes and long-term operating cost |

“Unestablished” does not necessarily mean absent from every user’s workflow. For example, a host agent or an organization’s own automation might supply activities outside the inspected product. That possibility is not evidence that DiffMind itself provides them.

## Future outcomes for the existing purposes

These are acceptance meanings for future evaluation, not designs or a delivery plan:

1. **The individual knows what their installation does.** They can distinguish a working connection, a prepared graph and a maintenance policy without learning internal orchestration.

2. **The agent supplies useful context at the right time.** A real task receives relevant graph evidence, while unsupported or stale information remains qualified.

3. **The joining developer reaches the intended company context.** Person and agent have understandable access, and missing authority does not look like missing architecture.

4. **The graph supports bounded claims.** Relevant relationships are verifiable, omitted scope and unsupported conventions are visible in the conclusion, and empty results are not treated as proof of absence.

5. **Exploration is usable across access needs.** Core questions can be answered with readable controls and equivalent interaction, with provenance preserved.

6. **PR feedback has a clear meaning.** The user understands what triggered analysis, which revision it describes, what is proved, what is a candidate and what has not been evaluated.

7. **Persistence has a clear continuity contract.** Reconnect, restart, refresh, upgrade, handoff and restore do not silently change what the user believes is available or authorized.

8. **Management can distinguish operational success from user value.** Successful jobs and larger graphs are supported by evidence of useful repeated tasks and appropriately bounded decisions.

The evidence does not select how any outcome should be achieved. Some may need only clearer behavior/documentation; others may expose larger coverage or operating boundaries. Determining that distinction requires the missing evidence below.

# Unanswered research questions and decision boundaries

## Individual developer

The unobserved portions are actual installation with the user’s existing client settings, natural selection of DiffMind tools during coding, understanding of agent versus query-only mode, reaction to omitted relationships and expectations after disconnect. Known backend acceptance gives a sound starting point but does not answer those user questions.

The important unknown outcome is **how much user attention is needed after successful setup**. It includes recognizing stale context, requesting refresh, resolving ownership conflicts and teaching company-specific conventions. No measured maintenance burden is available.

## Company joining developer

The unobserved portions are the handoff from operator to developer, browser versus agent credentials, identity/proxy compatibility, project appropriateness, administrator response to missing scope and offboarding execution. An access mechanism can work correctly while this organizational journey remains slow or unclear.

The important unknown outcome is **whether a developer can obtain useful context with ordinary rights and predictable support**. The product’s project-wide visibility, independent tokens and centralized configuration authority are deliberate boundaries the organization must understand.

## Architecture explorer

The unobserved portions are first-time orientation, comprehension of dependency direction and provenance, discovering filtered/limited information, keyboard/screen-reader tasks, explaining a snapshot comparison and interpreting a PR risk label. Screenshot inspection is insufficient to assess those tasks.

The important unknown outcome is **whether a person’s conclusion is both useful and appropriately qualified**. Attractive graphs, successful rendering and completed imports do not measure that outcome.

## Product and management decisions left open

The current evidence does not settle several choices:

- Which existing purpose is the primary repeated value: coding context, architectural exploration or change/PR investigation?

- What level of freshness is sufficient for each purpose, and who owns the refresh responsibility in each deployment?

- Which source/provider/framework boundaries are acceptable to the intended audience?

- What does a company consider its authorized architecture scope, and does project-wide visibility fit that scope?

- What weight should reviewers give a heuristic risk label before its outcome relationship is known?

- Who owns correction of identities, declarations and unsupported conventions?

- What continuity and recovery obligations arise when the service becomes part of a company’s daily work?

These are unresolved decisions, not recommendations to choose a direction now.

## Existing public evidence and its limits

There is already a useful separation of validation evidence in the repository:

| Evidence class | What it can establish | What it cannot establish alone |
| --- | --- | --- |
| Three-service controlled company acceptance | Exact known relationships, evidence, HTTP/MCP and persistence in a bounded fixture | Arbitrary company coverage |
| Six-repository Demo Shop | Controlled extraction/change story and individual scan checks | Natural user adoption or real PR performance |
| Pinned OpenTelemetry/Boutique scans | Compatibility observations and visible missing framework/configuration coverage | Assembled company graph precision/recall or review accuracy |
| Generated 150-service enterprise fixture | Reproducible navigation data and visual scope | Real extraction scale, runtime traffic, operating cost or latency |
| Historical release/readiness checkpoint | Recorded checks and explicit exclusions at a particular point | Current production certification or every real-client journey |

A broad claim such as “works with anything” needs a declared coverage meaning. Number of languages, number of repositories, successful scans, graph size, exact relationships and successful user decisions are different evidence. This audit does not convert any one of them into universal compatibility.

## Evidence still needed to resolve the hypotheses

No new experiments were performed in this research audit beyond the earlier session checks. The missing evidence categories are actual task observations for the three personas; representative company topology with known positive and negative relationships; real PR outcomes with explicitly reviewed misses and false alarms; actual host-client installation and repeat-use behavior; accessibility interaction; and deployment-specific continuity/recovery outcomes.

Relevant measurements are unknown, rather than zero: setup completion, first useful answer, repeated task success, user correction effort, administrator help burden, calibrated confidence, material PR misses/noise, queue waiting and recovery time. Establishing values or thresholds is outside the current report.

The present findings make these unknowns explicit without turning them into prescribed features, instrumentation or implementation tickets.

# Evidence index

All repository links below pin the inspected source revision. Source anchors identify useful entry points; some findings depend on more than one portion of a file.

| ID | Inspected evidence | Main use |
| --- | --- | --- |
| R01 | [README.md line 91](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/README.md#L91) | Product promise and agent first entry |
| R02 | [AGENT_SETUP.md line 34](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/AGENT_SETUP.md#L34) | Agent bootstrap scope and onboarding |
| R03 | [internal/workspace/agenthost/host.go line 27](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agenthost/host.go#L27) | Local runtime defaults |
| R04 | [docs/agent-operations.md line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-operations.md#L1) | Connection types and lifecycle |
| R05 | [docs/personal-setup.md line 176](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/personal-setup.md#L176) | Manual query only agent setup |
| R06 | [internal/workspace/ui/web/src/views/Projects.jsx line 16](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/Projects.jsx#L16) | Project selection and first project flow |
| R07 | [internal/workspace/ui/web/src/views/ProjectWorkspace.jsx line 246](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectWorkspace.jsx#L246) | Workspace controls progress and import |
| R08 | [internal/workspace/ui/web/src/components/Modal.jsx line 4](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/components/Modal.jsx#L4) | Dialog behavior |
| R09 | [internal/workspace/ui/web/src/views/GraphCanvas.jsx line 1218](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCanvas.jsx#L1218) | Graph controls and visible scope |
| R10 | [internal/workspace/ui/web/src/views/GraphDetails.jsx line 25](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphDetails.jsx#L25) | Evidence details and display limits |
| R11 | [internal/workspace/query/service.go line 117](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/service.go#L117) | Project resolution and completed snapshots |
| R12 | [docs/supported-patterns.md line 3](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/supported-patterns.md#L3) | Tested analysis coverage and limits |
| R13 | [examples/public-benchmarks/README.md line 30](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/public-benchmarks/README.md#L30) | Pinned public compatibility evidence |
| R14 | [docs/company-deployment.md line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/company-deployment.md#L1) | Company deployment and operations |
| R15 | [compose.yaml line 19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/compose.yaml#L19) | Shared deployment defaults |
| R16 | [docs/project-access.md line 3](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/project-access.md#L3) | Access modes permissions and trust boundaries |
| R17 | [docs/agent-tokens.md line 63](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/agent-tokens.md#L63) | Independent agent credentials |
| R18 | [docs/backup-recovery.md line 35](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/backup-recovery.md#L35) | Recovery scope and restored authority |
| R19 | [internal/workspace/ui/handlers_pull_requests.go line 24](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/handlers_pull_requests.go#L24) | Live PR projection scoring and revision checks |
| R20 | [internal/workspace/ui/webhook.go line 67](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/webhook.go#L67) | Push refresh event contract |
| R21 | [internal/workspace/query/contracts.go line 95](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/query/contracts.go#L95) | Bounded contract comparison |
| R22 | [docs/knowledge-packs.md line 91](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/knowledge-packs.md#L91) | Convention teaching provenance and integrity |
| R23 | [docs/readiness-verification.md line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/readiness-verification.md#L1) | Existing validation record and exclusions |
| R24 | [cmd/diffmind/agent_test.go line 27](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/cmd/diffmind/agent_test.go#L27) | Real stdio agent acceptance |
| R25 | [internal/workspace/ui/company_acceptance_test.go line 30](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/company_acceptance_test.go#L30) | Real binary company fixture acceptance |
| R26 | [internal/workspace/mcpserver/server.go line 35](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/mcpserver/server.go#L35) | Agent schemas descriptions and graph tools |
| R27 | [internal/workspace/agentapi/api.go line 132](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/agentapi/api.go#L132) | Management discovery and async workflow |
| R28 | [docs/operations.md line 19](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/operations.md#L19) | Queued work progress quotas and metrics |
| R29 | [internal/workspace/ui/web/src/views/GraphCompare.jsx line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/GraphCompare.jsx#L1) | Snapshot comparison UI |
| R30 | [internal/workspace/ui/web/src/views/ProjectTokens.jsx line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/internal/workspace/ui/web/src/views/ProjectTokens.jsx#L1) | Agent token UI |
| R31 | [testdata/extractor/eval/README.md line 4](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/testdata/extractor/eval/README.md#L4) | Obsolete evaluation instructions |
| R32 | [install.sh line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/install.sh#L1) | Release installation behavior |
| R33 | [scripts/agent-setup/main.go line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/scripts/agent-setup/main.go#L1) | Source bootstrap implementation |
| R34 | [examples/enterprise-showcase/README.md line 1](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/examples/enterprise-showcase/README.md#L1) | Synthetic enterprise fixture |
| R35 | [scripts/test-showcase.sh line 20](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/scripts/test-showcase.sh#L20) | Showcase verification scope |
| R36 | [docs/assets/readme/demo-shop-graph.jpg](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/demo-shop-graph.jpg) | Checked in six service screenshot |
| R37 | [docs/assets/readme/enterprise-overview.png](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/assets/readme/enterprise-overview.png) | Checked in enterprise screenshot |
| R38 | [docs/workspace/README.md line 35](https://github.com/safa-safakhou/diffmind/blob/5d2548514bfc5d920dfd40def5fcacd80162e9ae/docs/workspace/README.md#L35) | Ingestion completion contract |

External sources are cited beside their use in Comparative research and Research principles. They include seven product references, W3C accessibility guidance, MCP specifications, Nielsen Norman Group heuristics, Microsoft human-AI interaction research, DevEx research and Google HEART research. Vendor documentation is documentary evidence; the two research frameworks are lenses for identifying missing user-outcome evidence.

The [v0.1.0 release page](https://github.com/safa-safakhou/diffmind/releases/tag/v0.1.0) was checked separately to establish that a public release exists. This is not a new verification of every downloadable archive on every supported platform.

**Repository preservation:** the report was created outside the repository. No product changes, fixes, new features, imports or deployment/client reconfiguration were performed.

**Use of this report:** the 32 journey stages describe the current paths; F01–F42 provide a traceable findings inventory; the comparison and research sections explain the assessment lenses; the outcome and open-question sections state what remains unresolved. The next substantive decision belongs to the user after reviewing this evidence.

