# Historical evaluator fixtures

These labels belong to a removed extractor evaluator. The current CLI has no
`eval` command and the repository has no `internal/eval` package. The historical
`deterministic` flag and thresholds are not executable acceptance gates.

Retain these files as historical source material. Do not derive new ground truth
by copying candidate output, and do not report these labels as a passing current
accuracy benchmark. Current detector/protocol/workspace checks and their limits
are documented in [the validation guide](../../../docs/validation.md).

For a new corpus, record pinned source revisions, approved source scope and
expected relationships/callers independently before scoring the candidate.
Separate supported cases, unsupported patterns and unknown truth. Keep synthetic
rendering fixtures distinct from assembled-graph precision/recall and actual PR
review outcomes.
