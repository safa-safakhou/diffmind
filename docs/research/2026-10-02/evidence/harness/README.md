# Retained observer scripts

These are text copies of the live study's observer scripts, saved to explain the actual operations and selectors. They are research evidence, not product source or a supported automation suite. They contain the observer mistakes described in [the reproduction record](../../reproduction.md). A later script continued an earlier failed attempt; the final company script is the corrected version whose observation file is retained.

Paths were substituted with `<study-root>`, `<source-root>` and `<observer-playwright-module>`. The original helper derived the study root from its own script directory and used an existing Playwright installation. These `.txt` artifacts must not be executed as if they were ready-to-run code. A repeat requires recreating the disposable setup, supplying the browser runtime and respecting the state dependencies between scripts.

| Artifact | Purpose |
| --- | --- |
| [helpers.mjs.txt](helpers.mjs.txt) | JSON observation writer, real stdio client, loopback server lifecycle and Chromium helper. |
| [agent-study.mjs.txt](agent-study.mjs.txt) | Initial individual-agent journey. Includes the mistaken final contract tool call. |
| [agent-continuation.mjs.txt](agent-continuation.mjs.txt) | Correct contract comparison and unchanged ingestion. |
| [explorer-study.mjs.txt](explorer-study.mjs.txt) | Browser-first project/import and modal interaction. Includes the preview logger's wrong array key. |
| [explorer-continuation.mjs.txt](explorer-continuation.mjs.txt) | Failure inspection, PATH recovery, graph exploration and empty states. Includes the ineffective initial summary selector. |
| [company-study.mjs.txt](company-study.mjs.txt) | Final corrected scoped company trial using synthetic identities and temporary secrets. |
| [company-completion.mjs.txt](company-completion.mjs.txt) | Completed ingestion readback followed by an unsuccessful guessed GET route probe. |
| [public-pr-study.mjs.txt](public-pr-study.mjs.txt) | Public ingestion/listing; its generic helper initially misroutes the PR impact read. |
| [public-pr-continuation.mjs.txt](public-pr-continuation.mjs.txt) | Correct PR impact read, browser selection and refresh. |
| [provenance-study.mjs.txt](provenance-study.mjs.txt) | Source provenance and the PR company-context screenshot. |

The company script removes/revokes only prior study grants/tokens in the disposable fixture home to reset its own trials. Its actions are not instructions to reset a real company workspace. The generated secret values were never copied here; the artifact shows their random generation and use, not actual credentials.
