# tycs

Browse the [Teach Yourself CS](https://teachyourselfcs.com/) curriculum from the command line.

`tycs` is a single pure-Go binary. No API key required.

## Install

```bash
go install github.com/tamnd/tycs-cli/cmd/tycs@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/tycs-cli/releases), or run the container image:

```bash
docker run --rm ghcr.io/tamnd/tycs:latest --help
```

## Usage

```bash
# List all 9 CS subject guides
tycs subjects

# Table output
tycs subjects -o table

# JSON output
tycs subjects -o json

# Get subject URLs only
tycs subjects -o url

# Limit results
tycs subjects -n 3
```

## Commands

| Command | Description |
|---------|-------------|
| `subjects` | List all CS subject guides from Teach Yourself CS |
| `version` | Show version information |

## Global flags

```
-o, --output string    output format: table|json|jsonl|csv|tsv|url|raw (default "auto")
-n, --limit int        limit number of records (0 = all)
    --fields strings   comma-separated columns to include
    --no-header        omit header row
    --template string  Go text/template per record
    --timeout duration per-request timeout (default 30s)
    --delay duration   minimum spacing between requests
    --retries int      retry attempts on 429/5xx (default 3)
```

## License

Apache-2.0. See [LICENSE](LICENSE).
