# Backend

## Stack and schema

- Go, ent ORM, goose migrations.
- Any schema change = ent schema edit in `backend/internal/data/ent/schema/` + `task db:generate` (runs `go generate`) + a goose migration.
- Migrations live in `backend/internal/data/migrations/sqlite3/` **and** `postgres/`; both are required.
- Never edit an applied migration; every migration must be verified up **and** down.
- `backend/internal/data/migrations/migrations_test.go` asserts the resulting builds schema — extend it when you add a migration.

## Tests are the gate

- `cd backend && go build ./... && go test ./...` before any push or tag.
- Repo tests use in-memory SQLite (`backend/internal/data/repo/main_test.go`); tests must never run against the dev database.
- `go build` proves compilation only — it is not the gate.

## CI gating

- All four publish workflows (`docker-publish-rootless`, `docker-publish`, `docker-publish-hardened`, `publish-image`) run a test job (`go build ./... && go test ./...`, `pnpm typecheck`) with build/publish requiring it, plus `concurrency: cancel-in-progress` and `workflow_dispatch`.
- The `publish-image` gate is wired but not yet exercised by CI (tag/manual-trigger only) — verify it fires on the next tag.

## GHCR auth

- GHCR pushes use `GITHUB_TOKEN`, which requires the package to grant this repo Actions write access (Package settings → Manage Actions access).
- A push failing with `permission_denied: write_package` means that link is missing.

## Auth pepper

- `HBOX_AUTH_API_KEY_PEPPER` is required and must be ≥32 bytes (`backend/app/api/main.go`).
- Rotating it invalidates issued API keys, not user passwords.
- It must be backed up; never commit it.
