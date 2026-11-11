# WORKLOG

Operational notes for this fork (`antebudimir/homebox`). Not user-facing documentation.

## Fork sync strategy: merge + tag

Upstream is `sysadminsmedia/homebox`, tracked locally as the `upstream` remote.

1. `git fetch upstream`
2. `git merge upstream/main` (never rebase — fork commits stay intact and traceable)
3. Resolve conflicts, keeping fork features (builds, sport tagging, theme toggle, builds kill-switch)
4. Run the backend/frontend gate (see below)
5. Tag `vX.Y.Z-fork.N` **only after the gate is green**
6. `git push origin main && git push origin vX.Y.Z-fork.N` (push the single tag, never `--tags`, or upstream's tags land in the fork)

## Gate order (mandatory)

**merge → go test (CI or local) → tag → publish → pull to server**

- The image build only proves the code compiles. It does not prove behavior.
- A sync that touches shared code next to fork code (e.g. upstream's `RemoveMember`/`DeleteGroup`
  refactor next to the fork's `ensureBuildPartTemplates` call) must be covered by the Go tests.
- Local equivalent of the CI gate:

  ```bash
  cd backend && go build ./... && go test ./...
  cd ../frontend && pnpm install --frozen-lockfile && pnpm run typecheck
  ```

- Never pull a new image for the server until the tests have actually passed somewhere
  (CI run green, or a recorded local run).

## Image publishing policy

- GHCR only, this fork's namespace: `ghcr.io/antebudimir/homebox`.
- Never push to Docker Hub or to any `ghcr.io/sysadminsmedia/*` registry.
- Build cache lives in this fork's own package: `ghcr.io/antebudimir/homebox-devcache`.
- All four image workflows follow this policy (`docker-publish-rootless.yaml`, `publish-image.yaml`,
  `docker-publish.yaml`, `docker-publish-hardened.yaml`): Docker Hub env/logins/push/attest removed,
  upstream `binfmt`/`buildkit` images replaced with action defaults, `workflow_dispatch` and
  `concurrency: cancel-in-progress` added.
- Test gate (`go build ./... && go test ./...`, `pnpm install --frozen-lockfile && pnpm run typecheck`)
  runs in `docker-publish-rootless.yaml` and `publish-image.yaml` (the two that produce the images
  the server runs). The `docker-publish` and `docker-publish-hardened` variants publish only the
  plain/hardened variants and skip the gate to avoid tripling test cost on every main push.

## GHCR auth (pending manual step)

- The `ghcr.io/antebudimir/homebox` package was created by a PAT push (`publish-image.yaml` used
  `secrets.CR_PAT`), so it is not linked to the repository for `GITHUB_TOKEN` writes. Pushing with
  `GITHUB_TOKEN` fails with `denied: permission_denied: write_package`.
- Until the package grants this repo write access, all workflows log in with `secrets.CR_PAT`.
- To finish the cleanup: package settings → **Manage Actions access** → add repository
  `antebudimir/homebox` with **Write**, then replace `password: ${{ secrets.CR_PAT }}` with
  `password: ${{ secrets.GITHUB_TOKEN }}` in the four image workflows. Do the switch only after
  access is granted, otherwise the next push fails the same way.

## Incidents / lessons

- **2026-11-11 — tag pushed before any backend gate ran.** `v0.27.1-fork.2` was pushed right after
  the upstream merge, before `go test` had run anywhere. Same class of mistake as HouseOwl's
  "tests never ran", caught before it shipped a broken server rather than after.
  Correct order from now on: merge → go test → tag.
- **2026-11-11 — tag push produced no Actions run.** No workflow run was created for commit
  `ce9fcad4` (or the tag), and GHCR only had `0.27.0-fork.1` + `latest`, so no `v0.27.1` image
  existed. This was a GitHub-side delay, not a trigger bug: the `main` push at the same commit also
  produced no run even though it matched the `paths` filter. (For the record, the docs state
  "Path filters are not evaluated for pushes of tags", so path filters were never the cause and were
  kept for branch pushes.)
- **2026-11-11 — GHCR push denied.** Once runs started, both arch builds failed with
  `failed to push ghcr.io/antebudimir/homebox: denied: permission_denied: write_package`. Cause and
  fix: see "GHCR auth" above.
- **2026-11-11 — the gate paid off immediately.** `go test ./...` failed to compile
  `repo_builds_test.go` after the upstream merge changed `GroupDelete` to return
  `([]uuid.UUID, error)`. Exactly the class of breakage the gate exists to catch.

