# Upstream Sync

- This is a hard fork of Homebox (upstream: `sysadminsmedia/homebox`). Sync strategy: **merge, not cherry-pick** — the fork tracks upstream closely and adds features on top.
- Never use the GitHub "Sync fork" button — it creates parallel merge commits that fight the local workflow. Use `git fetch upstream && git merge upstream/main` locally.
- Fork features to protect on every merge: builds system (`repo_builds`, `build_positions`, usage sessions + schemas), sport tagging (`Sport` field, `DefaultSport` template inheritance), `ShowBuilds` kill-switch, `ensureBuildPartTemplates` group-creation seeding, theme toggle.
- Verify after merge: `git diff upstream/main HEAD` on overlapping files must show only fork additions; typecheck frontend; run Go tests. Upstream refactors (e.g. transactional `RemoveMember`/`DeleteGroup`) may sit next to fork hooks — verify semantically, not just textually. Real example: upstream's refactor changed `GroupDelete`'s signature and broke `repo_builds_test.go` — caught only because the test gate ran.
- Tag last, verify first. Tagging (`v<upstream>-fork.N`) triggers the image publish. Order: merge → go test → pnpm typecheck → then tag and push. Never tag on unverified backend code. Recorded in `WORKLOG.md` — cross-reference, don't duplicate.
- **Verify tag runs actually start.** After pushing a tag, check the Actions tab — both `v0.27.1` tags initially produced 0 runs (GitHub-side delay; path filters are NOT evaluated for tag pushes). If no run appears, re-run or re-tag rather than assuming it fired.
- Workflow policy: GHCR-only — no Docker Hub (`sysadminsmedia/homebox` is upstream's), no upstream devcache namespace. All four publish workflows follow this policy.
- Tag convention: `vX.Y.Z-fork.N` — upstream version + fork counter.
