# Quickstart: NATO Phonetic Alphabet CLI

## Installation

### Via Homebrew (Recommended)

```bash
brew install dreamiurg/tap/nato
```

### Via Go Install

```bash
go install github.com/dreamiurg/nato@latest
```

### From Binary

Download the latest release from [GitHub Releases](https://github.com/dreamiurg/nato/releases) and add to your PATH.

## Basic Usage

Convert text to NATO phonetic alphabet:

```bash
nato hello
# Output: Hotel Echo Lima Lima Oscar

nato "Hello World"
# Output:
# Hotel Echo Lima Lima Oscar
# Whiskey Oscar Romeo Lima Delta
```

## Options

| Flag | Short | Description |
|------|-------|-------------|
| `--alphabet` | `-a` | Select phonetic alphabet (default: nato) |
| `--no-color` | | Disable colored output |
| `--list-alphabets` | `-l` | List available alphabets |
| `--version` | `-v` | Show version |
| `--help` | `-h` | Show help |

## Examples

### Different Alphabets

```bash
# NATO (default)
nato -a nato abc
# Alfa Bravo Charlie

# LAPD/Police
nato -a lapd abc
# Adam Boy Charles

# Western Union
nato -a western-union abc
# Adams Boston Chicago

# German
nato -a german abc
# Anton Berta Cäsar
```

### Piped Input

```bash
echo "test" | nato
# Tango Echo Sierra Tango

cat file.txt | nato
```

### Disable Colors

```bash
nato --no-color hello
# Plain text output without ANSI colors

# Colors are automatically disabled when piping
nato hello | cat
```

### List Available Alphabets

```bash
nato --list-alphabets
# Available alphabets:
#   nato         - NATO/ICAO phonetic alphabet (international standard)
#   lapd         - LAPD/APCO alphabet (US law enforcement)
#   western-union - Western Union alphabet (telegraphy)
#   german       - German spelling alphabet (traditional)
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NO_COLOR` | Disable colors when set (any value) |

## Tips

- Colors are automatically disabled when output is piped or redirected
- Use quotes for multi-word input: `nato "multiple words"`
- The tool is case-insensitive: `nato ABC` and `nato abc` produce the same output
- Special characters and punctuation are silently skipped
