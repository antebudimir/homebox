# Core Behavior

Applies to every task. Fork invariants: `06`–`08`. Gates: `03-workflow.md`.

## Orient first

- Read `AGENTS.md` before your first edit: what this fork added matters as much as what upstream owns.
- Read the target file and its callers; grep for a sibling implementation before inventing a pattern.
- Check `git log --oneline -15` and `git status` before touching anything.
- State the plan when a task spans more than one file.

## Change minimally

- Make the smallest diff that solves the problem. No drive-by refactors, renames or reformatting.
- Match the surrounding code: naming, error handling, component structure, test layout.
- Reuse existing helpers, composables and components; do not add a second way to do the same thing.
- Do not add dependencies, env vars, config keys or API surface without a concrete need.
- Never hand-edit generated files: `backend/internal/data/ent/**`, `frontend/lib/api/types/data-contracts.ts`, `frontend/.nuxt/**`, lockfiles, build output.

## Verify

- Prove behaviour with a test or real command output, not by reading code.
- Run the gates in `03-workflow.md`; if a gate cannot run here, say so explicitly.
- Fix a failing test or explain precisely why it is unrelated — never leave the suite red.

## Guardrails

- Never commit secrets: the auth pepper, tokens, personal inventory data, `.db` files.
- Ask before irreversible calls: schema changes, migration edits, public API shape, image names.
- Never publish to upstream registries (Docker Hub, `ghcr.io/sysadminsmedia/*`); this fork is GHCR-only (`08-upstream-sync.md`).
- Do not remove fork features to satisfy a request; propose a fork-shaped alternative instead.
