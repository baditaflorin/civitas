# 0016 - Local Git Hooks

## Status

Accepted

## Context

The project does not use GitHub Actions, so quality gates must run locally.

## Decision

Use plain `.githooks/` wired with `core.hooksPath` via `make install-hooks`.

Hooks:

- `pre-commit`: format/lint/typecheck where tools are installed, plus gitleaks if available.
- `commit-msg`: Conventional Commits validator.
- `pre-push`: `make test`, `make build`, `make smoke`.
- `post-merge` and `post-checkout`: regenerate API types if dependencies are installed.

## Consequences

- Hooks are transparent shell scripts and can be run manually.
- Missing optional tools fail clearly when required by the hook.
- Contributors must opt in once with `make install-hooks`.

## Alternatives Considered

- Lefthook: rejected for v1 to keep bootstrapping minimal.
- No hooks: rejected by project requirements.
