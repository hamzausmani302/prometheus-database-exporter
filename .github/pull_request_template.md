## Description

<!-- What does this PR do? Why is the change needed? -->

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update
- [ ] Refactoring / code cleanup

## Changes Made

<!-- Brief bullet points of what changed -->

## Testing Done

- [ ] Unit tests pass: `go test ./...`
- [ ] No vet errors: `go vet ./...`
- [ ] Lint passes: `golangci-lint run`
- [ ] Integration tests pass (if applicable): `make test-integration`
- [ ] Manually verified against a running database

## Checklist

- [ ] No `fmt.Println` or `log.Fatal` calls in library code (`internal/`, `pkg/`)
- [ ] No real passwords or connection strings committed
- [ ] New datasource/store/scheduler types are registered in their factory
- [ ] Documentation updated if behaviour changed
