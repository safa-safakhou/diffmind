# DiffMind candidate test pack

Status: Codex owns all verification; no human sign-off is required. The [4 October autonomous verification report](../research/2026-10-04/autonomous-verification/README.md) records ten successful installed-model cases, three additional fixes, the runtime corpus's contract limit, and remaining technical coverage. Latest product is 66cf19a1e7c86223979bf284b06336117ff3ded6; each browser/model cohort retains its own exact identity. The [earlier broad system report](../research/2026-10-04/system-validation/README.md) remains historical evidence. [Automated journey verification](usability-study.md) and [source/runtime oracles](corpus.md) supersede the former human-review protocol.

1. Freeze a clean candidate; follow [runbook](runbook.md).
2. Freeze source/runtime oracle expectations before verification; follow [corpus protocol](corpus.md).
3. Run [scenarios](scenarios.md); keep every subcase result.
4. Execute [automated journey sessions](usability-study.md).
5. Copy [result template](result-template.json) and [session template](session-template.md).

| Task | Required scenarios | Verification dependency |
| --- | --- | --- |
| [MNI-180](https://linear.app/mnim/issue/MNI-180) | A01-A06, U01-U03 | Actual installed host/unprompted model turns |
| [MNI-192](https://linear.app/mnim/issue/MNI-192) | P01-P05 | Target Enterprise endpoint/matching-host credential |
| [MNI-197](https://linear.app/mnim/issue/MNI-197) | C01-C03, R01-R08 | Source/runtime oracle agreement and actual upstream PRs |
| [MNI-198](https://linear.app/mnim/issue/MNI-198) | R01-R08 | Automated revision/deletion corpus |
| [MNI-199](https://linear.app/mnim/issue/MNI-199) | R01-R08, U01-U03 | Frozen source oracles and observed model interpretation |
| [MNI-201](https://linear.app/mnim/issue/MNI-201) | U01-U03, A01-A06 | Three automated cases per original journey; no human sign-off |
| [MNI-202](https://linear.app/mnim/issue/MNI-202) | All/native matrix | Exact candidate; 42 audit/11 live findings reconciled |

PASS requires expected behavior plus retained evidence; FAIL records mismatch; BLOCKED identifies missing prerequisite; NOT_RUN means not executed. UNKNOWN is evidence classification, not a passing result. One passing subcase cannot hide a blocked subcase. Preparation closes no task.

Run baseline, access/PR controls, source-oracle corpus and installed-host cases, then final reconciliation. Product fixes require affected checks on the new candidate; later documentation commits may explicitly reference frozen product identity. Critical activation/access/evidence defects cannot be averaged away. Existing issue criteria define acceptance; report measured timing/accuracy/adoption without inventing thresholds after results.

Administrator operation supports, rather than replaces, the original individual, company-joiner and architecture-explorer personas.
