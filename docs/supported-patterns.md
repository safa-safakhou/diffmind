# Tested patterns and limitations

Diffmind is deterministic static analysis, not a runtime dependency inventory.
These are reproducible baselines, not a claim that every use of a framework is
supported. A registered detector is not equivalent to complete framework support.

| Surface | Tested baseline | Verification | Important limits |
| --- | --- | --- | --- |
| Go `net/http` | Literal package-level Get/Head/Post/PostForm calls; import aliases and shadowing negatives | `internal/extractor/detectors/languages/golang/httpclient/nethttp/nethttp_test.go`, company acceptance | Dynamic URLs, custom client wrappers and request construction are outside this detector |
| Python `requests` | Literal URLs and unshadowed same-file uppercase string constants; recognized calls with unresolved dynamic destinations remain visible | `internal/extractor/stage/discovery/calls_test.go` | Imported constants, expressions, wrappers and runtime configuration remain unresolved |
| Python Redis | Literal `redis.Redis`/`StrictRedis` host, port and database identity; anonymous clients are scoped to their owning service | `internal/extractor/stage/discovery/calls_test.go`, `internal/workspace/archgraph/route_match_test.go` | Dynamic constructor configuration does not prove shared cache identity; use a knowledge pack for intentional shared resources |
| Python Django | A simple service route used in a cross-repository HTTP relationship | `testdata/company/catalog`, company acceptance | Not a guarantee for arbitrary routing, middleware or dynamic URL construction |
| Java Spring MVC | A simple service endpoint used in a cross-repository HTTP relationship | `testdata/company/billing`, company acceptance | Not a guarantee for arbitrary annotation composition or runtime configuration |
| Helm identity (built in) | Name, ingress aliases, owned queues, port metadata | `packs/helm-values/testdata`, pack tests | Conventional `deploy/values.yaml`; no Helm rendering or chart evaluation |
| Explicit service manifest (opt in) | HTTP/RPC targets, aliases, regex service resolution, queue publishers/consumers, external HTTP | `packs/service-manifest`, pack graph acceptance over HTTP/MCP | Declarations are not proof of a call or entrypoint reachability; no resource creation or infrastructure execution |
| Spring Cloud OpenFeign configuration (opt in) | Literal client URLs in base application YAML/YML and single-line properties; comments/placeholders/credentials negative cases | `packs/spring-openfeign-config` | No profile selection, environment/config-server expansion, annotation precedence, load-balancing discovery or escaped/continued properties |
| OpenAPI 3.0 request contracts | Local `openapi`/`swagger` YAML or JSON; operation/path parameters, JSON request properties, required/nullable/type metadata, bounded local `#/components/schemas/*` references | `internal/extractor/stage/discovery/openapi_test.go`, query contract tests | Exact method/path matching only; 4 MiB/document, 20 documents/repository, 1,000 fields/operation and depth 8; no remote refs/fetching, OpenAPI 3.1, composition flattening, callbacks or runtime-truth claim |

The real-binary company acceptance test imports synthetic Go/Python/Java Git
repositories and checks exact service relationships and source evidence through
HTTP and MCP. Pack graph acceptance separately verifies installation/integrity,
the graph job, query transport, disabling packs and immutable historical graphs.
Neither substitutes for validating your company's known architecture.

## Version evidence and detector releases

Each analysis records dependency declarations, effective versions where they can
be established, module scope, source files and resolution limitations. Supported
inputs include npm/pnpm/Yarn locks, requirements and pyproject with uv/Poetry
locks, Maven POM inheritance and cached pinned parents/BOMs, literal Gradle
coordinates, Go module requirements, and Gemfile declarations/locks. These
inputs are static evidence; they do not prove the deployed environment.

Maven resolution uses matching local parents or exact pinned cache artifacts; it
does not fetch or choose the newest cached parent. Gradle variables/catalogs,
Maven profiles, conditional Python environments, non-semver versions and Go
replacements can remain unresolved. Module-specific versions are selected before
root declarations. Unknown/future versions still receive built-in source-backed
analysis, with explicit coverage limits.

