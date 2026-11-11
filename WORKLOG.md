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
- `publish-image.yaml`'s gate is wired but not yet exercised by CI (tag/manual-trigger only) — it
  proves itself on the next tag.
- `Publish Release Binaries` (SLSA) fails on tags; it is upstream release plumbing this fork does not
  use. Candidate for disabling at the next housekeeping pass.

## GHCR auth

- GHCR pushes authenticate with `GITHUB_TOKEN`. This requires the package
  `ghcr.io/antebudimir/homebox` to grant this repository Actions write access
  (package settings → **Manage Actions access** → add `antebudimir/homebox` with **Write**).
- History: the package was originally created by a PAT push (`publish-image.yaml` used
  `secrets.CR_PAT`), so it was not linked to the repo and `GITHUB_TOKEN` pushes failed with
  `denied: permission_denied: write_package`. Access was granted on 2026-11-11 and the four image
  workflows were switched to `GITHUB_TOKEN`; `CR_PAT` is no longer used by CI.
- If a push fails with `permission_denied: write_package`, that link is missing again.

## Incidents / lessons

- **2026-11-11 — tag pushed before any backend gate ran.** `v0.27.1-fork.2` was pushed right after
  the upstream merge, before `go test` had run anywhere. Same class of mistake as HouseOwl's
  "tests never ran", caught before it shipped a broken server rather than after.
  Correct order from now on: merge → go test → tag.
- **2026-11-11 — tag push produced no Actions run.** No workflow run was created for commit
  `ce9fcad4` (or the tag), and GHCR only had `0.27.0-fork.1` + `latest`, so no `v0.27.1` image
  existed. This was a GitHub-side delay, not a trigger bug: the `main` push at the same commit also
  produced no run even though it matched the `paths` filter.
- **Verify tag runs actually start.** After pushing a tag, check the Actions tab — both `v0.27.1`
  tags initially produced 0 runs (GitHub-side delay; path filters are NOT evaluated for tag pushes
  per GitHub docs). If no run appears, re-run or re-tag rather than assuming it fired.
- **2026-11-11 — GHCR push denied.** Once runs started, both arch builds failed with
  `failed to push ghcr.io/antebudimir/homebox: denied: permission_denied: write_package`. Cause and
  fix: see "GHCR auth" above.
- **2026-11-11 — the gate paid off immediately.** `go test ./...` failed to compile
  `repo_builds_test.go` after the upstream merge changed `GroupDelete` to return
  `([]uuid.UUID, error)`. Exactly the class of breakage the gate exists to catch.

