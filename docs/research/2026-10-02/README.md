# DiffMind journey research — 2 October 2026

Source revision: `5d2548514bfc5d920dfd40def5fcacd80162e9ae`.

This research examines three overlapping journeys: an individual developer using an agent, a developer joining an existing company workspace, and a person exploring architecture. The first audit inspected source, documentation and public practice. The follow-up then exercised the compiled application, Chromium, stdio MCP, HTTP MCP, company grants, refresh work and an actual public GitHub PR.

## Read the reports

| Record | Contents |
| --- | --- |
| [Journey audit](journey-audit.md) | The original broad report: 32 journey stages, 42 findings, cross-disciplinary assessment and external references. Its limitations describe the audit phase before live interaction. |
| [Live usability study](live-usability-study.md) | Executed tasks, observed outcomes, 11 live findings, strengths, interpretation boundaries and remaining research questions. |
| [Reproduction and evidence guide](reproduction.md) | Environment, actual execution sequence, MCP/browser protocol, evidence navigation and observer mistakes excluded from product conclusions. |
| [Evidence index](evidence/README.md) | Timestamped observation logs, fixture and real PR results, metadata, hashes and the retained observer scripts. |
| [Coordinated improvement plan](improvement-plan.md) | Subsequently requested proposed solutions, shared decisions, dependency graph and coverage of all 53 findings. |
| [Detailed work items](work-items.md) | Solutions, acceptance criteria and validation for 28 tasks grouped into seven outcome epics. |
| [Linear project and task index](linear-backlog.md) | Created project, three milestones, first four tasks, all issue links and verified hierarchy/dependencies. |

## What the live work established

The individual agent path produced a graph of all six Demo Shop repositories, persisted it across disconnect/reconnect, reused unchanged analysis and identified the expected contract changes. The company path enforced viewer restrictions, accepted an editor's queued refresh and removed the open graph after membership revocation. A real public PR loaded successfully and displayed file evidence, a risk explanation and explicit stale-graph caveats.

Concrete experience problems also appeared. The dashboard closed the dry-run preview without displaying candidates, discarded the import scope, placed an input error behind the modal and let keyboard focus escape the first-use dialog. The initial graph fit made a service label approximately six pixels high. A dashboard launched by absolute binary path required that binary's directory in the process PATH for analysis to work. In the real-repository trial, examples were attributed to the enclosing `diffmind` service; source evidence confirmed their Demo Shop origin.

The live study contains **143 retained timestamped events and 22 screenshots**. It is an AI-observer walkthrough with real runtime interactions, rather than a recruited-human usability study. No product source, personal agent configuration or existing company workspace was changed. Only research documentation and evidence were added to this repository. GitHub was read without posting comments, checks or reviews. Test processes were stopped.

The original audit and live-study reports describe the observed state without selecting fixes. At the user's subsequent request, the coordinated plan adds proposed solutions and a verified Linear backlog. The planning phase did not implement product changes; subsequent implementation is recorded below.

## Implementation

[Integrated implementation batch 1](implementation-batch-1/README.md) records the changes, task dispositions, exact candidate and browser/MCP/source-scope regression evidence.

[Integrated implementation batch 2](implementation-batch-2/README.md) records committed import consistency, shared live freshness, overdue reconnects, bounded failure recovery and access continuity, with remaining acceptance criteria.
