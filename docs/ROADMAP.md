# Product roadmap

DiffMind is a deterministic architecture analysis tool for individual developers
and shared teams. It runs on one server and exposes the same saved source
evidence through the browser, HTTP API and MCP.

## Current capabilities

- Local repository and GitHub organization import with reviewed source scope.
- Incremental analysis, saved graph history, contract comparison and caller tracing.
- Declarative knowledge packs with source evidence and positive/negative fixtures.
- Read-only queries, local agent management and scoped viewer/editor access.
- Bounded background refresh, durable jobs and JSON or SQLite queue storage.
- Offline backup/restore, managed retention and optional Linux systemd scheduling.
- Native macOS and Linux distribution on AMD64 and ARM64.

See [architecture](ARCHITECTURE.md), [supported patterns](supported-patterns.md),
[agent operations](agent-operations.md) and [contributor validation](validation.md)
for the contracts and checks behind these features.

## Remaining work

- Broader framework and protocol coverage, backed by source/runtime expectations
  and positive, negative and ambiguity tests.
- More reliable agent discovery of saved architecture under ordinary dependency
  questions, including multiple accessible projects and empty local directories.
- PR evidence for deleted surfaces and indirect changes beyond supported changed
  source locations. Current exact caller evidence is a bounded static-source claim.
- Qualification against actual Enterprise endpoints and production identity
  providers. Controlled provider tests do not replace deployment checks.
- Organization-scale metadata storage and distributed workers with fencing,
  shared artifacts, authorization and recovery compatibility.
- Identity-provider group synchronization and user-bound personal tokens with
  explicit provisioning and revocation semantics.
- Application-schema migration and workspace path relocation with dry runs and
  rollback. Current recovery restores original paths.
- Encrypted off-host backup storage and additional scheduler integrations.
- Longer production operation trials, physical-device accessibility checks and
  broader deployment coverage.

Current memberships and project tokens are independent grants. Queue indexing
and resource caps support a single server; they do not provide distributed
execution or OS isolation. Static evidence does not establish runtime traffic,
complete dependency coverage or safe-to-merge decisions.
