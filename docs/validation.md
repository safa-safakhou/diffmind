# Contributor validation

Use the Go version in `go.mod`, Git and a C compiler. Dashboard changes also
require the Node.js version supported by the frontend packages.

## Required checks

Run `make test` for every change. Run the additional checks that apply:

| Change | Checks |
| --- | --- |
| Extraction rules or knowledge packs | `make test-packs test-acceptance` |
| Ingestion, storage or graph pipeline | `make test-race` |
| Either dashboard | `make ui-test ui-build` |
| Frontend dependencies | `make ui-audit` |
| Go dependencies | `make vulncheck` |
| Installation or distribution | `make test-distribution` |
| Agent management or onboarding | `make test-agent` |

`make verify` runs the comprehensive local suite. See
[CONTRIBUTING](../CONTRIBUTING.md) and [distribution](distribution.md) for the
contribution and release requirements. For a built native archive, run
`make test-release-native ARCHIVE="$diffmind_archive" VERSION="$diffmind_version"`
with the archive location and its exact version.

The [company acceptance fixture](../testdata/company/README.md) builds the real
CLI and checks synthetic Go/Python/Java relationships, HTTP/MCP queries,
incremental refresh and recovery. The public demo has a separate check:

```sh
make test-showcase
```

To check an installed binary, set `DIFFMIND_BINARY` to its location and run
`sh scripts/test-showcase.sh`. Use a disposable workspace for manual checks;
never run fixture mutations against a company workspace.

## Optional integration checks

The scripts in `scripts/validation` cover boundaries beyond the normal suite.
Arguments are documented at the top of each script. Pass output directories
outside the checkout. Use new destinations and check each script's requirements
before running it.

| Scripts | Coverage |
| --- | --- |
| `prepare-boutique.py`, `public-correction.mjs` | Pinned public source, reviewed correction packs and graph rollback |
| `flask-pr.py`, `source-corpus.py`, `upstream-pr.py` | Runtime/source expectations and PR caller evidence |
| `provider-matrix.py` | Provider state, credential and redirect boundaries |
| `background-matrix.py` | Concurrent reads, refresh reuse and access checks on both queue backends |
| `recovery-upgrade.py` | Prior-to-current backup/restore and grant reconciliation |
| `systemd-native.py` | Backup service lifecycle on an isolated Linux host |
| `installed-host.py`, `company-host.py` | Installed agent capabilities, project scope and token revocation |
| `browser-audit.mjs`, `extractor-browser.mjs` | Browser layout, keyboard behavior and accessibility |

Browser checks require Playwright; accessibility checks also require
`@axe-core/playwright`. Set `DIFFMIND_PLAYWRIGHT_MODULE` and `DIFFMIND_AXE_MODULE`
to modules from your dependency installation. `DIFFMIND_BROWSER` selects
`chromium`, `firefox` or `webkit`; `DIFFMIND_CHROMIUM` optionally selects a
Chromium executable. Source/runtime checks document their pinned upstream
revisions and Python dependencies.

## Evidence and limits

Define expected relationships from source, runtime registration or independent
contract parsing before comparing analyzer output. Include near-match negatives,
unsupported cases and unknown truth. Test dirty/stale revisions, access loss,
missing patches and failed provider requests as well as successful cases.

Record the source revision, binary identity, environment, commands, results and
limitations in the PR or CI artifacts. Keep generated graphs, logs, transcripts,
screenshots, clones and backups outside Git. Reusable fixtures must be synthetic
and should contain only the inputs and expected behavior needed by a test.

Report static extraction, declared configuration, graph rendering and actual
agent behavior separately. Synthetic graphs and scripted model sessions do not
establish complete framework coverage, human usability or production reliability.
Native release checks must run on every supported target; cross-compilation alone
does not verify installation. PR heuristics and zero exact callers do not prove
that a change is safe to merge.
