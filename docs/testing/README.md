# DiffMind candidate test pack

Status: operator execution and fixes recorded in the [4 October system validation report](../research/2026-10-04/system-validation/README.md), with a [per-scenario result](../research/2026-10-04/system-validation/scenarios.md). Tested product: 2cd9d328cc24de955184fc8aab863892a6ae91b3. External host, independent-reviewer, human, Enterprise and other-native-platform gates remain open. The protocol was originally prepared against c02bf351b913bef6ddb84b72be55d56ca90d94cb; earlier reports are historical evidence.

1. Freeze a clean candidate; follow [runbook](runbook.md).
2. Independently freeze source labels before output; follow [corpus protocol](corpus.md).
3. Run [scenarios](scenarios.md); keep every subcase result.
4. Conduct [unprompted persona sessions](usability-study.md).
5. Copy [result template](result-template.json) and [session template](session-template.md).

| Task | Required scenarios | External gate |
| --- | --- | --- |
| [MNI-180](https://linear.app/mnim/issue/MNI-180) | A01-A06, U01-U03 | Actual installed host/unprompted model turns |
| [MNI-192](https://linear.app/mnim/issue/MNI-192) | P01-P05 | Target Enterprise endpoint/matching-host credential |
| [MNI-197](https://linear.app/mnim/issue/MNI-197) | C01-C03, R01-R08 | Independent blind labels/agreement/actual upstream PRs |
| [MNI-198](https://linear.app/mnim/issue/MNI-198) | R01-R08 | Independent revision/deletion corpus |
| [MNI-199](https://linear.app/mnim/issue/MNI-199) | R01-R08, U01-U03 | Independent labels/human comprehension |
| [MNI-201](https://linear.app/mnim/issue/MNI-201) | U01-U03, A01-A06 | Three observations per original persona minimum |
| [MNI-202](https://linear.app/mnim/issue/MNI-202) | All/native matrix | Exact candidate; 42 audit/11 live findings reconciled |

PASS requires expected behavior plus retained evidence; FAIL records mismatch; BLOCKED identifies missing prerequisite; NOT_RUN means not executed. UNKNOWN is evidence classification, not a passing result. One passing subcase cannot hide a blocked subcase. Preparation closes no task.

Run baseline, access/PR controls, independent corpus/users, then final reconciliation. Product fixes require affected checks on the new candidate; later documentation commits may explicitly reference frozen product identity. Critical activation/access/evidence defects cannot be averaged away. Existing issue criteria define acceptance; report measured timing/accuracy/adoption without inventing thresholds after results.

Administrator operation supports, rather than replaces, the original individual, company-joiner and architecture-explorer personas.
