# CLI Interface Contract: nato

**Version**: 1.0.0

## Command Signature

```
nato [flags] [text...]
```

## Input Sources

| Priority | Source | Behavior |
|----------|--------|----------|
| 1 | Command-line arguments | `nato hello world` → converts "hello world" |
| 2 | Standard input (pipe) | `echo "hello" \| nato` → converts piped text |
| 3 | Neither | Shows help message |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--alphabet` | `-a` | string | `nato` | Phonetic alphabet to use |
| `--no-color` | | bool | `false` | Disable colored output |
| `--list-alphabets` | `-l` | bool | `false` | List available alphabets and exit |
| `--version` | `-v` | bool | `false` | Print version and exit |
| `--help` | `-h` | bool | `false` | Print help and exit |

## Output Format

### Standard Conversion (TTY)

When outputting to a terminal with color support:

```
<colored_word> <colored_word> ...
```

Each input word produces one line of output. Phonetic words are space-separated.
Colors alternate for visual distinction.

### Standard Conversion (Non-TTY or --no-color)

Plain text without ANSI escape codes:

```
Word1 Word2 Word3
Word4 Word5
```

### List Alphabets Output

```
Available alphabets:
  <name>        - <description>
  <name>        - <description>
  ...
```

### Version Output

```
nato version <semver>
```

### Help Output

Standard Cobra help format with usage, flags, and examples.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Invalid alphabet name |
| 1 | Other error |

## Environment Variables

| Variable | Effect |
|----------|--------|
| `NO_COLOR` | When set (any value), disables colored output |
| `TERM` | Used for TTY detection (handled by go-isatty) |

## Examples

### Basic Usage

```bash
$ nato hello
Hotel Echo Lima Lima Oscar

$ nato "hello world"
Hotel Echo Lima Lima Oscar
Whiskey Oscar Romeo Lima Delta
```

### With Flags

```bash
$ nato -a lapd hello
Adam David Adam Mary

$ nato --no-color hello
Hotel Echo Lima Lima Oscar

$ nato -l
Available alphabets:
  nato          - NATO/ICAO phonetic alphabet (international standard)
  lapd          - LAPD/APCO alphabet (US law enforcement)
  western-union - Western Union alphabet (telegraphy)
  german        - German spelling alphabet (traditional)

$ nato -v
nato version 1.0.0
```

### Piped Input

```bash
$ echo "test" | nato
Tango Echo Sierra Tango
```

### Error Cases

```bash
$ nato -a invalid hello
Error: unknown alphabet "invalid"

Available alphabets: nato, lapd, western-union, german
exit status 1
```

## Compatibility Notes

- Windows: Requires Windows 10 or newer for ANSI color support
- Older terminals: Colors auto-disabled when TERM indicates no support
- CI environments: Typically no TTY, colors auto-disabled
