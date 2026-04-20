# AGENTS.md

## Build & Test Commands

- **Full CI**: `make ci` - runs verify, check-license, lint, build, unit-tests
- **Build**: `make build` - builds binary via Docker (uses Go 1.25)
- **Test**: `make test` - runs unit tests with race detector
- **Lint**: `make lint` - runs golangci-lint
- **Format**: `make fmt` - formats code with goimports/gofmt
- **Verify**: `make verify` - checks code generation and module sync

## Dependencies

- Uses vendored dependencies: always use `GOFLAGS="-mod=vendor"` or `-mod=vendor` flag
- Go version: 1.25

## Products

This repo automates releases for multiple AppsCode products: `voyager`, `kubevault`, `kubedb`, `ace`, `stash`, `kubestash`.

- **Bump minor version**: `make bump-release-minor PRODUCT=<product>` or `./hack/scripts/bump-release-minor.sh <product>`

## CLI Commands

The binary supports subcommands: `release`, `ace`, `kubedb`, `kubestash`, `kubevault`, `stash`, `virtualsecrets`, `voyager`, `list-versions`, `update-assets`, `update-bundles`, `update-envvars`.

## CI

- GitHub Actions uses Go 1.25 on `ubuntu-24.04`
- Runs `make ci` on PRs and pushes to master
