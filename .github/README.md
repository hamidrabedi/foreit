# .github — repository automation & agent hints

This folder is used for GitHub-specific automation and guidance for external tools.

Key files added for agents:

- `../AGENTS.md` — how to access the repo via GitHub API and `gh`.
- `../.github/copilot-instructions.md` — Copilot-specific guidance (if present).

If automated agents cannot read repository files via the API, ensure they are given a PAT with `repo` (private) or `public_repo` (public) scope, or use a GitHub App with read-only content permissions.
