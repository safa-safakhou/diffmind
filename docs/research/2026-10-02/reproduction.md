# Live study execution and reproduction record

This records the actual environment and task sequence behind [the live report](live-usability-study.md). Reproduction may produce different GitHub data, ports, timestamps and graph run IDs. The retained scripts are observer artifacts, not a maintained product test framework.

## Isolation and preparation

The original source checkout remained the working source, at `5d2548514bfc5d920dfd40def5fcacd80162e9ae`. All build products, cloned repositories, fixture commits, generated tokens, homes and runtime state were placed in an external disposable directory. The application sources and embedded dashboard were unchanged. Only this documentation tree was added to the checkout.

The following commands show the actual operations with portable placeholders. `<study-root>` was a fresh directory created outside the repository. Placeholders must be replaced before running; they are not shell variables.

```text
go build -o <study-root>/diffmind ./cmd/diffmind
sh scripts/prepare-showcase.sh <study-root>/demo-shop
env -u GITHUB_TOKEN -u GH_TOKEN go run ./scripts/agent-setup \
  --repo-root <source-root> --home <study-root>/bootstrap-home \
  --bin-dir <study-root>/installed-bin --name diffmind-live-study
```

The source bootstrap output is retained in [bootstrap-output.json](evidence/bootstrap-output.json). No generated configuration was merged into a personal agent host.

The driver stripped inherited `DIFFMIND_*`, `GITHUB_TOKEN` and `GH_TOKEN` variables from application processes, then supplied the designated disposable `DIFFMIND_HOME`. Public PR clients additionally used an empty temporary `GH_CONFIG_DIR`. The anonymous clone disabled system/global Git configuration and prompting:

```text
env GIT_TERMINAL_PROMPT=0 GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null \
  git clone --depth 1 https://github.com/safa-safakhou/diffmind.git \
  <study-root>/public-repositories/diffmind
```

The cloned commit was verified to match the reviewed source. No writes were made to GitHub.

## Runtime and browser

The driver used Node 24.12.0 and an already available Playwright installation, launching headless Chromium 153.0.8010.12. It did not install product dependencies or rebuild the frontend. Main viewport: 1440 × 900. Additional captured sizes: 1280 × 720 and 1920 × 1080.

Browser-first servers used the compiled binary with:

```text
<study-root>/diffmind ui --no-spa-rebuild --host 127.0.0.1 --port <free-port>
```

The first browser trial intentionally left the study binary directory out of PATH. A later process added it to PATH and reused the same explorer home, establishing the observed failure and recovery. That environment adjustment was confined to the study process.

The company process used scoped access and randomly generated temporary admin/proxy secrets. Each browser context represented a different synthetic subject using `X-DiffMind-User`, `X-DiffMind-Role` and the trusted proxy secret header. Secret values and request authorization headers were not logged. This is a test of DiffMind behind its trust boundary, not a production authentication deployment.

## Actual sequence

1. Build and bootstrap; prepare six synthetic Git repositories.
2. `agent-study.mjs`: connect an empty home, discover tools, create a project, preview/import, build/query the graph, apply the existing fixture change, observe five seconds without refresh, disconnect/reconnect and explicitly refresh.
3. `agent-continuation.mjs`: discover the correct contract tool schema, compare the two completed runs, repeat unchanged ingestion and retain run IDs.
4. `explorer-study.mjs`: use the empty dashboard, inspect first-use focus/labels, create a project, preview candidates, reopen the form, submit an invalid root and attempt ingestion.
5. `explorer-continuation.mjs`: inspect the failed ingestion, restart with corrected process PATH, refresh, measure initial/focused graph labels, inspect evidence, visit local PR and single-snapshot comparison states and capture extra viewports.
6. `company-study.mjs`: start the scoped service against completed fixture state; exercise ungranted viewer, UI grants, viewer prohibition, one-time agent token, HTTP MCP queries, editor refresh, snapshot comparison, verified subject removal and explicit token revocation. The retained observation file is the corrected final trial.
7. Clone the public repository anonymously.
8. `public-pr-study.mjs`: ingest the clone and retrieve live PRs through MCP and browser.
9. `public-pr-continuation.mjs`: use the advertised read operation for PR #9, select it in the browser and refresh the list.
10. `provenance-study.mjs`: query the public service/dependencies and capture PR company context to identify where the unexpected queue and example routes came from.
11. `company-completion.mjs`: read the final company ingestion back through MCP; inspect the persisted job record separately to verify the accepted refresh succeeded.
12. Copy path-normalized JSON and unchanged screenshots into this research directory, record hashes and stop owned application processes.

The retained [observer scripts](evidence/harness/README.md) document the calls and selectors. Some originally failed attempts remain visible in observation logs, with explanations below. They should not be treated as polished commands that succeed from a clean machine without their setup/state dependencies.

