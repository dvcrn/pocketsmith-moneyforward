- Repo: dvcrn/pocketsmith-moneyforward

# Repository Guidelines

This repository contains a small Go CLI that syncs MoneyForward accounts into Pocketsmith.

## Project Structure & Module Organization
- `main.go` holds the CLI entry point and sync flow.
- `sanitizer/` is a standalone package for payee normalization.
- `sanitizer/*_test.go` contains unit tests for sanitizer behavior.
- `Dockerfile` and `.github/workflows/docker-publish.yml` define container build and publish automation.

## Build, Test, and Development Commands
Prefer `mise` tasks (see `mise.toml`):
- `mise run build` builds the module.
- `mise run test` runs all Go tests.
- `mise run run` runs the CLI; pass config via flags or env:
  `MONEYFORWARD_COOKIE=... POCKETSMITH_TOKEN=... mise run run`.
- `mise run docker-build` builds and pushes the multi-arch image (requires Docker auth).

## Coding Style & Naming Conventions
- Use `gofmt` output (tabs, standard Go layout).
- Exported identifiers use `PascalCase`; locals use `camelCase`.
- Test files use the `*_test.go` convention with table-driven tests where practical.

## Testing Guidelines
- Tests use Go's `testing` package.
- Run `mise run test` or `go test ./...` before changes to sanitizer logic.

## Configuration & Secrets
- Required config: `MONEYFORWARD_COOKIE` and `POCKETSMITH_TOKEN` (or `-mf-cookie` / `-pocketsmith-token`).
- Do not commit secrets; store local values in `fnox.toml` when using fnox.

## Commit & Pull Request Guidelines
- Commit subjects in this repo are short, imperative, and capitalized (e.g., "Add Docker publish GitHub Action").
- PRs should include a concise summary, testing notes, and highlight any Docker publish impact (pushes to `main` build and push images).
