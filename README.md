# nato

Convert text to NATO phonetic alphabet.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install dreamiurg/tap/nato
```

### Go Install

```bash
go install github.com/dreamiurg/nato@latest
```

### Binary

Download from [Releases](https://github.com/dreamiurg/nato/releases).

## Usage

```bash
nato hello
# Hotel Echo Lima Lima Oscar

nato "Hello World"
# Hotel Echo Lima Lima Oscar
# Whiskey Oscar Romeo Lima Delta

nato ABC123
# Alfa Bravo Charlie One Two Three
```

### Piped Input

```bash
echo "test" | nato
# Tango Echo Sierra Tango
```

### Alternative Alphabets

```bash
nato -a lapd hello
# Henry Edward Lincoln Lincoln Ocean

nato -a german hello
# Heinrich Emil Ludwig Ludwig Otto

nato --list-alphabets
# german         - German spelling alphabet (traditional)
# lapd           - US law enforcement phonetic alphabet
# nato           - International standard phonetic alphabet
# western-union  - Telegraphy phonetic alphabet
```

### Options

| Flag | Short | Description |
|------|-------|-------------|
| `--alphabet` | `-a` | Phonetic alphabet (default: nato) |
| `--no-color` | | Disable colored output |
| `--list-alphabets` | `-l` | List available alphabets |
| `--version` | `-v` | Show version |
| `--help` | `-h` | Show help |

## Alphabets

- **NATO/ICAO**: International standard (Alfa, Bravo, Charlie...)
- **LAPD**: US law enforcement (Adam, Boy, Charles...)
- **Western Union**: Telegraphy (Adams, Boston, Chicago...)
- **German**: Traditional spelling (Anton, Berta, Cäsar...)

## License

MIT
