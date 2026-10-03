# Validating a DiffMind candidate

Run from the source checkout with Go 1.26.6 or newer, Git and a C compiler.

```sh
go test ./...
go test -race ./internal/workspace/archgraph ./internal/workspace/query ./internal/workspace/ui ./internal/workspace/agentapi ./internal/workspace/mcpserver
go build -o bin/diffmind ./cmd/diffmind
sh scripts/test-showcase.sh
```

The showcase script prepares six public-safe synthetic repositories in a new
temporary directory, runs the actual analyzer and applies the canonical checkout
change. It removes only its own temporary directory. It requires the preceding
`bin/diffmind` build. It validates extraction mechanisms and source-change
regressions, not universal semantic coverage or PR accuracy.

To check an already installed candidate without replacing a local build, use
`DIFFMIND_BINARY=/absolute/path/to/diffmind sh scripts/test-showcase.sh`.

For frontend changes, use a supported Node.js runtime in
`internal/workspace/ui/web`:

```sh
node --test src/lib/*.test.js
node scripts/test-components.mjs
node node_modules/vite/bin/vite.js build
```

Install frontend dependencies first if absent. Retain the resulting embedded
bundle with the candidate. Browser tests must use the candidate containing that
bundle; a screenshot from an older binary is not candidate verification.

For an isolated source installation, follow [AGENT_SETUP](../AGENT_SETUP.md).
Use a fresh private home and user-owned binary directory. Query-only `mcp`,
full-management `agent`, and company HTTP `/mcp` are distinct contracts.
Verify actual installed-client discovery and source-backed queries before
reporting a working connection; prescribed tool calls do not demonstrate
unprompted adoption or human comprehension. See [agent operations](agent-operations.md)
and [company access](project-access.md).

The current CLI has **no `eval` command**. References to `diffmind eval --mode
cheap`, `score-run`, `variance`, `floor-coverage`, `internal/eval` or `internal/floor`
in historical extractor designs and fixture notes are obsolete. Preserve those
records as history, not installation or validation instructions.

The [pinned public benchmarks](../examples/public-benchmarks/README.md) report
bounded compatibility observations. Detector counts, rendered graph scale,
assembled graph correctness, actual installed-agent behavior and human outcomes
must be reported separately. Expected labels should precede candidate scoring.
Unknown truth is not a false positive or false negative. Synthetic PR fixtures
do not establish real code-host review accuracy or calibrated scores.

Native support covers macOS/Linux on AMD64/ARM64; platform support and the
platform actually tested for a candidate are separate facts. Record the exact
binary, source revision, environment, commands, retained evidence and limitations.
The [release verifier](distribution.md) exercises installed archives; publishing a
release is a separate action. Background maintenance and access migration retain
their [operational](operations.md) and [shared deployment](company-deployment.md)
requirements.

## Remaining acceptance scenarios

Use the [test scenario pack](testing/README.md) for executable checks, independent
labels, unprompted persona sessions, Enterprise trials and evidence templates.
Preparation is not an execution or completion receipt.
