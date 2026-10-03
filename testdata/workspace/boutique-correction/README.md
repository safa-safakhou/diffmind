# Evaluation-only Boutique correction

This narrow pack is a reproducible correction trial, not an automatically
installed official pack or general Kubernetes support.

It matches only double-quoted literal targets in deployment.yaml under seven
explicit gRPC address variable names. Comments, renamed variables, environment
placeholders and unrelated values are excluded by negative fixtures. It does
not parse arbitrary YAML semantics, evaluate profiles, or prove runtime calls.
The input deployment documents are copied verbatim from a pinned public source
next to the corresponding source directory; this expanded scope is explicit.
Positive/negative pack tests and exact graph assertions run without upstream
service execution. No proprietary fixture is included.

Run candidate pack lint/test here, then use scripts/validation/prepare-boutique.py
and scripts/validation/public-correction.mjs. The observer reports source-first
labels, declared-subset accuracy, graph provenance, role boundaries and rollback;
it never treats source-extracted extras as independently labeled facts.
