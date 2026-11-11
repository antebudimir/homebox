# Git

- `main` is the only long-lived branch and must always be deployable.
- Never force-push `main`; never rewrite history that is already pushed.
- Take a backup branch before risky history work: `git branch backup/<what>-$(date +%F)` (e.g. `backup/pre-upstream-sync-2026-11-11`).
- Keep the backup branch until CI is green and the smoke test passes; delete it deliberately, not by accident.
- Run git write commands serially; never batch `git add`/`commit`/`rebase` in parallel.
- After any history rewrite or multi-command git operation, verify actual repo state (`git log`, `git show`) before acting on cached status output.
- Upstream work is merged, never cherry-picked — see `08-upstream-sync.md`. Never use GitHub's "Sync fork" button.
- Fork commit subjects use Conventional Commits (`feat:`, `fix(scope):`, `chore:`, `docs:`, `test:`, `ci:`). Upstream subjects stay verbatim.
- Use `git --no-pager` / `GIT_PAGER=cat` in scripts and agent commands.
- Quote rev-range arguments: zsh expands an unquoted `^` as glob negation, so `upstream/main ^HEAD` silently becomes `upstream/main`. Use `'^HEAD'` or `HEAD..upstream/main`.
- Never commit generated files, secrets, build output or a `.db` file.
- Tags (`vX.Y.Z-fork.N`) trigger the image publish: only tag a commit that passed CI.
- Push the single tag (`git push origin vX.Y.Z-fork.N`), never `--tags` — that would push upstream's tags into the fork.
