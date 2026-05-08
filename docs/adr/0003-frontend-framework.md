# 0003 - Frontend Framework and Build Tooling

## Status

Accepted

## Context

The frontend must be static, GitHub Pages-compatible, accessible, typed, and small enough to load quickly.

## Decision

Use Vite, React, TypeScript strict mode, Tailwind CSS, TanStack Query, Zod, Lucide React, and Playwright for smoke/e2e checks.

Build output goes to `docs/`, with `base: "/civitas/"`, hashed assets, and a `404.html` SPA fallback.

## Consequences

- The public app is easy to publish from `main` `/docs`.
- UI state and data fetching have proven libraries.
- Initial JavaScript must stay below the 200 KB gzip target by avoiding eager heavyweight modules.

## Alternatives Considered

- Plain HTML: rejected because the dashboard needs stateful controls, API fetching, and typed data models.
- Next.js: rejected because static Pages output would add complexity for little v1 benefit.
