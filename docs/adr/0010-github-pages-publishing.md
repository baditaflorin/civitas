# 0010 - GitHub Pages Publishing Strategy

## Status

Accepted

## Context

The live Pages URL is a first-class deliverable from the first commit. The frontend must be published without GitHub Actions.

## Decision

Publish GitHub Pages from the `main` branch `/docs` folder at:

https://baditaflorin.github.io/civitas/

The Vite build writes directly to `docs/`, and `docs/` is intentionally not gitignored. The Vite base path is `/civitas/`. Asset names are hashed for cache busting. `docs/404.html` mirrors the built SPA fallback. No custom domain is configured in v1.

## Consequences

- A push to `main` updates the public site after GitHub Pages rebuilds.
- Build artifacts are committed by design.
- Rollback is a normal git revert of the Pages publishing commit.

## Alternatives Considered

- `gh-pages` branch: rejected because it adds branch management overhead.
- GitHub Actions deployment: rejected because the project explicitly uses local hooks instead of Actions.
