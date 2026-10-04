## Demo Commerce context is usable

The project **Demo Commerce** (`demo-commerce`) now has a saved, queryable graph: run **`20261004T164917Z`**. All six approved repositories were imported and analyzed successfully; readiness reports all six fresh, with no stale or dirty repositories.

The preview and registered repository list both contained exactly **catalog, checkout, gateway, notification, payment, and storefront** under the approved directory. No other repositories were included.

## Service and queue connections

- **Gateway directly calls Catalog and Checkout.** Its checkout handler sends `GET /products` to Catalog and `POST /v1/checkout` to Checkout. [gateway/main.go](<TRIAL_ROOT>/verified-demo/repositories/gateway/main.go:14)
- **Storefront calls Gateway** to begin the checkout flow, so Gateway then makes those Catalog and Checkout calls. [storefront/main.go](<TRIAL_ROOT>/verified-demo/repositories/storefront/main.go:8)
- **Checkout and Notification share the Kafka topic `orders.created`.** Checkout’s publisher sends to that topic, and its `order-out-0` binding names it as the destination. Notification listens to the same topic. [OrderPublisher.java](<TRIAL_ROOT>/verified-demo/repositories/checkout/src/main/java/example/OrderPublisher.java:8) · [application.yml](<TRIAL_ROOT>/verified-demo/repositories/checkout/src/main/resources/application.yml:4) · [OrderNotificationConsumer.java](<TRIAL_ROOT>/verified-demo/repositories/notification/src/main/java/example/OrderNotificationConsumer.java:5)

## Maintenance and limits

Automatic refresh was initially enabled every 15 minutes. I disabled it for this agent-owned DiffMind runtime and verified that refresh is now disabled; future maintenance is manual. That setting applies to the local runtime, rather than only to this project. No source files were edited, and ingestion reports zero synced repositories.

The graph is based on static source evidence, not runtime traffic, and its dependency coverage is marked unverified. In particular, the inspected `CheckoutController` does not show a call to `OrderPublisher`, so the saved evidence establishes the publisher and consumer configuration on the shared topic, not that every checkout request publishes a message. The repositories are local-only, with no supported Git remote for PR queries.