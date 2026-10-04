The Demo Commerce context is updated. Refresh job `refresh-3996e9f93bd0dfdf3d49b26ca83fc801` completed and saved graph run `20261004T165022Z`.

- **Scope preserved:** all six existing repositories remain in the project. `catalog` was reanalyzed; the other five were reused.
- **Local edit preserved:** `catalog/app.py` remains modified in the working tree. The new run includes that dirty checkout.
- **Graph comparison:** five evidence records changed, all marked evidence-only. Service and edge counts remain six and nine, with no added or removed facts.
- **Configuration:** unchanged.

**Remaining limitations:** `catalog` is still marked dirty, so readiness reports `needs_update` and the graph warns it was built from a dirty working tree. Architecture coverage remains unverified static-source evidence. The local repositories have no supported Git remote for PR queries.