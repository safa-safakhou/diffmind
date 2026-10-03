# Execution runbook

Clean supported Linux/macOS checkout; Go 1.26.6+, Git/C compiler, supported Node for frontend/Playwright. Disposable private homes, no production tokens.

## Baseline

~~~sh
git rev-parse HEAD
git status --porcelain
go version
go env GOOS GOARCH
go test ./...
go test -race ./internal/workspace/archgraph ./internal/workspace/query ./internal/workspace/ui ./internal/workspace/agentapi ./internal/workspace/mcpserver
go test ./internal/workspace/ui -run TestExactPRCallerRejectsCrossProtocolAndUnknownSurface -v -count=1
~~~

For frontend changes, before freezing, in internal/workspace/ui/web:

~~~sh
npm ci
node --test src/lib/*.test.js
node scripts/test-components.mjs
node node_modules/vite/bin/vite.js build
~~~

Commit product/assets, freeze SHA; artifacts outside checkout, no frontend rebuild between archive/browser.

## Native candidate

~~~sh
diffmind_trial=$(mktemp -d)
chmod 700 "$diffmind_trial"
diffmind_revision=$(git rev-parse HEAD)
diffmind_platform="$(go env GOOS)_$(go env GOARCH)"
go build -trimpath -ldflags "-X main.version=0.0.0-validation -X main.commit=$diffmind_revision" -o "$diffmind_trial/diffmind" ./cmd/diffmind
"$diffmind_trial/diffmind" version --json
go version -m "$diffmind_trial/diffmind"
go run ./scripts/release-check --package-binary "$diffmind_trial/diffmind" --archive "$diffmind_trial/diffmind_0.0.0-validation_${diffmind_platform}.tar.gz" --version 0.0.0-validation
DIFFMIND_BINARY="$diffmind_trial/diffmind" sh scripts/test-showcase.sh
~~~

Hash binary/archive via sha256sum Linux/shasum -a 256 macOS. Local version unpublished. Verifier actual installer/company-agent acceptance; showcase mechanisms, neither human outcomes. Repeat native Linux/macOS AMD64/ARM64 same revision. [CI matrix](../../.github/workflows/ci.yml) receipts must match candidate. Cross-compile alone NOT_RUN; WSL Linux not Windows release.

## Browser/public trial

Install isolated Playwright/Chromium. Replace absolute path placeholders:

~~~sh
export DIFFMIND_PLAYWRIGHT_MODULE=/absolute/node_modules/playwright/index.mjs
export DIFFMIND_CHROMIUM=/absolute/chromium
python3 scripts/validation/prepare-boutique.py "$diffmind_trial/public"
"$diffmind_trial/diffmind" pack lint testdata/workspace/boutique-correction
"$diffmind_trial/diffmind" pack test testdata/workspace/boutique-correction
node scripts/validation/public-correction.mjs "$diffmind_trial/diffmind" "$diffmind_trial/public" "$PWD/testdata/workspace/boutique-correction/pack.json"
~~~

Preparation output absent; retain pins/derived commits/results/screenshots/hashes. Observer Chromium/HTTP-MCP/scope/access/correction/cadence/rollback with synthetic identities, not proxy/human/blind reviewer.

## Recovery

Build agreed prior in separate clean checkout; hash/version; replace PRIOR_BINARY; output absent:

~~~sh
python3 scripts/validation/recovery-upgrade.py /absolute/PRIOR_BINARY "$diffmind_trial/diffmind" "$diffmind_trial/recovery" "$diffmind_trial/public/repositories"
~~~

Follow [recovery](../backup-recovery.md); isolate ingress until independent offboarding reconciled. Private backups/rollback local, export redacted results.

## Actual host/company/provider

[AGENT_SETUP](../../AGENT_SETUP.md), [agent operations](../agent-operations.md), [access](../project-access.md), [deployment](../company-deployment.md). Register frozen binary absolute path in isolated host profile. Source helpers rebuild: record hash/version/reconcile identity explicitly.

Actual model turns required A/U, actual proxy B01, authorized Enterprise matching-host read credential P02/P03. Missing prerequisite BLOCKED. No PR comments/contributions/releases/exposed homes authorized. No secrets committed.

Copy result/session templates, save redacted dated evidence under docs/research/<date>/<trial>/ with manifest/limits. Private recordings/source/backups outside repo. Reconcile candidate/native/scenarios/Linear and 42 audit/11 live findings. Later docs may explicitly reference frozen product identity. See [validation](../validation.md).
