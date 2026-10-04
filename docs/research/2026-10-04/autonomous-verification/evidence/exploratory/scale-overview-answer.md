The saved graph identifies **`catalog-orchestrator`** as the direct candidate: it has an inbound RPC dependency on `checkout-api`. Following reverse dependency edges yields **58 additional transitive candidates**—services that depend on a service in that chain:

- **Catalog:** `catalog-writer`, `catalog-api`, `catalog-reader`, `catalog-events`, `catalog-worker`
- **Platform, payments, identity, growth, fulfillment, edge, data, and customer:** for each team, the `orchestrator`, `writer`, `api`, `reader`, `events`, and `worker` services
- **Checkout:** `checkout-orchestrator`, `checkout-writer`, `checkout-reader`, `checkout-events`, `checkout-worker`

These are topology-based candidates, not confirmed runtime effects. The dependency evidence is static and its coverage is **unverified**; it does not establish that the edges carry live traffic or that a particular change would break them. Missing or unsupported dependencies may also be absent from the graph.

The evidence is pinned to saved graph run **`20261004T150313Z`**, saved at **2026-10-04 15:03:27 UTC**. `checkout-api` was analyzed on `main` at revision `0000000000000000000000000000000001352c16` at **15:03:13 UTC** and marked clean then. But its current checkout state is unknown: readiness reports freshness **unknown for all 150 repositories**, and the revision basis is stored repository status. So this establishes the analyzed revision, not that it matches today’s checkout or a PR head.