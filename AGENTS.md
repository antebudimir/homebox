# AGENTS.md

Orientation for agents landing in this repo. Conventions live in `.clinerules/`; this file is the map.

## What this is

A hard fork of [Homebox](https://github.com/sysadminsmedia/homebox) (upstream: home inventory and organization) adding a builds system for tracking skateboard and bicycle parts, with sport tagging and usage sessions.

- **Added:** builds (formerly "setups"), build positions, usage sessions, sport field (schema-backed `sport` / `DefaultSport`), `ShowBuilds` kill-switch, light/dark theme toggle, build part template seeding on group creation.
- **Kept:** everything else upstream.

## Where things live

- `backend/` — Go, ent ORM, goose migrations. `app/` (API, CLI), `internal/data/` (`ent/schema`, `migrations/{sqlite3,postgres}`, `repo`), `internal/core/` (services), `pkgs/`.
- `frontend/` — Nuxt 4 / Vue 3, pnpm. `pages/`, `components/` (incl. vendored `components/ui`), `composables/`, `stores/`, `locales/` (i18n JSON), `test/e2e/` (Playwright).
- `docs/` — documentation site. `.github/workflows/` — CI/CD. `WORKLOG.md` — operating history.
- `Taskfile.yml` — dev commands (`task pr`, `task db:generate`, `task go:test`, `task ui:check`).

## Gates

- Backend: `cd backend && go build ./... && go test ./...`
- Frontend: `cd frontend && pnpm run typecheck`
- Full PR set: `task pr`

## Sync

- Merge from upstream (`git fetch upstream && git merge upstream/main`), never the GitHub "Sync fork" button.
- Tag `vX.Y.Z-fork.N` only after the gates pass.
- The incident (tag pushed before any backend gate ran) is recorded in `WORKLOG.md`.

## Non-negotiables

- The kill-switch must fully hide Builds when off: nav entry, `/builds` route, sport field.
- Sport grouping is schema-backed, never string-prefix matching.
- Migrations are goose, in both `sqlite3/` and `postgres/`, verified up and down.

## Audience split

- `README.md` is for people: user-facing changes only — no process meta-language, no links to this file or `.clinerules/`.
- `AGENTS.md` and `.clinerules/` are for agents: conventions, gates, footguns. Keep the two audiences separate.

## Rules index

- `01-core-behavior.md` — how to work in this repo
- `02-code-quality.md` — Go/Vue expectations, tests, hygiene
- `03-workflow.md` — gates, commits, what never to commit
- `04-ux.md` — theming, states, interaction, copy
- `05-git.md` — branch safety, commit style, tags
- `06-backend.md` — ent + goose, test gate, CI gating, GHCR auth
- `07-frontend.md` — theme state, builds kill-switch, i18n
- `08-upstream-sync.md` — merge playbook and fork invariants
