# go-cache Contributor Guide

`go-cache` is a Go 1.25 in-memory cache library (`github.com/pzentenoe/go-cache`). Keep changes small, preserve its public API, and add behavior-focused tests with each change.

## Project Map

| Area | Purpose |
|---|---|
| `cache.go`, `item.go` | Single-lock cache, expiration, callbacks, and snapshots. |
| `increment.go`, `decrement.go`, `numeric.go`, `constants.go` | Typed numeric mutations and overflow handling (shared generic helpers live in `numeric.go`). |
| `janitor.go` | Background expiration cleanup and lifecycle controls (single janitor shared by `Cache` and `shardedCache`). |
| `serialization.go` | Gob reader/writer and file persistence for `Cache`. |
| `sharded_cache*.go` | Sharded implementation, its interface, hashing, persistence, and janitor. |
| `*_test.go` | Same-package unit, concurrency, lifecycle, serialization, and benchmark coverage. |
| `examples/`, `docs/`, `README.md` | Consumer documentation and runnable examples. |

## Code Discovery: Mandatory

Use `codebase-memory-mcp` first for all structural code discovery, exact symbols, callers/callees, risk analysis, and dead-code checks:

1. `list_projects` and `index_status`
2. `search_graph`
3. `trace_path`
4. `get_code_snippet`
5. `query_graph`

Fall back to `grep` only for string literals, configuration, or non-code files. Current graph project ID: `Users-pzentenoe-Documents-projects-vikingcode-go-cache`.

Refresh the graph after structural source changes so symbol and call-path results remain current:

```text
codebase-memory-mcp_index_repository(
  repo_path="/Users/pzentenoe/Documents/projects/vikingcode/go-cache",
  mode="moderate",
  persistence=true,
)
```

The persistence artifact is `.codebase-memory/graph.db.zst`. Do not add it to a commit unless repository policy explicitly requires it.

## Validation

Run the narrowest relevant check first, then the repository commands below as applicable:

```bash
go test ./...
go test -race ./...
go test -cover ./...
golangci-lint run ./...
make test
```

`make check` installs tools and runs `go fmt ./...` plus `misspell -w .`; it can modify the worktree. Use it only when those edits are intended.

CI uses Go 1.25, runs tests with atomic coverage excluding `examples/`, and runs `golangci-lint run ./...`.

## Releases

Releases are automated with GoReleaser (`.goreleaser.yml`, library mode: no binaries). To cut a release: update `CHANGELOG.md` (Keep a Changelog style), commit, then push a lightweight tag — `git tag vX.Y.Z && git push origin vX.Y.Z`. The `release` workflow runs the race test suite as a gate and publishes the GitHub Release with a changelog grouped by conventional-commit type. Never move or delete a published tag: Go module proxies cache them immutably.

## Conventions

- Keep production Go files in package `cache`; tests currently use the same package and Testify assertions.
- Guard shared maps with the existing `sync.RWMutex` pattern. Invoke user callbacks outside locks, but capture any callback while synchronized first.
- Preserve expiration semantics: `DefaultExpiration` is `0`, `NoExpiration` is `-1`, and expired entries may remain until cleanup.
- Public APIs use `any`, `time.Duration`, and Go-style error returns. Document exported identifiers with concise comments.
- Do not add abstractions or dependencies without a demonstrated need. Profile before performance changes; use benchmarks with allocation reporting for hot paths.
- Prefer deterministic lifecycle/concurrency tests over fixed sleeps where practical; run race detection for synchronization changes.
