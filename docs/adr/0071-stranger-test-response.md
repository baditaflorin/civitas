# 0071 - Stranger-Test Findings and Response

## Status

Accepted

## Context

Phase 3 requires a fresh-user workflow test. No other human is available inside this autonomous run.

## Decision

Run a me-as-stranger test in a fresh Playwright browser context with local storage cleared and real fixture input. Record confusion in `docs/phase3/stranger-test.md` and fix the top three issues before completion.

Expected issues to watch:

- Can the user tell whether a backend is needed?
- Can the user create a case and load data without docs?
- Can the user take output away?
- Can the user recover/reset local UI state?

## Consequences

This is not as strong as a real external human test, but it is better than skipping the usability pass.

## Alternatives Considered

Skipping the test was rejected by Phase 3 constraints.