The UI service inspector and saved JSON metrics expose the same coverage states:

| State | Meaning |
| --- | --- |
| `validated` | Exact version has a registered static extraction fixture. |
| `compatible_unverified` | Version matches applicability, but has no exact fixture. |
| `unknown_version` | Version or environment condition could not be established. |
| `unsupported_version` | Version is outside the declared rule range; source fallback used. |
| `ambiguous_version` | Multiple effective versions conflict in this module. |
| `unversioned` | Detector has no version-specific coverage declaration. |

The initial exact fixture matrix covers Express 4.21.2/5.1.0, Spring Web
6.2.0/7.0.0, FastAPI 0.115.0 and Flask 3.1.0. This validates the bounded patterns
in the fixtures, not all framework behavior. Express 5 removed path syntax is
rejected when version 5 is established. ORM detectors have applicability and
version evidence too; applicability alone does not make them validated.

Detector revision, applied rule and coverage are retained independently of the
DiffMind application version. A new framework/library release requires upstream
review and fixtures, followed by a DiffMind release updating that matrix—even
when the extraction implementation is unchanged. There is no automatic claim
that the latest upstream release is supported. Company libraries can use
[version-gated custom configuration](extractor/docs/CONFIGURATION.md#patterns).

Flow review reports dependency/runtime/build-tool version changes separately.
Changing version provenance or detector metadata does not fabricate a changed
execution flow. A dependency-only PR can therefore have version changes and no
structurally changed flows; source analysis cannot establish runtime compatibility.

Relevant upstream specifications:
[npm locks](https://docs.npmjs.com/cli/v11/configuring-npm/package-lock-json/),
[Maven dependency management](https://maven.apache.org/guides/introduction/introduction-to-dependency-mechanism.html),
[uv project layout](https://docs.astral.sh/uv/concepts/projects/layout/), and
[Express 5 migration](https://expressjs.com/en/guide/migrating-5/).

## Full repository corpus validation

Use `scripts/validation/corpus-audit.py` with a private JSONL inventory containing
`repository`, `commit`, `path`, `archived` and `status` (`downloaded` or `empty`)
for pinned local snapshots:

```sh
go build -o /tmp/diffmind-audit ./cmd/diffmind
python3 scripts/validation/corpus-audit.py \
  --inventory /private/audit/snapshots.jsonl --binary /tmp/diffmind-audit \
  --out /private/audit/results --workers 2
```

Every non-empty repository runs extraction and Protocol validation. Results
retain commit IDs, errors, evidence counts, coverage and version-resolution
limits; active and archived repositories are counted separately. The script
does not build/install the repository's dependencies or execute its code. Keep
company sources and reports outside this repository. A successful run proves
extraction/schema validity, not recall or correctness of every relationship;
promote reviewed misses into synthetic positive and negative fixtures.

## Teaching a missing convention

1. Reduce it to a **synthetic** repository/configuration example.
2. If it is a declarative convention, start with `diffmind pack init ./my-pack`.
   Adapt a field path or narrow regex; add exact detection and graph assertions.
3. Include a near-match negative fixture so a rule cannot pass by overmatching.
4. Run `diffmind pack lint ./my-pack`, `diffmind pack test ./my-pack` and
   `diffmind pack explain ./my-pack --repo ./test-repo`.
5. Install the pack, rebuild the project graph, and compare it against known
   relationships. Contribute the pack, fixtures and documented limits together.

For semantics such as shadowing, call reachability or framework annotations,
extend the AST detector and its tests instead of relying on a broad code regex.
Packs cannot run code, execute a shell, call an LLM, or load an arbitrary plugin.
See [pack authoring](knowledge-packs.md) and [contributing](../CONTRIBUTING.md).

The OpenFeign pack follows the configured URL key described in the
[official Spring Cloud OpenFeign documentation](https://docs.spring.io/spring-cloud-openfeign/reference/spring-cloud-openfeign.html).
Those rules also explain why configuration alone does not establish the effective
runtime target: an annotation URL can take precedence.
