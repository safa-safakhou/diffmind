# Shared setup, lifecycle and recovery batch

Continues the user-pushed branch at dee0bc72c0d77363608b31401434b8e97bf96a60.
First implementation commit: d501000. No existing configured mode is migrated.

## Changes and evidence

- MNI-193: the new shared deployment environment example selects scoped mode.
  CLI and Compose fallback remain legacy to preserve existing authority.
  Project access displays the active mode, project-wide evidence exposure and
  identity verification checklist. Real-handler migration tests cover prepared
  grants, ungranted/guessed URLs, viewer/editor roles and recovery credentials.
- MNI-182: competing-controller feedback gives scoped-token/current-endpoint
  guidance and forbids forwarding the unauthenticated backend. The actual
  stdio acceptance launches a second controller and checks prompt rejection;
  existing acceptance continues through private token issuance/revocation,
  maintenance, disconnect, crash cleanup, reconnect and saved graph history.
- MNI-196: private disposable archive drill proves revoked token and membership
  resurrection after exact-path restore, then reconciles both independently
  while preserving service/admin access. Runs with JSON and SQLite queues;
  interrupted work recovers and SQLite integrity verifies. Recovery documentation
  makes ingress isolation and independent access reconciliation explicit.

## Verification and limits

Frontend: 20 library and 26 component tests passed; embedded production build
passed. UI/MCP/store tests passed. UI/agenthost/homelock race suites passed.
Actual TestAgentAcceptance passed using a newly built candidate executable.
The isolated restore drill passed in both queue modes (about 0.08 seconds for
the complete small fixture on this WSL host); this is not a production RTO.

No human participants, real identity-provider deployment or production-size
backup were exercised. Migration enforcement uses actual handlers, not a newly
recruited browser/MCP persona trial. MNI-193 and MNI-196 retain remaining
cross-interface setup / upgrade drill criteria; this report does not imply
completion of every backlog task. No changes were pushed automatically.
