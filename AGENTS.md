# Agents / Bots / Automation

This document explains how external AI agents or automation tools can access repository contents and metadata via GitHub APIs or standard Git clients.

Recommended access methods

- Git clone (best for full repo):

  git clone https://github.com/<owner>/<repo>.git

- GitHub REST API (file-level access):

  - Get a file's content (base64-encoded):

    curl -H "Authorization: token $GH_TOKEN" \
      -H "Accept: application/vnd.github.v3+json" \
      https://api.github.com/repos/<owner>/<repo>/contents/<path>

  - List a directory:

    curl -H "Authorization: token $GH_TOKEN" \
      https://api.github.com/repos/<owner>/<repo>/contents/<dir>

- GitHub GraphQL (preferred when querying many files/trees with low rate cost).

Authentication & scopes

- Use a Personal Access Token (PAT) or GitHub App credential. For private repos, include `repo` scope. For public repos, `public_repo` is sufficient for read-only access.
- For Actions or Apps, use the automatically provided token with appropriate permissions.

Practical notes for AI agents

- The REST content API returns files base64-encoded; decode before parsing.
- Large files or many requests may hit rate limits; prefer git clone or GraphQL to reduce REST calls.
- Use `gh` CLI for convenience: `gh repo clone owner/repo` or `gh api repos/owner/repo/contents/path`.

Example: fetch `forgego/pkg/admin/register.go` via REST

  curl -s -H "Authorization: token $GH_TOKEN" \
    https://api.github.com/repos/hamidrabedi/foreit/contents/forgego/pkg/admin/register.go | jq -r .content | base64 --decode

If you (or CI) need help enabling API access for this repository, grant the agent a PAT or create a minimal GitHub App with read-only repo content access.
