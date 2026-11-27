# CLI client

The `jokes-cli` binary retrieves public jokes from a Jokes API server and writes results for terminal or shell-pipeline use. It does not require a token because the joke retrieval API is public.

```bash
jokes-cli --help
jokes-cli --server "$JOKES_SERVER_URL" --search programming --count 5
jokes-cli --server "$JOKES_SERVER_URL" --category dadjokes
jokes-cli --server "$JOKES_SERVER_URL" --id "$JOKE_ID" --format json
jokes-cli --server "$JOKES_SERVER_URL" --count 3
```

Set `JOKES_SERVER_URL` to the deployed API base URL before running these examples. The client also accepts `JOKES_SERVER_PRIMARY` as the default endpoint. The client does not assume an official endpoint. HTTPS is required except for localhost.

Set `JOKE_ID` to an opaque UUID copied from the `id` field in a joke API response. Numeric IDs in the embedded dataset are internal and are not accepted by public lookup routes.

## Flags

| Flag | Description |
|---|---|
| `--server URL` | API base URL; defaults to the configured endpoint or `JOKES_SERVER_PRIMARY` |
| `--config PROFILE` | Use a CLI configuration profile or a configuration file path |
| `--token TOKEN` | Bearer token for owner-scoped operations |
| `--token-file PATH` | Read the Bearer token from a file |
| `--category NAME` | Restrict retrieval to one category |
| `--search TEXT` | Search joke text and categories |
| `--id UUID` | Retrieve a joke by public UUID |
| `--count N` | Number of random jokes, from 1 to 100 |
| `--limit N`, `--page N` | Page size and page for collection/search requests |
| `--format FORMAT`, `--output FORMAT` | Output as `table`, `json`, `yaml`, `plain`/`text`, or `csv` |
| `--lang CODE` | Select a supported CLI language |
| `--color auto|yes|no` | Color preference; `NO_COLOR` is honored in auto mode |
| `--debug[=BOOL]` | Enable request diagnostics on stderr |
| `--verbose[=BOOL]` | Enable verbose request diagnostics on stderr |
| `--update check|yes` | Check for a CLI update or install an available update |
| `--shell completions|init [SHELL]` | Print shell completions or shell initialization for bash, zsh, fish, sh, dash, ksh, or PowerShell |
| `--version` (`-v`), `--help` (`-h`) | Print version or usage information |

When both a search, category, and count are provided, the client sends the search request with category and result limit. Otherwise, count greater than one selects random collection retrieval. Use `--help` on the installed binary for the authoritative flags for that release.

The client creates `cli.yml` under its per-user configuration directory on first run. On Unix-like systems the default is `~/.config/apimgr/jokes/cli.yml`; on Windows it is `%APPDATA%\apimgr\jokes\cli.yml`. The `--config` flag accepts an absolute path, a path relative to the configuration directory, or a profile name. Configuration files may set the server URL/API version/timeout/retries, auth token or token-file, output format/color/pager, theme and TUI options, logging, cache, display mode, update policy, and default language/output/limit. Keep token-bearing profiles private; the CLI writes them with owner-only permissions on Unix-like systems.

The CLI checks `/api/autodiscover` on each invocation. `update.auto` defaults to `false`, so an interactive run asks before installing an available update and a non-interactive run leaves it alone. Set `update.auto: true` to allow a non-interactive invocation to install an available update without a prompt. `jokes-cli --update check` checks without installing, and `jokes-cli --update yes` explicitly installs an available version. A server-enforced `cli_min_version` blocks ordinary API requests when the installed client is too old.

Environment and terminal overrides:

| Variable | Effect |
|---|---|
| `JOKES_SERVER_PRIMARY` | Default API server URL |
| `JOKES_TOKEN` | Default Bearer token for owner-scoped operations |
| `JOKES_OUTPUT_FORMAT` | Default output format |
| `JOKES_VERBOSE`, `JOKES_DEBUG` | Enable verbose diagnostics or request debugging |
| `JOKES_COLOR`, `NO_COLOR` | Set color behavior; `NO_COLOR` disables automatic color |
| `LC_ALL`, `LC_MESSAGES`, `LANG` | Auto-detect the CLI language after `--lang` and `defaults.lang` in `cli.yml`; environment precedence is `LC_ALL`, then `LC_MESSAGES`, then `LANG`, then English |
| `PAGER` | Select the pager used for long output |

The client sends its resolved language in `Accept-Language` on API requests. Unsupported language values silently fall back to English.

## Server operator commands

