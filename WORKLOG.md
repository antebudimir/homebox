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
- `docker-publish-rootless.yaml` was cleaned up accordingly (Docker Hub env, logins, push and
  attest steps removed; upstream `binfmt`/`buildkit` images replaced with action defaults).

## Incidents / lessons

- **2026-11-11 — tag pushed before any backend gate ran.** `v0.27.1-fork.2` was pushed right after
  the upstream merge, before `go test` had run anywhere. Same class of mistake as HouseOwl's
  "tests never ran", caught before it shipped a broken server rather than after.
  Correct order from now on: merge → go test → tag.
- **2026-11-11 — tag push produced no Actions run.** No workflow run was created for commit
  `ce9fcad4` (or the tag), and GHCR only had `0.27.0-fork.1` + `latest`, so no `v0.27.1` image
  existed. Root causes addressed:
  - the push trigger used path filters that can also gate tag pushes (a tag on an already-pushed
    commit has no changed files) — path filters removed from `push`;
  - `workflow_dispatch` added so publishing can be started manually;
  - `concurrency: cancel-in-progress` added so nightly and tag/main runs cannot race the same tags.
  If pushes still produce no runs, check **Settings → Actions → General** ("Disable actions")
  and the Actions tab banner; re-enable, then run the workflow manually.
