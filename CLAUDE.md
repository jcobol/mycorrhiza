# CLAUDE.md

This repo (`jcobol/mycorrhiza`) is a fork of `bouncepaw/mycorrhiza`.

**Never reference, target, or interact with the upstream `bouncepaw/mycorrhiza` repo** — not in PRs, issues, branches, commit messages, or `gh`/`git` remote operations. Upstream does not accept AI-authored or AI-assisted submissions. All work here stays local to this fork: PRs and pushes go to `origin` (`jcobol/mycorrhiza`) only.

**Pitfall:** GitHub's own "Create a pull request" link/UI (the one `git push` prints, and the compare page it opens) defaults the base repository to the fork's parent (`bouncepaw/mycorrhiza`) instead of this fork, even when the URL is under `jcobol/mycorrhiza`. This already happened once (bouncepaw/mycorrhiza#274, since closed). Always create PRs with an explicit repo, e.g. `gh pr create --repo jcobol/mycorrhiza --base master ...`, and never just follow the printed link or accept the compare page's default base without checking it.
