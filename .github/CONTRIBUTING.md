# Contributing

## Local workflow

Use Docker for Go compilation and tests; do not install or run Go directly on the host.

```sh
make dev
make test
make build
python -m pip install -r docs/requirements.txt
python -m mkdocs build --strict
```

Run `./tests/run_tests.sh` for the Docker binary smoke tier. Browser beta checks are run explicitly with `./tests/e2e.sh`.

Create a topic branch from `main`, make focused changes, and open a pull request. Include tests for behavior changes and update README, API reference, or operator docs when their descriptions change. Keep changes to AI.md out of feature pull requests; it is the read-only implementation source for this repository.

## Security reports

Do not open a public issue for a vulnerability. Use [GitHub private vulnerability reporting](https://github.com/apimgr/jokes/security/advisories/new). See [the security policy](SECURITY.md) for the disclosure process.