These commands run on the `jokes` server binary, not `jokes-cli`:

```bash
jokes --baseurl /myproject
jokes --update check
jokes --update branch devel
jokes --update yes
jokes --maintenance update check
jokes --maintenance setup
jokes --maintenance mode production
jokes --maintenance backup [filename] [--include-data] [--include-ssl]
jokes --maintenance backup test [file]
jokes --maintenance email test --dry-run
jokes --maintenance email test [--to ADDRESS]
jokes --maintenance restore /path/to/backup.tar.gz
jokes --maintenance pgp generate
jokes --maintenance pgp rotate
jokes --maintenance pgp publish
jokes --maintenance pgp export public [file]
jokes --maintenance pgp export private <file>
jokes --maintenance pgp import <file>
jokes --maintenance pgp delete
jokes --maintenance secret rotate installation_secret
jokes --maintenance secret rotate encryption_key
jokes --maintenance security list
jokes --maintenance security show <id>
jokes --maintenance security export <id>
jokes --maintenance security update <id> [--status STATUS] [--severity LEVEL] [--comment TEXT] [--disclosure-date YYYY-MM-DD]
jokes --maintenance token list
jokes --maintenance token revoke <prefix> [--yes]
jokes --maintenance data export --owner <prefix> [--output FILE]
jokes --maintenance data delete --owner <prefix> [--yes]
jokes --maintenance compliance report
jokes --maintenance scheduler list
jokes --maintenance scheduler show <task-id>
jokes --maintenance scheduler run <task-id>
jokes --maintenance scheduler enable <task-id>
jokes --maintenance scheduler disable <task-id>
jokes --maintenance scheduler history <task-id>
jokes --maintenance tor status
jokes --maintenance tor validate
jokes --maintenance tor restart
jokes --maintenance tor regenerate --yes
jokes --maintenance tor vanity start --prefix abcd [--workers N]
jokes --maintenance tor vanity stop
jokes --maintenance tor vanity apply [address] --yes
jokes --maintenance tor import-keys /absolute/path/to/keys --yes
```

`--baseurl` sets the URL path prefix for a subpath deployment; trusted forwarded prefix headers take precedence. The update checker follows the configured `server.update.branch`; `--update branch` changes that setting. Installing an update requires an elevated process. Backups include the configuration directory and an existing local SQLite database, and may be encrypted with `server.backup.encryption_password`. `backup test` checks archive integrity and restores into a temporary directory without changing live data. Security-report commands list, inspect, export, and update report records. Tor control commands contact the running server over its private loopback endpoint; destructive commands require `--yes`.

`--maintenance update` accepts the same `check`, `yes`, and `branch <name>` subcommands as `--update`. `--maintenance setup` is first-run only unless run as root; `--maintenance mode` accepts `production`, `development`, or `debug` and applies after restart. PGP key management, secret rotation, token administration, owner data export/deletion, and compliance reporting require the operator token or root. PGP private-key export/import/deletion and secret rotation require typed confirmation; commands authorized with `server.token` also prompt for that token again, while root runs skip the token prompt. `token revoke` accepts a linked-token prefix or primary-token owner ID and requires confirmation unless `--yes` is supplied. Secret rotation supports only `installation_secret` and `encryption_key`; the previous installation secret is retained for seven days and the previous encryption key for 30 days while encrypted data is migrated. Data export writes a mode-0600 JSON file by default; `--output -` writes it to stdout.

Scheduler commands require the operator token or root. They list task state, show details, run a registered task immediately, enable or disable a task by updating `yaml.d/maintenance.yml`, and read recent task history. Scheduler state has no HTTP API.

`--maintenance email test --dry-run` renders every built-in or overridden email template with sample data, prints each subject and body, and reports all invalid overrides without requiring SMTP. Validation errors exit with status `65`. `--maintenance email test` sends a real test message after a successful SMTP connection check; it uses `server.notifications.email.reply_to` as the recipient unless `--to ADDRESS` is supplied, prefixes the subject with `[TEST]`, and records the send in the audit log.

Restore requires an explicit backup-file path. Before changing files, the command verifies the archive, prompts securely for an encrypted archive's password, and warns when the backup was created by another application version. A root run against an existing local database asks for confirmation; an empty local SQLite install can restore without that confirmation. The configured service account can restore an existing database after securely entering the operator token. Other unprivileged users are denied. An empty remote libSQL database can be seeded from the verified SQLite snapshot; a configured local cache can also receive the restored database file. Restart the service after a successful restore so it loads the restored configuration. Scheduled backup tasks apply retention automatically.
