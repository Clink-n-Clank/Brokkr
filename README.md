# Brokkr
![Coverage](https://img.shields.io/badge/Coverage-86.3%25-brightgreen)
![GitHub Super-Linter](https://github.com/Clink-n-Clank/Brokkr/actions/workflows/lint.yml/badge.svg)

Brokkr micro framework that helps you quickly write simple API's and test it

## Run tests

```bash
go test -v --race $(go list ./... | (grep -v /vendor/) | (grep -v internal/test/bdd/integration_tests)) -covermode=atomic -coverprofile=coverage.out
```

