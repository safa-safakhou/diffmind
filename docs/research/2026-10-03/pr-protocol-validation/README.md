# PR caller protocol validation

Follow-up to MNI-198 and MNI-199. Base commit: 65979b8.
This report and its implementation are committed together; logs were produced
from the working tree containing that implementation.

## Defect and correction

Exact PR caller matching previously compared IDs and HTTP-shaped operation names
across every graph edge type. A clean, source-extracted RPC fact could therefore
qualify as an exact HTTP caller just because its name or ID coincided. Missing
protocol and unknown surface kinds could also qualify.

A regression reproduced five false positives before the fix. IDs now match
within the known surface protocol. HTTP method/path fallback applies only to
HTTP surfaces and HTTP edges. Recognized HTTP route/endpoint/webhook and RPC
endpoint kinds can qualify; unknown kinds remain candidates. Existing clean
revision, source provenance and changed-head-line gates remain required.

Eight positive/negative cases now pass: HTTP operation, RPC operation lookalike,
unknown edge, cross-protocol ID collisions in both directions, valid HTTP ID,
valid RPC ID and unknown surface. Existing fixture metadata now reflects real
HTTP graph kinds/protocols rather than missing metadata.

## Verification and limits

- Complete UI/backend package: passed.
- Full Go suite: passed; retained in go-tests.txt.
- Race-enabled exact-caller, provenance, changed-line and revision tests: passed.
- New eight-case regression: passed; retained in regression.txt.
- No frontend changes, new live browser trial, independent reviewer, human
  participant, Enterprise credential trial or native platform certification
  occurred in this follow-up. Previous live-study results remain tied to their
  recorded older candidate, not this changed product revision.

This fixes a concrete defect, but does not complete MNI-198's independently
labeled real-source corpus requirement or MNI-199's human comprehension gate.
The seven remaining child tasks retain their external validation requirements.
Unknown output and absence of exact callers still do not establish merge safety.

Reproduce:

    go test ./...
    go test ./internal/workspace/ui -run TestExactPRCallerRejectsCrossProtocolAndUnknownSurface -v -count=1
    go test -race ./internal/workspace/ui -run 'Test(Exact|DeclaredOrUnknown|ChangedEntrypoints|PullRequestGraph|ServiceGraphRevision)' -count=1