## MCP and task protocol

Stdio sessions executed the real `diffmind agent` process and exchanged newline-delimited JSON-RPC. The initialization protocol version was `2025-11-25`, followed by `notifications/initialized` and `tools/list`. Named tool calls used the advertised schemas. Async ingestion was polled using `inspect_workspace` with `get_ingestion` until its status completed or failed.

The principal management sequence was:

```text
list_projects
get_graph_summary                        # empty-state check
agent_runtime {action: status}
manage_workspace {operation: create_project, body: {...}}
manage_workspace {operation: import_repositories,
                  selectors: {pid}, body: {provider: local, root, dry_run: true}}
manage_workspace {operation: start_ingestion, selectors: {pid}, body: {...}}
inspect_workspace {operation: get_ingestion, selectors: {pid}}
get_graph_summary / get_dependencies / get_service / get_contracts
compare_contracts {project, from, to, service: checkout}
inspect_workspace {operation: list_pull_requests, selectors: {pid}}
inspect_workspace {operation: pull_request_impact,
                   selectors: {pid, repo_id: diffmind, number: "9"}}
```

The company agent made actual HTTP MCP POSTs to `/mcp` with a project-scoped viewer bearer token, accepting JSON or SSE responses. Its initialization, tool discovery and graph queries are recorded. A direct viewer refresh attempt tested backend authorization independently of disabled UI controls.

The controlled source change was the existing repository script:

```text
sh scripts/apply-demo-shop-change.sh <study-root>/demo-shop
```

It modified and committed the disposable checkout repository only. The original source fixture in this repository was not edited.

## Retained outcomes and identifiers

| Evidence state | Identifier/result |
| --- | --- |
| Individual project | `live-agent-demo-shop` |
| Baseline fixture graph | `20261002T161959Z` |
| Changed fixture graph | `20261002T162005Z` |
| Unchanged repeat graph | `20261002T162104Z` |
| Browser-first project | `live-explorer-demo-shop` |
| Browser-first recovered graph | `20261002T162610Z` |
| Company final refresh | `refresh-4615f0ec8f9e6df5337ada4f9ae29f84` |
| Company final graph | `20261002T163243Z`; zero analyzed, six reused |
| Public trial project | `public-diffmind-pr-trial` |
| Public graph | `20261002T163244Z` |
| Public PR | `safa-safakhou/diffmind#9`, observed risk 19/low, graph context stale and score-ineligible |

## Observer mistakes and exclusions

These distinctions prevent tooling mistakes from becoming invented product defects:

- The initial agent script attempted `diff_contracts`, which is not an advertised tool. The product correctly rejected it. The continuation discovered and used `compare_contracts`; only that comparison supports the result.
- The generic driver helper initially routed `pull_request_impact` to `manage_workspace`. The product instructed the caller to use `inspect_workspace`. The corrected operation returned 200 and is the retained PR result.
- In an earlier company attempt, the observer removed the first membership row after the server had sorted rows, targeting the editor rather than the viewer. That attempt was discarded. The final run finds and verifies the viewer subject value before removal; the resulting graph disappearance is recorded.
- An earlier discarded company attempt treated native Fetch's `status` property as a function. This was corrected in the driver.
- The explorer continuation waited for a nonexistent ingestion-summary selector, introducing a 30-second wait. That wait is excluded from product timing. Its other screenshots/readbacks support the observed failure. The first explorer script's separate graph wait ended because analysis had actually failed; backend ingestion results corroborate that failure.
- The browser preview logger looked for names under a `repos` key, while the backend's candidate array is `results`. Its logged `candidate_names: []` must not be interpreted as the API returning no names. MCP preview evidence contains all six names. The UI's absent preview result is supported by its displayed body/screenshot and closed modal.
- The completion probe guessed an unadvertised `list_refresh_jobs` operation and then an unavailable GET route with an `/api/v1` prefix. Those observer calls are excluded from product conclusions. The completed ingestion MCP result and persisted job record independently verify the final refresh outcome.
- Preliminary driver parsing/selector corrections are not application crashes. No browser JavaScript exception was recorded in the retained browser runs. Expected 400/403/404/401 results represent input validation and access tests, not an assertion that the application emitted no failed requests.

The retained `company-completion.json` contains the successful ingestion readback and orderly disconnect preceding the final probe's JSON parsing failure. The saved job is copied from the disposable application's durable job storage. The report does not claim that the guessed GET route succeeded.

## Remaining limits

No real agent host registration, human participants, screen-reader session, production company authentication deployment, private company clone, cross-browser matrix, long-duration soak, restore exercise or automatic PR posting was performed. The report makes no measured claims about those paths. Every owning agent/UI process was stopped at the end of its retained trial.
