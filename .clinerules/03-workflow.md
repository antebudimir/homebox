# Workflow

## Gates, in order

```bash
# 1. backend
cd backend && go build ./... && go test ./...

# 2. frontend
cd frontend && pnpm install --frozen-lockfile && pnpm run typecheck

# 3. full PR set (regenerates swagger + TS types, lints, runs the suite)
task pr
```

- Run the narrowest failing test first (`go test ./internal/data/repo -run TestName`, `pnpm exec vitest run <file>`), then the full suite before push.
- CI mirrors the backend/frontend gates in all four publish workflows (a `test` job, with `build`/`publish` requiring it).
- No Go toolchain locally: say so explicitly and run the gate in a container (`podman run --rm -v "$PWD":/src -w /src/backend golang:1.27 bash -c 'go build ./... && go test ./...'`) or wait for CI. Never tag on unverified backend code.

## Commits

- One logical change per commit. Prefix: `feat:`, `fix:`, `fix(area):`, `chore:`, `docs:`, `test:`, `ci:`.
- API surface changed? Regenerate the derived files before committing: `task swag --force` (swagger/OpenAPI, also copied into `docs/public/api/`) and `task typescript-types --force` (`frontend/lib/api/types/data-contracts.ts`), or just `task generate`.
- Update `WORKLOG.md` when an operating rule, incident or CI policy changes.
- Update `docs/` in the same commit when behaviour, flags or env vars change.
- A schema change ships as one commit: ent schema + generated code + goose migration (`06-backend.md`).

## Never commit

- Generated or build output: `frontend/.nuxt/`, `node_modules/`, `backend/build/`, `backend/.data/`.
- Secrets or data: the auth pepper, tokens, `*.db`, `backend/.env`.
- Anything you have not run the gates against.

## When stuck

- Reproduce with the smallest command, then read the actual error.
- Check whether the behaviour is a fork feature or an upstream one (`AGENTS.md`).
- Report the blocker with the exact command and output instead of guessing.
