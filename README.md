# pocketsmith-moneyforward

CLI for syncing MoneyForward accounts into Pocketsmith.

## Run locally
- `MONEYFORWARD_COOKIE=... POCKETSMITH_TOKEN=... mise run run`
- `mise run mf-test` to verify MoneyForward connectivity.

## Build
- `mise run build` to compile the module.

## Container image
- Image: `ghcr.io/dvcrn/pocketsmith-moneyforward`
- Example: `docker pull ghcr.io/dvcrn/pocketsmith-moneyforward:latest`

## Docker publish
- Use `mise run docker-build` to build and push to GHCR.
- GitHub Actions only publishes when the workflow is manually triggered.
