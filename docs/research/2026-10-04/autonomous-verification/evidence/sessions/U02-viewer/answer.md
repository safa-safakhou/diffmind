Using saved graph **`20261004T164807Z`** for Demo Commerce:

- **Gateway calls both services.** Its saved source evidence records `GET http://catalog/products` and `POST http://checkout/v1/checkout` in `gateway/main.go` (lines 15 and 17). Both destinations resolve to their matching service names.
- **`orders.created` connects checkout to notification.** Checkout publishes to the Kafka topic, based on `checkout/src/main/resources/application.yml` line 6. Notification consumes the same topic in `OrderNotificationConsumer.notifyCustomer` (`notification/src/main/java/example/OrderNotificationConsumer.java`, lines 6–9).

**Can this connection update context or change repository scope?** The current connection exposes query and inspection operations, but readiness reports `refresh: false` and `configure: false`; no context-management operations are available here. Its saved scope contains six repositories: catalog, checkout, gateway, notification, payment, and storefront.

**Evidence limits:** This is static source evidence, not observed runtime traffic, so it doesn’t prove calls or messages actually occur or succeed. Dependency coverage is marked unverified. The saved graph reports five repositories fresh and `catalog` dirty, so its analysis may need updating. I made no changes.